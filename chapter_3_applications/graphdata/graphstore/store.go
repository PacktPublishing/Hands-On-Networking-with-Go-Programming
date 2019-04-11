package graphstore

import "errors"

// NewStore creates a new store.
func NewStore() *Store {
	return &Store{
		Nodes: make(map[string]*Node),
	}
}

// Store of a graph.
type Store struct {
	Nodes map[string]*Node
}

// Add a node.
func (s *Store) Add(n Node) {
	s.Nodes[n.ID] = &n
}

var errChildNodeNotFound = errors.New("child node not found")
var errParentNodeNotFound = errors.New("parent node not found")

// AddEdge between parent and child.
func (s *Store) AddEdge(parent, child string) (err error) {
	c, cok := s.Nodes[child]
	if !cok {
		err = errChildNodeNotFound
		return
	}
	p, pok := s.Nodes[parent]
	if !pok {
		err = errParentNodeNotFound
		return
	}
	p.Children = append(p.Children, child)
	c.Parent = &parent
	return
}

// Get a node.
func (s *Store) Get(id string) (n *Node) {
	if nn, ok := s.Nodes[id]; ok {
		n = nn
	}
	return
}

// Roots of the graph.
func (s *Store) Roots() []string {
	op := []string{}
	for k, n := range s.Nodes {
		if n.Parent == nil {
			op = append(op, k)
		}
	}
	return op
}
