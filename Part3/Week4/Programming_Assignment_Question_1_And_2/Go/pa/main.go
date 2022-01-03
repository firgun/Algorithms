package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"path"
)

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// Ideas to improve time and space
//
// Idea 1:
//
// As suggested in the problem definition, use a recursive approach with some kind of 
// memoization. This'll ensure that we don't unnecessarily solve any subproblems who's 
// solutions aren't actually used by bigger subproblems.
//
// The recursive solution takes a top-down approach and subproblems are solved lazily.
//
// We'll use memoization to ensure that a sub-problem is only solved once.
//
// * How will we solve the base cases? What are they?
//   - capacity empty
//   - prefix of items empty
//   - solution with these arguments already solved
// * What are the non-base cases branches?
//   - We need to compare the weight of the last item, to the capacity left to determine if this item can fit in
//     and be used itself or if we must resort to the other subproblem's solution.
// * How will we cache results to avoid redundant computation?
//   - hash table

type Task struct {
	cap, pref int
}

func knapsack3(capacity int, weights, values []int) int {
	prev, curr := make([]int, capacity), make([]int, capacity)

	for i := 0; i < len(values); i++ {
		for c := 1; c <= capacity; c++ {
			if i == 0 {
				if  weights[i] <= c {
					curr[c-1] = values[i]
				}
				continue
			}

			s1 := prev[c-1]

			if weights[i] <= c {
				s2 := values[i]
				if weights[i] != c {
					s2 += prev[c-weights[i]-1]
				}
				if s2 > s1 {
					curr[c-1] = s2
					continue
				}
			}

			curr[c-1] = s1
		}
		prev, curr = curr, prev
	}

	return prev[len(prev)-1]
}

func knapsack2(capacity int, weights, values []int) int {
	m := make(map[Task]int)
	s := _knapsack2(capacity, weights, values, m)
	return s
}

func _knapsack2(capacity int, weights, values []int, m map[Task]int) int {
	// FIXME: not sure why this is so slow ... it's seems to be working on MORE recursive subproblems than 
	// len(items) * capacity ... which doesn't really make sense to me. But I could be measuring 
	// incorrectly and generally it seems to be under and as input size increases it gets closer to being 
	// equal?

	if len(weights) != len(values) {
		panic("invalid arguments: len(weights) != len(values)")
	}
	
	if capacity < 0 {
		panic("invalid arguments: capacity < 0")
	}

	if len(values) == 0 || capacity < 0 {
		return 0
	}

	if s, ok := m[Task{capacity, len(values)}]; ok {
		return s
	}

	vn := values [len(values )-1]
	wn := weights[len(weights)-1]

	s1 := _knapsack2(capacity, weights[:len(weights)-1], values[:len(values)-1], m)
	s := s1
	if wn <= capacity {
		s2 := _knapsack2(capacity - wn, weights[:len(weights)-1], values[:len(values)-1], m) + vn
		if s2 > s {
			s = s2
		}
	}

	m[Task{capacity, len(values)}] = s
	return s
}

func knapsack(capacity int, weights, values []int) int {
	s := make([][]int, len(values))
	for i := range s {
		s[i] = make([]int, capacity)
	}

	for i := range s {
		for c := 1; c <= capacity; c++ {
			if i == 0 {
				if  weights[i] <= c {
					s[i][c-1] = values[i]
				}
				continue
			}

			s1 := s[i-1][c-1]

			if weights[i] <= c {
				s2 := values[i]
				if weights[i] != c {
					s2 += s[i-1][c-weights[i]-1]
				}
				if s2 > s1 {
					s[i][c-1] = s2
					continue
				}
			}

			s[i][c-1] = s1
		}
	}

	return s[len(s)-1][capacity-1]
}

func main() {
	if len(os.Args) == 2 {
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
			capacity, weights, values := readInput(iFile)
			expected := readOutput(oFile)
			actual := knapsack3(capacity, weights, values)
			if expected != actual { 
				fmt.Printf("failed %q  expected %d  got %d", entry.Name(), expected, actual)
				failed++
			} else {
				fmt.Printf("passed %q", entry.Name())
			}
			fmt.Println()
		}
		fmt.Printf("total: %d, failed: %d\n", len(dirEntries), failed)
		return
	}

	capacity, weights, values := readInput(os.Stdin)
	fmt.Println(knapsack(capacity, weights, values))
}

func readInput(r io.Reader) (capacity int, weights []int, values []int) {
	_, _ = fmt.Fscanf(r, "%d", &capacity)
	var numItems int
	_, _ = fmt.Fscanf(r, "%d\n", &numItems)
	values = make([]int, numItems)
	weights = make([]int, numItems)
	for i := 0; i < numItems; i++ {
		_, _ = fmt.Fscanf(r, "%d"  , &values [i])
		_, _ = fmt.Fscanf(r, "%d\n", &weights[i])
	}
	return
}

func readOutput(r io.Reader) int {
	var n int
	_, _  = fmt.Fscanf(r, "%d", &n)
	return n
}