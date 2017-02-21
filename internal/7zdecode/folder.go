package _zdecode

import "bufio"

type codecSpec struct {
	id         uint64
	numStreams uint64
	props      []byte
}

type folder struct {
	codecs        []codecSpec
	packStream    []uint64
}

func readFolder(reader *bufio.Reader) (*folder, error) {
	folder := new(folder)
	numCoders, err := readEncodedUInt64(reader)
	if err != nil {
		return nil, err
	}
	folder.codecs = make([]codecSpec, numCoders)
	numInStreams := uint64(0)
	for i := uint64(0); i < numCoders; i++ {
		mainByte, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if mainByte&0xC0 != 0 {
			return nil, errInvalidFile
		}
		codecIdSize := mainByte & 0xF
		if codecIdSize > 8 {
			return nil, errInvalidFile
		}
		codecId := uint64(0)
		for pos := 0; pos < int(codecIdSize); pos++ {
			codecIdPart, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}
			codecId = codecId<<8 | uint64(codecIdPart)
		}
		folder.codecs[i].id = codecId
		if mainByte&0x10 != 0 {
			numInStream, err := readEncodedUInt64(reader)
			if err != nil {
				return nil, err
			}
			if value, err := readEncodedUInt64(reader); err != nil || value != uint64(1) {
				return nil, errInvalidFile
			}
			folder.codecs[i].numStreams += numInStream
		} else {
			folder.codecs[i].numStreams = 1
		}
		if mainByte&0x20 != 0 {
			propsSize, err := readEncodedUInt64(reader)
			if err != nil {
				return nil, err
			}
			folder.codecs[i].props = make([]byte, propsSize)
			if value, err := reader.Read(folder.codecs[i].props); err != nil || value != int(propsSize) {
				return nil, errInvalidFile
			}
		}
		numInStreams += folder.codecs[i].numStreams
	}
	numBonds := numCoders - 1

	//TODO: what to do with this (create bonds)
	for i := 0; i < int(numBonds); i++ {
		_, err = readEncodedUInt64(reader)
		if err != nil {
			return nil, err
		}
		_, err = readEncodedUInt64(reader)
		if err != nil {
			return nil, err
		}
	}

	numPackStreams := numInStreams - numBonds
	folder.packStream = make([]uint64, numPackStreams)
	if numPackStreams == 1 {
		//TODO : should be more complex once bonds are managed
		folder.packStream[0] = 0
	} else {
		for i := 0; i < int(numPackStreams); i++ {
			folder.packStream[i], err = readEncodedUInt64(reader)
			if err != nil {
				return nil, err
			}
		}
	}

	return folder, nil
}
