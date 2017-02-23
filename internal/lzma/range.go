package lzma

import (
	"bytes"
)

const (
	kTopValue = uint32(1) << 24
)

type rangedecoder struct {
	nrange uint32
	code   uint32
	reader bytes.Reader
}

func (rd *rangedecoder) init() error {
	rd.nrange = 0xffffffff
	rd.code = 0

	b, err := rd.reader.ReadByte()
	if err != nil {
		return err
	}
	if b != 0 {
		return lzmaInvalidFormat
	}

	for i := 0; i < 4; i++ {
		if err := rd.updateCode(); err != nil {
			return err
		}
	}

	if rd.code >= rd.nrange {
		return lzmaInvalidFormat
	}

	return nil
}

func (rd *rangedecoder) isFinishedOK() bool {
	return rd.code == 0
}

func (rd *rangedecoder) updateCode() error {
	b, err := rd.reader.ReadByte()
	if err != nil {
		return err
	}
	rd.code = (rd.code << 8) | uint32(b)
	return nil
}

func (rd *rangedecoder) normalize() error {
	if rd.nrange < kTopValue {
		rd.nrange <<= 8
		if err := rd.updateCode(); err != nil {
			return err
		}
	}
	return nil
}

func (rd *rangedecoder) decodeDirectBits(numBits uint32) (uint32, error) {
	res := uint32(0)
	for {
		rd.nrange >>= 1
		rd.code -= rd.nrange
		t := 0 - (rd.code >> 31)
		rd.code += rd.nrange & t

		if err := rd.normalize(); err != nil {
			return 0, err
		}
		res <<= 1
		res += t + 1
		numBits--
		if numBits <= 0 {
			break
		}
	}
	return res, nil
}

func (rd *rangedecoder) decodeBit(prob *prob) (uint, error) {
	symbol := uint(0)
	bound := prob.bound(rd.nrange)
	if rd.code < bound {
		prob.inc()
		rd.nrange = bound
	} else {
		prob.dec()
		rd.code -= bound
		rd.nrange -= bound
		symbol = 1
	}
	if err := rd.normalize() ; err != nil {
		return 0, err
	}
	return symbol, nil
}

