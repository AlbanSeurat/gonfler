package _zdecode

import (
	"bufio"
	"encoding/binary"
)

func readPackInfo(reader *bufio.Reader) error {
	numPackStreams, err := readEncodedUInt64(reader)
	if err != nil {
		return nil
	}
	if err = readUntil(reader, kSize); err != nil {
		return err
	}
	packSizes := make([]uint64, numPackStreams)
	for idx, _ := range packSizes {
		if packSizes[idx], err = readEncodedUInt64(reader) ; err != nil {
			return err
		}
	}
	return readEndOrCrc(reader, func() error {
		return nil
	})
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
	for _, fold := range folders {
		for i := 0 ; i < int(fold.numOutStreams) ; i++ {
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