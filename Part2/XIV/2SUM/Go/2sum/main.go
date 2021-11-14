package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

// TODO
//
// When I have time and want... I should look at optimizing the solve routine.
//
// I was thinking that sorting once up front would allow me to do a binary search
// for a ~20001 element window that I'd have to inspect instead of the entire
// array.
//
// So everything should run as quickly as the 20000 element inputs, assuming there's
// not a crazy amount of duplicates?
//

func readIntArray(path string) ([]int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("readIntArray: failed: %v", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("readIntArray: failed: %v", err)
	}
	tokens := strings.Split(string(b), "\n")
	arr := make([]int, 0, len(tokens))
	for _, rt := range tokens {
		t := strings.Trim(rt, " \n\t")
		if t == "" {
			continue
		}
		n, err := strconv.Atoi(t)
		if err != nil {
			return nil, fmt.Errorf("readIntArray: failed to parse input line: %v", err)
		}
		arr = append(arr, n)
	}
	return arr, nil
}

func solve(arr []int) int {
	sums := 0
	c := make(map[int]bool)
	for _, n := range arr {
		c[n] = true
	}
	tStart, tEnd := -10000, 10000
	for t := tStart; t <= tEnd; t++ {
		for _, n := range arr {
			if c[t-n] && n != t/2 {
				sums++
				break
			}
		}
	}
	return sums
}

func main() {
	algoPath, ok := os.LookupEnv("ALGO_PATH")
	if !ok {
		fmt.Println("error: ALGO_PATH environment variable not set, cannot located generated tests, exiting.")
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Println("incorrect usage: 2sum [<filename>|-g]")
		os.Exit(1)
	}

	if os.Args[1] != "-g" {
		iPath := path.Join(algoPath, os.Args[1])
		iArr, err := readIntArray(iPath)
		if err != nil {
			fmt.Printf("error: failed to read input from %q, error: %v\n", iPath, err)
		}
		fmt.Println(solve(iArr))
		os.Exit(0)
	}

	basePath := path.Join(algoPath, "Tests/Part2/XIV/Generated")
	dirEntries, err := os.ReadDir(basePath)
	if err != nil {
		fmt.Printf("error: failed to read directory at path: %q, error: %v\n", basePath, err)
		os.Exit(1)
	}
	total, fails, maxFails := 0, 0, -1
	for _, e := range dirEntries {
		if !e.Type().IsRegular() || !strings.HasPrefix(e.Name(), "input") {
			continue
		}
		ifName, ofName := e.Name(), strings.Replace(e.Name(), "input", "output", 1)
		ifPath, ofPath := path.Join(basePath, ifName), path.Join(basePath, ofName)
		iArray, err := readIntArray(ifPath)
		if err != nil {
			fmt.Printf("error: failed to read input for path: %q, error: %v\n", ifPath, err)
			os.Exit(1)
		}
		oArray, err := readIntArray(ofPath)
		if err != nil {
			fmt.Printf("error: failed to read output for path: %q, error: %v\n", ifPath, err)
			os.Exit(1)
		}
		if len(oArray) != 1 {
			fmt.Println("error: output array unexpectedly has more than 1 element")
			os.Exit(1)
		}

		arr := iArray
		exp := oArray[0]
		act := solve(arr)

		total++
		if act != exp {
			fails++
			if maxFails != -1 && fails >= maxFails {
				break
			}
			fmt.Printf("FAILED test case at path: %q, expected: %d, actual: %d\n", ifPath, exp, act)
		} else {
			fmt.Printf("PASSED test case at path: %q\n", ifPath)
		}
	}
	if fails > 0 {
		fmt.Printf("failed, %d/%d test cases\n", fails, total)
	} else {
		fmt.Println("success, passed all test cases!")
	}
}
