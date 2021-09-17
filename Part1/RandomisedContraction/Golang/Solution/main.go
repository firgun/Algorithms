package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type node struct {
	id    int
	edges []int
}

type edge struct {
	a, b int
}

type graph struct {
	nodes []node
	edges map[int]edge
}

const (
	graphStringModeMatrix = iota
	graphStringModeList
)

var graphStringMode int = graphStringModeList

func (g *graph) String() string {
	var buf bytes.Buffer
	if graphStringMode == graphStringModeMatrix {
		for _, a := range g.nodes {
			for _, b := range g.nodes {
				if a.id == b.id {
					buf.WriteString("? ")
					continue
				}
				incident := false
				for _, ae := range a.edges {
					edge := g.edges[ae]
					if edge.a == b.id || edge.b == b.id {
						incident = true
						break
					}
				}
				if incident {
					buf.WriteString("% ")
				} else {
					buf.WriteString("_ ")
				}
			}
			buf.WriteString("\n")
		}
	} else if graphStringMode == graphStringModeList {
		for _, nd := range g.nodes {
			buf.WriteString(strconv.Itoa(nd.id+1) + " ")
			for _, index := range nd.edges {
				edge := g.edges[index]
				adjacent := edge.a
				if edge.a == nd.id {
					adjacent = edge.b
				}
				buf.WriteString(strconv.Itoa(adjacent+1) + " ")
			}
			buf.WriteRune('\n')
		}
	} else {
		panic(fmt.Sprintf("bad string mode: %d", graphStringMode))
	}
	return buf.String()
}

func load(path string) (*graph, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to load graph: %v\n", err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	g := &graph{nodes: make([]node, 0), edges: make(map[int]edge)}
	for s.Scan() {
		if s.Err() != nil {
			return nil, fmt.Errorf("error occurred while scanning: %v", err)
		}
		ss := strings.Split(s.Text(), " ")
		numbers := make([]int, 0, len(ss))
		for _, s := range ss {
			if s == "" {
				continue
			}
			n, err := strconv.Atoi(s)
			if err != nil {
				return nil, fmt.Errorf("invalid token: %v", err)
			}
			numbers = append(numbers, n)
		}
		if len(numbers) < 1 {
			continue
		}
		nodeId := numbers[0] - 1
		g.nodes = append(g.nodes, node{nodeId, make([]int, 0, len(numbers))})
		if len(numbers) > 1 {
			for _, n := range numbers[1:] {
				g.edges[len(g.edges)] = edge{nodeId, n - 1}
			}
		}
	}
	for i, ed := range g.edges {
		g.nodes[ed.a].edges = append(g.nodes[ed.a].edges, i)
		g.nodes[ed.b].edges = append(g.nodes[ed.b].edges, i)
	}
	return g, nil
}

func testGraphLoad() {
	tests := [...]string{
		"complete",
		"islands",
		"star",
		"cycle",
		"binary_tree_breadth_first",
		"binary_tree_depth_first",
	}
	for _, t := range tests {
		g, err := load(fmt.Sprintf("../../Tests/%s.input", t))
		if err != nil {
			fmt.Printf("failed to load graph, error: %v\n", err)
			return
		}
		fmt.Println(g)
	}
}

// contract fuses the two endpoints of a specified edge together creating a
// supernode. The resulting supernode is adjacent to all nodes either of the
// endpoints where adjacent to. Returns the supernode.
func contract(g *graph, edgeIdx int) node {
	if len(g.nodes) < 2 {
		panic("less than 2 nodes in graph")
	}

	if _, ok := g.edges[edgeIdx]; !ok {
		panic(fmt.Sprintf("no edge with index %s", edgeIdx))
	}

	e := g.edges[edgeIdx]
	na := &g.nodes[e.a]
	nb := &g.nodes[e.b]

	if na.id == nb.id {
		panic("trying to contract self-loop")
	}

	// for all edges incident on nb, swap the endpoint which points to B with
	// A
	for _, bEdgeIdx := range nb.edges {
		bEdge := g.edges[bEdgeIdx]
		if bEdge.a == e.b {
			bEdge.a = e.a
		} else {
			bEdge.b = e.a
		}
		g.edges[bEdgeIdx] = bEdge
		na.edges = append(na.edges, bEdgeIdx)
	}

	// remove all references to incident edges from b.
	nb.edges = nil

	return *na
}

