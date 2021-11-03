// Use two heaps
//
// H[High], stores the numbers greater than or equal to the current median,
// supports extract-min H[Low ], stores the numbers less    than the current
// median, supports extract-max (use same heap just use negation trick!)
//
// When handed a new number, insert into H[High or Low] according to wether
// it's greater than or equal to / less than the current median
//
// If heaps become unbalanced, then just balance them like this
// hLow.Insert(hHigh.ExtractMin()) / hHigh.Insert(-hLow.ExtractMin() /* keep in
// mind negation trick */)
//
// When someone asks for the median, we'll use hHigh.PeakMin()
//

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
)

func readArray(path string) ([]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %v", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %v", err)
	}
	lns := strings.Split(string(b), "\n")
	arr := make([]int, 0, len(lns))
	for _, rawLn := range lns {
		ln := strings.Trim(rawLn, " \n\t\r")
		if ln == "" {
			continue
		}
		n, err := strconv.Atoi(ln)
		if err != nil {
			fmt.Printf("warning: bad input line: %q, %v\n", ln, err)
			continue
		}
		arr = append(arr, n)
	}
	return arr, nil
}

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

//
// input is not large enough to overflow but in the general case this property
// of modulo might be useful.
//
// (a + b) mod n = [(a mod n) + (b mod n)] mod n.
//

func main() {
	if len(os.Args) != 2 {
		fmt.Println("bad usage")
		os.Exit(1)
	}
	arr, err := readArray(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var s int64
	var hHigh, hLow heap
	for _, n := range arr {
		if len(hHigh) == 0 || n >= hHigh.peek() {
			hHigh.insert(n, n)
		} else {
			hLow.insert(-n, n)
		}
		for len(hHigh) > len(hLow)+1 {
			e, _ := hHigh.extractMin()
			hLow.insert(-e.key, e.val)
		}
		for len(hLow) > len(hHigh) {
			e, _ := hLow.extractMin()
			hHigh.insert(-e.key, e.val)
		}
		s += int64(hHigh.peek())
	}
	fmt.Println(s % 10000)
}



