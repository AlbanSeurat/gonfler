package lzma

type distance struct {
	posSlotDecoder [kNumLenToPosStates]*treedecoder
	posDecoders    [1 + kNumFullDistances - kEndPosModelIndex]prob
	alignDecoder   *treedecoder
}


func newDistance() *distance {
	dist := new(distance)
	for i := 0 ; i < kNumLenToPosStates ; i++ {
		dist.posSlotDecoder[i] = newTree(6)
	}
	dist.alignDecoder = newTree(kNumAlignBits)
	initProbs(dist.posDecoders[:])
	return dist
}

func (ld *LzmaDecoder) decodeDistance(len uint) (uint, error) {
	lenState := len
	if lenState >= kNumLenToPosStates - 1 {
		lenState = kNumLenToPosStates - 1
	}
	if posSlot, err := ld.distance.posSlotDecoder[lenState].decode(&ld.rangeDec) ; err == nil {
		if posSlot < 4 {
			return posSlot, nil
		}
		numDirectBytes := (posSlot >> 1 ) - 1
		dist := (2 | (posSlot & 1)) << numDirectBytes
		if posSlot < kEndPosModelIndex {
			if tmp, err := reverseDecode(ld.distance.posDecoders[dist - posSlot:], &ld.rangeDec) ; err == nil {
				dist += tmp
			} else {
				return 0, err
			}
		} else {
			if tmp, err := ld.rangeDec.decodeDirectBits(uint32(numDirectBytes - kNumAlignBits)) ; err == nil {
				dist += uint(tmp) << kNumAlignBits
			} else {
				return 0, err
			}
			if tmp, err := ld.alignDecoder.reverseDecode(&ld.rangeDec) ; err == nil {
				dist += tmp
			} else {
				return 0, err
			}
		}
		return dist, nil
	} else {
		return 0, err
	}
}

