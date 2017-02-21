package _zdecode

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

const (
	sigPrefix     = "7z\xBC\xAF\x27\x1C"
	sigPrefixSize = 32
	maxOffset     = uint64(^uint(0) >> 1)
	maxLength     = uint64(1 << 62)
)

var (
	errNoSign        = errors.New("invalid 7z signature")
	errCrcCheck      = errors.New("checksum error")
	errNextHeader    = errors.New("invalid header")
	errInvalidFile   = errors.New("invalid file")
	errCodecNotFound = errors.New("codec not found")

	errUnsupported = errors.New("unsupported")
)

type version struct {
	Major byte
	Minor byte
}

type header struct {
	Signature  [6]byte
	Version    version
	CrcCheck   uint32
	NextHeader [20]byte
}

type nextHeader struct {
	EndHeaderOffset uint64
	EndHeaderLen    uint64
	CrcEndHeader    uint32
}

func headerCheck(reader io.Reader) (*header, error) {
	header := new(header)
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

func readDatabase(header *header, file *os.File) error {
	nextHeader := new(nextHeader)
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

	_, err = file.Seek(int64(nextHeader.EndHeaderOffset), 1)
	if err != nil {
		return errNextHeader
	}

	nextHeaderContent := make([]byte, nextHeader.EndHeaderLen)
	file.Read(nextHeaderContent)
	if nextHeader.CrcEndHeader != crc32.Checksum(nextHeaderContent, crc32.IEEETable) {
		return errCrcCheck
	}
	if nextHeaderContent[0] == kHeader {
		return errUnsupported
	} else if nextHeaderContent[0] != kEncodedHeader {
		return errInvalidFile
	}

	ssInfo, err := readStreamInfo(bytes.NewReader(nextHeaderContent[1:]))
	if err != nil {
		return err
	}

	fmt.Println(decodeStream(file, ssInfo, sigPrefixSize))

	return nil
}
