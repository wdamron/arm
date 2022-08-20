package arm

import (
	"math"
	"math/bits"
)

// encode writes the matched instruction to the code buffer.
func (a *Assembler) encode() bool {
	if int(a.PC)+4 > len(a.Code) {
		return false
	}
	a.flattenArgs()
	opcode, args, cursor := a.Opcode, a.Flat, uint8(0)
Scan:
	for _, cmd := range a.Commands() {
		switch cmd.Op {
		case CmdAdv:
			cursor++
			continue Scan
		case CmdBack:
			cursor--
			continue Scan
		case CmdRwidth30:
			switch a.SimdSize {
			case 0, 16:
				opcode |= 1 << 30
			}
			continue Scan
		}

		switch arg := args[cursor].(type) {
		case FlatReg:
			switch cmd.Op {
			default:
				return false
			case CmdR0:
				opcode |= uint32(arg)
			case CmdR5:
				opcode |= uint32(arg) << 5
			case CmdR10:
				opcode |= uint32(arg) << 10
			case CmdR16:
				opcode |= uint32(arg) << 16
			case CmdRNz16:
				if arg == 31 {
					return false
				}
				opcode |= uint32(arg) << 16
			case CmdRLo16:
				if arg >= 16 {
					return false
				}
				opcode |= uint32(arg) << 16
			case CmdREven:
				offset := cmd.X[0]
				if arg&1 != 0 {
					return false
				}
				opcode |= uint32(arg) << offset
			case CmdRNext:
				prev, ok := args[cursor-1].(FlatReg)
				if !ok || arg != (prev+1)%32 {
					return false
				}
			}

		case FlatMod:
			switch cmd.Op {
			case CmdRotates:
				switch uint8(arg) {
				default:
					return false
				case SymLSL:
					opcode |= 0b00 << 22
				case SymLSR:
					opcode |= 0b01 << 22
				case SymASR:
					opcode |= 0b10 << 22
				case SymROR:
					opcode |= 0b11 << 22
				}
			case CmdExtendsW, CmdExtendsX:
				switch uint8(arg) {
				default:
					return false
				case SymUXTB:
					opcode |= 0b000 << 13
				case SymUXTH:
					opcode |= 0b001 << 13
				case SymUXTW:
					opcode |= 0b010 << 13
				case SymUXTX:
					opcode |= 0b011 << 13
				case SymSXTB:
					opcode |= 0b100 << 13
				case SymSXTH:
					opcode |= 0b101 << 13
				case SymSXTW:
					opcode |= 0b110 << 13
				case SymSXTX:
					opcode |= 0b111 << 13
				case SymLSL:
					if cmd.Op == CmdExtendsW {
						opcode |= 0b010 << 13
					} else {
						opcode |= 0b011 << 13
					}
				}
			}

		case FlatImm:
			switch cmd.Op {
			default:
				return false

			// symbols:
			case CmdCond:
				offset := cmd.X[0]
				if Symbol(arg) < EQ || Symbol(arg) > NV {
					return false
				}
				opcode |= uint32(Symbol(arg)-EQ) << offset
			case CmdCondInv:
				offset := cmd.X[0]
				if Symbol(arg) < EQ || Symbol(arg) >= AL {
					return false
				}
				opcode |= (uint32(Symbol(arg)-EQ) ^ 1) << offset
			case CmdLitList:
				offset, listSym := cmd.X[0], cmd.X[1]
				if !symListContains(listSym, Symbol(arg)) {
					return false
				}
				if listSym == SymCONTROLREGS {
					arg -= FlatImm(C0)
				}
				opcode |= uint32(SymbolValue[arg]) << offset

			// arithmetic/bitwise:
			case CmdUAlt2:
				offset, listIdx := cmd.X[0], cmd.X[1]
				i, ok := checkAlt(Alts2[listIdx][:], arg)
				if !ok {
					return false
				}
				opcode |= uint32(i) << offset
			case CmdUAlt4:
				offset, listIdx := cmd.X[0], cmd.X[1]
				i, ok := checkAlt(Alts4[listIdx][:], arg)
				if !ok {
					return false
				}
				opcode |= uint32(i) << offset
			case CmdUbits:
				offset, bitlen := cmd.X[0], cmd.X[1]
				mask := (uint32(1) << bitlen) - 1
				if !unsignedRangeCheck(uint64(arg), 0, mask, 0) {
					return false
				}
				opcode |= (uint32(arg) & uint32(mask)) << offset
			case CmdSbits:
				const offset = 12
				const mask = (int32(1) << 9) - 1
				const half = (int32(1) << (9 - 1)) * -1
				if !signedRangeCheck(int64(arg), half, mask+half, 0) {
					return false
				}
				opcode |= (uint32(arg) & uint32(mask)) << offset
			case CmdUscaled:
				offset, bitlen, shift := cmd.X[0], cmd.X[1], cmd.X[2]
				mask := (uint32(1) << bitlen) - 1
				if !unsignedRangeCheck(uint64(arg), 0, mask, shift) {
					return false
				}
				opcode |= ((uint32(arg) >> shift) & mask) << offset
			case CmdSscaled:
				shift := cmd.X[0]
				const offset = 15
				const mask = (int32(1) << 7) - 1
				const half = (int32(1) << (7 - 1)) * -1
				if !signedRangeCheck(int64(arg), half, mask-half, shift) {
					return false
				}
				opcode |= ((uint32(arg) >> shift) & uint32(mask)) << offset
			case CmdUslice, CmdSslice:
				offset, bitlen, start := cmd.X[0], cmd.X[1], cmd.X[2]
				mask := (uint32(1) << bitlen) - 1
				opcode |= ((uint32(arg) >> start) & mask) << offset
			case CmdUrange:
				offset, min, max := cmd.X[0], uint32(cmd.X[1]), uint32(cmd.X[2])
				if !unsignedRangeCheck(uint64(arg), min, max, 0) {
					return false
				}
				opcode |= (uint32(arg) - min) << offset
			case CmdUsub:
				offset, bitlen, add := cmd.X[0], cmd.X[1], int64(cmd.X[2])
				mask := (int64(1) << bitlen) - 1
				if !unsignedRangeCheck(uint64(arg), uint32(add-mask), uint32(add), 0) {
					return false
				}
				opcode |= ((uint32(add) - uint32(arg)) & uint32(mask)) << offset
			case CmdUnegmod:
				offset, bitlen := cmd.X[0], cmd.X[1]
				mask := (uint64(1) << bitlen) - 1
				add := int64(1) << bitlen
				if !unsignedRangeCheck(uint64(arg), 0, uint32(mask), 0) {
					return false
				}
				opcode |= ((uint32(add - int64(arg))) & uint32(mask)) << offset
			case CmdUsumdec:
				offset, bitlen := cmd.X[0], cmd.X[1]
				if cursor == 0 {
					return false
				}
				mask := (uint64(1) << bitlen) - 1
				prev, ok := args[cursor-1].(FlatImm)
				if !ok {
					return false
				}
				opcode |= uint32((uint64(prev)+uint64(arg)-1)&mask) << offset
			case CmdUfields11:
				count := cmd.X[0]
				mask := (uint32(1) << count) - 1
				if !unsignedRangeCheck(uint64(arg), 0, mask, 0) {
					return false
				}
				fields := [...]uint8{20, 21, 11}
				for i, b := range fields[3-count:] {
					opcode |= ((uint32(arg) >> i) & 1) << b
				}
			case CmdUfields30:
				count := cmd.X[0]
				mask := (uint32(1) << count) - 1
				if !unsignedRangeCheck(uint64(arg), 0, mask, 0) {
					return false
				}
				fields := [...]uint8{10, 11, 12, 30}
				for i, b := range fields[4-count:] {
					opcode |= ((uint32(arg) >> i) & 1) << b
				}
			case CmdUfields21:
				if arg&1 != arg {
					return false
				}
				opcode |= (uint32(arg) & 1) << 21
			case CmdSpecial:
				offset, specialType := cmd.X[0], cmd.X[1]
				enc, ok := encSpecialImm(offset, specialType, uint64(arg))
				if !ok {
					return false
				}
				opcode |= enc
			case CmdOffset:
				relType := cmd.X[0]
				enc, ok := encOffset(relType, int64(arg))
				if !ok {
					return false
				}
				opcode |= enc

			// non-consuming:
			case CmdChkUbits:
				bitlen := cmd.X[0]
				mask := (uint32(1) << bitlen) - 1
				if !unsignedRangeCheck(uint64(arg), 0, mask, 0) {
					return false
				}
			case CmdChkUsum:
				shift := cmd.X[0]
				if cursor == 0 {
					return false
				}
				prev, ok := args[cursor-1].(FlatImm)
				if !ok {
					return false
				}
				max := uint32(1) << shift
				if arg == prev {
					max -= uint32(arg)
				}
				if !unsignedRangeCheck(uint64(arg), 1, max, 0) {
					return false
				}
			case CmdChkSscaled:
				const mask = (int32(1) << 10) - 1
				const half = (int32(1) << (10 - 1)) * -1
				if !signedRangeCheck(int64(arg), half, mask+half, 3) {
					return false
				}
			case CmdChkUrange1:
				max := uint32(cmd.X[0])
				if !unsignedRangeCheck(uint64(arg), 1, max, 0) {
					return false
				}
			}

		case FlatLabel:
			if cmd.Op != CmdOffset {
				return false
			}
			relType := cmd.X[0]
			a.Relocs = append(a.Relocs, Reloc{a.PC, relType, arg})

		case FlatDefault:
			switch cmd.Op {
			default:
			case CmdR0: // default to R31
				opcode |= 0b11111
			case CmdR5: // default to R31
				opcode |= 0b11111 << 5
			case CmdR10: // default to R31
				opcode |= 0b11111 << 10
			case CmdR16: // default to R31
				opcode |= 0b11111 << 16
			case CmdExtendsW: // default to LSL
				opcode |= 0b010 << 13
			case CmdExtendsX: // default to LSL
				opcode |= 0b011 << 13
			case CmdUAlt2: // default to 0
				offset, listIdx := cmd.X[0], cmd.X[1]
				i, ok := checkAlt(Alts2[listIdx][:], 0)
				if !ok {
					return false
				}
				opcode |= uint32(i) << offset
			case CmdUAlt4: // default to 0
				offset, listIdx := cmd.X[0], cmd.X[1]
				i, ok := checkAlt(Alts4[listIdx][:], 0)
				if !ok {
					return false
				}
				opcode |= uint32(i) << offset
			}
		}

		switch cmd.Op {
		default:
			cursor++
		case CmdUslice, CmdSslice, CmdChkUbits, CmdChkUsum, CmdChkSscaled, CmdChkUrange1:
			// non-consuming
		}
	}

	enc32(a.Code[a.PC:], opcode)
	a.PC += 4
	return true
}

