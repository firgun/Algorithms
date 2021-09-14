package main

import (
	"fmt"
	"math/rand"
	"bufio"
	"os"
	"strconv"
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

func sort(arr []int) ([]int, int) {
	if len(arr) == 1 || len(arr) == 0 {
		return arr, 0
	}
	left, leftInversions := sort(arr[:len(arr)/2]) 
	right, rightInversions := sort(arr[len(arr)/2:])
	sortedArr := make([]int, len(arr))
	var j, k int
	splitInversions := 0 
	for i, _ := range sortedArr {
		if k >= len(right) || (j < len(left) && left[j] <= right[k]) {
			sortedArr[i] = left[j]
			j++
		} else {
			sortedArr[i] = right[k]
			k++
			splitInversions += len(left) - j
		}	
	}
	return sortedArr, leftInversions + splitInversions + rightInversions
}

func readIntArray() ([]int, error) {
	f, err := os.Open("intarray.txt")
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(f)
	arr := make([]int, 100000)
	for s.Scan() {
		n, err := strconv.Atoi(s.Text())
		if err != nil {
			return nil, err
		}
		arr = append(arr, n)
	}
	return arr, nil
}

func main() {
	if false {
		numFails := 0
		var failedIn, failedOut [][]int
		var failedInInv, failedOutInv []int
		for i := 0; i < numSamples; i++ {
			failed := false
			arr := make([]int, randint(numMinElements, numMaxElements))
			var minVal, maxVal = 0, len(arr)
			for i, _ := range arr {
				arr[i] = randint(minVal, maxVal)
			}
			sarr, ninv := sort(arr)
			if len(arr) != len(sarr) {
				failed = true
			}
			for i := 0; i < len(sarr)-1; i++ {
				if sarr[i] > sarr[i+1] {
					failed = true
					break
				}
			}
			count := 0
			for i := 0; i < len(arr)-1; i++ {
				for j := i+1; j < len(arr); j++ {
					if arr[i] > arr[j] {
						count++
					}
				}
			}
			if count != ninv {
				failed = true
			}
			if failed {
				failedIn = append(failedIn, arr)
				failedOut = append(failedOut, sarr)
				failedInInv = append(failedInInv, count)
				failedOutInv = append(failedOutInv, ninv)
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
				fmt.Println("in :", failedIn[i], failedInInv[i])
				fmt.Println("out:", failedOut[i], failedOutInv[i])
				fmt.Println()
			}
			fmt.Printf("too many errors -- %d total -- stopping\n", numFails)
		} else {
			fmt.Println("succeeded")
		}
	} else {
		arr, err := readIntArray()
		if err != nil {
			fmt.Printf("error: failed to read int array: %v\n", err)
			return	
		}
		_, numInversions := sort(arr)
		fmt.Println("num inversions:", numInversions)
	}
}
