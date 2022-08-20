package instref

import (
	"fmt"
	"math"
	"math/bits"
	"sort"
	"strings"

	"github.com/wdamron/arm"
	"github.com/wdamron/arm/gen/inst/opmap"
)

type EncodingWidth struct {
	Size int
	Doc  string
}

// Get returns a map from instruction name to list of encoding docs, where each encoding doc contains
// a single width (0) or pair of SIMD widths (16/8). Encoding constraints will be padded to lineChars.
func Get(lineChars int) map[string][][2]EncodingWidth {
	names := make([]string, 0, len(opmap.EncMap))
	for name := range opmap.EncMap {
		names = append(names, name)
	}
	sort.Strings(names)

	mapping := make(map[string][][2]EncodingWidth, len(opmap.EncMap))

	for _, name := range names {
		encs := opmap.EncMap[name]
		mapping[name] = make([][2]EncodingWidth, len(encs))
		for encIdx, enc := range encs {
			encInfo := getEncodingInfo(enc.Match, enc.Cmds)
			for widthIdx, width := range encInfo.widths() { // 2 encodings for dual-width SIMD instructions
				// Name and arguments:

				line := new(strings.Builder)
				line.Grow(64)
				line.WriteString(name)
				line.WriteByte(' ')
				for i, matcher := range encInfo.matchers {
					if matcher.m.Op == arm.MatEnd {
						line.WriteString(" {")
						continue
					}
					if i > 0 {
						line.WriteString(", ")
					}
					switch matcher.m.Op {
					case arm.MatLitSymbol:
						line.WriteString(arm.SymbolName[matcher.m.X[0]])
					case arm.MatLitInt:
						fmt.Fprintf(line, "#%d", matcher.m.X[0])
					case arm.MatLitFloat:
						fmt.Fprintf(line, "#%.01f", float32(matcher.m.X[0]))
					case arm.MatSymbol:
						line.WriteString("<symbol>")
					case arm.MatCond:
						line.WriteString("<cond>")
					case arm.MatImm, arm.MatFloat:
						line.WriteByte('#')
						line.WriteString(immName(matcher.flat[0].suffix, encInfo.numImms, false))
					case arm.MatW:
						line.WriteString("W" + regSuffixes[matcher.flat[0].suffix])
					case arm.MatX:
						line.WriteString("X" + regSuffixes[matcher.flat[0].suffix])
					case arm.MatWSP:
						line.WriteString(fmt.Sprintf("W%s|WSP", regSuffixes[matcher.flat[0].suffix]))
					case arm.MatXSP:
						line.WriteString(fmt.Sprintf("X%s|SP", regSuffixes[matcher.flat[0].suffix]))
					case arm.MatB:
						line.WriteString("B" + regSuffixes[matcher.flat[0].suffix])
					case arm.MatH:
						line.WriteString("H" + regSuffixes[matcher.flat[0].suffix])
					case arm.MatS:
						line.WriteString("S" + regSuffixes[matcher.flat[0].suffix])
					case arm.MatD:
						line.WriteString("D" + regSuffixes[matcher.flat[0].suffix])
					case arm.MatQ:
						line.WriteString("Q" + regSuffixes[matcher.flat[0].suffix])
					case arm.MatV:
						size := arm.Size(matcher.m.X[0])
						fmt.Fprintf(line, "V%s.%d%s", regSuffixes[matcher.flat[0].suffix], width/sizeBytes(size), sizeName(size))
					case arm.MatVStatic:
						size, lanes := arm.Size(matcher.m.X[0]), int(matcher.m.X[1])
						fmt.Fprintf(line, "V%s.%d%s", regSuffixes[matcher.flat[0].suffix], lanes, sizeName(size))
					case arm.MatVElement:
						size := arm.Size(matcher.m.X[0])
						fmt.Fprintf(line, "V%s.%s[i]", regSuffixes[matcher.flat[0].suffix], sizeName(size))
					case arm.MatVElementStatic:
						size, idx := arm.Size(matcher.m.X[0]), int(matcher.m.X[1])
						fmt.Fprintf(line, "V%s.%s[%d]", regSuffixes[matcher.flat[0].suffix], sizeName(size), idx)
					case arm.MatVStaticElement:
						size, lanes := arm.Size(matcher.m.X[0]), int(matcher.m.X[1])
						fmt.Fprintf(line, "V%s.%d%s[i]", regSuffixes[matcher.flat[0].suffix], lanes, sizeName(size))
					case arm.MatRegList:
						count, size := int(matcher.m.X[0]), arm.Size(matcher.m.X[1])
						fmt.Fprintf(line, "{V%s.%d%s * %d}", regSuffixes[matcher.flat[0].suffix], width/sizeBytes(size), sizeName(size), count)
					case arm.MatRegListStatic:
						count, size, lanes := int(matcher.m.X[0]), arm.Size(matcher.m.X[1]), int(matcher.m.X[2])
						fmt.Fprintf(line, "{V%s.%d%s * %d}", regSuffixes[matcher.flat[0].suffix], lanes, sizeName(size), count)
					case arm.MatRegListElement:
						count, size := int(matcher.m.X[0]), arm.Size(matcher.m.X[1])
						fmt.Fprintf(line, "{V%s.%s * %d}[i]", regSuffixes[matcher.flat[0].suffix], sizeName(size), count)
					case arm.MatOffset:
						line.WriteString("<offset>")
					case arm.MatRefBase:
						line.WriteString(fmt.Sprintf("[X%s|SP]", regSuffixes[matcher.flat[0].suffix]))
					case arm.MatRefOffset:
						line.WriteString(fmt.Sprintf("[X%s|SP {, #%s }]", regSuffixes[matcher.flat[0].suffix], immName(matcher.flat[1].suffix, encInfo.numImms, false)))
					case arm.MatRefPre:
						line.WriteString(fmt.Sprintf("[X%s|SP, #%s]!", regSuffixes[matcher.flat[0].suffix], immName(matcher.flat[1].suffix, encInfo.numImms, false)))
					case arm.MatRefIndex:
						line.WriteString(fmt.Sprintf("[X%s|SP, W%s|X%s {, LSL|UXTW|SXTW|SXTX #%s }]",
							regSuffixes[matcher.flat[0].suffix], regSuffixes[matcher.flat[1].suffix], regSuffixes[matcher.flat[1].suffix], immName(matcher.flat[3].suffix, encInfo.numImms, false)))
					case arm.MatLitMod:
						line.WriteString(strings.TrimPrefix(arm.ModName[matcher.m.X[0]], "Sym"))
						line.WriteString(" #")
						line.WriteString(immName(matcher.flat[0].suffix, encInfo.numImms, false))
					case arm.MatMod:
						imm := immName(matcher.flat[1].suffix, encInfo.numImms, false)
						switch list := matcher.m.X[0]; list {
						case arm.SymExtends:
							line.WriteString("LSL|UXT[BHWX]|SXT[BHWX] #" + imm)
						case arm.SymExtendsW:
							line.WriteString("UXT[BHW]|SXT[BHW] #" + imm)
						case arm.SymExtendsX:
							line.WriteString("LSL|UXTX|SXTX #" + imm)
						case arm.SymShifts:
							line.WriteString("LSL|LSR|ASR #" + imm)
						case arm.SymRotates:
							line.WriteString("LSL|LSR|ASR|ROR #" + imm)
						}

					}
				}

				if encInfo.hasOptional {
					line.WriteString(" }")
				}

				instFmt := line.String()

				// Constraints:

				if encInfo.hasConstraints {
					line = new(strings.Builder)
					line.Grow(64)
					line.WriteByte('(')
					sep := ""
					for _, m := range encInfo.matchers {
						for _, flat := range m.flat {
							for _, c := range flat.constraints {
								line.WriteString(sep)
								line.WriteString(c)
								sep = ", "
							}
						}
					}
					line.WriteByte(')')
					constraintsFmt := line.String()

					// pad right to line width:
					line = new(strings.Builder)
					line.Grow(lineChars)
					line.Write([]byte(instFmt))
					line.WriteByte(' ')
					line.WriteByte(' ')
					for i := len(instFmt) + 2; i < lineChars-(2+len(constraintsFmt)); i++ {
						line.WriteString("·")
					}
					line.WriteByte(' ')
					line.WriteByte(' ')
					line.WriteString(constraintsFmt)

					instFmt = line.String()
				}

				mapping[name][encIdx][widthIdx] = EncodingWidth{Size: width, Doc: instFmt}
			}
		}
	}

	return mapping
}

