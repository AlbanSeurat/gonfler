package lzma


type literal struct {
	probs []prob
}

func newLiteral(lc, lp int) (*literal) {
	literal := new(literal)
	literal.probs = make([]prob, 0x300<<uint(lc+lp))
	for i := range literal.probs {
		literal.probs[i] = probInit
	}
	return literal
}


func (*literal) decode() {
	
}
