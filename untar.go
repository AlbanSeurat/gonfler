package gonfler

import (
	"archive/tar"
	"os"
)

type TarArchive struct {
	handle *tar.Reader
	file   *os.File
}

func (archive TarArchive) Close() error {
	return archive.file.Close()
}

func (archive TarArchive) Volumes() VolumeIterator {
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

func openTar(name string) (Archive, error) {
	file, e := os.Open(name)
	if file != nil {
		return TarArchive{tar.NewReader(file), file}, nil
	} else {
		return nil, e
	}
}
