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
	//TODO: move this outside the function, should only see a stream
	stream := make([]byte, info.packSize[0])
	if lenght, err := file.Read(stream) ; err != nil || lenght != len(stream) {
		return err
	}

	//TODO : packStreams init


	for _, folder := range info.folders {
		for _, codecSpec := range folder.codecs {
			codec, found := codecMap[int(codecSpec.id)]
			if !found {
				return errCodecNotFound
			}
			if err := codec.Props(codecSpec.props) ; err != nil {
				return err
			}
			fmt.Println(codec.Decode(stream))
		}
	}


	return nil
}
