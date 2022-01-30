package main

import (
	"fmt"
	"os"
	"io"
	"math"
)

type pos struct { x, y float64 }

func (u pos) distanceTo(v pos) float64 { return math.Sqrt((u.x-v.x)*(u.x-v.x) + (u.y-v.y)*(u.y-v.y)) }

func tsp(arr []pos, start pos, current pos) float64 {
	if len(arr) == 0 { 
		return start.distanceTo(current)
	}
	min := math.Inf(1)
	for i, next := range arr {
		l := len(arr)-1
		arr[i], arr[l] = arr[l], arr[i]
		cost := current.distanceTo(next) + tsp(arr[:l], start, next)
		if cost < min {
			min = cost
		}
		arr[i], arr[l] = arr[l], arr[i]
	}	
	return min
}

func main() {
	cities := readInput(os.Stdin)
	fmt.Println(tsp(cities[1:], cities[0], cities[0]))
	
	// fmt.Println(len(cities))
	// for _, pos := range cities {
	//     fmt.Printf("%.4f %.4f\n", pos.x, pos.y)
	// }
}

func readInput(r io.Reader) []pos {
	var cityCount int
	_, err := fmt.Fscanf(r, "%d\n", &cityCount)
	if err != nil { panic(err) }
	cities := make([]pos, cityCount)
	for i := 0; i < cityCount; i++ {
		var x, y float64
		_, err = fmt.Fscanf(r, "%f %f\n", &x, &y)
		if err != nil { panic(err) }
		cities[i] = pos{x, y}
	}
	return cities
}