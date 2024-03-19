package io

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
)

// Open a file and returns a reader
func Open(filename string, blockSize int) (*bufio.Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	fileBlocks := int(math.Ceil(float64(fileSize) / float64(blockSize)))
	if fileBlocks <= 1 {
		return nil, fmt.Errorf("at least 2 blocks are required")
	}

	return bufio.NewReader(file), nil
}

// Reads a block of blockSize from the reader
func GetBlock(reader *bufio.Reader, blockSize int) []byte {
	block := make([]byte, blockSize)
	bytesRead, err := reader.Read(block)
	if bytesRead == 0 || err == io.EOF {
		return nil
	}

	return block
}

// Reads a byte from the reader
func GetByte(reader *bufio.Reader) *byte {
	c, err := reader.ReadByte()
	if err == io.EOF || err != nil {
		return nil
	}

	return &c
}
