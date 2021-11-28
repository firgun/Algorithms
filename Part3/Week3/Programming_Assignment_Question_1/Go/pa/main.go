package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"io/ioutil"
)

type node struct {
	left, right *node
	index       int
	weight      int
}

func findMinNode(nodes []*node) int {
	min := 0
	for index, node := range nodes {
		if node.weight < nodes[min].weight {
			min = index
		}
	}
	return min
}

func removeNode(nodes []*node, index int) []*node {
	nodes[index], nodes[len(nodes)-1] = nodes[len(nodes)-1], nodes[index]
	return nodes[:len(nodes)-1]
}

func iterLeaves(n *node, depth int, fn func(*node, int)) {
	if n == nil {
		return
	} else if n.left == nil && n.right == nil {
		fn(n, depth)
	} else {
		iterLeaves(n.left , depth+1, fn)
		iterLeaves(n.right, depth+1, fn)
	}
}

// solve, given an array of integer frequencies computes the Huffman codes returning the
// max and min length codeword.
func solve(arr []int) (int, int) {
	nodes := make([]*node, len(arr))
	for i := range nodes {
		nodes[i] = &node{nil, nil, i, arr[i]}
	}
	for len(nodes) > 1 {
		min1 := findMinNode(nodes)
		n1 := nodes[min1]
		nodes = removeNode(nodes, min1)
		
		min2 := findMinNode(nodes)
		n2 := nodes[min2]
		nodes = removeNode(nodes, min2)
		
		nodes = append(nodes, &node{n1, n2, -1, n1.weight + n2.weight})
	}

	min, max := -1, -1
	iterLeaves(nodes[0], 0, func(_ *node, depth int) {
		if min == -1 || depth < min { 
			min = depth
		}
		if max == -1 || depth > max {
			max = depth
		}
	})
	
	return max, min
}

func main() {
	algoPath, ok := os.LookupEnv("ALGO_PATH")
	if !ok {
		fmt.Println("error: ALGO_PATH environment variable not set, cannot locate generated tests, exiting.")
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Println("incorrect usage: supply [<filename>|-g]")
		os.Exit(1)
	}

	if os.Args[1] != "-g" {
		p := path.Join(algoPath, os.Args[1])
		g, err := readInput(p)
		if err != nil {
			fmt.Printf("error: failed to read input from %q, error: %v\n", p, err)
			os.Exit(1)
		}
		fmt.Println(solve(g))
		os.Exit(0)
	}

	basePath := path.Join(algoPath, "Tests/Part3/Week3/Question_1/Generated")
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

		arr, err := readInput(ifPath)
		if err != nil {
			fmt.Printf("error: failed to read input for path: %q, error: %v\n", ifPath, err)
			os.Exit(1)
		}

		expMax, expMin, err := readOutput(ofPath)
		if err != nil {
			fmt.Printf("error: failed to read output for path: %q, error: %v\n", ifPath, err)
			os.Exit(1)
		}

		actMax, actMin := solve(arr)

		total++
		msg := "pass"
		if actMin != expMin || actMax != expMax {
			fails++
			if maxFails != -1 && fails >= maxFails {
				break
			}
			msg = "FAIL"
		}

		fmt.Printf("want %2d, %2d    got %2d, %2d    %s %s\n", expMin, expMax, actMin, actMax, msg, strings.Replace(ifName, ".txt", "", -1))
	}

	if fails > 0 {
		fmt.Printf("failed, %d/%d test cases\n", fails, total)
	} else {
		fmt.Println("success, passed all test cases!")
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
// I/O
//

func readInput(path string) ([]int, error) {
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
	return arr[1:], nil
}

func readOutput(path string) (int, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return -1, -1, fmt.Errorf("readOutput: failed: %v", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return -1, -1, fmt.Errorf("readOutput: failed: %v", err)
	}
	toks := strings.Split(string(b), "\n")
	if len(toks) < 2 {
		return -1, -1, fmt.Errorf("readOutput: want 2, got %d tokens", len(toks))
	}
	min, err := strconv.Atoi(toks[0])
	if err != nil {
		return -1, -1, fmt.Errorf("readOutput: failed, error: %v", err)
	}
	max, err := strconv.Atoi(toks[1])
	if err != nil {
		return -1, -1, fmt.Errorf("readOutput: failed, error: %v", err)
	}
	return min, max, nil
}






