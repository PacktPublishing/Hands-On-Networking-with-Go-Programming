package graphstore

// Node within a graph.
type Node struct {
	// ID of the node.
	ID string `json:"id"`
	// Parent node.
	Parent *string `json:"parent"`
	// Child nodes.
	Children []string `json:"children"`
}

// NewNode creates a new Node.
func NewNode(id string, parent *string) Node {
	return Node{
		ID:     id,
		Parent: parent,
	}
}
