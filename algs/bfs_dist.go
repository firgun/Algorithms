package main

import (
	"fmt"
)

type graph [][]int

func dist(g graph, s, t int) int {
	v := make([]bool, len(g))
	q := make([]int, 0, len(g))
	d := make([]int, len(g))
	v[s] = true
	q = append(q, s)
	for len(q) > 0 {
		n := q[0]
		if len(q) > 1 {
			q = q[1:]
		} else {
			q = nil
		}
		for _, a := range g[n] {
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
	g := graph([][]int{
			{ 1, 2 },
			{ 0, 3, 4 },
			{ 0, 5, 6 },
			{ 1 },
			{ 1 },
			{ 2 },
			{ 2 },
	})
	fmt.Println(dist(g, 0, 0))
	fmt.Println(dist(g, 0, 1))
	fmt.Println(dist(g, 0, 6))
	fmt.Println(dist(g, 6, 3))
}