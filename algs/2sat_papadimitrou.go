package main

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

// TODO
// [x] refresh yourself on the algorithm using the lecture slides
// [-] design small simple test cases
// [x] implement input parsing
// [x] implement the algorithm
// [ ] debug and test on all generated test cases

type clause struct {
	a, b int
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func sign(n int) int {
	if n < 0 {
		return -1
	}
	return 1
}

func papadimitrou(c []clause) bool {
	rand.Seed(time.Now().UnixNano())
	a := make([]bool, len(c)) // same amount of variables and clauses
	// repeat log2(n) times
	for i := 0; i < int(math.Log2(float64(len(c)))); i++ {
		// compute random assignment
		for i := range a {
			if rand.Float64() > 0.5 {
				a[i] = true
			} else {
				a[i] = false
			}
		}
		// repeat 2n^2 times
		for j := 0; j < 2*len(c)*len(c); j++ {
			found := false
			for _, cl := range c {
				vara, varb := abs(cl.a)-1, abs(cl.b)-1
				if (sign(cl.a) == 1) != a[vara] && (sign(cl.b) == 1) != a[varb] {
					// randomly flip the value of one of the variables
					if rand.Float64() > 0.5 {
						a[vara] = !a[vara]
					} else {
						a[varb] = !a[varb]
					}
					found = true
				}
			}
			if !found {
				return true
			}
		}
	}

	return false
}

func main() {
	dir := "tests/2sat"

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var total, fails int
	for _, dirent := range dirEntries {
		ifile, err := os.Open(path.Join(dir, dirent.Name()))

		if err != nil {
			panic(err)
		}
		defer ifile.Close()

		ofile, err := os.Open(path.Join(dir, strings.Replace(dirent.Name(), "input", "output", 1)))
		if err != nil {
			panic(err)
		}
		defer ofile.Close()

		in, want := readin(ifile), readout(ofile)
		got := papadimitrou(in)
		if want != got {
			fmt.Printf("fail: %q\n", dirent.Name())
			fails++
		} else {
			fmt.Printf("pass: %q\n", dirent.Name())
		}
	}

	fmt.Printf("ran %d tests, failed %d, passed %d (%f)\n", total, fails, total-fails, float64(total-fails)/float64(total))
}

func readin(r io.Reader) []clause {
	var n int
	if _, err := fmt.Fscanf(r, "%d\n", &n); err != nil {
		panic(err)
	}
	c := make([]clause, n)
	for i := range c {
		var a, b int
		if _, err := fmt.Fscanf(r, "%d %d\n", &a, &b); err != nil {
			panic(err)
		}
		c[i] = clause{a, b}
	}
	return c
}

func readout(r io.Reader) bool {
	var n int
	fmt.Fscanf(r, "%d", &n)
	return n == 1
}
