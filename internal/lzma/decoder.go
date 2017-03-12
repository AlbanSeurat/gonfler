package lzma


type LzmaDecoder struct {
	props    Props
	out      outWindow
	rangeDec rangedecoder
	litProbs []prob
}

func (ld *LzmaDecoder) initLiterals(lc, lp int) {
	ld .litProbs = make([]prob, 0x300<<uint(lc+lp))
	for i := range ld.litProbs {
		ld.litProbs[i] = probInit
	}
}

func (ld *LzmaDecoder) decodeLiterals(state uint, rep0 int) error {
	prevByte := byte(0)
	if !ld.out.isEmpty() {
		prevByte = ld.out.getByte(1)
	}
	symbol := uint(1)
	litState := ((ld.out.totalPos & ((1 << ld.props.lp) - 1)) << ld.props.lc) + (uint32(prevByte) >> (8 - ld.props.lc))
	probs := ld.litProbs[0x300 + litState:]
	if state >= 7 {
		matchByte := uint(ld.out.getByte(rep0 + 1))
		for {
			matchBit := (matchByte >> 7) & 1
			matchByte <<= 1
			bit, err := ld.rangeDec.decodeBit(&probs[((1 + matchBit) << 8) + symbol])
			if err != nil {
				return err
			}
			symbol = (symbol << 1) | bit
			if matchBit != bit {
				break
			}
			if symbol < 0x100 {
				break
			}
		}
	}

	for symbol < 0x100 {
		ret, err := ld.rangeDec.decodeBit(&probs[symbol])
		if err != nil {
			return err
		}
		symbol = (symbol << 1 ) | ret
	}
	return ld.out.putByte(byte(symbol - 0x100))
}


func decode(stream []byte) {

	for {


		break
	}
}