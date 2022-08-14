package arm

import "math"

// flattenArgs unnests matched arguments to an internal form before encoding.
func (a *Assembler) flattenArgs() {
	var cursor int
	for _, m := range a.Pattern() {
		if m.Op == MatEnd {
			continue
		}
		flatArgCount := int(MatcherFlatArgCounts[m.Op])
		argCount := len(a.Flat)
		if cursor < len(a.Args) {
			switch arg := a.Args[cursor].(type) {
			case Reg:
				a.appendFlat(FlatReg(arg.ID))
				if arg.HasElem() && m.Op != MatVElementStatic {
					a.appendFlat(FlatImm(arg.GetElem()))
				}
			case RegList:
				a.appendFlat(FlatReg(arg.First.ID))
				if arg.First.HasElem() && m.Op != MatVElementStatic {
					a.appendFlat(FlatImm(arg.First.GetElem()))
				}
			case Imm:
				a.appendFlat(FlatImm(arg))
			case Float:
				a.appendFlat(FlatImm(math.Float32bits(float32(arg))))
			case Wide:
				a.appendFlat(FlatImm(arg))
			case Ref:
				a.appendFlat(FlatReg(arg.Base.ID))
			case RefOffset:
				a.appendFlat(FlatReg(arg.Base.ID), FlatImm(arg.Offset))
			case RefPreIndexed:
				a.appendFlat(FlatReg(arg.Base.ID), FlatImm(arg.Offset))
			case RefIndexed:
				a.appendFlat(FlatReg(arg.Base.ID), FlatReg(arg.Idx.ID))
				if arg.Mod.ID != 0 {
					a.appendFlat(FlatMod(arg.Mod.ID))
					if arg.Mod.HasImm() {
						a.appendFlat(FlatImm(arg.Mod.GetImm()))
					}
				}
			case Mod:
				if flatArgCount >= 2 {
					a.appendFlat(FlatMod(arg.ID))
				}
				if arg.HasImm() {
					a.appendFlat(FlatImm(arg.GetImm()))
				}
			case Label:
				a.appendFlat(FlatLabel(arg))
			case Symbol:
				switch arg {
				case INVERTED, LOGICAL: // skip
				default:
					a.appendFlat(FlatImm(arg))
				}
			}
		}

		for added := len(a.Flat) - argCount; added < flatArgCount; added++ {
			a.appendFlat(FlatDefault{})
		}

		cursor++
	}
}

func (a *Assembler) appendFlat(flat ...Flat) {
	a.Flat = append(a.Flat, flat...)
}
