package arm

// Condition code symbols
const (
	EQ Symbol = 0 + 1
	NE Symbol = 1 + 1
	CS Symbol = 2 + 1
	HS Symbol = 2 + 1
	CC Symbol = 3 + 1
	LO Symbol = 3 + 1
	MI Symbol = 4 + 1
	PL Symbol = 5 + 1
	VS Symbol = 6 + 1
	VC Symbol = 7 + 1
	HI Symbol = 8 + 1
	LS Symbol = 9 + 1
	GE Symbol = 10 + 1
	LT Symbol = 11 + 1
	GT Symbol = 12 + 1
	LE Symbol = 13 + 1
	AL Symbol = 14 + 1
	NV Symbol = 15 + 1
)

const (
	_ Symbol = iota

	// Literals

	RCTX     // restriction by context
	IVAU     // instruction cache line invalidation
	INVERTED // inverted moves
	LOGICAL  // logical moves
	CSYNC    // profile/trace synchronization

	C0  // CONTROLREGS (Control Registers)
	C1  // CONTROLREGS (Control Registers)
	C2  // CONTROLREGS (Control Registers)
	C3  // CONTROLREGS (Control Registers)
	C4  // CONTROLREGS (Control Registers)
	C5  // CONTROLREGS (Control Registers)
	C6  // CONTROLREGS (Control Registers)
	C7  // CONTROLREGS (Control Registers)
	C8  // CONTROLREGS (Control Registers)
	C9  // CONTROLREGS (Control Registers)
	C10 // CONTROLREGS (Control Registers)
	C11 // CONTROLREGS (Control Registers)
	C12 // CONTROLREGS (Control Registers)
	C13 // CONTROLREGS (Control Registers)
	C14 // CONTROLREGS (Control Registers)
	C15 // CONTROLREGS (Control Registers)

	S1E1R  // ATOPS (Address Translation)
	S1E1W  // ATOPS (Address Translation)
	S1E0R  // ATOPS (Address Translation)
	S1E0W  // ATOPS (Address Translation)
	S1E2R  // ATOPS (Address Translation)
	S1E2W  // ATOPS (Address Translation)
	S12E1R // ATOPS (Address Translation)
	S12E1W // ATOPS (Address Translation)
	S12E0R // ATOPS (Address Translation)
	S12E0W // ATOPS (Address Translation)
	S1E3R  // ATOPS (Address Translation)
	S1E3W  // ATOPS (Address Translation)
	S1E1RP // ATOPS (Address Translation)
	S1E1WP // ATOPS (Address Translation)

	IALLUIS // ICOPS (Instruction Cache)
	IALLU   // ICOPS (Instruction Cache)

	IVAC  // DCOPS (Data Cache)
	ISW   // DCOPS (Data Cache)
	CSW   // DCOPS (Data Cache)
	CISW  // DCOPS (Data Cache)
	ZVA   // DCOPS (Data Cache)
	CVAC  // DCOPS (Data Cache)
	CVAU  // DCOPS (Data Cache)
	CIVAC // DCOPS (Data Cache)
	CVAP  // DCOPS (Data Cache)

	SY    // BARRIEROPS (Instruction/Data Synchronization Barriers)
	ST    // BARRIEROPS (Instruction/Data Synchronization Barriers)
	LD    // BARRIEROPS (Instruction/Data Synchronization Barriers)
	ISH   // BARRIEROPS (Instruction/Data Synchronization Barriers)
	ISHST // BARRIEROPS (Instruction/Data Synchronization Barriers)
	ISHLD // BARRIEROPS (Instruction/Data Synchronization Barriers)
	NSH   // BARRIEROPS (Instruction/Data Synchronization Barriers)
	NSHST // BARRIEROPS (Instruction/Data Synchronization Barriers)
	NSHLD // BARRIEROPS (Instruction/Data Synchronization Barriers)
	OSH   // BARRIEROPS (Instruction/Data Synchronization Barriers)
	OSHST // BARRIEROPS (Instruction/Data Synchronization Barriers)
	OSHLD // BARRIEROPS (Instruction/Data Synchronization Barriers)

	SPSEL        // MSRIMMOPS (System Registers)
	DAIFSET      // MSRIMMOPS (System Registers)
	DAIFCLR      // MSRIMMOPS (System Registers)
	UAO          // MSRIMMOPS (System Registers)
	PAN          // MSRIMMOPS (System Registers)
	DIT          // MSRIMMOPS (System Registers)
	VMALLE1IS    // MSRIMMOPS (System Registers)
	VAE1IS       // MSRIMMOPS (System Registers)
	ASIDE1IS     // MSRIMMOPS (System Registers)
	VAAE1IS      // MSRIMMOPS (System Registers)
	VALE1IS      // MSRIMMOPS (System Registers)
	VAALE1IS     // MSRIMMOPS (System Registers)
	VMALLE1      // MSRIMMOPS (System Registers)
	VAE1         // MSRIMMOPS (System Registers)
	ASIDE1       // MSRIMMOPS (System Registers)
	VAAE1        // MSRIMMOPS (System Registers)
	VALE1        // MSRIMMOPS (System Registers)
	VAALE1       // MSRIMMOPS (System Registers)
	IPAS2E1IS    // MSRIMMOPS (System Registers)
	IPAS2LE1IS   // MSRIMMOPS (System Registers)
	ALLE2IS      // MSRIMMOPS (System Registers)
	VAE2IS       // MSRIMMOPS (System Registers)
	ALLE1IS      // MSRIMMOPS (System Registers)
	VALE2IS      // MSRIMMOPS (System Registers)
	VMALLS12E1IS // MSRIMMOPS (System Registers)
	IPAS2E1      // MSRIMMOPS (System Registers)
	IPAS2LE1     // MSRIMMOPS (System Registers)
	ALLE2        // MSRIMMOPS (System Registers)
	VAE2         // MSRIMMOPS (System Registers)
	ALLE1        // MSRIMMOPS (System Registers)
	VALE2        // MSRIMMOPS (System Registers)
	VMALLS12E1   // MSRIMMOPS (System Registers)
	ALLE3IS      // MSRIMMOPS (System Registers)
	VAE3IS       // MSRIMMOPS (System Registers)
	VALE3IS      // MSRIMMOPS (System Registers)
	ALLE3        // MSRIMMOPS (System Registers)
	VAE3         // MSRIMMOPS (System Registers)
	VALE3        // MSRIMMOPS (System Registers)
	VMALLE1OS    // MSRIMMOPS (System Registers)
	VAE1OS       // MSRIMMOPS (System Registers)
	ASIDE1OS     // MSRIMMOPS (System Registers)
	VAAE1OS      // MSRIMMOPS (System Registers)
	VALE1OS      // MSRIMMOPS (System Registers)
	VAALE1OS     // MSRIMMOPS (System Registers)
	RVAE1IS      // MSRIMMOPS (System Registers)
	RVAAE1IS     // MSRIMMOPS (System Registers)
	RVALE1IS     // MSRIMMOPS (System Registers)
	RVAALE1IS    // MSRIMMOPS (System Registers)
	RVAE1OS      // MSRIMMOPS (System Registers)
	RVAAE1OS     // MSRIMMOPS (System Registers)
	RVALE1OS     // MSRIMMOPS (System Registers)
	RVAALE1OS    // MSRIMMOPS (System Registers)
	RVAE1        // MSRIMMOPS (System Registers)
	RVAAE1       // MSRIMMOPS (System Registers)
	RVALE1       // MSRIMMOPS (System Registers)
	RVAALE1      // MSRIMMOPS (System Registers)
	RIPAS2E1IS   // MSRIMMOPS (System Registers)
	RIPAS2LE1IS  // MSRIMMOPS (System Registers)
	ALLE2OS      // MSRIMMOPS (System Registers)
	VAE2OS       // MSRIMMOPS (System Registers)
	ALLE1OS      // MSRIMMOPS (System Registers)
	VALE2OS      // MSRIMMOPS (System Registers)

	VMALLS12E1OS // TLBIOPS (Translation Table)
	RVAE2IS      // TLBIOPS (Translation Table)
	RVALE2IS     // TLBIOPS (Translation Table)
	IPAS2E1OS    // TLBIOPS (Translation Table)
	RIPAS2E1     // TLBIOPS (Translation Table)
	RIPAS2E1OS   // TLBIOPS (Translation Table)
	IPAS2LE1OS   // TLBIOPS (Translation Table)
	RIPAS2LE1    // TLBIOPS (Translation Table)
	RIPAS2LE1OS  // TLBIOPS (Translation Table)
	RVAE2OS      // TLBIOPS (Translation Table)
	RVALE2OS     // TLBIOPS (Translation Table)
	RVAE2        // TLBIOPS (Translation Table)
	RVALE2       // TLBIOPS (Translation Table)
	ALLE3OS      // TLBIOPS (Translation Table)
	VAE3OS       // TLBIOPS (Translation Table)
	VALE3OS      // TLBIOPS (Translation Table)
	RVAE3IS      // TLBIOPS (Translation Table)
	RVALE3IS     // TLBIOPS (Translation Table)
	RVAE3OS      // TLBIOPS (Translation Table)
	RVALE3OS     // TLBIOPS (Translation Table)
	RVAE3        // TLBIOPS (Translation Table)
	RVALE3       // TLBIOPS (Translation Table)
)

