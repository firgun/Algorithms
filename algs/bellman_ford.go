package main

import (
	"fmt"
	"math"
	"os"
	"io"
	"strings"
	"path"
	"strconv"
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

	var i int
	for i = 1; i < len(a); i++ {
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

	return a[i]
}

func shortestPaths(g Graph) [][]float64 {
	a := make([][]float64, len(g))
	for i := range g {
		a[i] = bellmanFord(g, i)
		if a[i] == nil {
			// negative cycle
			return nil
		}
	}
	return a
}

func shortestShortestPath(g Graph) (float64, bool) {
	a := shortestPaths(g)
	
	if a == nil {
		return 0, false
	}
	
	min := math.Inf(1)
	for i := range a {
		for j := range a[i] {
			if a[i][j] < min {
				min = a[i][j]
			}
		}
	}

	if min == math.Inf(1) {
		return 0, false
	} 

	return min, true
}

func main() {
	dirEntries, err := os.ReadDir(os.Args[1])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	failed := 0
	for _, entry := range dirEntries {
		if !strings.HasPrefix(entry.Name(), "input") {
			continue
		}

		iFile, err := os.Open(path.Join(os.Args[1], entry.Name()))
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		oFile, err := os.Open(path.Join(os.Args[1], strings.Replace(entry.Name(), "input", "output", 1) ))
		if err != nil { 
			fmt.Printf("error: %v\n", err)
		}

		graph := readInput(iFile)
	
		var expected string
		_, err = fmt.Fscanf(oFile, "%s", &expected)
		if err != nil {
			panic(err)
		}

		actual, ok := shortestShortestPath(graph)
		if !ok && expected == "NULL" {
			fmt.Printf("passed %q want NULL got NULL\n", entry.Name())
			continue
		}

		expectedN, err := strconv.ParseFloat(expected, 64)
		if err != nil { 
			panic(err) 
		}

		if expectedN != actual { 
			fmt.Printf("failed %q  want %.2f  got %.2f", entry.Name(), expectedN, actual)
			failed++
		} else {
			fmt.Printf("passed %q want %.2f got %.2f", entry.Name(), expectedN, actual)
		}
		fmt.Println()
	}
	fmt.Printf("failed: %d\n", failed)
	return
}

func readInput(r io.Reader) Graph {
	var numVerts, numEdges int
	_, err := fmt.Fscanf(r, "%d %d\n", &numVerts, &numEdges)
	if err != nil { panic(err) }
	g := make([]Node, numVerts)
	for i := 0; i < numEdges; i++ {
		var tail, head int
		var cost float64
		_, err := fmt.Fscanf(r, "%d %d %f\n", &tail, &head, &cost)
		if err != nil { panic(err) }
		g[head-1] = append(g[head-1], Endpoint{tail-1, cost})
	}
	return g
}