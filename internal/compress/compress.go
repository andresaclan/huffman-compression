package compress

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"huffman-compression/internal/huffman"

	"github.com/icza/bitio"
)

func Compress(data []byte) ([]byte, error) {
	freqTable := huffman.BuildFrequencyTable(data)
	root := huffman.BuildHuffmanTree(freqTable)
	codes := huffman.GenerateHuffmanCodes(root)
	encodedText := huffman.EncodeText(data, codes)

	// create a buffer to store the compressed data
	compressed := new(bytes.Buffer)

	// encode the huffman codes
	encodedTable, err := huffman.EncodeHuffmanCodes(codes)
	if err != nil {
		return nil, err
	}

	// write the length of the encoded table to the buffer
	tableLength := uint32(len(encodedTable))
	err = binary.Write(compressed, binary.LittleEndian, tableLength)
	if err != nil {
		return nil, err
	}

	// create a bit writer to write the compressed data bit by bit
	b := &bytes.Buffer{}
	w := bitio.NewCountWriter(b)

	// write the data to the bitWriter we are compressing bit by bit
	for _, bit := range encodedText {
		if bit == '0' {
			err := w.WriteBool(false)
			if err != nil {
				fmt.Println("Error writing bit:", err)
				return nil, err
			}
		} else {
			err := w.WriteBool(true)
			if err != nil {
				fmt.Println("Error writing bit:", err)
				return nil, err
			}

		}
	}
	fmt.Println("Bits Written during compression: ", w.BitsCount)

	// write the number of bits written to the buffer
	err = binary.Write(compressed, binary.LittleEndian, uint32(w.BitsCount))
	if err != nil {
		return nil, err
	}

	// write the encoded table and the compressed data to the buffer
	compressed.Write(encodedTable)
	compressed.Write(b.Bytes())

	// compressed looks like this: [tableLength][bitsWritten][encodedTable][compressedData]
	return compressed.Bytes(), nil
}
