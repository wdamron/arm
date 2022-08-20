package arm

const (
	_ uint8 = iota

	// commands that advance the argument pointer

	CmdR0    // R0, encode a register or reference base into a 5-bit bitfield at bit 0
	CmdR5    // R5, encode a register or reference base into a 5-bit bitfield at bit 5
	CmdR10   // R10, encode a register or reference base into a 5-bit bitfield at bit 10
	CmdR16   // R16, encode a register or reference base into a 5-bit bitfield at bit 16
	CmdRLo16 // RLo16, encode a register in the range 0-15 into a 4-bit bitfield at bit 16
	CmdRNz16 // RNz16, encode a register (except 31) or reference base into a 5-bit bitfield at bit 16
	CmdREven // REven(offset), encode an even register or reference base into a 5-bit bitfield at bit $0
	CmdRNext // encode that this register should be the previous register, plus one

	CmdRwidth30 // Rwidth, SIMD 128-bit indicator at bit 30

	// unsigned immediate encodings

	CmdUbits     // Ubits(offset, bitlen), encode an unsigned immediate starting at bit $0, $1 bits long
	CmdUscaled   // Uscaled(offset, bitlen, shift), encode an unsigned immediate, starting at bit $0, $1 bits long, shifted $2 bits to the right before encoding
	CmdUAlt2     // UAlt2(offset, *[2]uint8), encode an immediate that can only be a limited amount of options
	CmdUAlt4     // UAlt4(offset, *[4]uint8), encode an immediate that can only be a limited amount of options
	CmdUrange    // Urange(offset, min, max), (loc, min, max) asserts the immediate is below or equal to max, encodes the value of (imm-min)
	CmdUsub      // Usub(offset, bitlen, val), encode at $0, $1 bits long, $2 - value. Checks if the value is in the range 0 .. value
	CmdUnegmod   // Unegmod(offset, bitlen), encode at $0, $1 bits long, -value % (1 << $1). Checks if the value is in the range 0 .. value
	CmdUsumdec   // Usumdec(offset, bitlen), encode at $0, $1 bits long, the value of the previous arg + the value of the current arg - 1
	CmdUfields11 // Ufields11(count), encode an immediate bitwise with $0 fields, into bits [11, 21, 20]
	CmdUfields30 // Ufields30(count), encode an immediate bitwise with $0 fields, into bits [30, 12, 11, 10]
	CmdUfields21 // Ufields21, encode an immediate bitwise with 1 field, into bit 21

	// signed immediate encodings

	CmdSbits   // Sbits, encode a signed immediate starting at bit 12, 9 bits long
	CmdSscaled // Sscaled(shift), encode a signed immediate, starting at bit 15, 7 bits long, shifted $0 bits to the right before encoding

	// bit slice encodings. These don't advance the current argument. Only the slice argument actually encodes anything

	CmdChkUbits   // ChkUbits(bitlen), checks if the pointed value fits in $0 (0, 6, 8)
	CmdChkUsum    // ChkUsum(shift), checks that the pointed value fits between 1 and (1 << $0) - prev
	CmdChkSscaled // ChkSscaled, with (offset 10, shift 3)
	CmdChkUrange1 // ChkUrange(max), // check if the pointed value is between 1 and $0
	CmdUslice     // Uslice(offset, bitlen, startoffset), encode at $0, $1 bits long, the bitslice starting at $2 from the current arg
	CmdSslice     // Sslice(offset, bitlen, startoffset), encodes at $0, $1 bits long, the bitslice starting at $2 from the current arg

	CmdSpecial // Special(offset, SpecialType)

	// Extend/Shift fields

	CmdRotates  // Rotates, 2-bits field encoding at bit 22 [LSL, LSR, ASR, ROR]
	CmdExtendsW // ExtendsW, 3-bits field encoding at bit 13 [UXTB, UXTH, UXTW, UXTX, SXTB, SXTH, SXTW, SXTX]. Additionally, LSL is interpreted as UXTW
	CmdExtendsX // ExtendsX, 3-bits field encoding at bit 13 [UXTB, UXTH, UXTW, UXTX, SXTB, SXTH, SXTW, SXTX]. Additionally, LSL is interpreted as UXTX

	// Condition encodings.

	CmdCond    // Cond(offset), normal condition code 4-bit encoding
	CmdCondInv // CondInv(offset), 4-bit encoding, but the last bit is inverted. No AL/NV allowed

	// Mapping of literal -> bitvalue

	CmdLitList // LitList(offset, *[]Symbol)

	// Offsets

	CmdOffset // Offset(RelType)

	// special commands

	CmdAdv  // advances the argument pointer, only needed to skip over an argument.
	CmdBack // moves the argument pointer back.

	// Relocation command types (CmdOffset)

	RelB     // b, bl 26 bits, dword aligned
	RelBCond // b.cond, cbnz, cbz, ldr, ldrsw, prfm: 19 bits, dword aligned
	RelAdr   // adr split 21 bit, byte aligned
	RelAdrp  // adrp split 21 bit, 4096-byte aligned
	RelTbz   // tbnz, tbz: 14 bits, dword aligned

	// Symbol groups (CmdLitList)

	SymATOPS
	SymDCOPS
	SymICOPS
	SymTLBIOPS
	SymBARRIEROPS
	SymMSRIMMOPS
	SymCONTROLREGS
)

// Arm Architecture Reference Manual for A-profile architecture, 4 Feb 2022 Issue H.a
// B2.6.2: Instruction endianness:
// A64 instructions have a fixed length of 32 bits and are always little-endian.

