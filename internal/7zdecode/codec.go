package _zdecode

const (
	copy = 0x0
	lzma = 0x030101
	lzma2 = 0x21
	_7zaes = 0x06F10701
)

var (
	codecs = map[int]codec {
		lzma: lzmaCodec{},
		lzma2: lzma2Codec{},
	}
)

type codec interface {
	decode() error
}

func findCodec(id int) (codec, error) {
	codec, ok := codecs[id]
	if ok {
		return codec, nil
	} else {
		return nil, errCodecNotFound
	}
}

type lzmaCodec struct {

}

type lzma2Codec struct {

}

func (lzmaCodec) decode() error {
	return nil
}

func (lzma2Codec) decode() error {
	return nil
}