package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

//
// TODO
// ===
// [_] Implement solve2, using heaps to store edges, so repeated minimum computations only take O(log(n)) time
// [_] Implement solve3, using heaps to store vertices, same asymptotic running time, but practically faster
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

func solve(g graph) int {
	x := make([]bool, len(g.nodes))
	x[0] = true
	t := make([]int, 0)
	for c := 1; c < len(g.nodes); c++ {
		s := -1
		for i, e := range g.edges {
			if (s == -1 || e.c < g.edges[s].c) && x[e.a] != x[e.b] {
				s = i
			}
		}
		t = append(t, s)
		e := g.edges[s]
		x[e.a] = true
		x[e.b] = true
	}
	sum := 0
	for _, eIndex := range t {
		sum += g.edges[eIndex].c
	}
	return sum
}

func solveFast(g graph) int {
	var h heap
	for i, e := range g.edges {
		h.insert(e.c, i)
	}

	for {
		if entry, ok := h.extractMin(); ok {
			i := entry.val.(int)
			k := entry.key
			e := g.edges[i]
			fmt.Printf("(%d, %d, %d) key: %d\n", e.a+1, e.b+1, e.c, k)
		} else {
			break
		}
	}

	panic("foo")

	x := make([]bool, len(g.nodes))
	x[0] = true
	t := make([]int, 0)
	for c := 1; c < len(g.nodes); c++ {
		s := -1
		ret := make([]int, 0)
		for {
			if entry, ok := h.extractMin(); ok {
				s = entry.val.(int)
				e := g.edges[s]
				if x[e.a] != x[e.b] {
					break
				} else {
					ret = append(ret, s)
				}
			} else {
				panic("heap empty")
			}
		}
		for _, r := range ret {
			e := g.edges[r]
			h.insert(e.c, r)
		}
		if s == -1 {
			panic("no good crossing edge found")
		}
		t = append(t, s)
		e := g.edges[s]
		fmt.Printf("(%d, %d, %d)\n", e.a+1, e.b+1, e.c)
		x[e.a] = true
		x[e.b] = true
	}
	sum := 0
	for _, eIndex := range t {
		sum += g.edges[eIndex].c
	}
	fmt.Println(sum)
	os.Exit(1)
	return sum
}

func solveFastest(g graph) int {
	return 0
}

func main() {
	algoPath, ok := os.LookupEnv("ALGO_PATH")
	if !ok {
		fmt.Println("error: ALGO_PATH environment variable not set, cannot located generated tests, exiting.")
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Println("incorrect usage: 2sum [<filename>|-g]")
		os.Exit(1)
	}

	if !strings.HasPrefix(os.Args[1], "-g") {
		iPath := path.Join(algoPath, os.Args[1])
		iGraph, err := readGraph(iPath)
		if err != nil {
			fmt.Printf("error: failed to read input from %q, error: %v\n", iPath, err)
		}
		fmt.Println(solve(iGraph))
		os.Exit(0)
	}

	var solveFn func(graph) int
	fnMode := strings.ToLower(strings.TrimPrefix(os.Args[1], "-g"))
	if fnMode == "" {
		solveFn = solve
	} else {
		switch fnMode {
		case "normal":
			solveFn = solve
		case "fast":
			solveFn = solveFast
		case "fastest":
			solveFn = solveFastest
		default:
			fmt.Printf("error: invalid function mode %s, must be normal, fast or fastest\n", fnMode)
			os.Exit(1)
		}
	}

	basePath := path.Join(algoPath, "Tests/Part3/Week1/Question_3")
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

		fmt.Println(ifName)

		iGraph, err := readGraph(ifPath)
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
			fmt.Println("error: output array must have exactly 1 elements")
			os.Exit(1)
		}

		exp := oArray[0]
		act := solveFn(iGraph)

		total++

		if act != exp {
			fails++
			if maxFails != -1 && fails >= maxFails {
				break
			}
			fmt.Printf("FAILED (question 1) test case at path: %q, expected: %d, actual: %d\n", ifName, exp, act)
		} else {
			fmt.Printf("PASSED (question 1) test case at path: %q\n", ifName)
		}
	}
	if fails > 0 {
		fmt.Printf("failed, %d/%d test cases\n", fails, total)
	} else {
		fmt.Println("success, passed all test cases!")
	}
}

// Heap data structure

type heapEntry struct {
	key int
	val interface{}
}

type heap []struct {
	key int
	val interface{}
}

func (h *heap) insert(key int, value interface{}) {
	*h = append(*h, heapEntry{key, value})
	h.siftUp(len(*h) - 1)
}

func (h *heap) extractMin() (heapEntry, bool) {
	if len(*h) == 0 {
		return heapEntry{}, false
	}
	r := (*h)[0]
	(*h)[0], (*h)[len(*h)-1] = (*h)[len(*h)-1], (*h)[0]
	*h = (*h)[:len(*h)-1]
	h.siftDown(0)
	return r, true
}

func (h *heap) siftUp(i int) {
	for ; i != 0 && (*h)[i].key < (*h)[i/2].key; i = i / 2 {
		(*h)[i], (*h)[i/2] = (*h)[i/2], (*h)[i]
	}
}

func (h *heap) siftDown(i int) {
	for ; i < len(*h) && 2*i < len(*h); i++ {
		l, r := 2*i, 2*i+1
		if (*h)[i].key > (*h)[l].key || (r < len(*h) && (*h)[i].key > (*h)[r].key) {
			var c int
			if r >= len(*h) || (*h)[l].key < (*h)[r].key {
				c = l
			} else {
				c = r
			}
			(*h)[i], (*h)[c] = (*h)[c], (*h)[i]
			i = c
		} else {
			break
		}
	}
}

func (h *heap) delete(i int) {
	if i < 0 || i >= len(*h) {
		panic(fmt.Errorf("index out of range: %d", i))
	}
	(*h)[i], (*h)[len(*h)-1] = (*h)[len(*h)-1], (*h)[i]
	*h = (*h)[:len(*h)-1]
	h.siftUp(i)
	h.siftDown(i)
}

func (h *heap) dump() {
	t := 1 // 2 to power 0
	i := 0
	for _, e := range *h {
		fmt.Print(e.key)
		i++
		if i == t {
			fmt.Println()
			i = 0
			t *= 2
		} else {
			fmt.Print(" ")
		}
	}
}

func (h *heap) peek() int {
	return (*h)[0].val.(int)
}
