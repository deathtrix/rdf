package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/deathtrix/rdf/delta"
	"github.com/deathtrix/rdf/io"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("usage:", os.Args[0], "<block-size> <oldfilename> <newfilename>")
		return
	}
	blockSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(fmt.Errorf("block-size must be integer"))
		return
	}
	oldFilename := os.Args[2]
	newFilename := os.Args[3]

	oldReader, err := io.Open(oldFilename, blockSize)
	if err != nil {
		fmt.Println(err)
		return
	}
	newReader, err := io.Open(newFilename, blockSize)
	if err != nil {
		fmt.Println(err)
		return
	}

	sig, err := delta.BuildSignature(oldReader, blockSize)
	if err != nil {
		fmt.Println(err)
		return
	}
	diff, err := delta.BuildDelta(newReader, sig)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(delta.FormatDelta(diff))
}
