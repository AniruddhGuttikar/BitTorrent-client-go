package bitfield

type BitField []byte

func (bf BitField) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	destByte := bf[byteIndex]

	return destByte>>(7-offset)&1 == 1
}

func (bf *BitField) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8
	
	(*bf)[byteIndex] |= 1 << (7-offset)
}