func encOffset(relType uint8, imm int64) (opcode uint32, ok bool) {
	switch relType {
	case RelB: // b, bl 26-bit (+/- 128 MB); DWORD aligned
		const mask = (int32(1) << 26) - 1
		const half = (int32(1) << (26 - 1)) * -1
		if !signedRangeCheck(imm, half, mask+half, 2) {
			return 0, false
		}
		return (uint32(imm) >> 2) & uint32(mask), true
	case RelBCond: // b.cond, cbnz, cbz, ldr, ldrsw, prfm 19-bit (+/- 1 MB); DWORD aligned
		const mask = (int32(1) << 19) - 1
		const half = (int32(1) << (19 - 1)) * -1
		if !signedRangeCheck(imm, half, mask+half, 2) {
			return 0, false
		}
		return ((uint32(imm) >> 2) & uint32(mask)) << 5, true
	case RelAdr: // adr split 21-bit (+/- 1 MB); BYTE aligned
		const mask = (int32(1) << 21) - 1
		const half = (int32(1) << (21 - 1)) * -1
		if !signedRangeCheck(imm, half, mask+half, 0) {
			return 0, false
		}
		low := ((uint32(imm) >> 2) & 0x7FFFF) << 5
		return low | (uint32(imm)&3)<<29, true
	case RelAdrp: // adrp split 21-bit (+/- 4 GB); page aligned (4KB)
		const mask = (int64(1) << 21) - 1
		const half = (int64(1) << (21 - 1)) * -1
		scaled := imm >> 12
		if scaled<<12 != imm || scaled < half && scaled > mask+half {
			return 0, false
		}
		low := ((uint32(uint64(imm) >> 14)) & 0x7FFFF) << 5
		return low | ((uint32(uint64(imm)>>12))&3)<<29, true
	case RelTbz: // tbnz, tbz 14-bit (+/- 32 KB); DWORD aligned
		const mask = (int32(1) << 14) - 1
		const half = (int32(1) << (14 - 1)) * -1
		if !signedRangeCheck(imm, half, mask+half, 2) {
			return 0, false
		}
		return ((uint32(imm) >> 2) & uint32(mask)) << 5, true
	}
	return 0, false
}

