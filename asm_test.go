package arm

import "testing"

func TestLabels(t *testing.T) {
	code := make([]byte, 256)
	var a Assembler
	a.Init(code)

	// Backward label jump:

	start := a.NewLabel()
	a.PC = 64
	if !a.Inst(B, start) {
		t.Fatalf("Failed to encode B: %v", a.Err)
	}
	// Patch instructions:
	if !a.ApplyRelocations() {
		t.Fatalf("Failed to apply relocs: %v", a.Err)
	}
	if actual, expected := dec32(code[a.PC-4:]), 0x17FFFFF0; actual != uint32(expected) {
		t.Fatalf("Invalid B %08X, expecting %08X", actual, expected)
	}
	// Ensure branch with immediate offset encodes to the same instruction:
	a.PC = 0
	a.Inst(B, Imm(-64))
	if actual, expected := dec32(code[0:]), 0x17FFFFF0; actual != uint32(expected) {
		t.Fatalf("Invalid B %08X, expecting %08X", actual, expected)
	}

	// Forward label jump, with extra extraOffset:

	const endPC, extraOffset = 96, -16
	code = make([]byte, 256)
	a.Init(code)
	end := a.NewLabel() // PC assigned later
	if !a.Inst(B, Label{end.ID, extraOffset}) {
		t.Fatalf("Failed to encode B: %v", a.Err)
	}
	// Assign the current PC to the label:
	a.PC = endPC
	a.SetLabel(end)
	// Patch instructions:
	if !a.ApplyRelocations() {
		t.Fatalf("Failed to apply relocs: %v", a.Err)
	}
	if actual, expected := dec32(code[0:]), 0x14000014; actual != uint32(expected) {
		t.Fatalf("Invalid B %08X, expecting %08X", actual, expected)
	}
	// Ensure branch with immediate offset encodes to the same instruction:
	a.PC = 0
	a.Inst(B, Imm(endPC+extraOffset))
	if actual, expected := dec32(code[0:]), 0x14000014; actual != uint32(expected) {
		t.Fatalf("Invalid B %08X, expecting %08X", actual, expected)
	}
}

