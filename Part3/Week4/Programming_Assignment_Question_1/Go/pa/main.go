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
			actual := knapsack(capacity, weights, values)
			if expected != actual { 
				fmt.Printf("failed %q  expected %d  got %d\n", entry.Name(), expected, actual)
				failed++
			} else {
				fmt.Printf("passed %q\n", entry.Name())
			}
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
		_, _ = fmt.Fscanf(r, "%d", &values[i])
		_, _ = fmt.Fscanf(r, "%d\n", &weights[i])
	}
	return
}

func readOutput(r io.Reader) int {
	var n int
	_, _  = fmt.Fscanf(r, "%d", &n)
	return n
}