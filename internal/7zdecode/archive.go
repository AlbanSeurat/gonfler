package _zdecode

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"os"
	"fmt"
)

const (
	//sigPrefix = []byte{'7', 'z', 0xBC, 0xAF, 0x27, 0x1C};
	sigPrefix = "7z\xBC\xAF\x27\x1C"
	maxOffset = uint64(^uint(0) >> 1)
	maxLength = uint64(1 << 62)
)

var (
	errNoSign     = errors.New("invalid 7z signature")
	errCrcCheck   = errors.New("checksum error")
	errNextHeader = errors.New("invalid header")
	errInvalidFile = errors.New("invalid file")

	errUnsupported = errors.New( "unsupported")
)

type _version struct {
	Major byte
	Minor byte
}

type _header struct {
	Signature  [6]byte
	Version    _version
	CrcCheck   uint32
	NextHeader [20]byte
}

type _nextHeader struct {
	EndHeaderOffset uint64
	EndHeaderLen    uint64
	CrcEndHeader    uint32
}

func headerCheck(reader io.Reader) (*_header, error) {
	header := new(_header)
	err := binary.Read(reader, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal([]byte(sigPrefix), header.Signature[:]) {
		return nil, errNoSign
	}
	if header.CrcCheck != crc32.Checksum(header.NextHeader[:], crc32.IEEETable) {
		return nil, errCrcCheck
	}
	return header, nil
}

func readDatabase(header *_header, file *os.File) error {
	nextHeader := new(_nextHeader)
	err := binary.Read(bytes.NewReader(header.NextHeader[:]), binary.LittleEndian, nextHeader)
	if err != nil {
		return err
	}
	if nextHeader.EndHeaderOffset > maxOffset || nextHeader.EndHeaderLen > maxLength {
		return errNextHeader
	}
	if nextHeader.EndHeaderLen == 0 && nextHeader.EndHeaderOffset != 0 {
		return errNextHeader
	}

	_, err = file.Seek(int64(nextHeader.EndHeaderOffset), io.SeekCurrent)
	if err != nil {
		return errNextHeader
	}

	nextHeaderContent := make([]byte, nextHeader.EndHeaderLen)
	file.Read(nextHeaderContent)
	if nextHeader.CrcEndHeader != crc32.Checksum(nextHeaderContent, crc32.IEEETable) {
		return errCrcCheck
	}
	fmt.Println(nextHeaderContent)
	switch nextHeaderContent[0] {
	case kHeader:
		return nil
	case kEncodedHeader:
		readStreamInfo(bytes.NewReader(nextHeaderContent[1:]))
	}

	return nil
}
