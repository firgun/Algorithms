package main

import (
	"fmt"
)

type graph struct {
	nodes [][]int
	edges [][2]int
}

func bfs(g graph, i int) []int {
	v := make([]bool, len(g.nodes))
	q := make([]int, 0, len(g.nodes))
	r := make([]int, 0, len(g.nodes))
	v[i] = true
	q = append(q, i)
	for len(q) > 0 {
		n := q[0]
		if len(q) > 1 {
			q = q[1:]
		} else {
			q = nil
		}
		r = append(r, n)
		for _, a := range g.nodes[n] {
			if !v[a] {
				v[a] = true
				q = append(q, a)
			}
		}
	}
	return r
}

func main() {
	g := graph{
		nodes: [][]int{
			{ 1, 2 },
			{ 0, 3, 4 },
			{ 0, 5, 6 },
			{ 1 },
			{ 1 },
			{ 2 },
			{ 2 },
		},
		edges: [][2]int{
			{0, 1}, 
			{0, 2}, 
			{1, 3}, 
			{1, 4}, 
			{2, 5}, 
			{2, 6},
		},
	}
	fmt.Println(bfs(g, 0))
}