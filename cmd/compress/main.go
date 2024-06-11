package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"huffman-compression/internal/bitio/bitwriter"
	"huffman-compression/internal/file"
	"huffman-compression/internal/huffman"
	"huffman-compression/pkg/cli"
	"os"
)

// takes already compressed data along with the non-serialized huffman codes
func compress(data []byte, codes map[rune]string) ([]byte, error) {
	fmt.Println("Original Codes:")
	for k, v := range codes {
		fmt.Printf("%c:%s\n", k, v)
	}
	compressed := new(bytes.Buffer)
	encodedTable, err := huffman.EncodeHuffmanCodes(codes)
	if err != nil {
		return nil, err
	}
	tableLength := uint32(len(encodedTable))
	fmt.Println("Table length:", tableLength)
	err = binary.Write(compressed, binary.LittleEndian, tableLength)
	if err != nil {
		return nil, err
	}

	// write the data bit by bit
	bitWriter, err := bitwriter.NewBitWriter()
	if err != nil {
		return nil, err
	}
	defer bitWriter.Close()

	var bitsWritten uint32
	for _, bit := range data {
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

	err = binary.Write(compressed, binary.LittleEndian, &bitsWritten)
	if err != nil {
		return nil, err
	}

	compressed.Write(encodedTable)
	compressed.Write(bitWriter.Bytes.Bytes())

	return compressed.Bytes(), nil
}

func main() {
	inputFilePath, outputFilePath := cli.ParseArgs(os.Args)
	data, err := file.ReadFile(inputFilePath)
	// uncompressedSize := len(data)
	// fmt.Println("Uncompressed file size:", uncompressedSize, "bytes")
	if err != nil { // handle error
		fmt.Println("Error reading file:", err)
		return
	}
	freqTable := huffman.BuildFrequencyTable(data)
	var totalChars uint32
	for _, v := range freqTable {
		totalChars += uint32(v)
	}

	root := huffman.BuildHuffmanTree(freqTable)
	codes := huffman.GenerateHuffmanCodes(root)
	encodedText := huffman.EncodeText(data, codes)
	compressedData, err := compress(encodedText, codes)
	if err != nil { // handle error
		fmt.Println("Error compressing data:", err)
		return
	}
	err = file.WriteFile(outputFilePath, compressedData)
	if err != nil { // handle error
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Println("File compressed successfully!")

}
