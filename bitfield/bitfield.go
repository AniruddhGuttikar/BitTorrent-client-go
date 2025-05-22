package bitfield

type BitField []byte

func (bf BitField) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8

	if byteIndex < 0 || byteIndex >= len(bf) {
		return false
	}

	destByte := bf[byteIndex]

	return destByte>>(7-offset)&1 == 1
}

func (bf BitField) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8

	// silently discard invalid bounded index
	if byteIndex < 0 || byteIndex >= len(bf) {
		return
	}
	bf[byteIndex] |= 1 << uint(7-offset)
}
