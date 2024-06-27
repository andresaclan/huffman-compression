package compress

import (
	"bytes"
	"huffman-compression/internal/huffman"
)

func Compress(data []byte) ([]byte, error) {
	freqTable := huffman.GetFrequencies(data)
	huffmanTree := huffman.NewHuffmanTreeFromFrequencies(freqTable)

	// fmt.Println("Displaying Codes:")
	// huffmanTree.CreateCodes()
	// huffmanTree.DisplayCodes()

	b := &bytes.Buffer{}
	w := huffman.EncodeHuffmanTree(huffmanTree.Htree.Root, b)

	huffmanTree.Encode(data, w)

	return b.Bytes(), nil

}