func checkAlt(alts []uint16, v FlatImm) (i uint8, ok bool) {
	for i := len(alts) - 1; i >= 0; i-- {
		if v == FlatImm(alts[i]) {
			return uint8(i), true
		}
	}
	return 0, false
}

func unsignedRangeCheck(v uint64, min, max uint32, scale uint8) bool {
	scaled := v >> scale
	return scaled<<scale == v && scaled >= uint64(min) && scaled <= uint64(max)
}

func signedRangeCheck(v int64, min, max int32, scale uint8) bool {
	scaled := v >> scale
	return scaled<<scale == v && scaled >= int64(min) && scaled <= int64(max)
}

func encSpecialImm(offset, op uint8, v uint64) (opcode uint32, ok bool) {
	switch op {
	case SpecialImmWideInv64:
		return encImmWide64(offset, ^v)
	case SpecialImmWideInv32:
		return encImmWide32(offset, v, true)
	case SpecialImmWide64:
		return encImmWide64(offset, v)
	case SpecialImmWide32:
		return encImmWide32(offset, v, false)
	case SpecialImmStretched:
		return encImmStretched(offset, v)
	case SpecialImmLogical32:
		return encImmLogical32(offset, v)
	case SpecialImmLogical64:
		return encImmLogical64(offset, v)
	case SpecialImmFloat:
		return encImmFloat(offset, v)
	case SpecialImmFloatSplit:
		return encImmFloatSplit(offset, v)
	}
	return
}

