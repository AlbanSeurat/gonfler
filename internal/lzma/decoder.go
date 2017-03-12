package lzma

const (
	kNumPosBitsMax     = 4
	kNumLenToPosStates = 4
	kEndPosModelIndex  = 14
	kNumFullDistances  = 1 << (kEndPosModelIndex >> 1)
	kNumAlignBits      = 4
)

type LzmaDecoder struct {
	literals
	distance
	props         Props
	out           outWindow
	rangeDec      rangedecoder
	lenDecoder    lendecoder
	repLenDecoder lendecoder
}

func decode(stream []byte) {

	for {

		break
	}
}
