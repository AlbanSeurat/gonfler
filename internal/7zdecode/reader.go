package _zdecode

import (
	"os"
)

type ReadCloser struct {
	file *os.File
}

// Close closes the rar file.
func (rc *ReadCloser) Close() error {
	return nil
}

func openVolume(name string) (*os.File, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	header, err := headerCheck(file)
	if err != nil {
		return nil, err
	}
	return file, readDatabase(header, file)
}

// OpenReader opens a 7z archive specified by the name and returns a ReadCloser.
func OpenReader(name string) (*ReadCloser, error) {
	file, err := openVolume(name)
	if err != nil {
		return nil, err
	}
	rc := new(ReadCloser)
	rc.file = file
	return rc, nil
}
