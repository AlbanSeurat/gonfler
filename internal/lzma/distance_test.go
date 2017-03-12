package lzma

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestNewDistance(t *testing.T) {
	RegisterTestingT(t)

	distance := newDistance()
	Expect(distance.posDecoders).Should(Equal(nil))
}

