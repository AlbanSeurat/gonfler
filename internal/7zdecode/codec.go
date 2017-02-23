package _zdecode

import "github.com/alkpone/gonfler/internal/7zdecode/codecs"

const (
	copy = 0x0
	lzma = 0x030101
	lzma2 = 0x21
	_7zaes = 0x06F10701
)

var (
	codecMap = map[int]codec {
		lzma: codecs.LzmaCodec{},
		lzma2: lzma2Codec{},
	}
)

type codec interface {
	Props(codedProps []byte) error
	Decode(stream []byte) ([]byte, error)
}

type lzma2Codec struct {

}

func (lzma2Codec) Props(codedProps []byte) error {
	return nil
}

func (lzma2Codec) Decode(stream []byte)  ([]byte, error) {
	return nil, nil
}