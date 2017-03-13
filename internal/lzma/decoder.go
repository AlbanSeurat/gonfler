package lzma

import "errors"

const (
	kNumPosBitsMax     = 4
	kNumLenToPosStates = 4
	kEndPosModelIndex  = 14
	kNumFullDistances  = 1 << (kEndPosModelIndex >> 1)
	kNumAlignBits      = 4
	kNumStates         = 12
	kMatchMinLen       = 2
)

var (
	errError = errors.New("decoder internal error")
	errNoMarker = errors.New("finished without marker")
	errWithMarker = errors.New("finished with marker")
)


type state struct {
	isMatch [kNumStates << kNumPosBitsMax]prob
	isRep [kNumStates]prob
	isRepG0 [kNumStates]prob
	isRepG1[kNumStates]prob
	isRepG2[kNumStates]prob
	isRep0Long[kNumStates << kNumPosBitsMax]prob

	lenDecoder lendecoder
	repLenDecoder lendecoder
}

type LzmaDecoder struct {
	literals
	distance
	state
	props         Props
	out           outWindow
	rangeDec      rangedecoder
	lenDecoder    lendecoder
	repLenDecoder lendecoder
}

func (ld LzmaDecoder) updateLiteral(state uint) uint {
	if state < 4 {
		return 0
	} else if state < 10 {
		return state - 3
	} else {
		return state - 6
	}
}

func (ld LzmaDecoder) updateRep(state uint) uint {
	if state < 7 {
		return 8
	} else {
		return 11
	}
}

func (ld LzmaDecoder) updateShortRep(state uint) uint {
	if state < 7 {
		return 9
	} else {
		return 11
	}
}

func (ld LzmaDecoder) updateMatch(state uint) uint {
	if state < 7 {
		return 7
	} else {
		return 10
	}
}


func NewDecoder(stream []byte) *LzmaDecoder {
	decoder := new(LzmaDecoder)

	return decoder
}

func (ld *LzmaDecoder) Decode(unpackSize uint64) error {
	rep0 := uint(0)
	rep1 := uint(0)
	rep2 := uint(0)
	rep3 := uint(0)
	state := uint(0)

	for {
		if unpackSize == 0 {
			if ld.rangeDec.isFinishedOK() {
				return errNoMarker
			}
		}

		posState := uint(ld.out.totalPos & ((1 << ld.props.pb) - 1))
		db, err := ld.rangeDec.decodeBit(&ld.state.isMatch[(state << kNumPosBitsMax) + posState])
		if err != nil {
			return err
		}
		if db == 0 {
			if unpackSize == 0 {
				return errError
			}
			if err := ld.decodeLiteral(state, int(rep0)); err != nil {
				return err
			}
			state = ld.updateLiteral(state)
			unpackSize--
			continue
		}

		var len uint
		db, err = ld.rangeDec.decodeBit(&ld.state.isRep[state])
		if err != nil {
			return err
		}
		if db != 0 {
			if unpackSize == 0 || ld.out.isEmpty() {
				return errError
			}
			db, err = ld.rangeDec.decodeBit(&ld.state.isRepG0[state])
			if err != nil {
				return err
			}
			if db == 0 {
				db, err = ld.rangeDec.decodeBit(&ld.state.isRep0Long[(state << kNumPosBitsMax) + posState])
				if err != nil {
					return err
				}
				if db == 0 {
					state = ld.updateShortRep(state)
					ld.out.putByte(ld.out.getByte(int(rep0 + 1)))
					unpackSize--
					continue
				}
			} else {
				var dist uint
				db, err = ld.rangeDec.decodeBit(&ld.state.isRepG1[state])
				if err != nil {
					return err
				}
				if db == 0 {
					dist = rep1
				} else {
					db, err = ld.rangeDec.decodeBit(&ld.state.isRepG2[state])
					if err != nil {
						return err
					}
					if db == 0 {
						dist = rep2
					} else {
						dist = rep3
						rep3 = rep2
					}
					rep2 = rep1
				}
				rep1 = rep0
				rep0 = dist
			}
			if len, err = ld.repLenDecoder.decode(&ld.rangeDec, posState) ; err != nil {
				return err
			}
			state = ld.updateRep(state)

		} else {
			rep3 = rep2
			rep2 = rep1
			rep1 = rep0
			if len, err = ld.lenDecoder.decode(&ld.rangeDec, posState) ; err != nil {
				return err
			}
			state = ld.updateMatch(state)
			if rep0, err = ld.decodeDistance(len) ; err != nil {
				return err
			}
			if rep0 == 0xFFFFFFFF {
				if ld.rangeDec.isFinishedOK() {
					return errWithMarker
				} else {
					return errError
				}
			}
			if unpackSize == 0 || rep0 >= uint(ld.props.dicSize) || !ld.out.checkDistance(int(rep0)) {
				return errError
			}
			len += kMatchMinLen
			if unpackSize < uint64(len) {
				return errError
			}
			if err = ld.out.copyMatch(int(rep0 + 1), len) ; err != nil {
				return err
			}
			unpackSize -= uint64(len)
		}
	}
}

