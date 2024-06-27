package compress

import (
	"huffman-compression/internal/huffman"
	"reflect"
	"testing"
)

func TestGetFrequencies(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []int
	}{
		{
			name:  "empty slice",
			input: []byte(""),
			want:  make([]int, 256),
		},
		{
			name:  "single character",
			input: []byte("a"),
			want:  func() []int { freq := make([]int, 256); freq['a'] = 1; return freq }(),
		},
		{
			name:  "multiple characters",
			input: []byte("banana"),
			want:  func() []int { freq := make([]int, 256); freq['b'] = 1; freq['a'] = 3; freq['n'] = 2; return freq }(),
		},
		{
			name:  "characters with ASCII values",
			input: []byte{0, 1, 255},
			want:  func() []int { freq := make([]int, 256); freq[0] = 1; freq[1] = 1; freq[255] = 1; return freq }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := huffman.GetFrequencies(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFrequencies() = %v, want %v", got, tt.want)
			}
		})
	}
}
