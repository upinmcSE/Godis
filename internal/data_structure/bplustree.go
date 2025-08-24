package data_structure

type Item struct {
	Score  float64
	Member string
}

func (i *Item) CompareTo(other *Item) int {
	if i.Score < other.Score {
		return -1
	}
	if i.Score > other.Score {
		return 1
	}
	// Scores are equal, use the member string for tie-breaking
	if i.Member < other.Member {
		return -1
	}
	if i.Member > other.Member {
		return 1
	}
	return 0
}

// Node of B+ Tree
type Node struct {
	Items    []*Item // A list of key-value pairs, max M-1 items
	Children []*Node // Pointers to child nodes, max M children
	IsLeaf   bool    // True if it's a leaf node
	Parent   *Node   // Pointer to the parent node
	Next     *Node   // For leaf nodes, a pointer to the next leaf in the sequence
}

type BPlusTree struct {
	Root   *Node
	Degree int // The maximum number of children a node can have, = M
}

func NewBPlusTree(degree int) *BPlusTree {
	return &BPlusTree{
		Root:   &Node{IsLeaf: true},
		Degree: degree,
	}
}

// O(N) => O(1) with HashMap
func (t *BPlusTree) Score(member string) (float64, bool) {
	return 0, false
}

func (t *BPlusTree) Add(score float64, member string) int {
	return 0
}

func (t *BPlusTree) splitNode(node *Node) {}

func (t *BPlusTree) splitLeaf(node *Node) {}

func (t *BPlusTree) splitInternal(node *Node) {}

func (t *BPlusTree) splitRoot() {}

func (t *BPlusTree) GetRank(member string) int {
	return 0
}
