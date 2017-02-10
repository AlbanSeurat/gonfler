package gonfler

import (
	"fmt"
	. "github.com/onsi/gomega"
	"testing"
	"bytes"
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
	for it := archive.Volumes(); it.next != nil; it = it.next() {
		fmt.Println(it.volume)
		buf := new(bytes.Buffer)
		buf.ReadFrom(it.volume)
		fmt.Println(buf.String())
	}
}

func TestOpenZip(t *testing.T) {
	RegisterTestingT(t)

	archive, err := Open("testdata/archive.zip")
	Expect(err).ShouldNot(HaveOccurred())

	defer archive.Close()
	for it := archive.Volumes(); it.next != nil; it = it.next() {
		fmt.Println(it.volume)
		buf := new(bytes.Buffer)
		buf.ReadFrom(it.volume)
		fmt.Println(buf.String())
	}
}


func TestOpenTar(t *testing.T) {
	RegisterTestingT(t)

	archive, err := Open("testdata/archive.tar")
	Expect(err).ShouldNot(HaveOccurred())

	defer archive.Close()
	for it := archive.Volumes(); it.next != nil; it = it.next() {
		fmt.Println(it.volume)
		buf := new(bytes.Buffer)
		buf.ReadFrom(it.volume)
		fmt.Println(buf.String())
	}
}
