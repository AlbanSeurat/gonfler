package lzma

type treedecoder struct {
	probs []prob
}

func newTree(numBits int) *treedecoder {
	tree := new(treedecoder)
	tree.probs = make([]prob, numBits)
	return tree
}

func (t *treedecoder) decode(rd *rangedecoder) (uint, error) {
	m := uint(1)
	for _, prob := range t.probs {
		bit, err := rd.decodeBit(&prob)
		if err != nil {
			return 0, err
		}
		m = (m << 1) + bit
	}
	return m - uint(1 << uint(len(t.probs))), nil
}

func (t *treedecoder) reverseDecode(rd *rangedecoder) (uint, error) {
	m := uint(1)
	symbol := uint(0)
	for pos, prob := range t.probs {
		bit, err := rd.decodeBit(&prob)
		if err != nil {
			return 0, err
		}
		m <<= 1
		m += bit
		symbol |= bit << uint(pos)
	}
	return symbol, nil
}
