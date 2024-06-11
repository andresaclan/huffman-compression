package bitwriter

import (
	"encoding/binary"
	"testing"
)

func TestWriteBit(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int
		input    []byte
		expected uint16
	}{
		{
			name:     "single byte",
			bytes:    1,
			input:    []byte{1, 0, 1, 0, 1, 0, 1, 0},
			expected: uint16(170),
		},
		{
			name:     "single byte with padding",
			bytes:    1,
			input:    []byte{1, 0, 0, 0, 0, 0, 0, 0},
			expected: uint16(128),
		},
		{
			name:     "two bytes",
			bytes:    2,
			input:    []byte{1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1},
			expected: uint16(43605),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bw, _ := NewBitWriter()
			for _, bit := range test.input {
				if err := bw.WriteBit(bit); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
			if err := bw.Close(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if test.bytes == 1 {
				if got := uint16(bw.Bytes.Bytes()[0]); got != test.expected {
					t.Errorf("expected %d, got %d", test.expected, got)
				}
			} else {
				if got := binary.BigEndian.Uint16(bw.Bytes.Bytes()); got != test.expected {
					t.Errorf("expected %d, got %d", test.expected, got)
				}
			}
		})
	}
}
