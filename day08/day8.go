
package main

import (
	"fmt"
	"os"
)

type Node struct {
	children []*Node
	metadata []int
}

// Try out an iterator pattern with a closure
func NewSymbolIterator(symbols []int) (func() (value int, valid bool)) {
	pos := 0
	return func() (int, bool) {
		if pos >= len(symbols) {
			return 0, false
		}
		ret := symbols[pos]
		pos += 1
		return ret, true
	}
}

func NewNode() *Node {
	var node Node
	node.children = make([]*Node, 0)
	node.metadata = make([]int, 0)
	return &node
}


func ReadSymbols(filepath string) []int {
	symbols := make([]int, 0)
	
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	for {
		var sym int
		_, err := fmt.Fscan(f, &sym)
		if err != nil {
			break
		}
		symbols = append(symbols, sym)
		
	}
	return symbols
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func ReadNode(node *Node, nextSym func()(int, bool)) {
	childCount, valid := nextSym()
	if !valid {
		panic("Out of symbols early")
	}
	metadataCount, valid := nextSym()
	if !valid {
		panic("Out of symbols early")
	}
	node.children = make([]*Node, childCount)
	node.metadata = make([]int, metadataCount)

	for i := 0; i < childCount; i += 1 {
		node.children[i] = NewNode()
		ReadNode(node.children[i], nextSym)
	}

	for i := 0; i < metadataCount; i += 1 {
		node.metadata[i], valid = nextSym()
		if !valid {
			panic("Out of symbols early")
		}
	}
}

func SumAllMetadata(node *Node) int {
	sum := 0
	for _, childNode := range node.children {
		sum += SumAllMetadata(childNode)
	}
	for _, metadata := range node.metadata {
		sum += metadata
	}
	return sum
}

func GetNodePart2Value(node *Node) int {
	sum := 0
	// If a node has no children, its value is the sum of its metadata
	if len(node.children) == 0 {
		for _, m := range node.metadata {
			sum += m
		}
		return sum
	}
	// Otherwise, the metadata are indices to child nodes to count towards value
	for _, m := range node.metadata {
		if m == 0 {
			continue
		}
		m -= 1
		if m >= len(node.children) {
			continue
		}
		sum += GetNodePart2Value(node.children[m])
	}
	return sum
}

func main() {
	symbols := ReadSymbols("day8_input.txt")
	fmt.Printf("Read %d symbols\n", len(symbols))

	rootNode := NewNode()
	iter_next := NewSymbolIterator(symbols)
	ReadNode(rootNode, iter_next)

	// PART 1
	metadataSum := SumAllMetadata(rootNode)
	fmt.Printf("Part 1\n------\n")
	fmt.Printf("Sum of all metadata: %d\n", metadataSum)
	// PART 2	
	part2Value := GetNodePart2Value(rootNode)
	fmt.Printf("Part 2\n------\n")
	fmt.Printf("Root node value: %d\n", part2Value)
}