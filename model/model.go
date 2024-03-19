package model

type Signature struct {
	BlockSize int
	Hashes    []uint32
}

type Delta struct {
	Blocks []Block
}

type Block struct {
	Op        string // O - original position, M - moved, C - changed
	FromStart int    // start index from the old version
	FromEnd   int    // end index from the old version
	ToStart   int    // start index from the new version
	ToEnd     int    // end index from the new version

	// for changed blocks
	// DiffBytes []byte
}
