package internal

import (
	"fmt"
	"huffman-compression/internal/compress"
	"huffman-compression/internal/decompress"
	"testing"
)

func TestCompressAndDecompress(t *testing.T) {
	// This test will compress a string, decompress it, and verify the result
	tests := [][]byte{
		[]byte("test"),
		[]byte("hello"),
		[]byte("this is a test of the huffman compression algorithm with a lot of data to compress"),
		[]byte("After a while Uncle came in, in a Cossack coat, blue trousers, and small top boots. And Natasha felt that this costume, the very one she had regarded with surprise and amusement at Otradnoe, was just the right thing and not at all worse than a swallow-tail or frock coat. Uncle too was in high spirits and far from being offended by the brother's and sister's laughter (it could never enter his head that they might be laughing at his way of life) he himself joined in the merriment."),
	}

	for i, test := range tests {
		// Compress the data
		fmt.Println("Beginning Test: ", i)
		compressedData, err := compress.Compress(test)
		if err != nil {
			t.Errorf("Error compressing data: %v", err)
		}

		// Decompress the data
		decompressedData, err := decompress.Decompress(compressedData)
		if err != nil {
			t.Errorf("Error decompressing data: %v", err)
		}

		// Verify the decompressed data matches the original
		if string(decompressedData) != string(test) {
			t.Errorf("Test %d failed: Decompressed data = %q, want %q", i, string(decompressedData), string(test))
		} else {
			t.Logf("Test %d passed", i)
		}
	}
}
