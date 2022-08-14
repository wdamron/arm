package arm

// Arg is any instruction argument.
//
// The following are argument types:
//   - Reg: integer, SP, SIMD scalar, or SIMD vector register (with optional element index)
//   - RegList: list of sequential registers
//   - Ref: memory reference with register base, optionally followed by X register or immediate for post-indexing
//   - RefOffset: memory reference with register base and immediate offset
//   - RefPreIndexed: pre-indexed memory reference with register base and immediate offset
//   - RefIndexed: memory index with register base, register index, and optional index modifier
//   - Imm: 32-bit immediate integer
//   - Float: 32-bit immediate float
//   - Wide: 64-bit immediate integer
//   - Mod: modifier with optional immediate shift/rotate
//   - Label: label reference with optional offset from label address
//   - Symbol: constant identifier
type Arg interface {
	arg()
}

// ----------------------------------------------------------------

// Reg is a scalar or vector register argument. Vector registers may include an element specifier.
type Reg struct {
	ID   uint8   // 0-31 integer, SP, or SIMD register
	Type RegType // element size and 34/64/128-bit indicator
	Elem uint8   // vector element index (bitwise complement, zero indicates unset)
}

func (r Reg) arg() {}

// RegList is a register list argument with sequentially numbered registers.
type RegList struct {
	First Reg
	Len   uint8
}

func (r RegList) arg() {}

// ----------------------------------------------------------------

// Ref is a memory reference argument with a register base.
// For post-indexing, a Ref may be followed by an X register or immediate offset.
type Ref struct {
	Base Reg // X|SP
}

func (r Ref) arg() {}

// RefOffset is a memory reference argument with a register base and immediate offset.
type RefOffset struct {
	Base   Reg // X|SP
	Offset int32
}

func (r RefOffset) arg() {}

// RefPreIndexed is a pre-indexed memory reference argument with a register base and immediate offset.
type RefPreIndexed struct {
	Base   Reg // X|SP
	Offset int32
}

func (r RefPreIndexed) arg() {}

// RefIndexed is an memory reference argument with a register base, register index, and optional index modifier.
type RefIndexed struct {
	Base Reg // X|SP
	Idx  Reg // X|W
	Mod  Mod // Idx=X: LSL|SXTX, Idx=W: SXTW|UXTW; LSL requires imm
}

func (r RefIndexed) arg() {}

// ----------------------------------------------------------------

// Imm is a 32-bit integer immediate argument.
type Imm int32

func (i Imm) arg() {}

// Wide is a 64-bit integer immediate argument.
type Wide uint64

func (i Wide) arg() {}

// Float is a 32-bit float immediate argument.
type Float float32

func (i Float) arg() {}

// ----------------------------------------------------------------

// Mod is a shift, rotate, or extension modifier argument.
// Shift and rotate modifiers require an immediate.
type Mod struct {
	ID     uint8
	ImmInv uint8 // bitwise complement, zero indicates unset
}

func (m Mod) arg() {}

// ----------------------------------------------------------------

// Label is a label reference argument with an optional offset from the label.
type Label struct {
	ID     uint32
	Offset int32 // optional offset from label
}

func (l Label) arg() {}

// ----------------------------------------------------------------

// Symbol is a symbol/identifier argument.
type Symbol uint8

func (l Symbol) arg() {}

// ----------------------------------------------------------------

// Flat is an argument flattened for encoding.
type Flat interface {
	flat()
}

type FlatReg uint8

func (r FlatReg) flat() {}

type FlatImm uint64

func (i FlatImm) flat() {}

type FlatMod uint8

func (m FlatMod) flat() {}

type FlatLabel Label

func (l FlatLabel) flat() {}

type FlatDefault struct{}

func (d FlatDefault) flat() {}
