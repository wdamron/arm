package arm

// matchPattern returns true if the encoding at the current iterator position matches.
func (a *Assembler) matchPattern() bool {
	a.SimdSize = 0
	args, pattern, required, optional := a.Args, a.Pattern(), int(a.patternLen), 0
	for i, m := range pattern {
		if m.Op == MatEnd {
			required = i
			optional = int(a.patternLen) - (i + 1) // skip MatEnd
			break
		}
	}
	if len(args) != required && len(args) != required+optional {
		return false
	}
	if !a.matchArgSlice(args[:required], pattern[:required]) {
		return false
	}
	if len(args) == required || optional == 0 {
		return true
	}
	return a.matchArgSlice(args[required:], pattern[required+1:]) // skip MatEnd
}

func (a *Assembler) matchArgSlice(args []Arg, pattern []EncOp) bool {
	for i, m := range pattern {
		if !a.matchArg(args[i], m) {
			return false
		}
	}
	return true
}

func (a *Assembler) matchArg(arg Arg, m EncOp) bool {
	switch arg := arg.(type) {
	case Reg:
		if !checkReg(arg) {
			return false
		}
		switch m.Op {
		// scalar
		case MatW:
			return arg.Type == RW
		case MatX:
			return arg.Type == RX
		case MatWSP:
			return arg.Type == RWSP || (arg.Type == RW && arg != WZR)
		case MatXSP:
			return arg.Type == RXSP || (arg.Type == RX && arg != XZR)
		// simd
		case MatB:
			return arg.Type == RB
		case MatH:
			return arg.Type == RH
		case MatS:
			return arg.Type == RS
		case MatD:
			return arg.Type == RD
		case MatQ:
			return arg.Type == RQ
		case MatV:
			return !arg.HasElem() && arg.IsVec() && arg.ElemSize() == Size(m.X[0]) && a.matchOrSetSimdSize(arg)
		case MatVStatic:
			return !arg.HasElem() && arg.IsVec() && arg.ElemSize() == Size(m.X[0]) && arg.Lanes() == m.X[1]
		case MatVStaticElement:
			return arg.HasElem() && arg.ElemSize() == Size(m.X[0]) && arg.Lanes() == m.X[1]
		case MatVElement:
			return arg.HasElem() && arg.ElemSize() == Size(m.X[0])
		case MatVElementStatic:
			return arg.HasElem() && arg.ElemSize() == Size(m.X[0]) && arg.GetElem() == m.X[1]
		}

	case RegList:
		if !checkReg(arg.First) {
			return false
		}
		switch m.Op {
		case MatRegList:
			return !arg.First.HasElem() && arg.Len == m.X[0] && arg.First.ElemSize() == Size(m.X[1]) && a.matchOrSetSimdSize(arg.First)
		case MatRegListStatic:
			return !arg.First.HasElem() && arg.Len == m.X[0] && arg.First.ElemSize() == Size(m.X[1]) && arg.First.Lanes() == m.X[2]
		case MatRegListElement:
			return arg.First.HasElem() && arg.Len == m.X[0] && arg.First.ElemSize() == Size(m.X[1])
		}

	case Imm:
		switch m.Op {
		case MatImm, MatOffset:
			return true
		case MatLitInt:
			return arg == Imm(m.X[0])
		}

	case Wide:
		switch m.Op {
		case MatImm, MatOffset:
			return true
		case MatLitInt:
			return arg == Wide(m.X[0])
		}

	case Float:
		switch m.Op {
		case MatFloat:
			return true
		case MatLitFloat:
			return arg == Float(m.X[0])
		}

	case Mod:
		switch m.Op {
		case MatMod:
			return checkMod(ModList[m.X[0]], arg.ID)
		case MatLitMod:
			return arg.ID == m.X[0]
		}

	case Ref:
		return (m.Op == MatRefBase || m.Op == MatRefOffset) && checkReg(arg.Base) && checkRefBase(arg.Base)

	case RefOffset:
		return m.Op == MatRefOffset && checkReg(arg.Base) && checkRefBase(arg.Base)

	case RefPreIndexed:
		return m.Op == MatRefPre && checkReg(arg.Base) && checkRefBase(arg.Base)

	case RefIndexed:
		return m.Op == MatRefIndex && checkReg(arg.Base) && checkReg(arg.Idx) && checkRefBase(arg.Base) && arg.Idx.Family() == RegInt

	case Label:
		if int(arg.ID) >= len(a.LabelPC) {
			return false
		}
		return m.Op == MatOffset

	case Symbol:
		switch m.Op {
		case MatSymbol, MatCond:
			return true
		case MatLitSymbol:
			return arg == Symbol(m.X[0])
		}
	}
	return false
}

func checkReg(r Reg) bool {
	switch r.Family() {
	case RegInt:
		switch r.Type {
		case RW, RX:
			return r.ID < 32 && !r.HasElem()
		}
	case RegSP:
		switch r.Type {
		case RWSP, RXSP:
			return r.ID == 31 && !r.HasElem()
		}
	case RegFloat:
		switch r.Type {
		case RB, RH, RS, RD, RQ:
			return r.ID < 32 && !r.HasElem()
		}
	case RegVec32:
		switch r.Type {
		case V4B, V2H:
			return r.ID < 32 && (!r.HasElem() || r.GetElem() < r.Lanes())
		}
	case RegVec64:
		switch r.Type {
		case V8B, V4H, V2S, V1D:
			return r.ID < 32 && (!r.HasElem() || r.GetElem() < r.Lanes())
		}
	case RegVec128:
		switch r.Type {
		case V16B, V8H, V4S, V2D, V1O:
			return r.ID < 32 && (!r.HasElem() || r.GetElem() < r.Lanes())
		}
	}
	return false
}

func checkRefBase(r Reg) bool { return r.Family() == RegInt || r.Family() == RegSP }

func (a *Assembler) matchOrSetSimdSize(reg Reg) bool {
	switch reg.Family() {
	case RegInt, RegSP, RegFloat:
		return true
	default:
		regSize := reg.Type.Bytes()
		if a.SimdSize != 0 {
			return a.SimdSize == regSize
		}
		a.SimdSize = regSize
		return true
	}
}
