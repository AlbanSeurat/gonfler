package _zdecode

import (
	"io"
	"bufio"
)

type streamInfo struct {
	dataOffset uint64
	folders []folder
}

func readStreamInfo(reader io.Reader) (*streamInfo, error) {
	bufReader := bufio.NewReader(reader)
	streamType, err := bufReader.ReadByte()
	if err != nil {
		return nil, err
	}
	ssInfo := new(streamInfo)
	if streamType == kPackInfo {
		ssInfo.dataOffset, err = readEncodedUInt64(bufReader)
		if err != nil {
			return nil, err
		}
		err = readPackInfo(bufReader)
		if err != nil {
			return nil, err
		}
		if streamType, err = bufReader.ReadByte() ; err != nil {
			return nil, err
		}
	}

	if streamType == kUnpackInfo {
		ssInfo.folders, err = readUnpackInfo(bufReader)
		if err != nil {
			return nil, err
		}
		if streamType, err = bufReader.ReadByte() ; err != nil {
			return nil, err
		}
	}

	if streamType != kEnd {
		return nil, errInvalidFile
	}

	return ssInfo, nil
}
