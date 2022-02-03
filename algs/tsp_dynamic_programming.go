package main

import (
	"fmt"
	"os"
	"io"
	"math"
)

type pos struct { x, y float64 }

func dist(u, v pos) float64 { return math.Sqrt((u.x-v.x)*(u.x-v.x) + (u.y-v.y)*(u.y-v.y)) }

func numbits(n int) int {
	var c int
	for n != 0 {
		n &= (n-1)
		c++
	}
	return c
}

func tsp(arr []pos) float64 {
	s := make([][]float64, 1<<len(arr))
	for i := range s {
		s[i] = make([]float64, len(arr))
		for j := range s[i] {
			s[i][j] = math.Inf(1)
		}
	}
	s[1][0] = 0

	for size := 2; size <= len(arr); size++ {
		for i := range s {
			if i & 1 == 0 || numbits(i) != size {
				continue
			}

			for j := 1; j < len(arr); j++ {
				if i & (1<<j) == 0 {
					continue
				}

				p := i & (0xffffffff ^ (1<<j))

				for k := 0; k < len(arr); k++ {
					if k != j {
						s[i][j] = math.Min(s[p][k] + dist(arr[k], arr[j]), s[i][j])
					}
				}
			}
		}
	}

	min := math.Inf(1)
	for d, c := range s[1<<len(arr) - 1] {
		min = math.Min(min, c + dist(arr[d], arr[0]))
	}

	return float64(int(min))
}

func main() {
	arr := readInput(os.Stdin)
	fmt.Println(tsp(arr))
}

func readInput(r io.Reader) []pos {
	var count int
	_, err := fmt.Fscanf(r, "%d\n", &count)
	if err != nil { panic(err) }
	arr := make([]pos, count)
	for i := 0; i < count; i++ {
		var x, y float64
		_, err = fmt.Fscanf(r, "%f %f\n", &x, &y)
		if err != nil { panic(err) }
		arr[i] = pos{x, y}
	}
	return arr
}
