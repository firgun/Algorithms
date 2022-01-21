package main

import (
	"fmt"
	"math"
)

type Endpoint struct { Node int; EdgeWeight float64 }
type Node     []Endpoint
type Graph    []Node

func bellmanFord(g Graph, s int) []float64 {
	a := make([][]float64, len(g)+1)
	for i := range a {
		a[i] = make([]float64, len(g))
	}

	a[0][s] = 0
	for i, _ := range g {
		if i != s {
			a[0][i] = math.Inf(1)
		}
	}

	var stable bool
	for i := 1; i < len(a); i++ {
		stable = true
		
		for j, n := range g {
			a[i][j] = a[i-1][j]
			for _, e := range n {
				if a[i-1][e.Node] + e.EdgeWeight < a[i][j] {
					a[i][j] = a[i-1][e.Node] + e.EdgeWeight
					stable = false
				}
			}
		}

		if stable {
			break
		}
	}

	if !stable {
		return nil
	}

	return a[len(a)-1]
}

func main() {
	fmt.Println(bellmanFord(Graph{
		[]Endpoint{},
		[]Endpoint{Endpoint{0, 2}},
		[]Endpoint{Endpoint{1, 1}, Endpoint{0, 5}},
	}, 0))
}