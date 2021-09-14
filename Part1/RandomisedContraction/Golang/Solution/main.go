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

/*
type graph struct {
	nodes map[int]node
	edges map[int]edge
}
*/

type graph struct {
	nodes []node
	edges []edge
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
	g := &graph{}
	for s.Scan() {
		if s.Err() != nil {
			return nil, fmt.Errorf("error occurred while scanning: %v", err)
		}
		ss := strings.Split(s.Text(), " ")
		nd := node{len(g.nodes), make([]int, 0, len(ss))}
		g.nodes = append(g.nodes, nd)
		for _, s := range ss {
			if s == "" {
				continue
			}
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("invalid token: %v", err)
			}
			g.edges = append(g.edges, edge{nd.id, n - 1})
		}
	}
	for i, ed := range g.edges {
		na := &g.nodes[ed.a]
		nb := &g.nodes[ed.b]
		na.edges = append(na.edges, i)
		nb.edges = append(nb.edges, i)
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
	/*
		e := &g.edges[edgeIdx]
		for _, bEdgeIdx := range e.b {
			bEdge := g.edges[bEdgeIdx]
			if bEdge.a == e.b {
				bEdge.a = e.a
			} else {
				bEdge.b = e.a
			}
		}
	*/
}

func main() {

}
