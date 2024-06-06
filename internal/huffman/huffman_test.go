package huffman

import (
	"reflect"
	"testing"
)

func TestBuildFrequencyTable(t *testing.T) {
	tests := []struct {
		text     []byte
		expected map[rune]int
	}{
		{[]byte("test"), map[rune]int{'t': 2, 'e': 1, 's': 1}},
		{[]byte("hello"), map[rune]int{'h': 1, 'e': 1, 'l': 2, 'o': 1}},
	}

	for _, test := range tests {
		result := BuildFrequencyTable(test.text)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("BuildFrequencyTable(%q) = %v, want %v", test.text, result, test.expected)
		}
	}
}

func TestBuildHuffmanTree(t *testing.T) {
	// This test assumes the existence of a function to verify the structure of the Huffman tree
	// Since the exact structure can vary but still be correct, this is a simplified example
	freqTable := map[rune]int{'a': 3, 'b': 2, 'c': 1}
	root := BuildHuffmanTree(freqTable)

	// Verify the tree has the correct total frequency
	expectedFreq := 6
	totalFreq := root.Freq
	if totalFreq != expectedFreq {
		t.Errorf("BuildHuffmanTree total frequency = %d, want %d", totalFreq, expectedFreq)
	}
}

func TestDecodeText(t *testing.T) {
	// This test will encode a string, decode it, and verify the result
	tests := []string{
		"test",
		"hello",
	}

	for _, test := range tests {
		freqTable := BuildFrequencyTable([]byte(test))
		codes := GenerateHuffmanCodes(BuildHuffmanTree(freqTable))
		encodedText := EncodeText([]byte(test), codes)
		decodedText := DecodeText([]byte(encodedText), codes)
		if !reflect.DeepEqual([]byte(test), decodedText) {
			t.Errorf("DecodeText(EncodeText(%q)) = %q, want %q", test, string(decodedText), test)
		}
	}
}