func (g *graph) randomEdgeIndex() int {
	for k, _ := range g.edges {
		return k
	}
	panic("no edges!")
}

type testCasePath struct {
	inputPath  string
	outputPath string
}

type test struct {
	path testCasePath
	in   *graph
	out  int
}

var testsPath string = "../../Tests/"

func listTests() []test {
	files, err := ioutil.ReadDir(testsPath)
	if err != nil {
		panic(err)
	}
	paths := make(map[string]testCasePath)
	for _, f := range files {
		if f.Mode().IsRegular() {
			fn := f.Name()
			if strings.HasPrefix(fn, "input") {
				key := strings.Replace(fn, "input_", "", 1)
				path := paths[key]
				path.inputPath = fn
				paths[key] = path
			} else if strings.HasPrefix(fn, "output") {
				key := strings.Replace(fn, "output_", "", 1)
				path := paths[key]
				path.outputPath = fn
				paths[key] = path
			}
		}
	}
	// validate test case paths
	toDelete := make(map[string]bool)
	for k, p := range paths {
		if p.inputPath == "" {
			toDelete[k] = true
			fmt.Println("warning: test case missing input path, omitting from list")
		}
		if p.outputPath == "" {
			toDelete[k] = true
			fmt.Println("warning: test case missing output path, omitting from list")
		}
	}
	for key, _ := range toDelete {
		delete(paths, key)
	}
	pathList := make([]testCasePath, 0)
	for _, p := range paths {
		pathList = append(pathList, p)
	}
	ts := make([]test, 0)
	for _, path := range paths {
		g, err := load(testsPath + path.inputPath)
		if err != nil {
			panic(err)
		}
		bytes, err := ioutil.ReadFile(testsPath + path.outputPath)
		if err != nil {
			panic(err)
		}
		s := strings.Trim(string(bytes), "\n\t ")
		a, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		ts = append(ts, test{path, g, a})
	}
	return ts
}

func runTests() {
	for _, t := range listTests() {
		fmt.Println(t.path, t.out)
	}
}

func testContract() {
	tests := [...]string{
		"complete",
		"islands",
		"star",
		"cycle",
		"binary_tree_breadth_first",
		"binary_tree_depth_first",
		"simple",
		"simple2",
	}
	for _, t := range tests {
		fmt.Println(t)
		g, err := load(fmt.Sprintf("../../Tests/%s.input", t))
		if err != nil {
			fmt.Printf("failed to load graph, error: %v\n", err)
			return
		}
		n := len(g.nodes)
		if len(g.edges) < 1 {
			continue
		}
		rand.Seed(time.Now().UnixNano())

		fmt.Println(g)

		for i := 0; i < 10; i++ {
			for n > 2 {
				eidx := g.randomEdgeIndex()

				theEdge := g.edges[eidx]
				fmt.Printf("choose {%d, %d}\n\n", theEdge.a+1, theEdge.b+1)

				sn := contract(g, eidx)

				idxToDel := make([]int, 0)
				for _, index := range sn.edges {
					edge := g.edges[index]
					if edge.a == edge.b {
						idxToDel = append(idxToDel, index)
					}
				}

				for _, index := range idxToDel {
					delete(g.edges, index)
					del := -1
					for refIndex, ref := range sn.edges {
						if index == ref {
							del = refIndex
							break
						}
					}
					if del != -1 {
						l := len(sn.edges) - 1
						sn.edges[l], sn.edges[del] = sn.edges[del], sn.edges[l]
						sn.edges = sn.edges[:l]
					}
				}

				g.nodes[sn.id] = sn
				fmt.Println(g)
				n--
			}
		}
	}
}

func main() {
	// testGraphLoad()
	// testContract()
	runTests()
}
