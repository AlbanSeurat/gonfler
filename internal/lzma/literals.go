package lzma

type literals struct {
	litProbs []prob
}

func (lit *literals) initLiterals(lc, lp int) {
	lit.litProbs = make([]prob, 0x300<<uint(lc+lp))
	for i := range lit.litProbs {
		lit.litProbs[i] = probInit
	}
}

func (ld *LzmaDecoder) decodeLiteral(state uint, rep0 int) error {
	prevByte := byte(0)
	if !ld.out.isEmpty() {
		prevByte = ld.out.getByte(1)
	}
	symbol := uint(1)
	litState := ((ld.out.totalPos & ((1 << ld.props.lp) - 1)) << ld.props.lc) + (uint32(prevByte) >> (8 - ld.props.lc))
	probs := ld.literals.litProbs[0x300 + litState:]
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

