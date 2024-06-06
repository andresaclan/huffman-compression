package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"huffman-compression/internal/file"
	"huffman-compression/internal/huffman"
	"huffman-compression/pkg/cli"
	"os"
)

// takes already compressed data along with the non-serialized huffman codes
func compress(data []byte, codes map[rune]string) ([]byte, error) {
	compressed := new(bytes.Buffer)
	encodedTable, err := huffman.EncodeHuffmanCodes(codes)
	if err != nil {
		return nil, err
	}
	tableLength := uint32(len(encodedTable))
	err = binary.Write(compressed, binary.LittleEndian, tableLength)
	if err != nil {
		return nil, err
	}
	compressed.Write(encodedTable)
	compressed.Write(data)

	return compressed.Bytes(), nil
}

func main() {
	inputFilePath, outputFilePath := cli.ParseArgs(os.Args)
	data, err := file.ReadFile(inputFilePath)
	uncompressedSize := len(data)
	fmt.Println("Uncompressed file size:", uncompressedSize, "bytes")
	if err != nil { // handle error
		fmt.Println("Error reading file:", err)
		return
	}
	freqTable := huffman.BuildFrequencyTable(data)
	root := huffman.BuildHuffmanTree(freqTable)
	codes := huffman.GenerateHuffmanCodes(root)
	encodedText := huffman.EncodeText(data, codes)
	encodedSize := len(encodedText)
	fmt.Println("Encoded file size:", encodedSize, "bytes")
	compressedData, err := compress(encodedText, codes)
	if err != nil { // handle error
		fmt.Println("Error compressing data:", err)
	}

	err = file.WriteFile(outputFilePath, compressedData)

	if err != nil { // handle error
		fmt.Println("Error writing file:", err)
	}

	fmt.Println("File compressed successfully!")

}
