package bitreader

import (
	"bytes"
	"io"
)

type BitReader struct {
	Bytes *bytes.Buffer
	buf   byte
	n     uint8
}

func NewBitReader(readFrom *bytes.Buffer) (*BitReader, error) {
	return &BitReader{Bytes: readFrom}, nil
}

func (br *BitReader) ReadBit() (byte, error) {
	if br.n == 0 {
		buf := make([]byte, 1)
		if _, err := br.Bytes.Read(buf); err == io.EOF {
			return 0, err
		}
		br.buf = buf[0]
		br.n = 8
	}
	// 7 - br.n
	bit := (br.buf >> (br.n - 1)) & 1
	br.n--
	return bit, nil
}

func (br *BitReader) Close() error {
	return nil
}
