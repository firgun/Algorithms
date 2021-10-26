package main

import (
	"fmt"
)

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
	i := len(*h) - 1
	for i != 0 && (*h)[i].key < (*h)[i/2].key {
		(*h)[i], (*h)[i/2] = (*h)[i/2], (*h)[i]
		i = i / 2
	}
}

func (h *heap) extractMin() (heapEntry, bool) {
	if len(*h) == 0 {
		return heapEntry{}, false
	}
	r := (*h)[0]
	(*h)[0], (*h)[len(*h)-1] = (*h)[len(*h)-1], (*h)[0]
	*h = (*h)[:len(*h)-1]
	for i := 0; i < len(*h) && 2*i < len(*h); i++ {
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
	return r, true
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

func main() {
	h := heap{}
	h.insert(2, "Bob")
	h.insert(3, "Cat")
	h.insert(1, "Sam")
	h.dump()

	fmt.Println()
	for len(h) > 0 {
		fmt.Println(h.extractMin())

		fmt.Println()
		h.dump()
		fmt.Println()
	}
}
