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
	ModSXTX = Mod{ID: SymSXTX}
	ModSXTW = Mod{ID: SymSXTW}
	ModSXTH = Mod{ID: SymSXTH}
	ModSXTB = Mod{ID: SymSXTB}
	ModUXTX = Mod{ID: SymUXTX}
	ModUXTW = Mod{ID: SymUXTW}
	ModUXTH = Mod{ID: SymUXTH}
	ModUXTB = Mod{ID: SymUXTB}

	ModLSL = Mod{ID: SymLSL}
	ModLSR = Mod{ID: SymLSR}
	ModASR = Mod{ID: SymASR}
	ModROR = Mod{ID: SymROR}
	ModMSL = Mod{ID: SymMSL}
)

func (m Mod) Imm(i uint8) Mod { return Mod{ID: m.ID, ImmInv: ^i} }
func (m Mod) HasImm() bool    { return m.ImmInv != 0 }
func (m Mod) GetImm() uint8   { return ^m.ImmInv }

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
