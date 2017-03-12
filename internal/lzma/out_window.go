package lzma

import "bufio"

type outWindow struct {
	bufio.Writer
	buf []byte
	pos int
	totalPos uint32
	isFull bool
}


func NewOutWindow(dictSite uint32) (*outWindow){
	window := new(outWindow)
	window.buf = make([]byte, dictSite)
	window.pos = 0
	window.totalPos = 0
	window.isFull = false

	return window
}

func (w *outWindow) putByte(aByte byte) error {
	w.totalPos++
	w.buf[w.pos] = aByte
	w.pos++
	if w.pos == len(w.buf) {
		w.isFull = true
		w.pos = 0
	}
	return w.Writer.WriteByte(aByte)
}

func (w *outWindow) getByte(dist int) (byte) {
	calcPos := len(w.buf) - dist + w.pos
	if dist <= w.pos {
		calcPos = w.pos - dist
	}
	return w.buf[calcPos]
}

func (w *outWindow) copyMatch(dist int, len uint) error {
	for ; len > 0 ; len-- {
		err := w.putByte(w.getByte(dist))
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *outWindow) checkDistance(dist int) bool {
	return dist <= w.pos || w.isFull
}

func (w *outWindow) isEmpty() bool {
	return w.pos == 0 && !w.isFull
}
