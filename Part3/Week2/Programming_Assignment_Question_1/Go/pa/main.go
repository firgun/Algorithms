package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

type sorter struct {
	edges   []edge
	indices []int
}

func (s sorter) Less(i, j int) bool {
	a := s.edges[s.indices[i]]
	b := s.edges[s.indices[j]]
	r := a.c <= b.c
	return r
}

func (s sorter) Swap(i, j int) {
	s.indices[i], s.indices[j] = s.indices[j], s.indices[i]
}

func (s sorter) Len() int {
	return len(s.indices)
}

func solve(g graph, k int) int {
	if len(g.nodes) < k {
		panic("we need at least k nodes in the graph")
	}
	leaders := make([]int, len(g.nodes))
	for i := range leaders {
		leaders[i] = i
	}
	ccs := make([][]int, len(g.nodes))
	for i := range ccs {
		ccs[i] = []int{i}
	}
	indices := make([]int, len(g.edges))
	for i := range indices {
		indices[i] = i
	}
	sort.Sort(sorter{g.edges, indices})
	count := len(g.nodes)
	spacing := -1
	for _, index := range indices {
		e := g.edges[index]
		if leaders[e.a] != leaders[e.b] {
			if count == k {
				spacing = g.edges[index].c
				break
			}
			lesser, greater := leaders[e.a], leaders[e.b]
			if len(ccs[lesser]) > len(ccs[greater]) {
				lesser, greater = greater, lesser
			}
			for _, index := range ccs[lesser] {
				leaders[index] = greater
				ccs[greater] = append(ccs[greater], index)
			}
			ccs[lesser] = nil
			count--
		}
	}
	return spacing
}

func main() {
	algoPath, ok := os.LookupEnv("ALGO_PATH")
	if !ok {
		fmt.Println("error: ALGO_PATH environment variable not set, cannot located generated tests, exiting.")
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Println("incorrect usage: supply [<filename>|-g]")
		os.Exit(1)
	}

	if os.Args[1] != "-g" {
		p := path.Join(algoPath, os.Args[1])
		g, err := readGraph(p)
		if err != nil {
			fmt.Printf("error: failed to read input from %q, error: %v\n", p, err)
		}
		fmt.Println(solve(g, 4))
		os.Exit(0)
	}

	basePath := path.Join(algoPath, "Tests/Part3/Week2/Question_1/Generated")
	dirEntries, err := os.ReadDir(basePath)
	if err != nil {
		fmt.Printf("error: failed to read directory at path: %q, error: %v\n", basePath, err)
		os.Exit(1)
	}
	total, fails, maxFails := 0, 0, -1
	for _, e := range dirEntries {
		if !e.Type().IsRegular() || !strings.HasPrefix(e.Name(), "input") {
			continue
		}

		ifName, ofName := e.Name(), strings.Replace(e.Name(), "input", "output", 1)
		ifPath, ofPath := path.Join(basePath, ifName), path.Join(basePath, ofName)

		graph, err := readGraph(ifPath)
		if err != nil {
			fmt.Printf("error: failed to read input for path: %q, error: %v\n", ifPath, err)
			os.Exit(1)
		}

		oArray, err := readIntArray(ofPath)
		if err != nil {
			fmt.Printf("error: failed to read output for path: %q, error: %v\n", ifPath, err)
			os.Exit(1)
		}

		if len(oArray) != 1 {
			fmt.Println("error: output array must have exactly 2 elements")
			os.Exit(1)
		}

		exp := oArray[0]
		act := solve(graph, 4)

		total++
		msg := "pass"
		if act != exp {
			fails++
			if maxFails != -1 && fails >= maxFails {
				break
			}
			msg = "FAIL"
		}

		fmt.Printf("want %5d    got %5d    %s %s\n", exp, act, msg, strings.Replace(ifName, ".txt", "", -1))
	}

	if fails > 0 {
		fmt.Printf("failed, %d/%d test cases\n", fails, total)
	} else {
		fmt.Println("success, passed all test cases!")
	}
}

//////////////////////////////////////////////////////////////////////////////
// Graph loading and representation
//

type node []int

type edge struct {
	a, b, c int
}

type graph struct {
	nodes []node
	edges []edge
}

func debugDumpGraph(g graph) {
	fmt.Printf("Nodes (%d)\n", len(g.nodes))
	for i, n := range g.nodes {
		fmt.Printf("%d: ", i+1)
		for j, eIndex := range n {
			e := g.edges[eIndex]
			fmt.Printf("(%d, %d, %d)", e.a+1, e.b+1, e.c)
			if j < len(n)-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println()
	}

	fmt.Println()

	fmt.Printf("Edges (%d)\n", len(g.edges))
	for _, e := range g.edges {
		fmt.Printf("(%d, %d, %d)\n", e.a+1, e.b+1, e.c)
	}
}

func readGraph(path string) (graph, error) {
	f, err := os.Open(path)
	if err != nil {
		return graph{}, fmt.Errorf("readIntArray: failed: %v", err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return graph{}, fmt.Errorf("readIntArray: failed: %v", err)
	}

	lines := strings.Split(string(b), "\n")
	edges := make([]edge, 0, len(lines))

	maxNode := -1

	if len(lines) > 1 {
		for lineNo, line := range lines[1:] {
			line = strings.Trim(line, " \n\r\t")
			if line == "" {
				continue
			}
			tokens := strings.Split(line, " ")
			for i := range tokens {
				tokens[i] = strings.Trim(tokens[i], " \n\r\t")
			}
			if len(tokens) != 3 {
				return graph{}, fmt.Errorf("readGraph: invalid line (no. %d), expected 3 tokens, found: %d", lineNo, len(tokens))
			}
			a, err := strconv.Atoi(tokens[0])
			if err != nil {
				return graph{}, fmt.Errorf("readGraph: on line %d, failed to parse first endpoint token error: %v", lineNo, err)
			}

			b, err := strconv.Atoi(tokens[1])
			if err != nil {
				return graph{}, fmt.Errorf("readGraph: on line %d, failed to parse second endpoint token, error: %v", lineNo, err)
			}

			c, err := strconv.Atoi(tokens[2])
			if err != nil {
				return graph{}, fmt.Errorf("readGraph: on line %d, failed to parse cost token, error: %v", lineNo, err)
			}

			e := edge{a - 1, b - 1, c}
			if e.a > e.b {
				e.a, e.b = e.b, e.a
			}

			if e.b > maxNode {
				maxNode = e.b
			}

			edges = append(edges, e)
		}
	}

	nodes := make([]node, maxNode+1)

	for i, e := range edges {
		nodes[e.a] = append(nodes[e.a], i)
		nodes[e.b] = append(nodes[e.b], i)
	}

	g := graph{nodes, edges}

	// DEBUGGING
	// debugDumpGraph(g)

	return g, nil
}

//////////////////////////////////////////////////////////////////////////////
// I/O util
//

func readIntArray(path string) ([]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("readIntArray: failed: %v", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("readIntArray: failed: %v", err)
	}
	tokens := strings.Split(string(b), "\n")
	arr := make([]int, 0, len(tokens))
	for _, rt := range tokens {
		t := strings.Trim(rt, " \n\t")
		if t == "" {
			continue
		}
		n, err := strconv.Atoi(t)
		if err != nil {
			return nil, fmt.Errorf("readIntArray: failed to parse input line: %v", err)
		}
		arr = append(arr, n)
	}
	return arr, nil
}