// casp* instructions have destination as the last register
var regSuffixes = [...]string{"d", "n", "m", "a", "b", "d"} // Xd, Xn, Xm, Xa, Xb, Xd

func immName(immNum, immCount int, isOffset bool) string { // imm, imm1-imm4, offset
	if isOffset {
		return "offset"
	}
	if immCount > 1 {
		return fmt.Sprintf("imm%d", immNum+1)
	}
	return "imm"
}

func sizeName(size arm.Size) string {
	switch size {
	case arm.BYTE:
		return "B"
	case arm.WORD:
		return "H"
	case arm.DWORD:
		return "S"
	case arm.QWORD:
		return "D"
	case arm.OWORD:
		return "Q"
	}
	return ""
}
func sizeBytes(size arm.Size) int {
	switch size {
	case arm.BYTE:
		return 1
	case arm.WORD:
		return 2
	case arm.DWORD:
		return 4
	case arm.QWORD:
		return 8
	case arm.OWORD:
		return 16
	}
	return 1
}

// Constraints
// ------------------------------------------------------------------------------

type encodingInfo struct {
	numImms        int
	hasOptional    bool
	hasConstraints bool
	dualWidth      bool
	matchers       []matcherInfo
}

var (
	dualWidth   = [2]int{16, 8}
	singleWidth = [1]int{0}
)

