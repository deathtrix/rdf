package delta

import (
	"bufio"
	"fmt"
	"sort"

	"github.com/deathtrix/rdf/hash"
	"github.com/deathtrix/rdf/model"

	"github.com/deathtrix/rdf/io"
)

// Builds signature of the old file version
func BuildSignature(reader *bufio.Reader, blockSize int) (*model.Signature, error) {
	signature := model.Signature{
		BlockSize: blockSize,
	}
	h := hash.Init(blockSize)

	// hash old file version block by block
	for {
		block := io.GetBlock(reader, blockSize)
		if block == nil {
			break
		}

		h.HashBlock(block)
		signature.Hashes = append(signature.Hashes, h.GetHash())
	}

	return &signature, nil
}

// Builds delta from signature and a new file version
func BuildDelta(reader *bufio.Reader, signature *model.Signature) (*model.Delta, error) {
	blockSize := signature.BlockSize
	block := []byte{}
	newSignature := model.Signature{}
	h := hash.Init(blockSize)
	newFileSize := 0
	first := true

	// hash new file with rolling window
	for {
		// read byte by byte and append to block
		b := io.GetByte(reader)
		if b == nil {
			break
		}
		block = append(block, *b)

		// if block is filled
		if len(block) >= blockSize {
			if first {
				// first block hash
				h.HashBlock(block)
				first = false
			} else {
				// rolling hash for the next bytes added
				h.HashRolling(block)
			}
			newSignature.Hashes = append(newSignature.Hashes, h.GetHash())
		}

		// remove first byte in block as you roll the window
		if len(block) > blockSize {
			block = block[1:]
		}
		newFileSize++
	}

	// hash end bytes in block if fileSize mod blockSize != 0
	for range blockSize - 1 {
		block = append(block, 0)
		h.HashRolling(block)
		newSignature.Hashes = append(newSignature.Hashes, h.GetHash())
		block = block[1:]
	}

	// calculate delta and indexes
	diff := model.Delta{}
	previ := 0
	for i, kn := range signature.Hashes {
		sd := model.Block{
			FromStart: i * blockSize,
			FromEnd:   (i + 1) * blockSize,
		}
		prevj := 0
		for j, ko := range newSignature.Hashes {
			// if hash found
			if ko == kn {
				if j == i*blockSize {
					// if same position
					sd.Op = "O"
					sd.ToStart = j
					sd.ToEnd = j + blockSize
					if len(diff.Blocks) < len(signature.Hashes) {
						diff.Blocks = append(diff.Blocks, sd)
					} else {
						break // exit for if enough blocks found
					}
				} else {
					// if different position
					sd.Op = "M"
					if len(diff.Blocks) < len(signature.Hashes) {
						// make sure the move is in another block otherwise if you have
						// a repeating character > blockSize, it will detect it at the next
						// position
						if i > previ || j >= prevj+blockSize {
							sd.ToStart = j
							sd.ToEnd = j + blockSize
							if sd.ToEnd > newFileSize {
								sd.ToEnd = newFileSize
							}
							diff.Blocks = append(diff.Blocks, sd)
							prevj = j
							previ = i
						}
					} else {
						break // exit for if enough blocks found
					}
				}
			}
			if len(diff.Blocks) >= len(signature.Hashes) {
				break
			}
		}
	}

	// fill in changed blocks with calculated indexes
	sd := model.Block{Op: "C"}
	if len(diff.Blocks) > 0 {
		prevBlock := model.Block{ToEnd: 2 * blockSize}

		for _, b := range diff.Blocks {
			if prevBlock.ToEnd < b.ToStart {
				sd.ToStart = prevBlock.ToEnd
				sd.ToEnd = b.ToStart
				if sd.ToEnd > newFileSize {
					sd.ToEnd = newFileSize
				}
				diff.Blocks = append(diff.Blocks, sd)
			}
			prevBlock = b
		}
	} else {
		// if no original or moved blocks are found
		for i := 0; i < newFileSize; i += blockSize {
			sd.ToStart = i
			sd.ToEnd = i + blockSize
			diff.Blocks = append(diff.Blocks, sd)
		}
	}

	return &diff, nil
}

// Formats delta for output
func FormatDelta(diff *model.Delta) string {
	comparePrice := func(i, j int) bool {
		return diff.Blocks[i].ToStart < diff.Blocks[j].ToStart
	}
	sort.Slice(diff.Blocks, comparePrice)

	ret := ""
	for _, d := range diff.Blocks {
		ret += fmt.Sprintf("%s %d:%d %d:%d\n", d.Op, d.ToStart, d.ToEnd, d.FromStart, d.FromEnd)
	}

	return ret
}