// Address Translation
var ATOPS = [...]Symbol{S1E1R, S1E1W, S1E0R, S1E0W, S1E2R, S1E2W, S12E1R, S12E1W, S12E0R, S12E0W, S1E3R, S1E3W, S1E1RP, S1E1WP}

// Instruction Cache
var ICOPS = [...]Symbol{IALLUIS, IALLU}

// Data Cache
var DCOPS = [...]Symbol{IVAC, ISW, CSW, CISW, ZVA, CVAC, CVAU, CIVAC, CVAP}

// Instruction/Data Synchronization Barriers
var BARRIEROPS = [...]Symbol{SY, ST, LD, ISH, ISHST, ISHLD, NSH, NSHST, NSHLD, OSH, OSHST, OSHLD}

// System Registers
var MSRIMMOPS = [...]Symbol{SPSEL, DAIFSET, DAIFCLR, UAO, PAN, DIT}

// Translation Table
var TLBIOPS = [...]Symbol{VMALLE1IS, VAE1IS, ASIDE1IS, VAAE1IS, VALE1IS, VAALE1IS, VMALLE1, VAE1, ASIDE1, VAAE1, VALE1, VAALE1, IPAS2E1IS, IPAS2LE1IS, ALLE2IS, VAE2IS, ALLE1IS, VALE2IS, VMALLS12E1IS, IPAS2E1, IPAS2LE1, ALLE2, VAE2, ALLE1, VALE2, VMALLS12E1, ALLE3IS, VAE3IS, VALE3IS, ALLE3, VAE3, VALE3, VMALLE1OS, VAE1OS, ASIDE1OS, VAAE1OS, VALE1OS, VAALE1OS, RVAE1IS, RVAAE1IS, RVALE1IS, RVAALE1IS, RVAE1OS, RVAAE1OS, RVALE1OS, RVAALE1OS, RVAE1, RVAAE1, RVALE1, RVAALE1, RIPAS2E1IS, RIPAS2LE1IS, ALLE2OS, VAE2OS, ALLE1OS, VALE2OS, VMALLS12E1OS, RVAE2IS, RVALE2IS, IPAS2E1OS, RIPAS2E1, RIPAS2E1OS, IPAS2LE1OS, RIPAS2LE1, RIPAS2LE1OS, RVAE2OS, RVALE2OS, RVAE2, RVALE2, ALLE3OS, VAE3OS, VALE3OS, RVAE3IS, RVALE3IS, RVAE3OS, RVALE3OS, RVAE3, RVALE3}

