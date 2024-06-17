package main

import (
	"fmt"
	"huffman-compression/internal/compress"
	"huffman-compression/internal/decompress"
	"huffman-compression/internal/file"
	"huffman-compression/pkg/cli"
	"os"
)

func main() {
	function, inputFilePath, outputFilePath := cli.ParseArgs(os.Args)

	if function == "" {
		fmt.Println("Invalid arguments")
		return
	} else if function == "compress" {
		// compress logic
		data, err := file.ReadFile(inputFilePath)
		if err != nil { // handle error
			fmt.Println("Error reading file:", err)
			return
		}
		compressedData, err := compress.Compress(data)
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
		return
	} else if function == "decompress" {
		// decompress logic
		data, err := file.ReadFile(inputFilePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		decompressedData, err := decompress.Decompress(data)
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
		return
	} else {
		fmt.Println("Invalid function")
		return
	}
}
