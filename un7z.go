package gonfler

import (
	"github.com/alkpone/gonfler/internal/7zdecode"
)

type _7zArchive struct {
	handle interface{}
}

func (archive _7zArchive) Close() error {
	return nil
}

func (archive _7zArchive) Volumes() VolumeIterator {
	return VolumeIterator{nil, nil}
}

func open7z(name string) (Archive, error) {
	handle, e := _zdecode.OpenReader(name)
	if handle != nil {
		return _7zArchive{handle}, nil
	} else {
		return nil, e
	}
}
