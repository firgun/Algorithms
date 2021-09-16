package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type node struct {
	id    int
	edges []int
}

type edge struct {
	a, b int
}

type graph struct {
	nodes []node
	edges map[int]edge
}

func (g *graph) String() string {
	var buf bytes.Buffer
	for _, a := range g.nodes {
		for _, b := range g.nodes {
			if a.id == b.id {
				buf.WriteString("? ")
				continue
			}
			incident := false
			for _, ae := range a.edges {
				edge := g.edges[ae]
				if edge.a == b.id || edge.b == b.id {
					incident = true
					break
				}
			}
			if incident {
				buf.WriteString("% ")
			} else {
				buf.WriteString("_ ")
			}
		}
		buf.WriteString("\n")
	}
	return buf.String()
}

func load(path string) (*graph, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load graph: %v\n", err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	g := &graph{nodes: make([]node, 0), edges: make(map[int]edge)}
	for s.Scan() {
		if s.Err() != nil {
			return nil, fmt.Errorf("error occurred while scanning: %v", err)
		}
		ss := strings.Split(s.Text(), " ")
		numbers := make([]int, 0, len(ss))
		for _, s := range ss {
			if s == "" {
				continue
			}
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("invalid token: %v", err)
			}
			numbers = append(numbers, n)
		}
		if len(numbers) < 1 {
			continue
		}
		nodeId := numbers[0]-1
		g.nodes = append(g.nodes, node{nodeId, make([]int, 0, len(numbers))})
		if len(numbers) > 1 {
			for _, n := range numbers[1:] {
				g.edges[len(g.edges)] = edge{nodeId, n - 1}
			}
		}
	}
	for i, ed := range g.edges {
		g.nodes[ed.a].edges = append(g.nodes[ed.a].edges, i)
		g.nodes[ed.b].edges = append(g.nodes[ed.b].edges, i)
	}
	return g, nil
}

func testGraphLoad() {
	tests := [...]string{
		"complete",
		"islands",
		"star",
		"cycle",
		"binary_tree_breadth_first",
		"binary_tree_depth_first",
	}
	for _, t := range tests {
		g, err := load(fmt.Sprintf("../../Tests/%s.input", t))
		if err != nil {
			fmt.Printf("failed to load graph, error: %v\n", err)
			return
		}
		fmt.Println(g)
	}
}

// contract fuses the two endpoints of a specified edge together creating a
// supernode. The resulting supernode is adjacent to all nodes either of the
// endpoints where adjacent to.
func contract(g *graph, edgeIdx int) {
	// TODO: Handle some edge cases
	// ---
	//  [ ] there's not two nodes in the graph
	//  [ ] there's not edge with id edgeIdx in graph 
	//  [ ] the edge we are given to contract is itself a self-loop
	e := g.edges[edgeIdx]
	na := &g.nodes[e.a]
	nb := &g.nodes[e.b]
	for _, bEdgeIdx := range nb.edges {
		bEdge := g.edges[bEdgeIdx]
		if (bEdge.a == e.a && bEdge.b == e.b) || (bEdge.a == e.b && bEdge.b == e.a) {
			for i, eidx := range na.edges {
				if eidx == bEdgeIdx {
					last := len(na.edges)-1
					na.edges[i], na.edges[last] = na.edges[last], na.edges[i]
					na.edges = na.edges[:len(na.edges)]
				}
			}
			// we're never going to be accessing b again, so we can just forget it exists? For our purposes??
			delete(g.edges, bEdgeIdx)
			continue
		}
		if bEdge.a == e.b {
			bEdge.a = e.a
		} else {
			bEdge.b = e.a
		}
		g.edges[bEdgeIdx] = bEdge
		g.nodes[e.a].edges = append(g.nodes[e.a].edges, bEdgeIdx)
	}
	nb.edges = nil
	delete(g.edges, edgeIdx)
}

func (g *graph) randomEdgeIndex() int {
	for k, _ := range g.edges {
		return k
	}
	panic("no edges!") 
}

func testContract() {
        tests := [...]string{
                "complete",
                "islands",
                "star",
                "cycle",
                "binary_tree_breadth_first",
                "binary_tree_depth_first",
        }
        for _, t := range tests {
                g, err := load(fmt.Sprintf("../../Tests/%s.input", t))
                if err != nil {
                        fmt.Printf("failed to load graph, error: %v\n", err)
                        return
                }
		fmt.Println(g)
		for len(g.edges) > 1 {
			eidx := g.randomEdgeIndex()
			contract(g, eidx)
			fmt.Println(g)
		}
		break
        }	
}

func main() {
	// testGraphLoad()
	testContract()
}
