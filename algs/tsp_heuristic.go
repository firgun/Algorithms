package main

import (
	"fmt"
	"io"
	"os"
	"math"
)

type position struct { 
	x, y float64 
}

func dist(a, b position) float64 {
	xd := b.x - a.x
	yd := b.y - a.y
	return math.Sqrt(xd*xd + yd*yd)
}

func nearestNeighbour(arr []position) float64 {
	var p int
	var cost float64
	visited := make(map[int]bool)
	visited[p] = true
	for i := 0; i < len(arr)-1; i++ {
		var minDist float64 = math.Inf(1)
		var idx int
		for j := 0; j < len(arr); j++ {
			if visited[j] { 
				continue
			}
			d := dist(arr[p], arr[j])
			if d < minDist {
				minDist = d
				idx = j
			}
		}
		p = idx
		cost += minDist
		visited[p] = true
	}
	cost += dist(arr[0], arr[p])
	return cost
}

func main() {
	i := readInput(os.Stdin)
	fmt.Println(int64(nearestNeighbour(i)))
}

func readInput(r io.Reader) []position {
	var count int
	if _, err := fmt.Fscanf(r, "%d\n", &count); err != nil {
		panic(err)
	}
	arr := make([]position, count)
	for i := 0; i < count; i++ {
		var n int
		var x, y float64
		if _, err := fmt.Fscanf(r, "%d %f %f\n", &n, &x, &y); err != nil {
			panic(err)
		}
		arr[i] = position{x,y}
	}
	return arr
}
