package main

import "fmt"

type graph [][]int

func dfs(g graph, i int) []int {
	v := make([]bool, len(g))
	r := make([]int, 0, len(g))
	v[i] = true
	var dfsHelper func(int)
	dfsHelper = func(n int) {
		r = append(r, n)
		for _, a := range g[n] {
			if !v[a] {
				v[a] = true
				dfsHelper(a)
			}
		}
	}
	dfsHelper(i)
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
	fmt.Println(dfs(g, 0))
}