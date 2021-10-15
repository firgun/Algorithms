package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type graph []node // 0 incoming, 1 outgoing

type node struct {
	o []int
	i []int
}

func loadGraph(p string) (*graph, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	g := graph{}
	prev := -1
	for s.Scan() {
		parts := strings.Split(s.Text(), " ")
		hs, ts := parts[0], parts[1]
		h, err := strconv.Atoi(hs)
		if err != nil {
			return nil, err
		}
		t, err := strconv.Atoi(ts)
		if err != nil {
			return nil, err
		}
		if prev == -1 || prev != h {
			g = append(g, node{})
			prev = h
		}
		g[len(g)-1].o = append(g[len(g)-1].o, t-1)
	}
	for ni, n := range g {
		for _, a := range n.o {
			g[a].i = append(g[a].i, ni)
		}
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	return &g, nil
}

func dfs1(g graph, n int, v []bool, f []int, t *int) {
	for _, a := range g[n].i {
		if !v[a] {
			v[a] = true
			dfs1(g, a, v, f, t)
		}
	}
	*t++
	if f[n] != 0 {
		panic("assigning finishing time to a node twice!!!")
	}
	f[n] = *t
}

func times(g graph) []int {
	t := 0
	tp := &t
	f := make([]int, len(g))
	v := make([]bool, len(g))
	for n := len(g) - 1; n >= 0; n-- {
		if !v[n] {
			v[n] = true
			dfs1(g, n, v, f, tp)
		}
	}
	return f
}

func leaders(g graph) []int {
	
}

func sccs(g graph) {
	f := times(g)
	fmt.Println(f)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: scc <input>.txt")
		os.Exit(1)
	}
	fileName := os.Args[1]
	resourcesPath, ok := os.LookupEnv("ALGO_RESOURCE_DIRECTORY")
	if !ok {
		fmt.Println("error: cannot locate resource directory, exiting ...")
		os.Exit(1)
	}
	filePath := path.Join(resourcesPath, fileName)
	g, err := loadGraph(filePath)
	if err != nil {
		fmt.Printf("error: failed to load graph: %v\n", err)
		os.Exit(1)
	}

	sccs(*g)

	/*
		for l, n := range *g {
			fmt.Print(l+1, " -> ")
			fmt.Print("i: ")
			for _, i := range n.i {
				fmt.Print(i+1, " ")
			}
			fmt.Print("o: ")
			for _, o := range n.o {
				fmt.Print(o+1, " ")
			}
			fmt.Println()
		}
	*/

}
