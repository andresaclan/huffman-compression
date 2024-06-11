package bitreader

import (
	"bytes"
	"testing"
)

func TestReadBit(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "single byte",
			input:    []byte{0b10101010},
			expected: []byte{1, 0, 1, 0, 1, 0, 1, 0},
		},
		{
			name:     "two bytes",
			input:    []byte{0b10101010, 0b01010101},
			expected: []byte{1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			br, _ := NewBitReader(bytes.NewBuffer(test.input))
			for i, expectedBit := range test.expected {
				bit, err := br.ReadBit()
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if bit != expectedBit {
					t.Errorf("at position %d: expected %d, got %d", i, expectedBit, bit)
				}
			}
			if _, err := br.ReadBit(); err == nil {
				t.Error("expected error when no more bits to read, got nil")
			}
		})
	}
}
