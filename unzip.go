package gonfler

import "archive/zip"

type ZipArchive struct {
	handle *zip.ReadCloser
}

func (archive ZipArchive) Close() error {
	return archive.handle.Close()
}

func (archive ZipArchive) Volumes() VolumeIterator {
	pos := 0
	var next func() VolumeIterator
	next = func() VolumeIterator {
		if len(archive.handle.File) == pos {
			return VolumeIterator{nil, nil}
		} else {
			file := archive.handle.File[pos]
			pos++
			fileHandle, _ := file.Open()
			return VolumeIterator{
				volume: &Volume{fileHandle, file.Name},
				next:   next,
			}
		}
	}
	return next()
}

func openZip(name string) (Archive, error) {
	handle, e := zip.OpenReader(name)
	if handle != nil {
		return ZipArchive{handle}, nil
	} else {
		return nil, e
	}
}
