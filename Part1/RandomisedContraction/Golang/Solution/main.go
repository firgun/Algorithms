package main

//
// todo
//
// If we per chance pick an edge that happens to be a self loop, then we get
// rid of it and never have to deal with it again.
//
// What's an upper bound on the number of self-loops that can form throughout
// the algorithm's execution?
//
// If the number of self-loops we have to deal with is linear over the whole
// execution of the algorithm then we only do linear work in addition to our
// already linear work and the running time of a single trial stays linear.
//
// there obviously can't be more than m self-loops, because all self-loops must
// come from an edge in the initial state of the graph.
//
// the number of self-loops is actually m - k? Isn't it? Because all edges that
// don't end up as self-loops at the end of the execution end up as crossing
// edges.
//
// Well, that's linear in the number of edges, and we have to read each edge
// while reading the input.
//
// So it wouldn't take more time than reading the input -- asymptotically
// speaking.
//
// so we don't need to delete a self-loop until we randomly sample it.
//
// we can have a shuffled array of m edges, then pick at random, if the edge
// picked at random is a self-loop, then we discard it and try again. We'll end
// up discarding at most m - k edges over the whole execution which is
// obviously only linear in the number of edges.

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
	Nodes []*Node
	Edges []*Edge
}

type Node struct {
	Edges [][]*Edge
	Label int
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
	fmt.Printf("contracting %d <- %d\n", g.Nodes[e[0]].Label+1, g.Nodes[e[1]].Label+1)
	for _, edgeList := range g.Nodes[e[1]].Edges {
		g.Nodes[e[0]].Edges = append(g.Nodes[e[0]].Edges, edgeList)
	}
	g.Nodes[e[1]].Edges = nil
	g.Nodes[e[1]] = g.Nodes[e[0]]

	g.dump()
	fmt.Println()
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
		for g.Nodes[g.Edges[edgeIndex][0]].Label == g.Nodes[g.Edges[edgeIndex][1]].Label {
			next++
			edgeIndex = ids[next]
		}
		g.contract(g.Edges[edgeIndex])
	}
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
	fmt.Println(g.minCut())
	fmt.Println(g)
}
