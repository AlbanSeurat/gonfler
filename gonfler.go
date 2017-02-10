package gonfler

import (
	"fmt"
	"github.com/h2non/filetype"
	"io"
)

type Volume struct {
	io.Reader
	name string
}

func (Volume) Close() error {
	return nil
}

type VolumeIterator struct {
	volume *Volume
	next   func() VolumeIterator
}

type Archive interface {
	volumes() VolumeIterator
	Close() error
}

func Open(name string) (Archive, error) {
	fileType, e := filetype.MatchFile(name)
	if e != nil {
		return nil, e
	}
	switch fileType.MIME.Value {
	case "application/x-rar-compressed":
		return openRar(name)
	case "application/zip":
		return openZip(name)
	default:
		return nil, fmt.Errorf("%s is not a recognized file (%s)", name, fileType.MIME.Value)

	}
	return nil, e
}