func (encInfo encodingInfo) widths() []int {
	if encInfo.dualWidth {
		return dualWidth[:]
	}
	return singleWidth[:]
}

type matcherInfo struct {
	m    arm.EncOp
	flat []flatArgInfo
}

func (matcher *matcherInfo) hasConstraints() bool {
	for _, f := range matcher.flat {
		if len(f.constraints) != 0 {
			return true
		}
	}
	return false
}

type flatArgInfo struct {
	suffix      int
	immMin      int64
	immMax      int64
	immSumDec   int64 // imm == immSumDec - prevImm
	immAlt      bool  // ignore min/max
	cmds        []arm.EncOp
	constraints []string
}

func getEncodingInfo(pattern, commands []arm.EncOp) encodingInfo {
	var encInfo encodingInfo
	// Index flat args to matchers and arg index within matcher:
	var flat2matcher []int
	var flat2matcherFlat []int
IndexCommands:
	for mi, m := range pattern {
		minfo := matcherInfo{m: m}
		if m.Op == arm.MatEnd {
			encInfo.hasOptional = true
			encInfo.matchers = append(encInfo.matchers, minfo)
			continue IndexCommands
		}
		switch m.Op {
		case arm.MatV, arm.MatRegList:
			encInfo.dualWidth = true
		case arm.MatRefOffset, arm.MatRefPre, arm.MatRefIndex, arm.MatImm, arm.MatFloat, arm.MatMod, arm.MatLitMod:
			encInfo.numImms++
		}
		fc := int(arm.MatcherFlatArgCounts[m.Op])
		minfo.flat = make([]flatArgInfo, fc)
		for fi := 0; fi < fc; fi++ {
			flat2matcherFlat = append(flat2matcherFlat, fi)
			flat2matcher = append(flat2matcher, mi)
		}
		encInfo.matchers = append(encInfo.matchers, minfo)
	}
	// Group commands by matcher->flat:
	var flatArgIdx int
GroupCommands:
	for _, c := range commands {
		switch c.Op {
		case arm.CmdAdv:
			flatArgIdx++
			continue GroupCommands
		case arm.CmdBack:
			flatArgIdx--
			continue GroupCommands
		case arm.CmdRwidth30:
			continue GroupCommands
		}
		matcherIdx := flat2matcher[flatArgIdx]
		matcherFlatIdx := flat2matcherFlat[flatArgIdx]
		encInfo.matchers[matcherIdx].flat[matcherFlatIdx].cmds = append(encInfo.matchers[matcherIdx].flat[matcherFlatIdx].cmds, c)
		switch c.Op {
		default:
			flatArgIdx++
		case arm.CmdUslice, arm.CmdSslice, arm.CmdChkUbits, arm.CmdChkUsum, arm.CmdChkSscaled, arm.CmdChkUrange1:
			// non-consuming
		}
	}
	// Collect constraint string lists for each flat arg:
	nextRegSuffix, nextImmSuffix := 1, 0
	updateReg := func(matcherIdx, flatArgIdx int) {
		matcher := &encInfo.matchers[matcherIdx]
		if matcherIdx > 0 || (len(matcher.flat) == 1 && len(matcher.flat[0].cmds) == 1 && matcher.flat[0].cmds[0].Op == arm.CmdREven) {
			matcher.flat[flatArgIdx].suffix = nextRegSuffix // n, m, a, b, d (casp* instructions have destination as the last register)
			nextRegSuffix++
		} else {
			matcher.flat[flatArgIdx].suffix = 0 // d (if first argument is a register type)
		}
		matcher.flat[flatArgIdx].fmtRegArgConstraints()
	}
	updateImm := func(matcherIdx, flatArgIdx int) {
		matcher := &encInfo.matchers[matcherIdx]
		matcher.flat[flatArgIdx].suffix = nextImmSuffix
		matcher.flat[flatArgIdx].fmtImmArgConstraints(encInfo.numImms, false)
		nextImmSuffix++
	}
	for mi := range encInfo.matchers {
		matcher := &encInfo.matchers[mi]
		switch op := matcher.m.Op; op {
		case arm.MatW, arm.MatX, arm.MatWSP, arm.MatXSP, arm.MatB, arm.MatH, arm.MatS, arm.MatD, arm.MatQ,
			arm.MatV, arm.MatVStatic, arm.MatVElement, arm.MatVStaticElement, arm.MatVElementStatic,
			arm.MatRegList, arm.MatRegListStatic, arm.MatRegListElement, arm.MatRefBase:
			updateReg(mi, 0)
		case arm.MatRefOffset, arm.MatRefPre:
			updateReg(mi, 0)
			updateImm(mi, 1)
		case arm.MatRefIndex:
			updateReg(mi, 0)
			updateReg(mi, 1)
			updateImm(mi, 3)
		case arm.MatImm, arm.MatFloat, arm.MatLitMod:
			updateImm(mi, 0)
		case arm.MatMod:
			updateImm(mi, 1)
		case arm.MatOffset:
			matcher.flat[0].fmtImmArgConstraints(encInfo.numImms, true)
		case arm.MatSymbol, arm.MatLitSymbol:
			// no constraints
		}
		if matcher.hasConstraints() {
			encInfo.hasConstraints = true
		}
	}

	return encInfo
}

