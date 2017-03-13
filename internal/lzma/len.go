package lzma

type lendecoder struct {
	choice    prob
	choice2   prob
	lowCoder  [1 << kNumPosBitsMax]*treedecoder
	midCoder  [1 << kNumPosBitsMax]*treedecoder
	highCoder *treedecoder
}

func newLenDecoder() *lendecoder {
	ret := new(lendecoder)
	ret.choice = probInit
	ret.choice2 = probInit
	ret.highCoder = newTree(8)
	for i := 0 ; i < (1 << kNumPosBitsMax) ; i++ {
		ret.lowCoder[i] = newTree(3)
		ret.midCoder[i] = newTree(3)
	}
	return ret
}

func (ld *lendecoder) decode(rd *rangedecoder, posState uint) (uint, error) {
	if bit, err := rd.decodeBit(&ld.choice) ; err == nil {
		if bit == 0 {
			return ld.lowCoder[posState].decode(rd)
		}
	} else {
		return 0, err
	}
	if bit, err := rd.decodeBit(&ld.choice2) ; err == nil {
		if bit == 0 {
			if ret, err := ld.midCoder[posState].decode(rd); err == nil {
				return 8 + ret, nil
			} else {
				return 0, err
			}
		}
	} else {
		return 0, err
	}
	if ret, err := ld.highCoder.decode(rd) ; err == nil {
		return 16 + ret, nil
	} else {
		return 0, nil
	}
}