func dec32(code []byte) uint32 {
	_ = code[3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(code[0]) | uint32(code[1])<<8 | uint32(code[2])<<16 | uint32(code[3])<<24
}

func enc32(code []byte, opcode uint32) {
	_ = code[3] // bounds check hint to compiler; see golang.org/issue/14808
	code[0] = byte(opcode)
	code[1] = byte(opcode >> 8)
	code[2] = byte(opcode >> 16)
	code[3] = byte(opcode >> 24)
}

// Special immediate types (CmdSpecial)
const (
	_ uint8 = iota

	SpecialImmWideInv32
	SpecialImmWideInv64
	SpecialImmWide32
	SpecialImmWide64
	SpecialImmStretched
	SpecialImmLogical32
	SpecialImmLogical64
	SpecialImmFloat
	SpecialImmFloatSplit
)

var SpecialName = [...]string{
	SpecialImmWideInv32:  "SpecialImmWideInv32",
	SpecialImmWideInv64:  "SpecialImmWideInv64",
	SpecialImmWide32:     "SpecialImmWide32",
	SpecialImmWide64:     "SpecialImmWide64",
	SpecialImmStretched:  "SpecialImmStretched",
	SpecialImmLogical32:  "SpecialImmLogical32",
	SpecialImmLogical64:  "SpecialImmLogical64",
	SpecialImmFloat:      "SpecialImmFloat",
	SpecialImmFloatSplit: "SpecialImmFloatSplit",
}

var Alts2 = [...][2]uint16{
	{0, 0},
	{0, 1},
	{0, 2},
	{0, 3},
	{0, 4},
	{0, 8},
	{0, 12},
	{0, 16},
	{8, 16},
	{90, 270},
}

var Alts4 = [...][4]uint16{
	{0, 8, 16, 24},
	{0, 16, 32, 48},
	{0, 90, 180, 270},
}

var CmdArgCounts = [...]uint8{
	CmdR0:         0,
	CmdR5:         0,
	CmdR10:        0,
	CmdR16:        0,
	CmdRLo16:      0,
	CmdRNz16:      0,
	CmdREven:      1,
	CmdRNext:      0,
	CmdRwidth30:   0,
	CmdUbits:      2,
	CmdUscaled:    3,
	CmdUAlt2:      2,
	CmdUAlt4:      2,
	CmdUrange:     3,
	CmdUsub:       3,
	CmdUnegmod:    2,
	CmdUsumdec:    2,
	CmdUfields11:  1,
	CmdUfields30:  1,
	CmdUfields21:  0,
	CmdSbits:      0,
	CmdSscaled:    1,
	CmdChkUbits:   1,
	CmdChkUsum:    1,
	CmdChkSscaled: 0,
	CmdChkUrange1: 1,
	CmdUslice:     3,
	CmdSslice:     3,
	CmdSpecial:    2,
	CmdRotates:    0,
	CmdExtendsW:   0,
	CmdExtendsX:   0,
	CmdCond:       1,
	CmdCondInv:    1,
	CmdLitList:    2,
	CmdOffset:     1,
	CmdAdv:        0,
	CmdBack:       0,
}

var CmdName = [...]string{
	CmdR0:          "CmdR0",
	CmdR5:          "CmdR5",
	CmdR10:         "CmdR10",
	CmdR16:         "CmdR16",
	CmdRLo16:       "CmdRLo16",
	CmdRNz16:       "CmdRNz16",
	CmdREven:       "CmdREven",
	CmdRNext:       "CmdRNext",
	CmdRwidth30:    "CmdRwidth30",
	CmdUbits:       "CmdUbits",
	CmdUscaled:     "CmdUscaled",
	CmdUAlt2:       "CmdUAlt2",
	CmdUAlt4:       "CmdUAlt4",
	CmdUrange:      "CmdUrange",
	CmdUsub:        "CmdUsub",
	CmdUnegmod:     "CmdUnegmod",
	CmdUsumdec:     "CmdUsumdec",
	CmdUfields11:   "CmdUfields11",
	CmdUfields30:   "CmdUfields30",
	CmdUfields21:   "CmdUfields21",
	CmdSbits:       "CmdSbits",
	CmdSscaled:     "CmdSscaled",
	CmdChkUbits:    "CmdChkUbits",
	CmdChkUsum:     "CmdChkUsum",
	CmdChkSscaled:  "CmdChkSscaled",
	CmdChkUrange1:  "CmdChkUrange1",
	CmdUslice:      "CmdUslice",
	CmdSslice:      "CmdSslice",
	CmdSpecial:     "CmdSpecial",
	CmdRotates:     "CmdRotates",
	CmdExtendsW:    "CmdExtendsW",
	CmdExtendsX:    "CmdExtendsX",
	CmdCond:        "CmdCond",
	CmdCondInv:     "CmdCondInv",
	CmdLitList:     "CmdLitList",
	CmdOffset:      "CmdOffset",
	CmdAdv:         "CmdAdv",
	CmdBack:        "CmdBack",
	RelB:           "RelB",
	RelBCond:       "RelBCond",
	RelAdr:         "RelAdr",
	RelAdrp:        "RelAdrp",
	RelTbz:         "RelTbz",
	SymATOPS:       "SymATOPS",
	SymDCOPS:       "SymDCOPS",
	SymICOPS:       "SymICOPS",
	SymTLBIOPS:     "SymTLBIOPS",
	SymBARRIEROPS:  "SymBARRIEROPS",
	SymMSRIMMOPS:   "SymMSRIMMOPS",
	SymCONTROLREGS: "SymCONTROLREGS",
}
