package main

import (
	"fmt"
	"math/rand"
)

const (
	maxFailsToReport = 5
)

var (
	numSamples = 100
	numMinElements = 0
	numMaxElements = 10
)

func randint(min, max int) int {
	if min > max {
		panic("invalid args")
	}
	return min + int(float64(max - min + 1) * rand.Float64())
}

func test(name string, sortfn func([]int) []int) {
	fmt.Printf("testing %s ... ", name)
	numFails := 0
	var failedIn, failedOut [][]int
	for i := 0; i < numSamples; i++ {
		failed := false
		arr := make([]int, randint(numMinElements, numMaxElements))
		var minVal, maxVal = 0, len(arr)
		for i, _ := range arr {
			arr[i] = randint(minVal, maxVal)
		}
		sarr := sortfn(arr)
		if len(arr) != len(sarr) {
			failed = true
		}
		for i := 0; i < len(sarr)-1; i++ {
			if sarr[i] > sarr[i+1] {
				failed = true
				break
			}
		}
		if failed {
			failedIn = append(failedIn, arr)
			failedOut = append(failedOut, sarr)
			numFails++
		}
	}
	if numFails > 0 {
		fmt.Println("failed")
		n := numFails
		if n > maxFailsToReport {
			n = maxFailsToReport
		}
		fmt.Println("\nfailed test cases:")
		fmt.Println("---")
		for i := 0; i < maxFailsToReport; i++ {
			fmt.Println("in :", failedIn[i])
			fmt.Println("out:", failedOut[i])
			fmt.Println()
		}
		fmt.Printf("too many errors -- %d total -- stopping\n", numFails)
	} else {
		fmt.Println("succeeded")
	}
}

// sort1: bubble sort
func sort1(arr []int) []int {
	for {
		numSwaps := 0
		for i := 0; i < len(arr)-1; i++ {
			if arr[i] > arr[i+1] {
				arr[i], arr[i+1] = arr[i+1], arr[i]
				numSwaps++
			}
		}
		if numSwaps == 0 {
			break
		}
	}
	return arr
}

// sort2: merge sort
func sort2(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	} 
	l := sort2(arr[:len(arr)/2])
	r := sort2(arr[len(arr)/2:])
	sarr := make([]int, len(arr))
	var i, j, k int
	for  {
		if i < len(l) && j < len(r) {
			if l[i] < r[j] {
				sarr[k] = l[i]
				i++
			} else {
				sarr[k] = r[j]
				j++
			}
			k++
		} else if i < len(l) {
			sarr[k] = l[i]
			k, i = k+1, i+1
		} else if j < len(r) {
			sarr[k] = r[j]
			k, j = k+1, j+1
		} else {
			break
		}
	}
	return sarr
}

func main() {
	test("bubble-sort", sort1)
	test("merge-sort", sort2)
}


