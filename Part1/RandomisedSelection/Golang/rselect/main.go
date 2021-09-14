package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// s1 returns ith smallest element of x (by reduction to sorting)
func s1(x []int, i int) int {
	if len(x) < 1 {
		panic("invalid arg")
	}
	if i < 0 || i >= len(x) {
		panic("invalid arg")
	}
	y := make([]int, len(x))
	copy(y, x)
	sort.Ints(y)
	return y[i]
}

// partition splits the elements of the array such that all elements less than
// the pivot appear before the pivot in the array and all elements greater than
// the pivot appear after the pivot in the array. Returns the resting place
// of the pivot.
func partition(x []int) int {
	if len(x) == 0 {
		panic("invalid args")
	}
	if len(x) == 1 {
		return 0
	}
	j := 1
	for i, _ := range x {
		if x[i] < x[0] {
			x[i], x[j] = x[j], x[i]
			j++
		}
	}
	x[0], x[j-1] = x[j-1], x[0]
	return j-1
}

// s2 returns ith smallest element of x (by randomised selection)
func s2(x []int, i int) int {
	if len(x) < 1 { 
		panic("invalid arg")
	}
	if i < 0 || i >= len(x) {
		panic("invalid arg")
	}
	if len(x) == 1 {
		return x[0]
	}
	y := make([]int, len(x))
	copy(y, x)
	p := rand.Intn(len(y))
	y[0], y[p] = y[p], y[0]
	p = partition(y)
	if i < p {
		return s2(y[:p], i)
	} else if i == p {
		return y[p]
	} else { // i > p
		return s2(y[p+1:], i-(p+1))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		x := make([]int, i+1)
		for j, _ := range x {
			x[j] = j+1
		}
		rand.Shuffle(len(x), func(i, j int) { x[i], x[j] = x[j], x[i] })
		t := rand.Intn(len(x))
		want := s1(x, t)
		got  := s2(x, t)
		if want != got {
			m := 25
			if m > len(x) {
				m = len(x)
			}
			fmt.Printf("test case failed, want: %d, got: %d, looking for %dth order stat. in:\n%v\n", want, got, t, x[:m])
		}
	}
}
