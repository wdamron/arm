// Package asm implements an ARMv8 (AArch64) instruction assembler in Go.
//
// This library is mostly adapted from the CensoredUsername/dynasm-rs (Rust) project, and is not heavily tested.
// See https://github.com/CensoredUsername/dynasm-rs. SVE instructions are not yet supported.
//
// The Assembler type encodes executable instructions to a code buffer.
//
// Some instructions support label offset arguments, which may be resolved by the Assembler
// and encoded after all label addresses are assigned.
//
// The following are argument types:
//   - [Reg]: integer, SP, SIMD scalar, or SIMD vector register (with optional element index)
//   - [RegList]: list of sequential registers
//   - [Ref]: memory reference with register base, optionally followed by X register or immediate for post-indexing
//   - [RefOffset]: memory reference with register base and immediate offset
//   - [RefPreIndexed]: pre-indexed memory reference with register base and immediate offset
//   - [RefIndexed]: memory index with register base, register index, and optional index modifier
//   - [Imm]: 32-bit immediate integer
//   - [Float]: 32-bit immediate float
//   - [Wide]: 64-bit immediate integer
//   - [Mod]: modifier with optional immediate shift/rotate
//   - [Label]: label reference with optional offset from label address
//   - [Symbol]: constant identifier
package arm

// Assembler encodes executable instructions to a code buffer.
//
// Some instructions support [Label] offset arguments, which may be resolved
// and encoded after all label addresses are assigned.
type Assembler struct {
	Code    []byte   // code buffer indexed by PC
	LabelPC []uint32 // label PC by ID
	Relocs  []Reloc  // label references
	Args    []Arg    // arguments for the current instruction
	Flat    []Flat   // flattened arguments for the current matched instruction
	PC      uint32   // current code offset

	CurrentInst Inst   // current instruction mnemonic, offset into the Patterns array
	Count       uint8  // available encodings for the current instruction
	Idx         int8   // encoding index for the current instruction
	SimdSize    uint8  // SIMD width for the current instruction when applicable
	Opcode      uint32 // opcode (without arguments) for the current matched instruction
	Err         error  // most recent error

	patternLen  uint8  // argument-matcher count for the current instruction
	patsOffset  uint16 // current offset within the Patterns array
	cmdsOffset  uint16 // current offset within the Commands array
	cmdsLen     uint8  // encoding-command count for the current matched instruction
	scratchArgs [6]Arg
	scratchFlat [12]Flat
	pattern     [6]EncOp // current argument-matcher list unpacked from the Patterns array
	cmds        [8]EncOp // current encoding-command list unpacked from the Commands array
}

// Reloc is a [Label] reference deferred for encoding after all relocations are being applied.
// Relocations are used internally, and exposed for debugging.
type Reloc struct {
	InstPC uint32    // instruction with label offset argument
	Op     uint8     // relocation type
	Jump   FlatLabel // label ID with optional offset
}

// EncOp is a matching or encoding operator decoded from the [Patterns] or [Commands] arrays.
// Operators are used internally, and exposed for debugging.
type EncOp struct {
	Op uint8
	X  [3]uint8
}

// Initialize or re-initialize the assembler with a new code buffer, resetting the PC and all state.
func (a *Assembler) Init(mem []byte) {
	a.Code, a.PC, a.LabelPC, a.Relocs, a.Err = mem, 0, nil, nil, nil
	a.CurrentInst = 0
	a.Args = a.scratchArgs[:0]
	a.Flat = a.scratchFlat[:0]
	a.Count = 0
	a.Opcode = 0
	a.Idx = -1
	a.SimdSize = 0
	a.patsOffset = 0
	a.patternLen = 0
	a.cmdsOffset = 0
	a.cmdsLen = 0
}

// NewLabel registers a new label identifier at the current PC. The label may be used as an offset argument,
// and the PC for the label may be reassigned by calling SetLabel at the target PC. Label offset
// arguments must be processed through ApplyRelocations once all labels can be resolved.
func (a *Assembler) NewLabel() Label {
	a.LabelPC = append(a.LabelPC, a.PC)
	return Label{ID: uint32(len(a.LabelPC) - 1)}
}

