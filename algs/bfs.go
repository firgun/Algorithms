package main

import (
	"fmt"
)

type graph [][]int

func bfs(g graph, i int) []int {
	v := make([]bool, len(g))
	q := make([]int, 0, len(g))
	r := make([]int, 0, len(g))
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
		for _, a := range g[n] {
			if !v[a] {
				v[a] = true
				q = append(q, a)
			}
		}
	}
	return r
}

func main() {
	g := graph([][]int{
			{ 1, 2 },
			{ 0, 3, 4 },
			{ 0, 5, 6 },
			{ 1 },
			{ 1 },
			{ 2 },
			{ 2 },
	})
	fmt.Println(bfs(g, 0))
}