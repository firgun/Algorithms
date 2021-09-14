package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
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
		nd := node{len(g.nodes), make([]int, len(ss))}
		for i, s := range ss {
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
	g, err := load("test.graph")
	if err != nil {
		fmt.Printf("failed to load graph, error: %v\n", err)
		return
	}
	fmt.Println(g)
}
