package gonfler

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
)

func TestOpen(t *testing.T) {
	RegisterTestingT(t)

	_, err := Open("testdata/lorempixel.jpg")
	Expect(err).Should(HaveOccurred())
}

func TestOpenRar(t *testing.T) {
	RegisterTestingT(t)

	archive, err := Open("testdata/archive.rar")
	Expect(err).ShouldNot(HaveOccurred())
	defer archive.Close()
	for it := archive.volumes(); it.next != nil; it = it.next() {
		fmt.Println(it.volume)
		it.volume.Close()
	}
}

func TestOpenZip(t *testing.T) {
	RegisterTestingT(t)

	archive, err := Open("testdata/archive.zip")
	Expect(err).ShouldNot(HaveOccurred())

	defer archive.Close()
	for it := archive.volumes(); it.next != nil; it = it.next() {
		fmt.Println(it.volume)
		it.volume.Close()
	}
}
