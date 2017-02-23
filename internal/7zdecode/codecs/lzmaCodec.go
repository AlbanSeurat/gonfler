package codecs

import (
	"github.com/alkpone/gonfler/internal/lzma"
)

type LzmaCodec struct {
	props *lzma.Props
}

func (LzmaCodec) Decode(stream []byte) ([]byte, error) {

	return nil, nil
}

func (codec LzmaCodec) Props(codecProps []byte) error {
	var err error
	if codec.props, err = lzma.NewProps(codecProps) ; err != nil {
		return err
	} else {
		return nil
	}
}


