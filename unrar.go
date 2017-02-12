package gonfler

import (
	rar "github.com/nwaples/rardecode"
)

type RarArchive struct {
	handle *rar.ReadCloser
}

func (archive RarArchive) Close() error {
	return archive.handle.Close()
}

func (archive RarArchive) Volumes() VolumeIterator {

	var next func() VolumeIterator
	next = func() VolumeIterator {
		header, err := archive.handle.Next()
		if err != nil {
			return VolumeIterator{nil, nil}
		} else {
			return VolumeIterator{
				volume: &Volume{archive.handle, header.Name},
				next:   next,
			}
		}
	}
	return next()
}

func openRar(name string) (Archive, error) {
	handle, e := rar.OpenReader(name, "")
	if handle != nil {
		return RarArchive{handle}, nil
	} else {
		return nil, e
	}
}
