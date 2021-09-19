package main

// todo
// ---
// for each endpoint of an edge, store a pointer to a pointer to a node
//
// for each node, store an array of array of edges
//
// a supernode is just a node with the top-level array having more than one
// element
//
// if you want to contract a supernode and a node, you just add the one element
// in the node to the array in the supernode O(1)
//
// if you want to contract two regular ndoes, just pretend that the one with
// more incident edges is the supernode, and do as above O(1)
//
// if you want to contract two supernodes together, just concatenate the
// smaller to the larger, so the larger supernode becomes the new supernode
// O(n) time
//
// make a shuffled array of edges up front, then for each contraction pick pop
// an edge of the array.
//
// the assignments format specifies a single edge "from the point of view" of
// both endpoints
//
// handle parallel edges in graph loading

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Graph struct {
	Nodes []*Node
	Edges []*Edge
}

type Node struct {
	Edges [][]*Edge
	Label int
}

type Edge [2]int

func (e Edge) String() string {
	return fmt.Sprintf("{%d, %d}", e[0]+1, e[1]+1)
}

func NewEdge(a, b int) Edge {
	e := Edge([2]int{a, b})
	if e[0] > e[1] {
		e[0], e[1] = e[1], e[0]
	}
	return e
}

func readGraph(path string) (*Graph, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read file at %q: %v", path, err)
	}
	g := Graph{make([]*Node, 0), make([]*Edge, 0)}
	for _, ln := range strings.Split(string(bytes), "\n") {
		fields := strings.Fields(ln)
		if len(fields) < 1 {
			continue
		}
		indices := make([]int, len(fields))
		for i, field := range fields {
			n, err := strconv.Atoi(field)
			if err != nil {
				return nil, fmt.Errorf("cannot convert field %q to int: %v", field, err)
			}
			indices[i] = n
		}
		v := indices[0]
		edges := make([]*Edge, len(indices)-1)
		edgeCounts := make(map[Edge]int)
		if len(indices) > 1 {
			for i, adj := range indices[1:] {
				e := NewEdge(v-1, adj-1)
				ep := &e
				edges[i] = ep
				if edgeCounts[*ep] == 0 {
					g.Edges = append(g.Edges, ep)
				}
				edgeCounts[*ep]++
			}
		}
		n := &Node{make([][]*Edge, 1), len(g.Nodes)}
		n.Edges[0] = edges
		g.Nodes = append(g.Nodes, n)
	}
	return &g, nil
}

func (g *Graph) dump() {
	for vi, v := range g.Nodes {
		if vi != v.Label {
			fmt.Println("deleted")
			continue
		}
		fmt.Printf("%d ", v.Label+1)
		if v == nil {
			continue
		}
		if v.Edges != nil {
			for listIndex, edgeList := range v.Edges {
				for eIndex, e := range edgeList {
					endpoint := 0
					if g.Nodes[e[0]].Label == v.Label {
						endpoint = 1
					}
					fmt.Printf("%d", g.Nodes[e[endpoint]].Label+1)
					if eIndex < len(v.Edges[listIndex])-1 {
						fmt.Print(" ")
					}
				}
				if listIndex < len(v.Edges)-1 {
					fmt.Print(" ")
				}
			}
		}
		fmt.Println()
	}
}

func (g *Graph) contract(e *Edge) {
	// fmt.Println(e)
	for _, edgeList := range g.Nodes[e[1]].Edges {
		g.Nodes[e[0]].Edges = append(g.Nodes[e[0]].Edges, edgeList)
	}
	g.Nodes[e[1]].Edges = nil
	g.Nodes[e[1]] = g.Nodes[e[0]]
}

func (g *Graph) minCut() int {
	// *old_ptr_ptr = *last_ptr_ptr
	// pop(node_array)

	/*
		for nodesLeft := len(g.Nodes); nodesLeft > 2; nodesLeft-- {

		}
	*/

	g.contract(g.Edges[1])

	return -1000
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("error: bad args, specify path to graph input file")
		os.Exit(1)
	}

	g, err := readGraph(os.Args[1])
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	g.dump()
	g.contract(g.Edges[0])
	fmt.Println()
	g.dump()
}
