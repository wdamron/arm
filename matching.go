package arm

const (
	_ uint8 = iota

	MatLitSymbol // Lit(Symbol)
	MatLitInt    // LitInt(int32), immediate literal
	MatLitFloat  // LitFloat(float32), float literal
	MatSymbol    // Symbol
	MatCond      // condition code Symbol
	MatImm       // 32/64-bit integer immediate
	MatFloat     // 32-bit float immediate (FMOV)

	// W and X registers (except SP)

	MatW // scalar 32-bit integer register (except WSP)
	MatX // scalar 64-bit integer register (except XSP)

	// WSP and XSP registers (except WZR and XZR)

	MatWSP // 32-bit stack pointer register
	MatXSP // 64-bit stack pointer register

	// scalar simd regs

	MatB // scalar simd 8-bit register
	MatH // scalar simd 16-bit register
	MatS // scalar simd 32-bit register
	MatD // scalar simd 64-bit register
	MatQ // scalar simd 128-bit register

	// vector simd regs

	MatV              // V(Size), vector register with elements of the specified size. Accepts a lane count of either 64 or 128 total bits
	MatVStatic        // VStatic(Size, lanes), vector register with elements of the specifized size, with the specified lane count
	MatVElement       // VElement(Size), vector register with element specifier, with the element of the specified size. The lane count is unchecked.
	MatVElementStatic // VElementStatic(Size, idx), vector register with element specifier, with the element of the specified size and the element index set to the provided value
	MatVStaticElement // VStaticElement(Size, lanes), vector register with elements of the specified size, with the specified lane count, with an element specifier

	MatRegList        // RegList(len, Size), register list with $0 items, with the elements of size $1
	MatRegListStatic  // RegListStatic(len, Size, lanes), register list with $0 items, with the elements of size $1 and a lane count of $2
	MatRegListElement // RegListElement(len, Size), register list with element specifier. It has $0 items with a size of $1

	MatOffset // jump offset

	// references

	MatRefBase   // memory reference with base register (integer or SP register), optionally followed by a register or immediate offset argument for post-indexing
	MatRefOffset // memory reference with base register (integer or SP register) and immediate offset
	MatRefPre    // pre-indexed memory reference with base register (integer or SP register) and immediate offset
	MatRefIndex  // memory reference with base register (integer or SP register), index register (integer register), and optional index modifier

	MatLitMod // LitMod(Mod), a single modifier
	MatMod    // Mod(*[]Mod), a set of allowed modifiers

	MatEnd // possible op mnemnonic end (everything after this point uses the default encoding)
)

var MatcherArgCounts = [...]uint8{
	MatLitSymbol:      1,
	MatLitInt:         1,
	MatLitFloat:       1,
	MatSymbol:         0,
	MatCond:           0,
	MatImm:            0,
	MatFloat:          0,
	MatW:              0,
	MatX:              0,
	MatWSP:            0,
	MatXSP:            0,
	MatB:              0,
	MatH:              0,
	MatS:              0,
	MatD:              0,
	MatQ:              0,
	MatV:              1,
	MatVStatic:        2,
	MatVElement:       1,
	MatVElementStatic: 2,
	MatVStaticElement: 2,
	MatRegList:        2,
	MatRegListStatic:  3,
	MatRegListElement: 2,
	MatOffset:         0,
	MatRefBase:        0,
	MatRefOffset:      0,
	MatRefPre:         0,
	MatRefIndex:       0,
	MatLitMod:         1,
	MatMod:            1,
	MatEnd:            0,
}

var MatchName = [...]string{
	MatLitSymbol:      "MatLitSymbol",
	MatLitInt:         "MatLitInt",
	MatLitFloat:       "MatLitFloat",
	MatSymbol:         "MatSymbol",
	MatCond:           "MatCond",
	MatImm:            "MatImm",
	MatFloat:          "MatFloat",
	MatW:              "MatW",
	MatX:              "MatX",
	MatWSP:            "MatWSP",
	MatXSP:            "MatXSP",
	MatB:              "MatB",
	MatH:              "MatH",
	MatS:              "MatS",
	MatD:              "MatD",
	MatQ:              "MatQ",
	MatV:              "MatV",
	MatVStatic:        "MatVStatic",
	MatVElement:       "MatVElement",
	MatVElementStatic: "MatVElementStatic",
	MatVStaticElement: "MatVStaticElement",
	MatRegList:        "MatRegList",
	MatRegListStatic:  "MatRegListStatic",
	MatRegListElement: "MatRegListElement",
	MatOffset:         "MatOffset",
	MatRefBase:        "MatRefBase",
	MatRefOffset:      "MatRefOffset",
	MatRefPre:         "MatRefPre",
	MatRefIndex:       "MatRefIndex",
	MatLitMod:         "MatLitMod",
	MatMod:            "MatMod",
	MatEnd:            "MatEnd",
}

var MatcherFlatArgCounts = [...]uint8{
	MatLitSymbol:      0,
	MatLitInt:         0,
	MatLitFloat:       0,
	MatSymbol:         1,
	MatCond:           1,
	MatImm:            1,
	MatFloat:          1,
	MatW:              1,
	MatX:              1,
	MatWSP:            1,
	MatXSP:            1,
	MatB:              1,
	MatH:              1,
	MatS:              1,
	MatD:              1,
	MatQ:              1,
	MatV:              1,
	MatVStatic:        1,
	MatVElement:       2,
	MatVElementStatic: 1,
	MatVStaticElement: 2,
	MatRegList:        1,
	MatRegListStatic:  1,
	MatRegListElement: 2,
	MatOffset:         1,
	MatRefBase:        1,
	MatRefOffset:      2,
	MatRefPre:         2,
	MatRefIndex:       4,
	MatLitMod:         1,
	MatMod:            2,
	MatEnd:            0,
}
