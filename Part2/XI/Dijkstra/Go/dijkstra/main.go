package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type edge struct {
	a, b, w int
}

type graph struct {
	nodes [][]int
	edges []edge
}

func getFullPath(rel string) string {
	val, ok := os.LookupEnv("ALGO_PATH")
	if !ok {
		panic("ALGO_PATH not set")
	}
	fullPath := path.Join(val, rel)
	// fmt.Println(fullPath)
	return fullPath
}

func readGraph(path string) (graph, error) {
	fullPath := getFullPath(path)
	f, err := os.Open(fullPath)
	if err != nil {
		return graph{}, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	edges := make([]edge, 0)
	edgeSet := make(map[edge]bool)
	largest := -1
	for s.Scan() {
		toks := strings.Split(strings.Trim(s.Text(), " \t\n\r"), "\t")
		if len(toks) < 1 {
			continue
		}
		t := strings.Trim(toks[0], " \n\t\r")
		if t == "" && len(toks) == 1 {
			continue
		}
		label, err := strconv.Atoi(t)
		if err != nil {
			return graph{}, err
		}
		for _, t := range toks[1:] {
			tok := strings.Trim(t, " \n\t\r")
			if tok == "" {
				continue
			}
			parts := strings.Split(tok, ",")
			aStr := parts[0]
			wStr := ""
			if len(parts) > 1 {
				wStr = parts[1]
			}
			a, err := strconv.Atoi(aStr)
			if err != nil {
				return graph{}, err
			}
			w, err := strconv.Atoi(wStr)
			if err != nil {
				return graph{}, err
			}
			min, max := label-1, a-1
			if min > max {
				min, max = max, min
			}
			e := edge{min, max, w}
			if !edgeSet[e] {
				edgeSet[e] = true
				edges = append(edges, e)
			}
			if max > largest {
				largest = max
			}
		}
	}
	if largest == -1 {
		panic("no nodes?")
	}
	nodes := make([][]int, largest+1)
	for idx, e := range edges {
		nodes[e.a] = append(nodes[e.a], idx)
		nodes[e.b] = append(nodes[e.b], idx)
	}
	// for _, theNode := range nodes {
	// 	for idx, edgeIdx := range theNode {
	// 		fmt.Print(edges[edgeIdx])
	// 		if idx < len(nodes[0])-1 {
	// 			fmt.Print(", ")
	// 		} else {
	// 			fmt.Println()
	// 		}
	// 	}
	// }
	if s.Err() != nil {
		return graph{}, s.Err()
	}
	return graph{nodes, edges}, nil
}

func dijkstraShortestPath(g graph) {
	x := make([]bool, len(g.nodes))
	a := make([]int, len(g.nodes))
	for idx := range a {
		a[idx] = 1000000
	}
	x[0] = true
	a[0] = 0
	for {
		minScore, minNode := -1, -1
		for _, edge := range g.edges {
			v := edge.a
			w := edge.b
			if x[v] == x[w] {
				continue
			}
			if x[w] {
				v, w = w, v
			}
			// v is the node in X, w is the node not in X
			score := a[v] + edge.w
			if minScore == -1 || score < minScore {
				minScore = score
				minNode = w
			}
		}
		if minNode == -1 {
			break
		}
		x[minNode] = true
		a[minNode] = minScore
	}
	nodes := []int{7, 37, 59, 82, 99, 115, 133, 165, 188, 197}
	for idx, n := range nodes {
		fmt.Print(a[n-1])
		if idx < len(nodes)-1 {
			fmt.Print(",")
		}
	}
	fmt.Println()
}

func main() {
	g, err := readGraph("Tests/XI/DijkstraShortestPath.txt")
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	dijkstraShortestPath(g)
}
