package main

import (
	"fmt"
)

type graph struct {
	nodes [][]int
	edges [][2]int
}

func dist(g graph, s, t int) int {
	v := make([]bool, len(g.nodes))
	q := make([]int, 0, len(g.nodes))
	d := make([]int, len(g.nodes))
	v[s] = true
	q = append(q, s)
	for len(q) > 0 {
		n := q[0]
		if len(q) > 1 {
			q = q[1:]
		} else {
			q = nil
		}
		for _, a := range g.nodes[n] {
			if !v[a] {
				v[a] = true
				d[a] = d[n] + 1
				q = append(q, a)
			}
		}
	}
	return d[t]
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
	fmt.Println(dist(g, 0, 0))
	fmt.Println(dist(g, 0, 1))
	fmt.Println(dist(g, 0, 6))
	fmt.Println(dist(g, 6, 3))
}