func encImmLogical32(offset uint8, v uint64) (uint32, bool) {
	if v > math.MaxUint32 {
		return 0, false
	}
	transitions := uint32(v) ^ (bits.RotateLeft32(uint32(v), -1))
	div := uint32(bits.OnesCount32(transitions))
	if div == 0 {
		return 0, false
	}
	elemSize := uint32(64) / div
	if uint32(v) != bits.RotateLeft32(uint32(v), int(elemSize)) {
		return 0, false
	}
	elem := uint32(v) & ((1 << elemSize) - 1)
	ones := uint32(bits.OnesCount32(elem))
	imms := (^((elemSize << 1) - 1) & 0x3F) | (ones - 1)
	var immr uint32
	if elem&1 != 0 {
		immr = ones - uint32(bits.TrailingZeros32(^elem))
	} else {
		immr = elemSize - uint32(bits.TrailingZeros32(elem))
	}
	enc := (uint16(immr) << 6) | uint16(imms)
	return uint32(enc) << offset, true
}

func encImmLogical64(offset uint8, v uint64) (uint32, bool) {
	transitions := v ^ (bits.RotateLeft64(v, -1))
	div := uint64(bits.OnesCount64(transitions))
	if div == 0 {
		return 0, false
	}
	elemSize := uint64(128) / div
	if v != bits.RotateLeft64(v, int(elemSize)) {
		return 0, false
	}
	elem := v & ((1 << elemSize) - 1)
	ones := uint64(bits.OnesCount64(elem))
	imms := (^((elemSize << 1) - 1) & 0x7F) | (ones - 1)
	var immr uint64
	if elem&1 != 0 {
		immr = ones - uint64(bits.TrailingZeros64(^elem))
	} else {
		immr = elemSize - uint64(bits.TrailingZeros64(elem))
	}
	var n uint16
	if imms&0x40 == 0 {
		n = 1
	}
	imms &= 0x3F
	enc := (n << 12) | (uint16(immr) << 6) | uint16(imms)
	return uint32(enc) << offset, true
}

func encImmStretched(offset uint8, v uint64) (uint32, bool) {
	chk := v & 0x0101010101010101
	chk |= chk << 1
	chk |= chk << 2
	chk |= chk << 4
	if v != chk {
		return 0, false
	}
	masked := v & 0x8040201008040201
	masked |= masked >> 32
	masked |= masked >> 16
	masked |= masked >> 8
	enc := uint32(masked) & 0xFF
	opcode := (enc & 0x1F) << offset
	opcode |= (enc & 0xE0) << (offset + 6)
	return opcode, true
}

func encImmWide64(offset uint8, v uint64) (uint32, bool) {
	pos := uint32(bits.TrailingZeros64(v)) & 0b110000
	masked := uint32(v>>pos) & 0xFFFF
	if uint64(masked)<<pos != v {
		return 0, false
	}
	enc := masked | (pos << 12)
	return enc << offset, true
}

func encImmWide32(offset uint8, v uint64, invert bool) (uint32, bool) {
	if v > math.MaxUint32 {
		return 0, false
	}
	if invert {
		v = uint64(uint32(^v))
	}
	pos := uint32(bits.TrailingZeros32(uint32(v))) & 0b10000
	masked := uint32(v>>pos) & 0xFFFF
	if uint64(masked)<<pos != v {
		return 0, false
	}
	enc := masked | (pos << 12)
	return enc << offset, true
}

func encImmFloat(offset uint8, v uint64) (uint32, bool) {
	enc := uint8((((v >> 24) & 0x80) | ((v >> 19) & 0x7F)))
	if chk := (uint32(v) >> 25) & 0x3F; (chk == 0b100000 || chk == 0b011111) && uint32(v)&0x7FFFF == 0 {
		return uint32(enc) << offset, true
	}
	return 0, false
}

func encImmFloatSplit(offset uint8, v uint64) (uint32, bool) {
	enc := uint8((((v >> 24) & 0x80) | ((v >> 19) & 0x7F)))
	if chk := (uint32(v) >> 25) & 0x3F; (chk == 0b100000 || chk == 0b011111) && uint32(v)&0x7FFFF == 0 {
		opcode := uint32(enc&0x1F) << offset
		opcode |= uint32(enc&0xE0) << (offset + 6)
		return opcode, true
	}
	return 0, false
}
