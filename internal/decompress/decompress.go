package decompress

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"huffman-compression/internal/huffman"

	"github.com/icza/bitio"
)

func Decompress(data []byte) ([]byte, error) {
	// create a buffer to read the compressed data
	buf := bytes.NewReader(data)

	// read the length of the encoded table from the data
	var tableLength uint32
	err := binary.Read(buf, binary.LittleEndian, &tableLength)
	if err != nil {
		return nil, err
	}

	// read the number of bits of data that were written
	var bitsWritten uint32
	err = binary.Read(buf, binary.LittleEndian, &bitsWritten)
	if err != nil {
		return nil, err
	}
	fmt.Println("Bits to decode during decompression: ", bitsWritten)

	// byte slice to store the encoded table (not yet decoded)
	encodedTable := make([]byte, tableLength)
	_, err = buf.Read(encodedTable)
	if err != nil {
		return nil, err
	}

	// decode the huffman codes
	codes, err := huffman.DecodeHuffmanCodes(encodedTable)
	if err != nil {
		return nil, err
	}

	// isolate the remaining data that was written
	data = data[8+tableLength:]
	// at this point data only contains the huffman encoded data

	// read each bit of the data and store it as a byte in a byte slice
	encodedData := make([]byte, 0)
	r := bitio.NewCountReader(bytes.NewBuffer(data))
	for {
		bit, err := r.ReadBool()
		if err != nil {
			break
		}
		if bit {
			encodedData = append(encodedData, '1')
		} else {
			encodedData = append(encodedData, '0')
		}
	}

	fmt.Println("Bits actually read during decompression:", r.BitsCount)
	// takes care of padded 0s at the end of the encoded data
	// if int(r.BitsCount) < int(bitsWritten) {
	// 	for i := 0; i < int(bitsWritten)-int(r.BitsCount); i++ {
	// 		encodedData = append(encodedData, '0')
	// 	}
	// }
	// decode the huffman encoded data
	decodedText := huffman.DecodeText(encodedData, codes)
	return decodedText, nil

}