func (flat *flatArgInfo) fmtRegArgConstraints() {
	name := regSuffixes[flat.suffix]
	var list []string
	for _, c := range flat.cmds {
		switch c.Op {
		case arm.CmdRLo16:
			list = append(list, name+" < 16")
		case arm.CmdRNz16:
			list = append(list, name+" != 31")
		case arm.CmdREven:
			list = append(list, name+" is even")
		case arm.CmdRNext:
			prev := regSuffixes[flat.suffix-1]
			list = append(list, fmt.Sprintf("%s == %s + 1", name, prev))
		}
	}
	flat.constraints = list
}

func (flat *flatArgInfo) setMin(i int64) {
	if i > flat.immMin {
		flat.immMin = i
	}
}
func (flat *flatArgInfo) setMax(i int64) {
	if i < flat.immMax {
		flat.immMax = i
	}
}

func (flat *flatArgInfo) fmtImmArgConstraints(immCount int, isOffset bool) {
	name := immName(flat.suffix, immCount, isOffset)

	flat.immMin, flat.immMax = math.MinInt64, math.MaxInt64
	var misc []string
	for _, c := range flat.cmds {
		switch c.Op {
		case arm.CmdUbits:
			bitlen := c.X[1]
			flat.setMin(0)
			flat.setMax((1 << bitlen) - 1)
		case arm.CmdSbits:
			half := int64(1) << (9 - 1)
			flat.setMin(-1 * half)
			flat.setMax(half - 1)
		case arm.CmdChkUbits:
			bitlen := c.X[0]
			flat.setMin(0)
			flat.setMax((1 << bitlen) - 1)
		case arm.CmdUscaled:
			bitlen, scale := c.X[1], c.X[2]
			flat.setMin(0)
			flat.setMax((1 << (bitlen + scale)) - 1)
			misc = append(misc, fmt.Sprintf("%s >> %d", name, scale))
		case arm.CmdSscaled:
			const bits = uint8(7)
			scale := c.X[0]
			half := int64(1) << (bits + scale - 1)
			flat.setMin(-1 * half)
			flat.setMax(half - 1)
			misc = append(misc, fmt.Sprintf("%s >> %d", name, scale))
		case arm.CmdChkSscaled:
			const bits, scale = uint8(10), uint8(3)
			half := int64(1) << (bits + scale - 1)
			flat.setMin(-1 * half)
			flat.setMax(half - 1)
			misc = append(misc, fmt.Sprintf("%s >> %d", name, scale))
		case arm.CmdUrange:
			flat.setMin(int64(c.X[1]))
			flat.setMax(int64(c.X[2]))
		case arm.CmdChkUrange1:
			flat.setMin(1)
			flat.setMax(int64(c.X[0]))
		case arm.CmdUsub:
			bitlen, add := c.X[1], int64(c.X[2])
			flat.setMin(add + 1 - (1 << bitlen))
			flat.setMax(add)
		case arm.CmdUnegmod:
			bitlen := c.X[1]
			flat.setMin(0)
			flat.setMax((1 << bitlen) - 1)
		case arm.CmdUsumdec:
			bitlen := c.X[1]
			flat.setMin(1)
			flat.setMax(1 << bitlen)
			if sumDec := int64(1) << bitlen; flat.immSumDec == 0 || flat.immSumDec > sumDec {
				flat.immSumDec = sumDec
			}
		case arm.CmdChkUsum:
			shift := c.X[0]
			flat.setMin(1)
			flat.setMax(1 << shift)
			if sumDec := int64(1) << shift; flat.immSumDec == 0 || flat.immSumDec > sumDec {
				flat.immSumDec = sumDec
			}
		case arm.CmdUfields30, arm.CmdUfields11:
			count := c.X[0]
			flat.setMin(0)
			flat.setMax((1 << count) - 1)
		case arm.CmdUfields21:
			flat.setMin(0)
			flat.setMax(1)
		case arm.CmdUAlt2, arm.CmdUAlt4:
			listIdx := c.X[1]
			if c.Op == arm.CmdUAlt2 && arm.Alts2[listIdx] == ([2]uint16{0, 0}) {
				flat.setMin(0)
				flat.setMax(0)
				continue
			}
			var alts []uint16
			if c.Op == arm.CmdUAlt2 {
				alts = arm.Alts2[listIdx][:]
			} else {
				alts = arm.Alts4[listIdx][:]
			}
			flat.immAlt = true // ignore min/max
			seen := make(map[uint16]struct{})
			var sb strings.Builder
			sb.Grow(32)
			sb.WriteString(name)
			sb.WriteString(" in [")
			for _, v := range alts {
				if _, ok := seen[v]; ok {
					continue
				}
				if len(seen) != 0 {
					sb.WriteString(", ")
				}
				seen[v] = struct{}{}
				sb.WriteString(fmt.Sprintf("%d", v))
			}
			sb.WriteByte(']')
			misc = append(misc, sb.String())
		case arm.CmdSpecial:
			switch specialType := c.X[1]; specialType {
			case arm.SpecialImmWide32:
				misc = append(misc, name+" is 32-bit wide")
			case arm.SpecialImmWide64:
				misc = append(misc, name+" is 64-bit wide")
			case arm.SpecialImmWideInv32:
				misc = append(misc, name+" is 32-bit inverted wide")
			case arm.SpecialImmWideInv64:
				misc = append(misc, name+" is 64-bit inverted wide")
			case arm.SpecialImmLogical32:
				misc = append(misc, name+" is 32-bit logical")
			case arm.SpecialImmLogical64:
				misc = append(misc, name+" is 64-bit logical")
			case arm.SpecialImmFloat:
				misc = append(misc, name+" is float")
			case arm.SpecialImmFloatSplit:
				misc = append(misc, name+" is split float")
			case arm.SpecialImmStretched:
				misc = append(misc, name+" is stretched")
			}
		case arm.CmdOffset:
			switch relType := c.X[0]; relType {
			case arm.RelB:
				misc = append(misc, fmt.Sprintf("%s >> 2 is 26-bit (+/- 128 MB)", name))
			case arm.RelBCond:
				misc = append(misc, fmt.Sprintf("%s >> 2 is 19-bit (+/- 1 MB)", name))
			case arm.RelAdr:
				misc = append(misc, name+" is 21-bit (+/- 1 MB)")
			case arm.RelAdrp:
				misc = append(misc, fmt.Sprintf("%s >> 12 is 21-bit (+/- 4 GB)", name))
			case arm.RelTbz:
				misc = append(misc, fmt.Sprintf("%s >> 2 is 14-bit (+/- 32 KB)", name))
			}
		}
	}

	// Format min/max range constraint(s) then remaining constraints:

	if flat.immAlt || (flat.immMin == math.MinInt64 && flat.immMax == math.MaxInt64) { // ignore min/max
		flat.constraints = misc
		return
	}

	if flat.immSumDec != 0 {
		misc = append([]string{
			fmt.Sprintf("%s + %s <= %d", immName(flat.suffix-1, immCount, false), name, flat.immSumDec),
		}, misc...)
	}

	if flat.immMin == flat.immMax {
		flat.constraints = append([]string{fmt.Sprintf("%s == %d", name, flat.immMin)}, misc...)
		return
	}

	var rangeExpr strings.Builder
	rangeExpr.Grow(32)
	if flat.immMin != math.MinInt64 {
		if flat.immMin == 1 {
			rangeExpr.WriteString("0 < ")
		} else {
			fmt.Fprintf(&rangeExpr, "%d <= ", flat.immMin)
		}
	}
	rangeExpr.WriteString(name)
	if flat.immMax != math.MaxInt64 {
		powerOf2 := func(v int64) bool { return bits.OnesCount64(uint64(v)) == 1 }
		if powerOf2(flat.immMax + 1) {
			fmt.Fprintf(&rangeExpr, " < %d", flat.immMax+1)
		} else {
			fmt.Fprintf(&rangeExpr, " <= %d", flat.immMax)
		}
	}

	flat.constraints = append([]string{rangeExpr.String()}, misc...)
}
