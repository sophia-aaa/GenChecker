

import "fmt"

type Node struct {
	children []*Node
	isEnd bool
}

type Trie struct {
	root *Node
}

func InitTrie() *Trie {
	result := &Trie {root: &Node{}}
	return result
}

// Insert


// Search

func main