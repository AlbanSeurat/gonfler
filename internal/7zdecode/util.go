package _zdecode

import (
	"bufio"
)

const (
	kEnd                   = iota
	kHeader                = iota
	kArchiveProperties     = iota
	kAdditionalStreamsInfo = iota
	kMainStreamsInfo       = iota
	kFilesInfo             = iota
	kPackInfo              = iota
	kUnpackInfo            = iota
	kSubStreamsInfo        = iota
	kSize                  = iota
	kCRC                   = iota
	kFolder                = iota
	kCodersUnpackSize      = iota
	kNumUnpackStream       = iota
	kEmptyStream           = iota
	kEmptyFile             = iota
	kAnti                  = iota
	kName                  = iota
	kCTime                 = iota
	kATime                 = iota
	kMTime                 = iota
	kWinAttrib             = iota
	kComment               = iota
	kEncodedHeader         = iota
	kStartPos              = iota
	kDummy                 = iota

// kNtSecure,
// kParent,
// kIsAux

)
type functionWithError func() error

func readEncodedUInt64(reader *bufio.Reader) (uint64, error) {
	first, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}

	if (first & 0x80) == 0 {
		return uint64(first), nil
	}

	b, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	next := uint64(b)

	for i := uint64(1) ; i < 8 ; i++ {
		mask := byte(0x80) >> i
		if (first & mask) == 0 {
			return next | uint64(first & (mask - 1)) << (i *8), nil
		}
		b, err = reader.ReadByte()
		if err != nil {
			return 0, err
		}
		next |= uint64(b) << (i *8)
	}
	return next, nil
}


func readUntil(reader *bufio.Reader, id int) error {

	for {
		next, err := reader.ReadByte()
		if err != nil {
			return err
		}
		if int(next) == id {
			return nil
		}
		if int(next) == kEnd {
			return errInvalidFile
		}
		_, err = reader.ReadByte()
		if err != nil {
			return err
		}
	}
}

func readEndOrCrc(reader *bufio.Reader, crcFunction functionWithError) error {

	for {
		if packType, err := reader.ReadByte() ; err != nil {
			return err
		} else {
			if packType == kEnd {
				return nil
			}
			if packType == kCRC {
				err := crcFunction()
				if err != nil {
					return err
				}
				continue
			}
		}
	}
	return nil
}


func readBooleans(reader * bufio.Reader, numItems int) ([]bool, int, error) {

	bools := make([]bool, numItems)
	mask := byte(0)
	var b byte
	var count int
	var err error
	for i := 0; i < numItems; i++ {
		if mask == 0 {
			b, err = reader.ReadByte()
			if err != nil {
				return nil, 0, err
			}
			mask = 0x80
		}
		bools[i] = (b & mask) != 0
		if bools[i] {
			count++
		}
		mask >>= 1
	}
	return bools, count, nil
}

func readDefinedBooleans(reader *bufio.Reader, numItems int) ([]bool, int, error) {
	allAreDefined, err := reader.ReadByte()
	if err != nil {
		return nil, 0, err
	}
	if allAreDefined == 0 {
		return readBooleans(reader, numItems)
	} else {
		bools := make([]bool, numItems)
		for index, _ := range bools { bools[index] = true}
		return bools, numItems, nil
	}
}