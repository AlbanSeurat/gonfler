package lzma

import "testing"
import (
	. "github.com/onsi/gomega"
	"os"
	"encoding/binary"
)

func TestOpen(t *testing.T) {
	RegisterTestingT(t)

	file, _ := os.Open("testdata/data.eos.l3.lzma")
	properties := make([]byte, 5)
	file.Read(properties)
	var fileSize uint64
	binary.Read(file, binary.LittleEndian, &fileSize)


}

