package lzma

type prob uint16

const (
	kNumMoveBits = 5
	kNumBitModelTotalBits = 11
	probInit prob = 1 << (kNumBitModelTotalBits - 1)
)

// Dec decreases the probability. The decrease is proportional to the
// probability value.
func (p *prob) dec() {
	*p -= *p >> kNumMoveBits
}

// Inc increases the probability. The Increase is proportional to the
// difference of 1 and the probability value.
func (p *prob) inc() {
	*p += ((1 << kNumBitModelTotalBits) - *p) >> kNumMoveBits
}

// Computes the new bound for a given range using the probability value.
func (p prob) bound(r uint32) uint32 {
	return (r >> kNumBitModelTotalBits) * uint32(p)
}