func symListContains(listSym uint8, arg Symbol) bool {
	var list []Symbol
	switch listSym {
	default:
		return false
	case SymCONTROLREGS:
		if arg < C0 || arg > C15 {
			return false
		}
		return true
	case SymATOPS:
		list = ATOPS[:]
	case SymDCOPS:
		list = DCOPS[:]
	case SymICOPS:
		list = ICOPS[:]
	case SymTLBIOPS:
		list = TLBIOPS[:]
	case SymBARRIEROPS:
		list = BARRIEROPS[:]
	case SymMSRIMMOPS:
		list = MSRIMMOPS[:]
	}
	for _, x := range list {
		if arg == x {
			return true
		}
	}
	return false
}

var SymbolValue = [...]uint16{
	RCTX:         uint16(RCTX),
	IVAU:         uint16(IVAU),
	INVERTED:     uint16(INVERTED),
	LOGICAL:      uint16(LOGICAL),
	CSYNC:        uint16(CSYNC),
	C0:           uint16(C0),
	C1:           uint16(C1 - C0),
	C2:           uint16(C2 - C0),
	C3:           uint16(C3 - C0),
	C4:           uint16(C4 - C0),
	C5:           uint16(C5 - C0),
	C6:           uint16(C6 - C0),
	C7:           uint16(C7 - C0),
	C8:           uint16(C8 - C0),
	C9:           uint16(C9 - C0),
	C10:          uint16(C10 - C0),
	C11:          uint16(C11 - C0),
	C12:          uint16(C12 - C0),
	C13:          uint16(C13 - C0),
	C14:          uint16(C14 - C0),
	C15:          uint16(C15 - C0),
	S1E1R:        0b00001111000000,
	S1E1W:        0b00001111000001,
	S1E0R:        0b00001111000010,
	S1E0W:        0b00001111000011,
	S1E2R:        0b10001111000000,
	S1E2W:        0b10001111000001,
	S12E1R:       0b10001111000100,
	S12E1W:       0b10001111000101,
	S12E0R:       0b10001111000110,
	S12E0W:       0b10001111000111,
	S1E3R:        0b11001111000000,
	S1E3W:        0b11001111000001,
	S1E1RP:       0b00001111001000,
	S1E1WP:       0b00001111001001,
	IALLUIS:      0b00001110001000,
	IALLU:        0b00001110101000,
	IVAC:         0b00001110110001,
	ISW:          0b00001110110010,
	CSW:          0b00001111010010,
	CISW:         0b00001111110010,
	ZVA:          0b01101110100001,
	CVAC:         0b01101111010001,
	CVAU:         0b01101111011001,
	CIVAC:        0b01101111110001,
	CVAP:         0b01101111100001,
	SY:           0b1111,
	ST:           0b1110,
	LD:           0b1101,
	ISH:          0b1011,
	ISHST:        0b1010,
	ISHLD:        0b1001,
	NSH:          0b0111,
	NSHST:        0b0110,
	NSHLD:        0b0101,
	OSH:          0b0011,
	OSHST:        0b0010,
	OSHLD:        0b0001,
	SPSEL:        0b00001000000101,
	DAIFSET:      0b01101000000110,
	DAIFCLR:      0b01101000000111,
	UAO:          0b00001000000011,
	PAN:          0b00001000000100,
	DIT:          0b01101000000010,
	VMALLE1IS:    0b00010000011000,
	VAE1IS:       0b00010000011001,
	ASIDE1IS:     0b00010000011010,
	VAAE1IS:      0b00010000011011,
	VALE1IS:      0b00010000011101,
	VAALE1IS:     0b00010000011111,
	VMALLE1:      0b00010000111000,
	VAE1:         0b00010000111001,
	ASIDE1:       0b00010000111010,
	VAAE1:        0b00010000111011,
	VALE1:        0b00010000111101,
	VAALE1:       0b00010000111111,
	IPAS2E1IS:    0b10010000000001,
	IPAS2LE1IS:   0b10010000000101,
	ALLE2IS:      0b10010000011000,
	VAE2IS:       0b10010000011001,
	ALLE1IS:      0b10010000011100,
	VALE2IS:      0b10010000011101,
	VMALLS12E1IS: 0b10010000011110,
	IPAS2E1:      0b10010000100001,
	IPAS2LE1:     0b10010000100101,
	ALLE2:        0b10010000111000,
	VAE2:         0b10010000111001,
	ALLE1:        0b10010000111100,
	VALE2:        0b10010000111101,
	VMALLS12E1:   0b10010000111110,
	ALLE3IS:      0b11010000011000,
	VAE3IS:       0b11010000011001,
	VALE3IS:      0b11010000011101,
	ALLE3:        0b11010000111000,
	VAE3:         0b11010000111001,
	VALE3:        0b11010000111101,
	VMALLE1OS:    0b00010000001000,
	VAE1OS:       0b00010000001001,
	ASIDE1OS:     0b00010000001010,
	VAAE1OS:      0b00010000001011,
	VALE1OS:      0b00010000001101,
	VAALE1OS:     0b00010000001111,
	RVAE1IS:      0b00010000010001,
	RVAAE1IS:     0b00010000010011,
	RVALE1IS:     0b00010000010101,
	RVAALE1IS:    0b00010000010111,
	RVAE1OS:      0b00010000101001,
	RVAAE1OS:     0b00010000101011,
	RVALE1OS:     0b00010000101101,
	RVAALE1OS:    0b00010000101111,
	RVAE1:        0b00010000110001,
	RVAAE1:       0b00010000110011,
	RVALE1:       0b00010000110101,
	RVAALE1:      0b00010000110111,
	RIPAS2E1IS:   0b10010000000010,
	RIPAS2LE1IS:  0b10010000000110,
	ALLE2OS:      0b10010000001000,
	VAE2OS:       0b10010000001001,
	ALLE1OS:      0b10010000001100,
	VALE2OS:      0b10010000001101,
	VMALLS12E1OS: 0b10010000001110,
	RVAE2IS:      0b10010000010001,
	RVALE2IS:     0b10010000010101,
	IPAS2E1OS:    0b10010000100000,
	RIPAS2E1:     0b10010000100010,
	RIPAS2E1OS:   0b10010000100011,
	IPAS2LE1OS:   0b10010000100100,
	RIPAS2LE1:    0b10010000100110,
	RIPAS2LE1OS:  0b10010000100111,
	RVAE2OS:      0b10010000101001,
	RVALE2OS:     0b10010000101101,
	RVAE2:        0b10010000110001,
	RVALE2:       0b10010000110101,
	ALLE3OS:      0b11010000001000,
	VAE3OS:       0b11010000001001,
	VALE3OS:      0b11010000001101,
	RVAE3IS:      0b11010000010001,
	RVALE3IS:     0b11010000010101,
	RVAE3OS:      0b11010000101001,
	RVALE3OS:     0b11010000101101,
	RVAE3:        0b11010000110001,
	RVALE3:       0b11010000110101,
}

