package main

// todo
//
// we need an arbitrarily long trail of pointers to the supernode that a node has previously been contracted into
//

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Graph struct {
	Nodes []Node
	Edges []Edge
}

type Node struct {
	Edges  []int
	Label  int
	Parent int
}

type Edge [2]int

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
	g := Graph{Nodes: make([]Node, 0), Edges: make([]Edge, 0)}
	edgeToIndex := make(map[Edge]int)
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
		edges := make([]int, len(indices)-1)
		if len(indices) > 1 {
			for i, adj := range indices[1:] {
				e := NewEdge(v-1, adj-1)
				ei, ok := edgeToIndex[e]
				if !ok {
					ei = len(g.Edges)
					g.Edges = append(g.Edges, e)
					edgeToIndex[e] = ei 
				}
				edges[i] = ei
			}
		}
		n := Node{Edges: edges, Label: len(g.Nodes), Parent: -1}
		g.Nodes = append(g.Nodes, n)
	}
	return &g, nil
}

func (g *Graph) copy() *Graph {
	c := &Graph{make([]Node, len(g.Nodes)), make([]Edge, len(g.Edges))}
	for i, e := range g.Edges {
		c.Edges[i] = e
	}
	for i, n := range g.Nodes {
		edges := make([]int, len(n.Edges))
		copy(edges, n.Edges)
		c.Nodes[i] = Node{edges, n.Label, n.Parent}
	}
	return c
}

func (g *Graph) dump() {
	for _, n := range g.Nodes {
		if n.Parent != -1 {
			fmt.Println("deleted")
			continue
		}
		fmt.Print(n.Label+1)
		fmt.Print(" ")
		for j, ei := range n.Edges {
			e := g.Edges[ei]
			end := 0
			if g.Nodes[e[end]].Label == n.Label {
				end = 1
			}
			a := g.Nodes[e[end]]
			fmt.Print(a.Label+1)
			if j < len(n.Edges)-1 {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func (g *Graph) findRoot(i int) *Node {
	// visited := make(map[int]bool)
	path := []int{}
	for i != -1 {
		/*
		if visited[i] {
			panic("already visited, reference cycle?")
		}
		visited[i] = true
		*/
		n := g.Nodes[i]
		path = append(path, i)
		i = n.Parent
	}
	if len(path) == -1 {
		panic("expected at least one node in path")
	}
	root := &g.Nodes[path[len(path)-1]]
	for _, i := range path[:len(path)-1] {
		g.Nodes[i].Parent = root.Label
	}
	return root
}

func (g *Graph) contract(e Edge) {
	a, b := g.findRoot(e[0]), g.findRoot(e[1])
	b.Parent = a.Label
}

func (g *Graph) minCut() int {
	ids := make([]int, len(g.Edges))
	for i, _ := range ids {
		ids[i] = i
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})
	next := 0
	edgeIndex := ids[next]
	for nodesLeft := len(g.Nodes); nodesLeft > 2; nodesLeft-- {
		for g.findRoot(g.Edges[edgeIndex][0]).Label == g.findRoot(g.Edges[edgeIndex][1]).Label {
			next++
			edgeIndex = ids[next]
		}
		edge := g.Edges[edgeIndex]
		g.contract(edge)
	}

	// pick any root and count crossing edges, edges where the other endpoint's root is different.
	k := 0
	for _, e := range g.Edges {
		if g.findRoot(e[0]).Label != g.findRoot(e[1]).Label {
			k++
		}
	}
	return k
}

func (g *Graph) findMinCut() int {
	m := -1
	for i := 0; i < len(g.Nodes) * len(g.Nodes); i++ {
		k := g.copy().minCut()
		if k < m || m == -1 {
			m = k
		}
	}
	return m
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
	fmt.Println("read")

	fmt.Println(g.findMinCut())
}
