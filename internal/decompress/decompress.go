package decompress

import (
	"bytes"
	"huffman-compression/internal/huffman"

	"github.com/icza/bitio"
)

func Decompress(data []byte) ([]byte, error) {
	b := bytes.NewBuffer(data)
	r := bitio.NewCountReader(b)

	// Construct the HuffmanTree from the Reader
	root := huffman.GetHuffmanTree(r)

	huffmanTree := &huffman.HuffmanTree{
		Htree: &huffman.Tree{Root: root},
	}

	// fmt.Println("Displaying Codes:")
	// huffmanTree.CreateCodes()
	// huffmanTree.DisplayCodes()

	// Decode the data
	output := huffmanTree.Decode(r)

	return output, nil
}
