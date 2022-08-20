package arm

// Size represents the width of a register elememt.
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

// RegFamily represents the family of a register (integer, SP, scalar, vector).
type RegFamily uint8

const (
	_ RegFamily = iota

	RegInt    // X and W registers (except stack pointer)
	RegSP     // stack pointer registers
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

// W constructs a 32-bit integer register, with type RW and family RegInt.
func W(id uint8) Reg { return Reg{ID: id, Type: RW} }

// X constructs a 64-bit integer register, with type RX and family RegInt.
func X(id uint8) Reg { return Reg{ID: id, Type: RX} }

// ScalarB constructs an 8-bit scalar SIMD register, with type RB and family RegFloat.
func ScalarB(id uint8) Reg { return Reg{ID: id, Type: RB} }

// ScalarH constructs a 16-bit scalar SIMD register, with type RH and family RegFloat.
func ScalarH(id uint8) Reg { return Reg{ID: id, Type: RH} }

// ScalarS constructs a 32-bit scalar SIMD register, with type RS and family RegFloat.
func ScalarS(id uint8) Reg { return Reg{ID: id, Type: RS} }

// ScalarD constructs a 64-bit scalar SIMD register, with type RD and family RegFloat.
func ScalarD(id uint8) Reg { return Reg{ID: id, Type: RD} }

// ScalarQ constructs a 16-bit scalar SIMD register, with type RQ and family RegFloat.
func ScalarQ(id uint8) Reg { return Reg{ID: id, Type: RQ} }

// Vec4B constructs a 4x8-bit vector SIMD register, with type V4B and family RegVec32.
func Vec4B(id uint8) Reg { return Reg{ID: id, Type: V4B} }

// Vec8B constructs an 8x8-bit vector SIMD register, with type V8B and family RegVec64.
func Vec8B(id uint8) Reg { return Reg{ID: id, Type: V8B} }

// Vec16B constructs a 16x8-bit vector SIMD register, with type V16B and family RegVec128.
func Vec16B(id uint8) Reg { return Reg{ID: id, Type: V16B} }

// Vec2H constructs a 2x16-bit vector SIMD register, with type V2H and family RegVec32.
func Vec2H(id uint8) Reg { return Reg{ID: id, Type: V2H} }

// Vec4H constructs a 4x16-bit vector SIMD register, with type V4H and family RegVec64.
func Vec4H(id uint8) Reg { return Reg{ID: id, Type: V4H} }

// Vec8H constructs a 8x16-bit vector SIMD register, with type V8H and family RegVec128.
func Vec8H(id uint8) Reg { return Reg{ID: id, Type: V8H} }

// Vec2S constructs a 2x32-bit vector SIMD register, with type V2S and family RegVec64.
func Vec2S(id uint8) Reg { return Reg{ID: id, Type: V2S} }

// Vec4S constructs a 4x32-bit vector SIMD register, with type V4S and family RegVec128.
func Vec4S(id uint8) Reg { return Reg{ID: id, Type: V4S} }

// Vec1D constructs a 1x64-bit vector SIMD register, with type V1D and family RegVec64.
func Vec1D(id uint8) Reg { return Reg{ID: id, Type: V1D} }

// Vec2D constructs a 2x64-bit vector SIMD register, with type V2D and family RegVec128.
func Vec2D(id uint8) Reg { return Reg{ID: id, Type: V2D} }

// Vec1Q constructs a 1x128-bit vector SIMD register, with type V1Q and family RegVec128.
func Vec1Q(id uint8) Reg { return Reg{ID: id, Type: V1Q} }

// I selects a vector element from r.
func (r Reg) I(idx uint8) Reg { return Reg{ID: r.ID, Type: r.Type, ElemInv: ^idx} }

// List constructs a register list with sequential registers starting from r.
func (r Reg) List(length uint8) RegList { return RegList{First: r, Len: length} }

// I selects a vector element from all registers in r.
func (r RegList) I(idx uint8) RegList { return RegList{First: r.First.I(idx), Len: r.Len} }

// Family returns the register family for r (integer, SP, scalar, vector).
func (r Reg) Family() RegFamily { return r.Type.Family() }

// ElemSize returns the element size for r. For vector registers, the size represents a single lane.
func (r Reg) ElemSize() Size { return r.Type.Elem() }

// Lanes returns the number of lanes for the type of r.
func (r Reg) Lanes() uint8 { return r.Type.Lanes() }

// HasElem returns true if a vector element is selected for r.
func (r Reg) HasElem() bool { return r.ElemInv != 0 }

// GetElem returns the vector element selected for r. The element is only valid if HasElem returns true.
func (r Reg) GetElem() uint8 { return ^r.ElemInv }

// IsVec returns true is the type of r is a vector type.
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

const (
	// Register element size, family, and lane count indicators

	RW RegType = RegType(DWORD) | RegType(RegInt<<4) // RW represents 32-bit integer registers
	RX RegType = RegType(QWORD) | RegType(RegInt<<4) // RX represents 64-bit integer registers

	RWSP RegType = RegType(DWORD) | RegType(RegSP<<4) // RWSP represents 32-bit stack pointer registers
	RXSP RegType = RegType(QWORD) | RegType(RegSP<<4) // RXSP represents 64-bit stack pointer registers

	RB RegType = RegType(BYTE) | RegType(RegFloat<<4)  // RB represents 8-bit scalar SIMD registers
	RH RegType = RegType(WORD) | RegType(RegFloat<<4)  // RH represents 16-bit scalar SIMD registers
	RS RegType = RegType(DWORD) | RegType(RegFloat<<4) // RS represents 32-bit scalar SIMD registers
	RD RegType = RegType(QWORD) | RegType(RegFloat<<4) // RD represents 64-bit scalar SIMD registers
	RQ RegType = RegType(OWORD) | RegType(RegFloat<<4) // RQ represents 128-bit scalar SIMD registers

	V4B  RegType = RegType(BYTE) | RegType(RegVec32<<4)   // V4B represents 4x8-bit vector SIMD registers
	V8B  RegType = RegType(BYTE) | RegType(RegVec64<<4)   // V8B represents 8x8-bit vector SIMD registers
	V16B RegType = RegType(BYTE) | RegType(RegVec128<<4)  // V16B represents 16x8-bit vector SIMD registers
	V2H  RegType = RegType(WORD) | RegType(RegVec32<<4)   // V2H represents 2x16-bit vector SIMD registers
	V4H  RegType = RegType(WORD) | RegType(RegVec64<<4)   // V4H represents 4x16-bit vector SIMD registers
	V8H  RegType = RegType(WORD) | RegType(RegVec128<<4)  // V8H represents 8x16-bit vector SIMD registers
	V2S  RegType = RegType(DWORD) | RegType(RegVec64<<4)  // V2S represents 2x32-bit vector SIMD registers
	V4S  RegType = RegType(DWORD) | RegType(RegVec128<<4) // V4S represents 4x32-bit vector SIMD registers
	V1D  RegType = RegType(QWORD) | RegType(RegVec64<<4)  // V1D represents 1x64-bit vector SIMD registers
	V2D  RegType = RegType(QWORD) | RegType(RegVec128<<4) // V2D represents 2x64-bit vector SIMD registers
	V1Q  RegType = RegType(OWORD) | RegType(RegVec128<<4) // V1Q represents 1x128-bit vector SIMD registers
)

// Family returns the register family for sz (integer, SP, scalar, vector).
func (sz RegType) Family() RegFamily { return RegFamily(sz >> 4) }

// Elem returns the element size for sz. For vector registers, the size represents a single lane.
func (sz RegType) Elem() Size { return Size(sz) & 0xF }

// ElemBytes returns the element size in bytes for sz. For vector registers, the size represents a single lane.
func (sz RegType) ElemBytes() uint8 { return 1 << (sz.Elem() - BYTE) }

// ElemBytes returns the full size in bytes for sz. For vector registers, the size represents all lanes combined.
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

// Lanes returns the lane count for sz. If sz is not from a vector register family, the lane count is 0.
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
