package _zdecode

import (
	"io"
	"bufio"
)

func readStreamInfo(reader io.Reader) (uint64, error) {
	bufReader := bufio.NewReader(reader)
	streamType, err := bufReader.ReadByte()
	if err != nil {
		return 0, err
	}
	offset := uint64(0)
	if streamType == kPackInfo {
		offset, err = readEncodedUInt64(bufReader)
		if err != nil {
			return 0, err
		}
		readPackInfo(bufReader)
		if streamType, err = bufReader.ReadByte() ; err != nil {
			return 0, err
		}
	}

	if streamType == kUnpackInfo {
		readUnpackInfo(bufReader)
		if streamType, err = bufReader.ReadByte() ; err != nil {
			return 0, err
		}
	}

	if streamType != kEnd {
		return 0, errInvalidFile
	}

	return offset, nil
}
