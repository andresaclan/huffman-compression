package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"huffman-compression/internal/bitio/bitreader"
	"huffman-compression/internal/file"
	"huffman-compression/internal/huffman"
	"huffman-compression/pkg/cli"
	"io"
	"os"
)

func decompress(data []byte) ([]byte, error) {
	var tableLength uint32
	buf := bytes.NewReader(data)

	err := binary.Read(buf, binary.LittleEndian, &tableLength)
	if err != nil {
		return nil, err
	}
	fmt.Println("Table length:", tableLength)
	var bitsWritten uint32
	err = binary.Read(buf, binary.LittleEndian, &bitsWritten)
	if err != nil {
		return nil, err
	}
	fmt.Println("Total Bits to be Read:", bitsWritten)
	encodedTable := make([]byte, tableLength)
	_, err = buf.Read(encodedTable)
	if err != nil {
		return nil, err
	}
	codes, err := huffman.DecodeHuffmanCodes(encodedTable)
	if err != nil {
		return nil, err
	}

	// fmt.Println("Decoded Codes:")
	// for k, v := range codes {
	// 	fmt.Printf("%c:%s\n", k, v)
	// }

	data = data[8+tableLength:]

	encodedData := make([]byte, 0)
	bitReader, _ := bitreader.NewBitReader(bytes.NewBuffer(data))
	var bitsRead int
	for {
		bit, err := bitReader.ReadBit()
		if errors.Is(err, io.EOF) {
			break
		}
		if bit == 0 {
			bitsRead++
			encodedData = append(encodedData, '0')
		} else {
			encodedData = append(encodedData, '1')
			bitsRead++
		}
	}
	fmt.Println("Bits read:", bitsRead)
	if bitsRead != int(bitsWritten) {
		for i := 0; i < int(bitsWritten)-bitsRead; i++ {
			encodedData = append(encodedData, '0')
		}
	}
	decodedText := huffman.DecodeText(encodedData, codes)
	return decodedText, nil

}
func main() {
	inputFilePath, outputFilePath := cli.ParseArgs(os.Args)
	data, err := file.ReadFile(inputFilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	decompressedData, err := decompress(data)
	if err != nil {
		fmt.Println("Error decompressing data:", err)
		return
	}
	err = file.WriteFile(outputFilePath, decompressedData)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Println("File decompressed successfully!")

}
