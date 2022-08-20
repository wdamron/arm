package arm

// Modifier identifiers
const (
	_ uint8 = iota
	SymLSL
	SymLSR
	SymASR
	SymROR
	SymMSL
	SymSXTX
	SymSXTW
	SymSXTH
	SymSXTB
	SymUXTX
	SymUXTW
	SymUXTH
	SymUXTB
)

// Modifier lists
const (
	_ uint8 = iota
	SymExtends
	SymExtendsW
	SymExtendsX
	SymShifts
	SymRotates
)

// Modifiers
var (
	ModSXTX = Mod{ID: SymSXTX} // ModSXTX is an extension modifier argument with ID SymSXTX
	ModSXTW = Mod{ID: SymSXTW} // ModSXTW is an extension modifier argument with ID SymSXTW
	ModSXTH = Mod{ID: SymSXTH} // ModSXTH is an extension modifier argument with ID SymSXTH
	ModSXTB = Mod{ID: SymSXTB} // ModSXTB is an extension modifier argument with ID SymSXTB
	ModUXTX = Mod{ID: SymUXTX} // ModUXTX is an extension modifier argument with ID SymUXTX
	ModUXTW = Mod{ID: SymUXTW} // ModUXTW is an extension modifier argument with ID SymUXTW
	ModUXTH = Mod{ID: SymUXTH} // ModUXTH is an extension modifier argument with ID SymUXTH
	ModUXTB = Mod{ID: SymUXTB} // ModUXTB is an extension modifier argument with ID SymUXTB

	ModLSL = Mod{ID: SymLSL} // ModLSL is a shift modifier argument with ID SymLSL.
	ModLSR = Mod{ID: SymLSR} // ModLSR is a shift modifier argument with ID SymLSR.
	ModASR = Mod{ID: SymASR} // ModASR is a shift modifier argument with ID SymASR.
	ModROR = Mod{ID: SymROR} // ModROR is a rotate modifier argument with ID SymROR.
	ModMSL = Mod{ID: SymMSL} // ModMSL is a shift modifier argument with ID SymMSL.
)

// Imm constructs a modifier argument from m with an immediate shift or rotate amount.
func (m Mod) Imm(i uint8) Mod { return Mod{ID: m.ID, ImmInv: ^i} }

// HasImm returns true if an immediate shift or rotate amount is set for m.
func (m Mod) HasImm() bool { return m.ImmInv != 0 }

// GetImm returns the immediate shift or rotate amount for m. The amount is only valid if HasImm returns true.
func (m Mod) GetImm() uint8 { return ^m.ImmInv }

// ModList contains grouped modifier symbols.
var ModList = [...][]uint8{
	SymExtends:  {SymUXTB, SymUXTH, SymUXTW, SymUXTX, SymSXTB, SymSXTH, SymSXTW, SymSXTX, SymLSL},
	SymExtendsW: {SymUXTB, SymUXTH, SymUXTW, SymSXTB, SymSXTH, SymSXTW},
	SymExtendsX: {SymUXTX, SymSXTX, SymLSL},
	SymShifts:   {SymLSL, SymLSR, SymASR},
	SymRotates:  {SymLSL, SymLSR, SymASR, SymROR},
}

var ModListName = [...]string{
	SymExtends:  "SymExtends",
	SymExtendsW: "SymExtendsW",
	SymExtendsX: "SymExtendsX",
	SymShifts:   "SymShifts",
	SymRotates:  "SymRotates",
}

var ModName = [...]string{
	SymLSL:  "SymLSL",
	SymLSR:  "SymLSR",
	SymASR:  "SymASR",
	SymROR:  "SymROR",
	SymSXTX: "SymSXTX",
	SymSXTW: "SymSXTW",
	SymSXTH: "SymSXTH",
	SymSXTB: "SymSXTB",
	SymUXTX: "SymUXTX",
	SymUXTW: "SymUXTW",
	SymUXTH: "SymUXTH",
	SymUXTB: "SymUXTB",
	SymMSL:  "SymMSL",
}

var ModRequiresImm = [...]bool{
	SymLSL: true,
	SymLSR: true,
	SymASR: true,
	SymROR: true,
	SymMSL: true,
}

func checkMod(list []uint8, id uint8) bool {
	for _, x := range list {
		if id == x {
			return true
		}
	}
	return false
}
