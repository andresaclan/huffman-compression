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

func decompress(data []byte) ([]byte, error) {
	var tableLength uint32
	buf := bytes.NewReader(data)

	err := binary.Read(buf, binary.LittleEndian, &tableLength)
	if err != nil {
		return nil, err
	}

	encodedTable := make([]byte, tableLength)
	_, err = buf.Read(encodedTable)
	if err != nil {
		return nil, err
	}
	codes, err := huffman.DecodeHuffmanCodes(encodedTable)
	if err != nil {
		return nil, err
	}
	data = data[4+tableLength:]
	decodedText := huffman.DecodeText(data, codes)
	return decodedText, nil

}
func main() {
	inputFilePath, outputFilePath := cli.ParseArgs(os.Args)
	data, err := file.ReadFile(inputFilePath)
	if err != nil { // handle error
		fmt.Println("Error reading file:", err)
		return
	}
	decompressedData, err := decompress(data)
	if err != nil { // handle error
		fmt.Println("Error decompressing data:", err)
		return
	}
	err = file.WriteFile(outputFilePath, decompressedData)
	if err != nil { // handle error
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Println("File decompressed successfully!")

}
