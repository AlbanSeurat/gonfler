package lzma

import "errors"

const (
	maxPropsValue   = 225
	propsSize       = 5
	dictMinimumSize = 1 << 12
)

var (
	invalidProps = errors.New("Incorrect LZMA properties")
)

type Props struct {
	lc, lp, pb uint64
	dicSize    uint32
}

func NewProps(codedProps []byte) (*Props, error) {
	props := new(Props)
	if len(codedProps) < propsSize {
		return nil, invalidProps
	}
	d := uint64(codedProps[0])
	if d >= maxPropsValue {
		return nil, invalidProps
	}
	props.lc = d % 9
	d /= 9
	props.pb = d / 5
	props.lp = d % 5

	props.dicSize = uint32(codedProps[1]) | uint32(codedProps[2])<<8 | uint32(codedProps[3])<<16 | uint32(codedProps[4])<<24
	if props.dicSize < dictMinimumSize {
		props.dicSize = dictMinimumSize
	}
	return props, nil
}


