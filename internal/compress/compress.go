package compress

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"huffman-compression/internal/huffman"
	"huffman-compression/pkg/bitio/bitwriter"
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
	bitWriter, err := bitwriter.NewBitWriter()
	if err != nil {
		return nil, err
	}
	defer bitWriter.Close()

	// write the data to the bitWriter we are compressing bit by bit
	var bitsWritten uint32
	for _, bit := range encodedText {
		if bit == '0' {
			err := bitWriter.WriteBit(0)
			if err != nil {
				fmt.Println("Error writing bit:", err)
				return nil, err
			}
			bitsWritten++
		} else {
			err := bitWriter.WriteBit(1)
			if err != nil {
				fmt.Println("Error writing bit:", err)
				return nil, err
			}
			bitsWritten++

		}
	}
	fmt.Println("Bits Written during compression: ", bitsWritten)

	// write the number of bits written to the buffer
	err = binary.Write(compressed, binary.LittleEndian, &bitsWritten)
	if err != nil {
		return nil, err
	}

	// write the encoded table and the compressed data to the buffer
	compressed.Write(encodedTable)
	compressed.Write(bitWriter.Bytes.Bytes())

	// compressed looks like this: [tableLength][bitsWritten][encodedTable][compressedData]
	return compressed.Bytes(), nil
}
