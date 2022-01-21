package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
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
	prev := 0
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
		if prev == 0 || prev != h {
			for i := prev + 1; i <= h; i++ {
				g = append(g, node{})
			}
			prev = h
		}
		g[len(g)-1].o = append(g[len(g)-1].o, t-1)
	}
	if s.Err() != nil {
		return nil, s.Err()
	}
	for ni, n := range g {
		for _, a := range n.o {
			g[a].i = append(g[a].i, ni)
		}
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

func dfs2(g graph, v []bool, l []int, n int, f int) {
	l[n] = f
	for _, a := range g[n].o {
		if !v[a] {
			v[a] = true
			dfs2(g, v, l, a, f)
		}
	}
}

func leaders(g graph, f []int) []int {
	o := make([]int, len(g))
	for i, t := range f {
		o[t-1] = i
	}
	l := make([]int, len(g))
	for i := range l {
		l[i] = -1
	}
	v := make([]bool, len(g))
	for i := len(o) - 1; i >= 0; i-- {
		n := o[i]
		if !v[n] {
			v[n] = true
			dfs2(g, v, l, n, f[n])
		}
	}
	return l
}

func sccs(g graph) []int {
	f := times(g)
	l := leaders(g, f)
	c := make([]int, len(l))
	for _, t := range l {
		c[t-1]++
	}
	sort.Ints(c)
	r := make([]int, 5)
	for i := 0; i < 5; i++ {
		if i < 5 {
			r[i] = c[len(c)-1-i]
		}
	}
	return r
}

func loadSlice(path string) ([]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	toks := strings.Split(string(b), ",")
	r := make([]int, 0, 5)
	for _, t := range toks {
		n, err := strconv.Atoi(strings.Trim(t, " \n\t\r"))
		if err != nil {
			return nil, err
		}
		r = append(r, n)
	}
	return r, nil
}

func runGeneratedTestCases() {
	resourcesPath, ok := os.LookupEnv("ALGO_PATH")
	if !ok {
		fmt.Println("error: cannot locate resource directory, exiting ...")
		os.Exit(1)
	}
	basePath := path.Join(resourcesPath, "Tests/X/Generated")
	fInfos, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Printf("error: cannot read test directory: %v\n", err)
		os.Exit(1)
	}
	testCases := make([]struct {
		In  string
		Out string
	}, len(fInfos)/2)
	for i, info := range fInfos {
		if info.Mode().IsRegular() && strings.HasPrefix(info.Name(), "input") {
			testCases[i].In = info.Name()
			testCases[i].Out = strings.Replace(info.Name(), "input", "output", 1)
		}
	}
	failCount := 0
	for _, tc := range testCases {
		g, err := loadGraph(path.Join(basePath, tc.In))
		if err != nil {
			fmt.Printf("failed to load graph: %v\n", err)
			os.Exit(1)
		}
		expected, err := loadSlice(path.Join(basePath, tc.Out))
		if err != nil {
			fmt.Printf("failed to load output: %v\n", err)
			os.Exit(1)
		}
		actual := sccs(*g)
		fmt.Print("expected: ")
		for i, n := range expected {
			fmt.Print(n)
			if i != len(expected)-1 {
				fmt.Print(",")
			}
		}
		fmt.Print("  actual: ")
		for i, n := range actual {
			fmt.Print(n)
			if i != len(actual)-1 {
				fmt.Print(",")
			}
		}
		fail := false
		if len(expected) != len(actual) {
			fail = true
		}
		if !fail {
			for idx := range expected {
				if expected[idx] != actual[idx] {
					fail = true
					break
				}
			}
		}
		if fail {
			fmt.Println("  FAIL")
			failCount++
		} else {
			fmt.Println("  PASS")
		}
	}
	fmt.Println("===\nTOTAL TEST CASES FAILED:", failCount, "\n")

}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: scc <input>.txt")
		os.Exit(1)
	}
	arg := os.Args[1]

	if strings.HasPrefix(arg, "-") {
		if arg == "-g" || arg == "--generated" {
			runGeneratedTestCases()
		} else {
			fmt.Printf("error: unknown arg %v\n", arg)
		}
		return
	}
	fileName := arg

	resourcesPath, ok := os.LookupEnv("ALGO_PATH")
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
	fmt.Println(sccs(*g))
}
