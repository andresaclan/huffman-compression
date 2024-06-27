package huffman

import (
	"bytes"
	"container/heap"
	"fmt"
	"io"

	"github.com/icza/bitio"
)

type HuffmanNode struct {
	Index       uint8
	Freq        int
	HuffmanCode int
	HuffmanLen  int
	Left        *HuffmanNode
	Right       *HuffmanNode
}

type Tree struct {
	Root *HuffmanNode
}

type HuffmanTree struct {
	Htree          *Tree
	huffmanCodes   [256]int
	huffmanLengths [256]int
}

func NewHuffmanTreeFromFrequencies(frequencies [256]int) *HuffmanTree {
	htree := buildHuffmanTree(frequencies)
	h := &HuffmanTree{
		Htree: htree,
	}
	h.CreateCodes()
	return h
}

// Encodes the data to the writer
func (h *HuffmanTree) Encode(data []byte, w *bitio.CountWriter) {
	for _, b := range data {
		fmt.Println("Encoding byte:", b, "of length", h.huffmanLengths[b], "Code:", toBinaryString(h.huffmanCodes[b], h.huffmanLengths[b]))
		w.WriteBits(uint64(h.huffmanCodes[b]), uint8(h.huffmanLengths[b]))
	}

}

// Decodes the data from the reader
// returns the decoded data as a byte slice (each byte represents a character)
func (h *HuffmanTree) Decode(r *bitio.CountReader) []byte {
	output := make([]byte, 0)
	for {
		index, err := h.getCode(r)
		if err == io.EOF {
			break
		}
		output = append(output, index)
	}
	return output
}

func (h *HuffmanTree) getCode(r *bitio.CountReader) (uint8, error) {
	curr := h.Htree.Root
	for !isLeaf(curr) {
		bit, err := r.ReadBool()
		if err != nil {
			return curr.Index, err
		}
		if bit {
			curr = curr.Right
		} else {
			curr = curr.Left
		}
	}
	return curr.Index, nil
}

func buildHuffmanTree(frequencies [256]int) *Tree {
	pq := make(PriorityQueue, len(frequencies))
	for i := 0; i < len(frequencies); i++ {
		pq[i] = &HuffmanNode{
			Index: uint8(i),
			Freq:  frequencies[i],
		}
	}
	heap.Init(&pq)

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*HuffmanNode)
		right := heap.Pop(&pq).(*HuffmanNode)
		merged := join(left, right)
		heap.Push(&pq, merged)
	}
	root := heap.Pop(&pq).(*HuffmanNode)
	return &Tree{Root: root}
}

// Encodes the HuffmanTree to the writer
func EncodeHuffmanTree(root *HuffmanNode, buf *bytes.Buffer) *bitio.CountWriter {
	if buf == nil {
		panic("buf is nil")
	}
	w := bitio.NewCountWriter(buf)
	encodeHuffmanTreeHelper(root, w)
	return w
}
func encodeHuffmanTreeHelper(node *HuffmanNode, w *bitio.CountWriter) {
	if w == nil {
		panic("CountWriter is nil")
	}
	if isLeaf(node) {
		w.WriteBool(true)
		w.WriteBits(uint64(node.Index), 8)
	} else {
		w.WriteBool(false)
		encodeHuffmanTreeHelper(node.Left, w)
		encodeHuffmanTreeHelper(node.Right, w)
	}
}

// Creates HuffmanTree from the encoded tree
func GetHuffmanTree(r *bitio.CountReader) *HuffmanNode {
	if b, err := r.ReadBool(); b {
		if err != nil {
			panic(err)
		}
		value, err := r.ReadByte()
		if err != nil {
			panic(err)
		}
		return &HuffmanNode{
			Index: uint8(value),
		}
	} else {
		leftChild := GetHuffmanTree(r)
		rightChild := GetHuffmanTree(r)
		return &HuffmanNode{
			Left:  leftChild,
			Right: rightChild,
		}
	}
}

func (h *HuffmanTree) CreateCodes() {
	if h.Htree.Root == nil {
		return
	}
	h.createCodesRecursive(h.Htree.Root)
}

func (h *HuffmanTree) createCodesRecursive(node *HuffmanNode) {
	if !isLeaf(node) {
		node.Left.HuffmanCode = node.HuffmanCode << 1
		node.Left.HuffmanLen = node.HuffmanLen + 1
		h.createCodesRecursive(node.Left)

		node.Right.HuffmanCode = (node.HuffmanCode << 1) | 1
		node.Right.HuffmanLen = node.HuffmanLen + 1
		h.createCodesRecursive(node.Right)
	} else {
		h.huffmanCodes[node.Index] = node.HuffmanCode
		h.huffmanLengths[node.Index] = node.HuffmanLen
	}
}

func GetFrequencies(data []byte) [256]int {
	var freqTable [256]int
	for i := 0; i < len(data); i++ {
		b := uint8(data[i])
		freqTable[b]++
	}
	return freqTable

}
func (h *HuffmanTree) DisplayCodes() {
	for i := 0; i < 256; i++ {
		if h.huffmanLengths[i] > 0 {
			fmt.Println(i, ":", toBinaryString(h.huffmanCodes[i], h.huffmanLengths[i]))
		}
	}
}

func toBinaryString(value int, length int) string {
	retval := ""
	for i := length; i >= 0; i-- {
		if (value >> i & 1) == 1 {
			retval += "1"
		} else {
			retval += "0"
		}
	}
	return retval
}
func join(a, b *HuffmanNode) *HuffmanNode {
	return &HuffmanNode{
		Freq:  a.Freq + b.Freq,
		Left:  a,
		Right: b,
	}
}

func isLeaf(node *HuffmanNode) bool {
	return node.Left == nil && node.Right == nil
}

type PriorityQueue []*HuffmanNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Freq < pq[j].Freq
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	huffmanNode := x.(*HuffmanNode)
	*pq = append(*pq, huffmanNode)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	huffmanNode := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return huffmanNode
}

// func main() {
// 	// testing Huffman Algorithm
// 	// text := "baacccc"
// 	// fmt.Println("original text: " + text)
// 	// freqTable := BuildFrequencyTable(text)
// 	// root := BuildHuffmanTree(freqTable)
// 	// codes := GenerateHuffmanCodes(root)
// 	// fmt.Println("codes:")
// 	// for char, code := range codes {
// 	// 	fmt.Printf("%c:%s\n", char, code)
// 	// }
// 	// encodedText := EncodeText(text, codes)
// 	// fmt.Println("encodedText: " + encodedText)
// 	// testing End

// 	// testing the EncodeHuffmanCodes() and DecodeHuffmanCodes()
// 	// codes := map[rune]string{
// 	// 	'a': "00",
// 	// 	'b': "01",
// 	// 	'c': "10",
// 	// }

// 	// encodedData, err := EncodeHuffmanCodes(codes)
// 	// if err != nil {
// 	// 	fmt.Println("Error encoding Huffman codes:", err)
// 	// 	return
// 	// }

// 	// fmt.Println("Encoded data:", encodedData)

// 	// decodedData, err := DecodeHuffmanCodes(encodedData)
// 	// if err != nil {
// 	// 	fmt.Println("Error decoding Huffman codes:", err)
// 	// 	return
// 	// }

// 	// for key, value := range decodedData {
// 	// 	fmt.Printf("%c:%s\n", key, value)
// 	// }
// 	// testing End

// 	// testing WriteFile
// 	filePath := "./output.txt"
// 	err := file.WriteFile(filePath, "test")
// 	if err != nil {
// 		fmt.Println("Error writing file:", err)
// 	}
// 	// testing End
// }
