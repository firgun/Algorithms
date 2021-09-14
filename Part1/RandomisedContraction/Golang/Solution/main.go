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
		nd := node{len(g.nodes)+1, make([]int, 0, len(ss))}
		for i, s := range ss {
			if s == "" {
				continue
			}
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("invalid token: %v", err)
			}
			e := edge{i, n}
			g.edges = append(g.edges, e)
			nd.edges = append(nd.edges, len(g.edges)-1)
		}
		g.nodes = append(g.nodes, nd)
	}
	return g, nil
}

func main() {
	tests := [...]string{"complete", "islands"}

	for _, t := range tests {
		g, err := load(fmt.Sprintf("../../Tests/%s.input", t))
		if err != nil {
			fmt.Printf("failed to load graph, error: %v\n", err)
			return
		}
		fmt.Println(g)
	}
}
