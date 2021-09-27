package main

import (
	"fmt"
)

type graph [][]int

func topoSort(g graph) []int {
	v := make([]bool, len(g))
	f := make([]int , len(g))
	label := len(g)
	var helper func(int)
	helper = func(n int) {
		for _, a := range g[n] {
			if !v[a] {
				v[a] = true
				helper(a)
			}
		}
		f[n] = label
		label--
	}
	for n, _ := range g {
		if !v[n] {
			helper(n)
		}
	}
	return f
}

func main() {
	g := graph([][]int{
		{1, 3},
		{2},
		{4},
		{1, 4, 6},
		{5},
		{},
		{},
	})
	f := topoSort(g)
	for i, o := range f {
		fmt.Printf("%d: %d\n", i, o)
	}
}