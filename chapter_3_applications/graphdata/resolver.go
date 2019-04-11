package graphdata

import (
	"context"

	"github.com/PacktPublishing/Hands-On-Networking-with-Go-Programming/chapter_3_applications/graphdata/graphstore"
)

type Resolver struct {
	Store *graphstore.Store
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Node() NodeResolver {
	return &nodeResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateNode(ctx context.Context, node NewNode) (string, error) {
	r.Store.Add(graphstore.NewNode(node.ID, node.Parent))
	return node.ID, nil
}

func (r *mutationResolver) CreateEdge(ctx context.Context, edge NewEdge) (string, error) {
	err := r.Store.AddEdge(edge.Parent, edge.Child)
	return edge.Parent, err
}

type nodeResolver struct{ *Resolver }

func (r *nodeResolver) Parent(ctx context.Context, obj *graphstore.Node) (parent *graphstore.Node, err error) {
	if obj.Parent != nil {
		parent = r.Store.Get(*obj.Parent)
	}
	return
}
func (r *nodeResolver) Children(ctx context.Context, obj *graphstore.Node) (children []graphstore.Node, err error) {
	children = make([]graphstore.Node, len(obj.Children))
	for i, c := range obj.Children {
		if n := r.Store.Get(c); n != nil {
			children[i] = *n
		}
	}
	return
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) List(ctx context.Context, skip *int, take *int) (op []graphstore.Node, err error) {
	for _, id := range r.Store.Roots() {
		if n := r.Store.Get(id); n != nil {
			op = append(op, *n)
		}
	}
	return
}
func (r *queryResolver) Get(ctx context.Context, id string) (n *graphstore.Node, err error) {
	n = r.Store.Get(id)
	return
}
