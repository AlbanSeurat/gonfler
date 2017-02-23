package _zdecode

import (
	"bufio"
	"encoding/binary"
)

func readPackInfo(reader *bufio.Reader) ([]uint64, error) {
	numPackStreams, err := readEncodedUInt64(reader)
	if err != nil {
		return nil, err
	}
	if err = readUntil(reader, kSize); err != nil {
		return nil, err
	}
	packSizes := make([]uint64, numPackStreams)
	for idx := range packSizes {
		if packSizes[idx], err = readEncodedUInt64(reader) ; err != nil {
			return nil, err
		}
	}
	//TODO - manage CRC block when present
	err = readEndOrCrc(reader, func() error {
		return nil
	})
	if err != nil {
		return nil, err
	}
	return packSizes, nil
}


func readUnpackInfo(reader *bufio.Reader) ([]folder, error) {
	if err := readUntil(reader, kFolder) ; err != nil {
		return nil,err
	}
	numFolders, err := readEncodedUInt64(reader)
	if err !=  nil {
		return nil,err
	}
	folders := make([]folder, numFolders)

	external, err := readEncodedUInt64(reader)
	if err != nil {
		return nil,err
	}
	switch external {
	case 0:
		for i := uint64(0) ; i < numFolders ; i++ {
			folder, err := readFolder(reader)
			if err != nil {
				return nil,err
			} else {
				folders[i] = *folder
			}
		}
	case 1:
		return nil, errUnsupported

	}
	err = readUntil(reader, kCodersUnpackSize)
	if err != nil {
		return nil,err
	}

	//TODO - look at what to do with unpackSize
	for _, fold := range folders {
		for i := 0 ; i < len(fold.codecs) ; i++ {
			_, err = readEncodedUInt64(reader)
			if err != nil {
				return nil,err
			}
		}
	}
	err = readEndOrCrc(reader, func() error {
		return readHashDigests(reader, int(numFolders))
	})
	if err != nil {
		return nil,err
	}
	return folders, nil
}

func readHashDigests(reader *bufio.Reader, nbElem int) error {
	bools, count, err := readDefinedBooleans(reader, nbElem)
	if err != nil {
		return err
	}
	crcs := make([]uint32, count)
	for pos := 0 ; pos < nbElem ; pos++ {
		if bools[pos] {
			binary.Read(reader, binary.LittleEndian, &crcs[pos])
		}
	}
	return nil
}