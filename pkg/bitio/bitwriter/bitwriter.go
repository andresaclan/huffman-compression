package bitwriter

import (
	"bytes"
	"fmt"
)

type BitWriter struct {
	Bytes *bytes.Buffer
	buf   byte
	n     uint8
}

func NewBitWriter() (*BitWriter, error) {
	return &BitWriter{Bytes: bytes.NewBuffer(make([]byte, 0))}, nil
}

// WriteBit writes a single bit to the buffer
func (bw *BitWriter) WriteBit(bit byte) error {
	bw.buf |= (bit & 1) << (7 - bw.n)
	bw.n++
	if bw.n == 8 {
		if err := bw.Flush(); err != nil {
			return err
		}
	}
	return nil
}

// Flush writes the buffer to the file if it contains any bits
func (bw *BitWriter) Flush() error {
	if bw.n > 0 {
		if _, err := bw.Bytes.Write([]byte{bw.buf}); err != nil {
			fmt.Println("Error writing bit buffer to bytes buffer:", err)
			return err
		}
		bw.buf = 0
		bw.n = 0
	}
	return nil
}

// Closes the file and flushes the buffer
func (bw *BitWriter) Close() error {
	if err := bw.Flush(); err != nil {
		return err
	}
	return nil
}
