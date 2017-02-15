package _zdecode

import "bufio"

type folder struct {
	numInStreams  uint64
	numOutStreams uint64
	codecs	[]uint64
}


func readFolder(reader *bufio.Reader) (*folder, error) {
	folder := new(folder)
	numCoders, err := readEncodedUInt64(reader)
	if err != nil {
		return nil, err
	}
	folder.numOutStreams = numCoders
	folder.codecs = make([]uint64, numCoders)

	for i := uint64(0) ; i < numCoders ; i++ {
		mainByte, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if mainByte & 0xC0 != 0 {
			return nil, errUnsupported
		}
		codecIdSize := mainByte & 0xF
		if codecIdSize > 8 {
			return nil, errUnsupported
		}
		codecId := uint64(0)
		for pos := 0 ; pos < int(codecIdSize) ; pos++ {
			codecIdPart, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}
			codecId = codecId << 8 | uint64(codecIdPart)
		}
		folder.codecs[i] = codecId
		if mainByte & 0x10 != 0 {
			numInStream, err := readEncodedUInt64(reader)
			if err != nil {
				return nil, err
			}
			if value, err := readEncodedUInt64(reader) ; err != nil || value != uint64(1) {
				return nil, errInvalidFile
			}
			folder.numInStreams += numInStream
		}
		if mainByte & 0x20 != 0 {
			propsSize, err := readEncodedUInt64(reader)
			if err != nil {
				return nil, err
			}
			props := make([]byte, propsSize)
			if value , err := reader.Read(props) ; err != nil || value != int(propsSize) {
				return nil, errInvalidFile
			}
		}
	}
	//TODO: what to do with this
	for i := 0 ; i < int(numCoders - 1) ; i++ {
		_ , err = readEncodedUInt64(reader)
		if err != nil {
			return nil, err
		}
		_ , err = readEncodedUInt64(reader)
		if err != nil {
			return nil, err
		}
	}

	//TODO: what to do with this
	numPackStreams := folder.numInStreams - (numCoders - 1)
	for i:= 0 ; i < int(numPackStreams) ; i++ {
		_ , err = readEncodedUInt64(reader)
		if err != nil {
			return nil, err
		}
	}
	return folder, nil
}
