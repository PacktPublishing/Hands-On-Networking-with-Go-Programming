package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/machinebox/graphql"
)

func main() {
	c := graphql.NewClient("http://localhost:11111/query")
	// c.Log = func(s string) { fmt.Println(s) }
	// Create parent node.
	id, err := createNode(context.Background(), c, newNode{
		ID:     "node1",
		Parent: nil,
	})
	if err != nil {
		fmt.Printf("failed to create first node: %v\n", err)
		return
	}
	fmt.Printf("created node '%s'\n", id)
	// Create child node 1.
	node1Parent := "node1"
	id, err = createNode(context.Background(), c, newNode{
		ID:     "childNode1",
		Parent: &node1Parent,
	})
	if err != nil {
		fmt.Printf("failed to create first child node: %v\n", err)
		return
	}
	fmt.Printf("created node '%s'\n", id)
	// Create child node 2.
	id, err = createNode(context.Background(), c, newNode{
		ID:     "childNode2",
		Parent: &node1Parent,
	})
	if err != nil {
		fmt.Printf("failed to create second child node: %v\n", err)
		return
	}
	fmt.Printf("created node '%s'\n", id)
	// Create edges.
	id, err = createEdge(context.Background(), c, newEdge{
		Parent: "node1",
		Child:  "childNode1",
	})
	if err != nil {
		fmt.Printf("failed to create edge between parent and child 1: %v\n", err)
		return
	}
	id, err = createEdge(context.Background(), c, newEdge{
		Parent: "node1",
		Child:  "childNode2",
	})
	if err != nil {
		fmt.Printf("failed to create edge between parent and child 2: %v\n", err)
		return
	}
	fmt.Printf("created node '%s'\n", id)
	// List nodes.
	nr, err := listParents(context.Background(), c)
	if err != nil {
		fmt.Printf("failed to list nodes: %v\n", err)
		return
	}
	fmt.Println("retrieved nodes:")
	d, err := json.MarshalIndent(nr, "", "  ")
	fmt.Println(string(d))
	if err != nil {
		fmt.Printf("failed to marshal return value: %v\n", err)
		return
	}
}

const createNodeQuery = `mutation($node: NewNode!) {
  createNode(node: $node)
}`

type newNode struct {
	ID     string  `json:"id"`
	Parent *string `json:"parent"`
}

func createNode(ctx context.Context, client *graphql.Client, node newNode) (id string, err error) {
	req := graphql.NewRequest(createNodeQuery)
	req.Var("node", node)
	req.Header.Set("Cache-Control", "no-cache")
	var result struct {
		CreateNode string `json:"createNode"`
	}
	err = client.Run(ctx, req, &result)
	id = result.CreateNode
	return
}

const createEdgeQuery = `mutation($edge: NewEdge!) {
  createEdge(edge: $edge)
}`

type newEdge struct {
	Parent string `json:"parent"`
	Child  string `json:"child"`
}

func createEdge(ctx context.Context, client *graphql.Client, edge newEdge) (id string, err error) {
	req := graphql.NewRequest(createEdgeQuery)
	req.Var("edge", edge)
	req.Header.Set("Cache-Control", "no-cache")
	var result struct {
		CreateEdge string `json:"createEdge"`
	}
	err = client.Run(ctx, req, &result)
	id = result.CreateEdge
	return
}

const listParentsQuery = `query {
  list(take:10, skip: 0) {
    id
    parent {
      id
    }
    children {
      id
      children {
        id
      }
    }
  }
}`

type nodeResult struct {
	ID       string        `json:"id"`
	Parent   *nodeResult   `json:"parent"`
	Children []*nodeResult `json:"children"`
}

func listParents(ctx context.Context, client *graphql.Client) (result []nodeResult, err error) {
	req := graphql.NewRequest(listParentsQuery)
	req.Header.Set("Cache-Control", "no-cache")
	var nr struct {
		List []nodeResult `json:"list"`
	}
	err = client.Run(ctx, req, &nr)
	result = nr.List
	return
}
