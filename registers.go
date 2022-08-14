package arm

type Size uint8

const (
	_ Size = iota

	BYTE  // 8-bit register elements
	WORD  // 16-bit register elements
	DWORD // 32-bit register elements
	QWORD // 64-bit register elements
	OWORD // 128-bit register element
)

var SizeName = [...]string{BYTE: "BYTE", WORD: "WORD", DWORD: "DWORD", QWORD: "QWORD", OWORD: "OWORD"}

type RegFamily uint8

const (
	_ RegFamily = iota

	RegInt    // X and W registers (except SP)
	RegSP     // stack-pointer registers
	RegFloat  // scalar simd registers
	RegVec32  // 32-bit vector registers
	RegVec64  // 64-bit vector registers
	RegVec128 // 128-bit vector registers
)

var (
	WZR = Reg{ID: 31, Type: RW} // 32-bit zero register
	XZR = Reg{ID: 31, Type: RX} // 64-bit zero register

	WSP = Reg{ID: 31, Type: RWSP} // 32-bit stack pointer register
	XSP = Reg{ID: 31, Type: RXSP} // 64-bit stack pointer register
)

func W(id uint8) Reg { return Reg{ID: id, Type: RW} }
func X(id uint8) Reg { return Reg{ID: id, Type: RX} }

func ScalarB(id uint8) Reg { return Reg{ID: id, Type: RB} }
func ScalarH(id uint8) Reg { return Reg{ID: id, Type: RH} }
func ScalarS(id uint8) Reg { return Reg{ID: id, Type: RS} }
func ScalarD(id uint8) Reg { return Reg{ID: id, Type: RD} }
func ScalarQ(id uint8) Reg { return Reg{ID: id, Type: RQ} }

func Vec4B(id uint8) Reg  { return Reg{ID: id, Type: V4B} }
func Vec8B(id uint8) Reg  { return Reg{ID: id, Type: V8B} }
func Vec16B(id uint8) Reg { return Reg{ID: id, Type: V16B} }
func Vec2H(id uint8) Reg  { return Reg{ID: id, Type: V2H} }
func Vec4H(id uint8) Reg  { return Reg{ID: id, Type: V4H} }
func Vec8H(id uint8) Reg  { return Reg{ID: id, Type: V8H} }
func Vec2S(id uint8) Reg  { return Reg{ID: id, Type: V2S} }
func Vec4S(id uint8) Reg  { return Reg{ID: id, Type: V4S} }
func Vec1D(id uint8) Reg  { return Reg{ID: id, Type: V1D} }
func Vec2D(id uint8) Reg  { return Reg{ID: id, Type: V2D} }
func Vec1O(id uint8) Reg  { return Reg{ID: id, Type: V1O} }

func (r Reg) I(idx uint8) Reg           { return Reg{ID: r.ID, Type: r.Type, Elem: ^idx} }
func (r Reg) List(length uint8) RegList { return RegList{First: r, Len: length} }

func (r RegList) I(idx uint8) RegList { return RegList{First: r.First.I(idx), Len: r.Len} }

func (r Reg) Family() RegFamily { return r.Type.Family() }
func (r Reg) ElemSize() Size    { return r.Type.Elem() }
func (r Reg) Lanes() uint8      { return r.Type.Lanes() }
func (r Reg) HasElem() bool     { return r.Elem != 0 }
func (r Reg) GetElem() uint8    { return ^r.Elem }
func (r Reg) IsVec() bool {
	switch r.Family() {
	case RegVec32, RegVec64, RegVec128:
		return true
	default:
		return false
	}
}

// RegType indicates an element size, family, and lane count for a register.
type RegType uint8

// Register element size, family, and lane count indicators
const (
	RW RegType = RegType(DWORD) | RegType(RegInt<<4)
	RX RegType = RegType(QWORD) | RegType(RegInt<<4)

	RWSP RegType = RegType(DWORD) | RegType(RegSP<<4)
	RXSP RegType = RegType(QWORD) | RegType(RegSP<<4)

	RB RegType = RegType(BYTE) | RegType(RegFloat<<4)
	RH RegType = RegType(WORD) | RegType(RegFloat<<4)
	RS RegType = RegType(DWORD) | RegType(RegFloat<<4)
	RD RegType = RegType(QWORD) | RegType(RegFloat<<4)
	RQ RegType = RegType(OWORD) | RegType(RegFloat<<4)

	V4B  RegType = RegType(BYTE) | RegType(RegVec32<<4)
	V8B  RegType = RegType(BYTE) | RegType(RegVec64<<4)
	V16B RegType = RegType(BYTE) | RegType(RegVec128<<4)
	V2H  RegType = RegType(WORD) | RegType(RegVec32<<4)
	V4H  RegType = RegType(WORD) | RegType(RegVec64<<4)
	V8H  RegType = RegType(WORD) | RegType(RegVec128<<4)
	V2S  RegType = RegType(DWORD) | RegType(RegVec64<<4)
	V4S  RegType = RegType(DWORD) | RegType(RegVec128<<4)
	V1D  RegType = RegType(QWORD) | RegType(RegVec64<<4)
	V2D  RegType = RegType(QWORD) | RegType(RegVec128<<4)
	V1O  RegType = RegType(OWORD) | RegType(RegVec128<<4)
)

func (sz RegType) Family() RegFamily { return RegFamily(sz >> 4) }
func (sz RegType) Elem() Size        { return Size(sz) & 0xF } // BYTE..OWORD
func (sz RegType) ElemBytes() uint8  { return 1 << (sz.Elem() - BYTE) }
func (sz RegType) Bytes() uint8 {
	switch sz.Family() {
	case RegVec32:
		return 4
	case RegVec64:
		return 8
	case RegVec128:
		return 16
	default:
		return sz.ElemBytes()
	}
}
func (sz RegType) Lanes() uint8 {
	elemSize := sz.ElemBytes()
	if elemSize == 0 {
		return 0
	}
	switch sz.Family() {
	case RegVec32:
		return 4 / elemSize
	case RegVec64:
		return 8 / elemSize
	case RegVec128:
		return 16 / elemSize
	default:
		return 0
	}
}
