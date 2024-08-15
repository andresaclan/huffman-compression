package huffman

import (
	"reflect"
	"testing"
)

func TestBuildFrequencyTable(t *testing.T) {
	tests := []struct {
		text     []byte
		expected [256]int
	}{
		{[]byte("test"), [256]int{'t': 2, 'e': 1, 's': 1}},
		{[]byte("hello"), [256]int{'h': 1, 'e': 1, 'l': 2, 'o': 1}},
	}

	for _, test := range tests {
		result := GetFrequencies(test.text)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("BuildFrequencyTable(%q) = %v, want %v", test.text, result, test.expected)
		}
	}
}
func TestGetFrequencies(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  [256]int
	}{
		{
			name:  "empty slice",
			input: []byte(""),
			want:  func() [256]int { var freq [256]int; return freq }(),
		},
		{
			name:  "single character",
			input: []byte("a"),
			want:  func() [256]int { var freq [256]int; freq['a'] = 1; return freq }(),
		},
		{
			name:  "multiple characters",
			input: []byte("banana"),
			want:  func() [256]int { var freq [256]int; freq['b'] = 1; freq['a'] = 3; freq['n'] = 2; return freq }(),
		},
		{
			name:  "characters with ASCII values",
			input: []byte{0, 1, 255},
			want:  func() [256]int { var freq [256]int; freq[0] = 1; freq[1] = 1; freq[255] = 1; return freq }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFrequencies(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFrequencies() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestHuffman(t *testing.T) {
// 	testString := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

// 	freqTable := GetFrequencies([]byte(testString))

// 	root := NewHuffmanTreeFromFrequencies(freqTable)

//		root.DisplayCodes()
//	}
func TestBuildHuffmanTree(t *testing.T) {
	// This test assumes the existence of a function to verify the structure of the Huffman tree
	// Since the exact structure can vary but still be correct, this is a simplified example
	var freqTable [256]int
	freqTable['a'] = 3
	freqTable['b'] = 2
	freqTable['c'] = 1

	root := NewHuffmanTreeFromFrequencies(freqTable)
	// Verify the tree has the correct total frequency
	expectedFreq := 6
	totalFreq := root.Htree.Root.Freq
	if totalFreq != expectedFreq {
		t.Errorf("BuildHuffmanTree total frequency = %d, want %d", totalFreq, expectedFreq)
	}
}

// func TestDecodeText(t *testing.T) {
// 	// This test will encode a string, decode it, and verify the result
// 	tests := []string{
// 		"test",
// 		"hello",
// 	}

// 	for _, test := range tests {
// 		freqTable := BuildFrequencyTable([]byte(test))
// 		codes := GenerateHuffmanCodes(BuildHuffmanTree(freqTable))
// 		encodedText := EncodeText([]byte(test), codes)
// 		decodedText := DecodeText([]byte(encodedText), codes)
// 		if !reflect.DeepEqual([]byte(test), decodedText) {
// 			t.Errorf("DecodeText(EncodeText(%q)) = %q, want %q", test, string(decodedText), test)
// 		}
// 	}
// }