// SetLabel sets the PC for a label to the current PC.
func (a *Assembler) SetLabel(label Label) { a.LabelPC[label.ID] = a.PC }

// Pattern returns the list of matching operators for the most recent matching iteration, useful for debugging.
func (a *Assembler) Pattern() []EncOp { return a.pattern[:a.patternLen] }

// Commands returns the list of encoding operators for the most recent matching iteration, useful for debugging.
func (a *Assembler) Commands() []EncOp { return a.cmds[:a.cmdsLen] }

// ApplyRelocations patches all instructions containing label offset arguments, with the currently assigned
// PC value for each label.
func (a *Assembler) ApplyRelocations() bool {
	if a.Err != nil {
		return false
	}
	for _, rel := range a.Relocs {
		opcode := dec32(a.Code[rel.InstPC:])
		targetPC := int64(a.LabelPC[rel.Jump.ID]) + int64(rel.Jump.Offset)
		enc, ok := encOffset(rel.Op, targetPC-int64(rel.InstPC))
		if !ok {
			a.Err = ErrInvalidEncoding
			return false
		}
		enc32(a.Code[rel.InstPC:], opcode|enc)
	}
	a.Relocs = nil
	return true
}

// Inst advances to the first matched encoding for inst and args, then writes
// the matched instruction to the code buffer if one was found.
//
// If no matching instruction was found or arguments could not be encoded,
// the call will return false and the Err field will be set.
func (a *Assembler) Inst(inst Inst, args ...Arg) bool {
	if a.Err != nil {
		return false
	}
	if inst == 0 || inst > ZIP2 {
		a.Err = ErrInvalidInst
		return false
	}
	a.CurrentInst = inst
	a.Args = append(a.scratchArgs[:0], args...)
	a.Flat = a.scratchFlat[:0]
	a.SimdSize = 0
	a.cmdsOffset = 0
	a.cmdsLen = 0

	a.patsOffset = uint16(inst)
	a.Count = uint8(Patterns[a.patsOffset])
	a.patsOffset++

	for a.Idx = 0; a.Idx < int8(a.Count); a.Idx++ { // each encoding pattern
		a.patternLen = uint8(Patterns[a.patsOffset])
		a.patsOffset++
		for m := uint8(0); m < a.patternLen; m++ { // each matcher for pattern
			op := Patterns[a.patsOffset]
			a.patsOffset++
			a.pattern[m].Op = op
			xs := MatcherArgCounts[op]
			copy(a.pattern[m].X[:], Patterns[a.patsOffset:a.patsOffset+uint16(xs)])
			a.patsOffset += uint16(xs)
		}

		cmdsOffset := uint16(Patterns[a.patsOffset])<<8 | uint16(Patterns[a.patsOffset+1])
		a.patsOffset += 2

		if !a.matchPattern() {
			continue
		}

		a.cmdsOffset = cmdsOffset
		opcode := Commands[a.cmdsOffset : a.cmdsOffset+4]
		a.Opcode = uint32(opcode[0])<<24 | uint32(opcode[1])<<16 | uint32(opcode[2])<<8 | uint32(opcode[3])
		a.cmdsOffset += 4
		a.cmdsLen = Commands[a.cmdsOffset]
		a.cmdsOffset++
		for i := uint8(0); i < a.cmdsLen; i++ {
			op := Commands[a.cmdsOffset]
			a.cmdsOffset++
			a.cmds[i].Op = op
			xs := CmdArgCounts[op]
			copy(a.cmds[i].X[:], Commands[a.cmdsOffset:a.cmdsOffset+uint16(xs)])
			a.cmdsOffset += uint16(xs)
		}
		if !a.encode() {
			a.Err = ErrInvalidEncoding
			return false
		}
		return true
	}

	a.Err = ErrNoMatch
	return false
}