var SymbolName = [...]string{
	RCTX:         "RCTX",
	IVAU:         "IVAU",
	INVERTED:     "INVERTED",
	LOGICAL:      "LOGICAL",
	CSYNC:        "CSYNC",
	C0:           "C0",
	C1:           "C1",
	C2:           "C2",
	C3:           "C3",
	C4:           "C4",
	C5:           "C5",
	C6:           "C6",
	C7:           "C7",
	C8:           "C8",
	C9:           "C9",
	C10:          "C10",
	C11:          "C11",
	C12:          "C12",
	C13:          "C13",
	C14:          "C14",
	C15:          "C15",
	S1E1R:        "S1E1R",
	S1E1W:        "S1E1W",
	S1E0R:        "S1E0R",
	S1E0W:        "S1E0W",
	S1E2R:        "S1E2R",
	S1E2W:        "S1E2W",
	S12E1R:       "S12E1R",
	S12E1W:       "S12E1W",
	S12E0R:       "S12E0R",
	S12E0W:       "S12E0W",
	S1E3R:        "S1E3R",
	S1E3W:        "S1E3W",
	S1E1RP:       "S1E1RP",
	S1E1WP:       "S1E1WP",
	IALLUIS:      "IALLUIS",
	IALLU:        "IALLU",
	IVAC:         "IVAC",
	ISW:          "ISW",
	CSW:          "CSW",
	CISW:         "CISW",
	ZVA:          "ZVA",
	CVAC:         "CVAC",
	CVAU:         "CVAU",
	CIVAC:        "CIVAC",
	CVAP:         "CVAP",
	SY:           "SY",
	ST:           "ST",
	LD:           "LD",
	ISH:          "ISH",
	ISHST:        "ISHST",
	ISHLD:        "ISHLD",
	NSH:          "NSH",
	NSHST:        "NSHST",
	NSHLD:        "NSHLD",
	OSH:          "OSH",
	OSHST:        "OSHST",
	OSHLD:        "OSHLD",
	SPSEL:        "SPSEL",
	DAIFSET:      "DAIFSET",
	DAIFCLR:      "DAIFCLR",
	UAO:          "UAO",
	PAN:          "PAN",
	DIT:          "DIT",
	VMALLE1IS:    "VMALLE1IS",
	VAE1IS:       "VAE1IS",
	ASIDE1IS:     "ASIDE1IS",
	VAAE1IS:      "VAAE1IS",
	VALE1IS:      "VALE1IS",
	VAALE1IS:     "VAALE1IS",
	VMALLE1:      "VMALLE1",
	VAE1:         "VAE1",
	ASIDE1:       "ASIDE1",
	VAAE1:        "VAAE1",
	VALE1:        "VALE1",
	VAALE1:       "VAALE1",
	IPAS2E1IS:    "IPAS2E1IS",
	IPAS2LE1IS:   "IPAS2LE1IS",
	ALLE2IS:      "ALLE2IS",
	VAE2IS:       "VAE2IS",
	ALLE1IS:      "ALLE1IS",
	VALE2IS:      "VALE2IS",
	VMALLS12E1IS: "VMALLS12E1IS",
	IPAS2E1:      "IPAS2E1",
	IPAS2LE1:     "IPAS2LE1",
	ALLE2:        "ALLE2",
	VAE2:         "VAE2",
	ALLE1:        "ALLE1",
	VALE2:        "VALE2",
	VMALLS12E1:   "VMALLS12E1",
	ALLE3IS:      "ALLE3IS",
	VAE3IS:       "VAE3IS",
	VALE3IS:      "VALE3IS",
	ALLE3:        "ALLE3",
	VAE3:         "VAE3",
	VALE3:        "VALE3",
	VMALLE1OS:    "VMALLE1OS",
	VAE1OS:       "VAE1OS",
	ASIDE1OS:     "ASIDE1OS",
	VAAE1OS:      "VAAE1OS",
	VALE1OS:      "VALE1OS",
	VAALE1OS:     "VAALE1OS",
	RVAE1IS:      "RVAE1IS",
	RVAAE1IS:     "RVAAE1IS",
	RVALE1IS:     "RVALE1IS",
	RVAALE1IS:    "RVAALE1IS",
	RVAE1OS:      "RVAE1OS",
	RVAAE1OS:     "RVAAE1OS",
	RVALE1OS:     "RVALE1OS",
	RVAALE1OS:    "RVAALE1OS",
	RVAE1:        "RVAE1",
	RVAAE1:       "RVAAE1",
	RVALE1:       "RVALE1",
	RVAALE1:      "RVAALE1",
	RIPAS2E1IS:   "RIPAS2E1IS",
	RIPAS2LE1IS:  "RIPAS2LE1IS",
	ALLE2OS:      "ALLE2OS",
	VAE2OS:       "VAE2OS",
	ALLE1OS:      "ALLE1OS",
	VALE2OS:      "VALE2OS",
	VMALLS12E1OS: "VMALLS12E1OS",
	RVAE2IS:      "RVAE2IS",
	RVALE2IS:     "RVALE2IS",
	IPAS2E1OS:    "IPAS2E1OS",
	RIPAS2E1:     "RIPAS2E1",
	RIPAS2E1OS:   "RIPAS2E1OS",
	IPAS2LE1OS:   "IPAS2LE1OS",
	RIPAS2LE1:    "RIPAS2LE1",
	RIPAS2LE1OS:  "RIPAS2LE1OS",
	RVAE2OS:      "RVAE2OS",
	RVALE2OS:     "RVALE2OS",
	RVAE2:        "RVAE2",
	RVALE2:       "RVALE2",
	ALLE3OS:      "ALLE3OS",
	VAE3OS:       "VAE3OS",
	VALE3OS:      "VALE3OS",
	RVAE3IS:      "RVAE3IS",
	RVALE3IS:     "RVALE3IS",
	RVAE3OS:      "RVAE3OS",
	RVALE3OS:     "RVALE3OS",
	RVAE3:        "RVAE3",
	RVALE3:       "RVALE3",
}
