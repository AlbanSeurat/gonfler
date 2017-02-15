package _zdecode

import (
	"os"
	"fmt"
)

func decodeStream(file *os.File, info *streamInfo, offset uint64) error {

	//TODO: move this outside the function, should only see a stream
	_, err := file.Seek(int64(offset + info.dataOffset), 0)
	if err != nil {
		return err
	}

	//TODO : packStreams init



	for _, folder := range info.folders {
		for _, codec := range folder.codecs {
			fmt.Println(findCodec(int(codec)))
		}
	}


	return nil
}
