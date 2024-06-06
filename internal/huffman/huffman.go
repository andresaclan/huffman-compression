package huffman

import (
	"bytes"
	"container/heap"
	"encoding/gob"
)

type HuffmanNode struct {
	Char  rune
	Freq  int
	Left  *HuffmanNode
	Right *HuffmanNode
}

func join(a, b *HuffmanNode) *HuffmanNode {
	return &HuffmanNode{
		Char:  0,
		Freq:  a.Freq + b.Freq,
		Left:  a,
		Right: b,
	}
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

// BuildFrequencyTable builds a frequency table from input text
func BuildFrequencyTable(text []byte) map[rune]int {
	m := make(map[rune]int)

	for _, c := range text {
		m[rune(c)]++
	}

	return m

}

func BuildHuffmanTree(freqTable map[rune]int) *HuffmanNode {
	pq := make(PriorityQueue, len(freqTable))
	i := 0
	for char, f := range freqTable {
		pq[i] = &HuffmanNode{
			Char: char,
			Freq: f,
		}
		i++
	}
	heap.Init(&pq)

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*HuffmanNode)
		right := heap.Pop(&pq).(*HuffmanNode)
		merged := join(left, right)
		heap.Push(&pq, merged)
	}
	return heap.Pop(&pq).(*HuffmanNode)
}

func GenerateHuffmanCodes(root *HuffmanNode) map[rune]string {
	codes := make(map[rune]string)
	generateHuffmanCodesHelper(root, "", codes)
	return codes
}

func generateHuffmanCodesHelper(node *HuffmanNode, code string, codes map[rune]string) {
	if node.Left == nil && node.Right == nil {
		codes[node.Char] = code
		return
	}
	if node.Left != nil {
		generateHuffmanCodesHelper(node.Left, code+"1", codes)
	}
	if node.Right != nil {
		generateHuffmanCodesHelper(node.Right, code+"0", codes)
	}
}

func EncodeText(text []byte, codes map[rune]string) []byte {
	var encodedText bytes.Buffer
	for _, char := range text {
		encodedText.WriteString(codes[rune(char)])
	}
	return encodedText.Bytes()

}

func DecodeText(text []byte, codes map[rune]string) []byte {
	decodedText := make([]byte, 0)
	lookup := make(map[string]rune)
	for char, code := range codes {
		lookup[code] = char
	}
	code := ""
	for _, bit := range text {
		code += string(bit)
		if char, ok := lookup[code]; ok {
			decodedText = append(decodedText, byte(char))
			code = ""
		}
	}
	return decodedText
}

func EncodeHuffmanCodes(codes map[rune]string) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(codes)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeHuffmanCodes(data []byte) (map[rune]string, error) {
	var codes map[rune]string
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(&codes)
	if err != nil {
		return nil, err
	}
	return codes, nil
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