func TestEncoding(t *testing.T) {
	code := make([]byte, 256)
	var a Assembler

	test := func(enc uint32, inst Inst, args ...Arg) {
		a.Init(code)
		if !a.Inst(inst, args...) {
			t.Logf("Failed to encode inst %v for enc %08X -- err: %v\n\targs:\n%#+v", inst, enc, a.Err, args)
			t.Fail()
		} else if actual := dec32(code); actual != enc {
			t.Logf("Invalid inst=%v:\n%032b (expected) = %08X\n%032b (actual)   = %08X\n%032b (opcode)\n%032b (args", inst, enc, enc, actual, actual, a.Opcode, a.Opcode^actual)
			t.Fail()
		}
	}

	// See https://github.com/CensoredUsername/dynasm-rs/blob/e6248cc55917e96b7ce7a590c0c93f61c609471e/testing/tests/gen_aarch64/aarch64_tests_0.rs.gen
	// (v1.2.3)

	test(0x5EE0B85F, ABS, ScalarD(31), ScalarD(2))
	test(0x5EE0B875, ABS, ScalarD(21), ScalarD(3))
	test(0x4E20B968, ABS, Vec16B(8), Vec16B(11))
	test(0x4E20BB55, ABS, Vec16B(21), Vec16B(26))
	test(0x0E20B990, ABS, Vec8B(16), Vec8B(12))
	test(0x0E20BBEA, ABS, Vec8B(10), Vec8B(31))
	test(0x4E60B8DA, ABS, Vec8H(26), Vec8H(6))
	test(0x4E60B896, ABS, Vec8H(22), Vec8H(4))
	test(0x0E60B980, ABS, Vec4H(0), Vec4H(12))
	test(0x0E60BA24, ABS, Vec4H(4), Vec4H(17))
	test(0x4EA0BB0A, ABS, Vec4S(10), Vec4S(24))
	test(0x4EA0B9AD, ABS, Vec4S(13), Vec4S(13))
	test(0x0EA0B934, ABS, Vec2S(20), Vec2S(9))
	test(0x0EA0BBDD, ABS, Vec2S(29), Vec2S(30))
	test(0x4EE0B8BE, ABS, Vec2D(30), Vec2D(5))
	test(0x4EE0B85F, ABS, Vec2D(31), Vec2D(2))

	test(0x1A0903BE, ADC, W(30), W(29), W(9))
	test(0x1A0903B6, ADC, W(22), W(29), W(9))
	test(0x9A190280, ADC, X(0), X(20), X(25))
	test(0x9A1502B5, ADC, X(21), X(21), X(21))

	test(0x3A0A0325, ADCS, W(5), W(25), W(10))
	test(0x3A1503EA, ADCS, W(10), W(31), W(21))
	test(0xBA1F008E, ADCS, X(14), X(4), XZR)
	test(0xBA0B03A8, ADCS, X(8), X(29), X(11))

	test(0x0B150065, ADD, W(5), W(3), W(21))
	test(0x0B882E6A, ADD, W(10), W(19), W(8), ModASR.Imm(11))
	test(0x8B8AEC81, ADD, X(1), X(4), X(10), ModASR.Imm(59))
	test(0x8B02F261, ADD, X(1), X(19), X(2), ModLSL.Imm(60))
	test(0x0B1D0BC8, ADD, W(8), W(30), W(29), ModLSL.Imm(2))
	test(0x0B3DCF99, ADD, W(25), W(28), W(29), ModSXTW.Imm(3))
	test(0x8B2B84C6, ADD, X(6), X(6), W(11), ModSXTB.Imm(1))
	test(0x8B2DA2F1, ADD, X(17), X(23), W(13), ModSXTH)
	test(0x8B1B004E, ADD, X(14), X(2), X(27))
	test(0x8B21E944, ADD, X(4), X(10), X(1), ModSXTX.Imm(2))
	test(0x1128E8B8, ADD, W(24), W(5), Imm(2618))
	test(0x115BE3AC, ADD, W(12), W(29), Imm(1784), ModLSL.Imm(12))
	test(0x9131F3D4, ADD, X(20), X(30), Imm(3196))
	test(0x91120031, ADD, X(17), X(1), Imm(1152))
	test(0x5EFC8632, ADD, ScalarD(18), ScalarD(17), ScalarD(28))
	test(0x5EF08594, ADD, ScalarD(20), ScalarD(12), ScalarD(16))
	test(0x4E288564, ADD, Vec16B(4), Vec16B(11), Vec16B(8))
	test(0x4E2B86E2, ADD, Vec16B(2), Vec16B(23), Vec16B(11))
	test(0x0E2F8575, ADD, Vec8B(21), Vec8B(11), Vec8B(15))
	test(0x0E278665, ADD, Vec8B(5), Vec8B(19), Vec8B(7))
	test(0x4E7F865C, ADD, Vec8H(28), Vec8H(18), Vec8H(31))
	test(0x4E618780, ADD, Vec8H(0), Vec8H(28), Vec8H(1))
	test(0x0E6C8559, ADD, Vec4H(25), Vec4H(10), Vec4H(12))
	test(0x0E6C85C3, ADD, Vec4H(3), Vec4H(14), Vec4H(12))
	test(0x4EA6849B, ADD, Vec4S(27), Vec4S(4), Vec4S(6))
	test(0x4EA2867F, ADD, Vec4S(31), Vec4S(19), Vec4S(2))
	test(0x0EB28601, ADD, Vec2S(1), Vec2S(16), Vec2S(18))
	test(0x0EBB8752, ADD, Vec2S(18), Vec2S(26), Vec2S(27))
	test(0x4EF48626, ADD, Vec2D(6), Vec2D(17), Vec2D(20))
	test(0x4EFF851C, ADD, Vec2D(28), Vec2D(8), Vec2D(31))

	test(0x2B542CAE, ADDS, W(14), W(5), W(20), ModLSR.Imm(11))
	test(0x2B515BED, ADDS, W(13), W(31), W(17), ModLSR.Imm(22))
	test(0xAB462899, ADDS, X(25), X(4), X(6), ModLSR.Imm(10))
	test(0xAB160C5F, ADDS, XZR, X(2), X(22), ModLSL.Imm(3))
	test(0x2B27ED10, ADDS, W(16), W(8), W(7), ModSXTX.Imm(3))
	test(0x2B202C0D, ADDS, W(13), W(0), W(0), ModUXTH.Imm(3))
	test(0xAB292CFF, ADDS, X(31), X(7), W(9), ModUXTH.Imm(3))
	test(0xAB304419, ADDS, X(25), X(0), W(16), ModUXTW.Imm(1))
	test(0xAB0D025B, ADDS, X(27), X(18), X(13))
	test(0xAB3C6027, ADDS, X(7), X(1), X(28), ModUXTX)
	test(0x31759288, ADDS, W(8), W(20), Imm(3428), ModLSL.Imm(12))
	test(0x31518F61, ADDS, W(1), W(27), Imm(1123), ModLSL.Imm(12))
	test(0xB10363E9, ADDS, X(9), XSP, Imm(216))
	test(0xB1019D0F, ADDS, X(15), X(8), Imm(103))

	test(0x4E31B89A, ADDV, ScalarB(26), Vec16B(4))
	test(0x4E31BAC1, ADDV, ScalarB(1), Vec16B(22))
	test(0x0E31B92D, ADDV, ScalarB(13), Vec8B(9))
	test(0x0E31BA2A, ADDV, ScalarB(10), Vec8B(17))
	test(0x4E71BB31, ADDV, ScalarH(17), Vec8H(25))
	test(0x4E71B9D9, ADDV, ScalarH(25), Vec8H(14))
	test(0x0E71BBE0, ADDV, ScalarH(0), Vec4H(31))
	test(0x0E71B948, ADDV, ScalarH(8), Vec4H(10))
	test(0x4EB1BBFE, ADDV, ScalarS(30), Vec4S(31))
	test(0x4EB1BB4F, ADDV, ScalarS(15), Vec4S(26))

	test(0x30497080, ADR, X(0), Imm(601617))
	test(0x300E940C, ADR, X(12), Imm(119425))

	test(0xF00FAFEA, ADRP, X(10), Imm(526381056))
	test(0xF00FFFF1, ADRP, X(17), Imm(536866816))

	test(0x54094645, B, PL, Imm(75976))
	test(0x54E8EDEF, B, NV, Imm(-188996))
	test(0x15F232F9, B, Imm(130599908))
	test(0x17396BCC, B, Imm(-52056272))

	test(0x330703F0, BFC, W(16), Imm(25), Imm(1))

	test(0x9713AA2A, BL, Imm(-61953880))
	test(0xD63F03E0, BLR, X(31))

	test(0x88AC7CA4, CAS, W(12), W(4), Ref{X(5)})
	test(0xC8A27FA0, CAS, X(2), X(0), Ref{X(29)})
	test(0xC8BC7FBF, CAS, X(28), XZR, Ref{X(29)})
	test(0x082C7FE0, CASP, W(12), W(13), W(0), W(1), Ref{XSP})
	test(0x483C7D2E, CASP, X(28), X(29), X(14), X(15), Ref{X(9)})

	test(0x351FE73A, CBNZ, W(26), Imm(261348))
	test(0x351B5997, CBNZ, W(23), Imm(224048))
	test(0xB5ED6C3C, CBNZ, X(28), Imm(-152188))
	test(0xB5F5663C, CBNZ, X(28), Imm(-86844))
	test(0x34105DC2, CBZ, W(2), Imm(134072))
	test(0x34F87A35, CBZ, W(21), Imm(-61628))
	test(0xB4ED32EE, CBZ, X(14), Imm(-154020))
	test(0xB41462FF, CBZ, X(31), Imm(167004))
	test(0x3A49F8E9, CCMN, W(7), Imm(9), Imm(9), NV)
	test(0x3A577BAB, CCMN, W(29), Imm(23), Imm(11), VC)
	test(0xBA508B44, CCMN, X(26), Imm(16), Imm(4), HI)

	test(0xD500401F, CFINV)

	test(0x2B582E3F, CMN, W(17), W(24), ModLSR.Imm(11))
	test(0x2B9E297F, CMN, W(11), W(30), ModASR.Imm(10))
	test(0x2B24E99F, CMN, W(12), W(4), ModSXTX.Imm(2))
	test(0xAB31A29F, CMN, X(20), W(17), ModSXTH)
	test(0xAB33037F, CMN, X(27), W(19), ModUXTB)

	test(0x5A93066F, CNEG, W(15), W(19), NE)
	test(0xDA90A60B, CNEG, X(11), X(16), LT)

	test(0xD50B7B3A, DC, CVAU, X(26))
	test(0xD4B08861, DCPS1, Imm(33859))
	test(0xD5033BBF, DMB, ISH)

	test(0x5E0706E3, DUP, ScalarB(3), Vec8B(23).I(3))
	test(0x5E1A04F5, DUP, ScalarH(21), Vec8H(7).I(6))

	test(0x5E30DAE3, FADDP, ScalarH(3), Vec2H(23))
	test(0x6E5F1663, FADDP, Vec8H(3), Vec8H(19), Vec8H(31))
	test(0x4E6DE6CD, FCMEQ, Vec2D(13), Vec2D(22), Vec2D(13))
	test(0x5EF8DAB3, FCMEQ, ScalarH(19), ScalarH(21), Float(0))
	test(0x4EF8DAFD, FCMEQ, Vec8H(29), Vec8H(23), Float(0))
	test(0x2F5C51E1, FCMLA, Vec4H(1), Vec4H(15), Vec4H(28).I(0), Imm(180))
	test(0x2E9BDE94, FCMLA, Vec2S(20), Vec2S(20), Vec2S(27), Imm(270))

	test(0x4F05FE97, FMOV, Vec8H(23), Float(-20.0))
	test(0x0F00FD1B, FMOV, Vec4H(27), Float(3.0))
	test(0x4F03F50A, FMOV, Vec4S(10), Float(0.75))
	test(0x0F06F634, FMOV, Vec2S(20), Float(-0.265625))
	test(0x6F04F60A, FMOV, Vec2D(10), Float(-4.0))
	test(0x9EAF01CE, FMOV, Vec2D(14).I(1), X(14))
	test(0x9EAE03E5, FMOV, X(5), Vec2D(31).I(1))
	test(0x1EE7501B, FMOV, ScalarH(27), Float(26.0))
	test(0x1EFB1009, FMOV, ScalarH(9), Float(-0.375))
	test(0x1E2FB000, FMOV, ScalarS(0), Float(1.8125))
	test(0x1E66B004, FMOV, ScalarD(4), Float(21.0))

	test(0x4C4073E1, LD1, Vec16B(1).List(1), Ref{XSP})
	test(0x4C407368, LD1, Vec16B(8).List(1), Ref{X(27)})
	test(0x0C407672, LD1, Vec4H(18).List(1), Ref{X(19)})
	test(0x4C40A0C5, LD1, Vec16B(5).List(2), Ref{X(6)})
	test(0x0C40A8AE, LD1, Vec2S(14).List(2), Ref{X(5)})
	test(0x0C402FFB, LD1, Vec1D(27).List(4), Ref{XSP})
	test(0x0CDF7053, LD1, Vec8B(19).List(1), Ref{X(2)}, Imm(8))
	test(0x4CDA752F, LD1, Vec8H(15).List(1), Ref{X(9)}, X(26))
	test(0x4CDF6CC7, LD1, Vec2D(7).List(3), Ref{X(6)}, Imm(48))
	test(0x0CD523CE, LD1, Vec8B(14).List(4), Ref{X(30)}, X(21))
	test(0x4D401D10, LD1, Vec16B(16).List(1).I(15), Ref{X(8)})
	test(0x0DDF80F8, LD1, Vec2S(24).List(1).I(0), Ref{X(7)}, Imm(4))
	test(0x4D4079F0, LD3, Vec8H(16).List(3).I(7), Ref{X(15)})

	test(0xD956B2AF, LDAPUR, X(15), RefOffset{X(21), -149})
	test(0x1945A2C5, LDAPURB, W(5), RefOffset{X(22), 90})
	test(0x6D4054E6, LDP, ScalarD(6), ScalarD(21), Ref{X(7)})
	test(0xAD74C6D0, LDP, ScalarQ(16), ScalarQ(17), RefOffset{X(22), -368})
	test(0x28EC4612, LDP, W(18), W(17), Ref{X(16)}, Imm(-160))
	test(0x29D05D12, LDP, W(18), W(23), RefPreIndexed{X(8), 128})
	test(0xA9EE109D, LDP, X(29), X(4), RefPreIndexed{X(4), -288})
	test(0xA9DEF2FE, LDP, X(30), X(28), RefPreIndexed{X(23), 488})
	test(0x29407BB1, LDP, W(17), W(30), Ref{X(29)})
	test(0x3C6158CB, LDR, ScalarB(11), RefIndexed{X(6), W(1), ModUXTW})
	test(0x3C67D85E, LDR, ScalarB(30), RefIndexed{X(2), W(7), ModSXTW})
	test(0x7C6C69C1, LDR, ScalarH(1), RefIndexed{X(14), X(12), Mod{}})
	test(0x7C6F6A88, LDR, ScalarH(8), RefIndexed{X(20), X(15), ModLSL})
	test(0xB876DA28, LDR, W(8), RefIndexed{X(17), W(22), ModSXTW.Imm(2)})
	test(0xF865CBEF, LDR, X(15), RefIndexed{XSP, W(5), ModSXTW})
	test(0xF869DA07, LDR, X(7), RefIndexed{X(16), W(9), ModSXTW.Imm(3)})

	test(0x5E150478, MOV, ScalarB(24), Vec16B(3).I(10))
	test(0x5E0E06BB, MOV, ScalarH(27), Vec4H(21).I(3))
	test(0x6E1F15C8, MOV, Vec16B(8).I(15), Vec8B(14).I(2))
	test(0x6E1F15C8, MOV, Vec16B(8).I(15), Vec16B(14).I(2))
	test(0x12B204B4, MOV, INVERTED, W(20), Imm(1876623359))
	test(0x12843223, MOV, INVERTED, W(3), Wide(4294958702))
	test(0x92A7BC98, MOV, INVERTED, X(24), Wide(18446744072671199231))
	test(0x52ADAFA1, MOV, W(1), Imm(1836908544))
	test(0xD2C83A36, MOV, X(22), Wide(72365903970304))
	test(0x4EA31C6C, MOV, Vec16B(12), Vec16B(3))
	test(0x3200F3FD, MOV, LOGICAL, W(29), Imm(1431655765))
	test(0xB201EBF1, MOV, LOGICAL, X(17), Wide(13527612320720337851))
	test(0x0E1C3DF2, MOV, W(18), Vec4S(15).I(3))

	test(0x4F01E47C, MOVI, Vec16B(28), Imm(35), ModLSL)
	test(0x4F02E493, MOVI, Vec16B(19), Imm(68))
	test(0x4F0267E1, MOVI, Vec4S(1), Imm(95), ModLSL.Imm(24))
	test(0x2F05E65F, MOVI, ScalarD(31), Wide(18374967950353432320))
	test(0x6F06E77A, MOVI, Vec2D(26), Wide(18446463698227757055))

	test(0xD503201F, NOP)

	test(0xD8E477B8, PRFM, Imm(24), Imm(-225548))
	test(0xF8A369AB, PRFM, Imm(11), RefIndexed{X(13), X(3), ModLSL})
	test(0xF89501AA, PRFUM, Imm(10), RefOffset{X(13), -176})
	test(0xD503223F, PSB, CSYNC)

	test(0xD65F02E0, RET, X(23))
	test(0xD65F0100, RET, X(8))
	test(0xD65F03C0, RET)

	test(0x0FA6E1DB, SDOT, Vec2S(27), Vec8B(14), Vec4B(6).I(1))
	test(0x4F9FEAC2, SDOT, Vec4S(2), Vec16B(22), Vec4B(31).I(2))
	test(0x0E979582, SDOT, Vec2S(2), Vec8B(12), Vec8B(23))
	test(0x4E86965F, SDOT, Vec4S(31), Vec16B(18), Vec16B(6))

	test(0x4C00725A, ST1, Vec16B(26).List(1), Ref{X(18)})
	test(0x4C007BFC, ST1, Vec4S(28).List(1), Ref{XSP})
	test(0x4C00AF7C, ST1, Vec2D(28).List(2), Ref{X(27)})
	test(0x0C0066DE, ST1, Vec4H(30).List(3), Ref{X(22)})
	test(0x4C002549, ST1, Vec8H(9).List(4), Ref{X(10)})
	test(0x4C9F7347, ST1, Vec16B(7).List(1), Ref{X(26)}, Imm(16))
	test(0x4C9F7FF0, ST1, Vec2D(16).List(1), Ref{XSP}, Imm(16))
	test(0x4C85A1D4, ST1, Vec16B(20).List(2), Ref{X(14)}, X(5))
	test(0x0C9F64ED, ST1, Vec4H(13).List(3), Ref{X(7)}, Imm(24))
	test(0x4C9D6B6F, ST1, Vec4S(15).List(3), Ref{X(27)}, X(29))
	test(0x0D0012FC, ST1, Vec8B(28).List(1).I(4), Ref{X(23)})
	test(0x0D00867D, ST1, Vec2D(29).List(1).I(0), Ref{X(19)})
	test(0x0D8E0123, ST1, Vec8B(3).List(1).I(0), Ref{X(9)}, X(14))
	test(0x4D8B82A6, ST1, Vec4S(6).List(1).I(2), Ref{X(21)}, X(11))
	test(0x0C9F82E3, ST2, Vec8B(3).List(2), Ref{X(23)}, Imm(16))

	test(0x6C9554B3, STP, ScalarD(19), ScalarD(21), Ref{X(5)}, Imm(336))
	test(0xACA38C9D, STP, ScalarQ(29), ScalarQ(3), Ref{X(4)}, Imm(-912))
	test(0x2D81A119, STP, ScalarS(25), ScalarS(8), RefPreIndexed{X(8), 12})
	test(0xADA07567, STP, ScalarQ(7), ScalarQ(29), RefPreIndexed{X(11), -1024})
	test(0x2D00272B, STP, ScalarS(11), ScalarS(9), Ref{X(25)})
	test(0xAD040D1C, STP, ScalarQ(28), ScalarQ(3), RefOffset{X(8), 128})
	test(0x28B151BC, STP, W(28), W(20), Ref{X(13)}, Imm(-120))
	test(0xA9AA431E, STP, X(30), X(16), RefPreIndexed{X(24), -352})
	test(0x3C87E46C, STR, ScalarQ(12), Ref{X(3)}, Imm(126))
	test(0x3C1B9C4B, STR, ScalarB(11), RefPreIndexed{X(2), -71})
	test(0xBC22FB16, STR, ScalarS(22), RefIndexed{X(24), X(2), ModSXTX.Imm(2)})
	test(0xB83F7991, STR, W(17), RefIndexed{X(12), X(31), ModLSL.Imm(2)})

	test(0xD50AA775, SYS, Imm(2), C10, C7, Imm(3), X(21))
	test(0xD52B8349, SYSL, X(9), Imm(3), C8, C3, Imm(2))

	test(0x4E1023DE, TBL, Vec16B(30), Vec16B(30).List(2), Vec16B(16))
	test(0x0E0642A3, TBL, Vec8B(3), Vec16B(21).List(3), Vec8B(6))

	test(0x3742E27B, TBNZ, W(27), Imm(8), Imm(23628))
	test(0x37F1471B, TBNZ, W(27), Imm(30), Imm(10464))
	test(0x375005E6, TBNZ, X(6), Imm(10), Imm(188))
	test(0x37B04418, TBNZ, X(24), Imm(22), Imm(2176))
	test(0x36C844CE, TBZ, W(14), Imm(25), Imm(2200))
	test(0x36F84370, TBZ, W(16), Imm(31), Imm(2156))
	test(0xB6A5833B, TBZ, X(27), Imm(52), Imm(-20380))
	test(0x363DA928, TBZ, X(8), Imm(7), Imm(-19164))

}
