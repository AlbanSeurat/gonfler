package codecs

type LzmaCodec struct {

}

func (LzmaCodec) Decode(stream []byte) error {
	//_, err := lzma.Uncompress(stream)
	return nil
}


