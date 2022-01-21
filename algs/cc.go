package main

import (
	"fmt"
)

type graph [][]int

func components(g graph) int {
	v := make([]bool, len(g))
	q := make([]int, 0, len(g))
	count := 0
	for i, _ := range g {
		if !v[i] {
			q = q[:0]
			count++
			v[i] = true
			q = append(q, i)
			for len(q) > 0 {
				n := q[0]
				if len(q) > 1 {
					q = q[1:]
				} else {
					q = q[:0]
				}
				for _, a := range g[n] {
					if !v[a] {
						v[a] = true
						q = append(q, a)
					}
				}
			}
		}
	}
	return count
}

func main() {
	g := graph([][]int{
			{ 1, /* 2 */ },
			{ 0, 3, 4 },
			{ /* 0, */ 5, 6 },
			{ 1 },
			{ 1 },
			{ 2 },
			{ 2 },
	})
	fmt.Println(components(g))
}