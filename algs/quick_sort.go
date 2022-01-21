package main

import (
	"os"
	"io/ioutil"
	"bufio"
	"strconv"
	"fmt"
	"strings"
)

func partition(arr []int) int {
	if len(arr) <= 1 {
		return 0
	}
	p := arr[0]
	i := 1
	for j, n := range arr {
		if n < p {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[0], arr[i-1] = arr[i-1], arr[0]
	return i-1
}

func qsort(arr []int, choose func(arr []int) int) int {
	if len(arr) <= 1 {
		return 0
	}
	i := choose(arr)
	arr[0], arr[i] = arr[i], arr[0]
	i = partition(arr)
	lcomps := qsort(arr[:i], choose)
	rcomps := qsort(arr[i+1:], choose)
	return len(arr)-1 + lcomps + rcomps
}

func first(arr []int) int {
	return 0
}

func last(arr []int) int {
	return len(arr)-1
}

func isMedian(a, b, c int) bool {
	return !(a < b && a < c) && !(a > b && a > c)
}

func medianOfThree(arr []int) int {
	fi, li := 0, len(arr)-1
	mi := len(arr)/2
	if len(arr)%2 == 0 {
		mi--
	}
	// fmt.Println("mi: %d\n, len(arr): %d", mi, len(arr))
	f, l, m := arr[fi], arr[li], arr[mi]
	if isMedian(f, l, m) {
		// fmt.Println("chose %d out of %d %d %d\n", f, f, l, m)
		return fi
	} 
	if isMedian(l, f, m) {
		// fmt.Println("chose %d out of %d %d %d\n", l, f, l, m)
		return li
	} 
	// fmt.Println("chose %d out of %d %d %d\n", m, f, l, m)
	if isMedian(m, f, l) {
		return mi
	}
	panic(fmt.Errorf("invalid code path"))
}

func readIntArray(path string) ([]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	arr := make([]int, 0)
	for s.Scan() {
		n, err := strconv.Atoi(s.Text())
		if err != nil {
			return nil, err
		}
		arr = append(arr, n)
	}
	return arr, nil
}

type testCaseInfo struct {
	inFile  string
	outFile string
}

var testCases []testCaseInfo

func init() {
	files, err := ioutil.ReadDir("TestCases")
	if err != nil {
		panic(err)
	}
	for _, fInfo := range files {
		if strings.HasPrefix(fInfo.Name(), "input") {
			outputName := strings.Replace(fInfo.Name(), "input", "output", 1)
			tc := testCaseInfo{"TestCases/"+fInfo.Name(), "TestCases/"+outputName}
			testCases = append(testCases, tc)
		}
	}
}

func copyOfIntSlice(a []int) []int {
	c := make([]int, len(a))
	copy(c, a)
	return c
}

func test() {
	failures, firstFailures, lastFailures, medianFailures := 0, 0, 0, 0
	for _, tc := range testCases {
		arr, err := readIntArray(tc.inFile)
		if err != nil {
			fmt.Printf("failed to read input file for test case: %q, skipping: %v\n", tc.inFile, err)
			continue
		}
		cmp, err := readIntArray(tc.outFile)
		if err != nil {
			fmt.Printf("failed to read output file for test case: %q, skipping: %v\n", tc.inFile, err)
			continue
		}
		res := []int{
			qsort(copyOfIntSlice(arr), first),
			qsort(copyOfIntSlice(arr), last),
			qsort(copyOfIntSlice(arr), medianOfThree),
		}
		success := true
		for i, c := range cmp {
			if c != res[i] {
				success = false
				break
			}
		}
		if !success {
			fmt.Println("!!!")
			fmt.Printf("FAILED: test case, input file: %q\n", tc.inFile)
			if cmp[0] != res[0] {
				fmt.Printf("\t > \"first\" failed, expected: %d, actual: %d\n", cmp[0], res[0])
				firstFailures++
			}
			if cmp[1] != res[1] {
				fmt.Printf("\t > \"last\" failed, expected: %d, actual: %d\n", cmp[1], res[1])
				lastFailures++
			}
			if cmp[2] != res[2] {
				fmt.Printf("\t > \"median-of-three\" failed, expected: %d, actual: %d\n", cmp[2], res[2])
				medianFailures++
			}
			fmt.Println("!!!")
			fmt.Println()
			failures++
		} else {
			fmt.Printf("PASSED: %v\n", tc)
		}
	}
	fmt.Println("===")
	if failures == 0 {
		fmt.Println("TEST SUITE SUCCEEDED")
	} else {
		fmt.Printf("TEST SUITE FAILED: failed %d/%d\n", failures, len(testCases))
		fmt.Printf("\t > first: %d\n", firstFailures)
		fmt.Printf("\t > last: %d\n", lastFailures)
		fmt.Printf("\t > median: %d\n", medianFailures)
	}
	fmt.Println("---")
}

func main() {
	if false {
		test()
	} else {
		arr, err := readIntArray("QuickSort.txt")
		if err != nil {
			panic(err)
		}
		fmt.Printf("first: %d\n" , qsort(copyOfIntSlice(arr), first))
		fmt.Printf("last: %d\n"  , qsort(copyOfIntSlice(arr), last))
		fmt.Printf("median: %d\n", qsort(copyOfIntSlice(arr), medianOfThree))
	}
}
