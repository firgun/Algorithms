package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

type job struct {
	weight, length int
}

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

func readJobsArray(path string) ([]job, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("readJobsArray: failed: %v", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("readJobsArray: failed: %v", err)
	}
	lines := strings.Split(string(b), "\n")
	arr := make([]job, 0, len(lines))
	if len(lines) > 1 {
		for _, ln := range lines[1:] {
			ln := strings.Trim(ln, " \n\t\r")
			if ln == "" {
				continue
			}
			toks := strings.Split(ln, " ")
			if len(toks) != 2 {
				return nil, fmt.Errorf("readJobsArray: failed to parse input line: unexpected number of tokens: %d", len(toks))
			}
			weight, err := strconv.Atoi(strings.Trim(toks[0], " \n\t\r"))
			if err != nil {
				return nil, fmt.Errorf("readJobsArray: failed to parse input token: %v", err)
			}
			length, err := strconv.Atoi(strings.Trim(toks[1], " \n\t\r"))
			if err != nil {
				return nil, fmt.Errorf("readJobsArray: failed to parse input token: %v", err)
			}
			j := job{weight, length}
			arr = append(arr, j)
		}
	}
	return arr, nil
}

type sortJobs struct {
	jobs    []job
	indices []int
	lessFn  func(indices []int, i, j int) bool
}

func (s sortJobs) Less(i, j int) bool {
	return s.lessFn(s.indices, i, j)
}

func (s sortJobs) Swap(i, j int) {
	s.indices[i], s.indices[j] = s.indices[j], s.indices[i]
}

func (s sortJobs) Len() int {
	return len(s.indices)
}

func solveQuestion1(js []job) int {
	return solve(js, func(is []int, i, j int) bool {
		j1 := js[is[i]]
		j2 := js[is[j]]
		d1 := j1.weight - j1.length
		d2 := j2.weight - j2.length
		if d1 == d2 {
			return j1.weight >= j2.weight
		}
		return d1 >= d2
	})
}

func solveQuestion2(js []job) int {
	return solve(js, func(is []int, i, j int) bool {
		j1, j2 := js[is[i]], js[is[j]]
		r1, r2 := float64(j1.weight)/float64(j1.length), float64(j2.weight)/float64(j2.length)
		return r1 >= r2
	})
}

func solve(js []job, lessFn func(is []int, i, j int) bool) int {
	s := make([]int, len(js))
	for i := range s {
		s[i] = i
	}
	sort.Sort(sortJobs{js, s, lessFn})
	t, n := 0, 0
	for _, index := range s {
		j := js[index]
		t += j.length
		n += j.weight * t
	}
	return n
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
		iArr, err := readJobsArray(iPath)
		if err != nil {
			fmt.Printf("error: failed to read input from %q, error: %v\n", iPath, err)
		}
		fmt.Println("Question 1:", solveQuestion1(iArr))
		fmt.Println("Question 2:", solveQuestion2(iArr))
		os.Exit(0)
	}

	basePath := path.Join(algoPath, "Tests/Part3/Week1/Question_1_and_2")
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

		iArray, err := readJobsArray(ifPath)
		if err != nil {
			fmt.Printf("error: failed to read input for path: %q, error: %v\n", ifPath, err)
			os.Exit(1)
		}

		oArray, err := readIntArray(ofPath)
		if err != nil {
			fmt.Printf("error: failed to read output for path: %q, error: %v\n", ifPath, err)
			os.Exit(1)
		}

		if len(oArray) != 2 {
			fmt.Println("error: output array must have exactly 2 elements")
			os.Exit(1)
		}

		arr := iArray

		exp1 := oArray[0]
		exp2 := oArray[1]

		act1 := solveQuestion1(arr)
		act2 := solveQuestion2(arr)

		total += 2

		if act1 != exp1 {
			fails++
			if maxFails != -1 && fails >= maxFails {
				break
			}
			fmt.Printf("FAILED (question 1) test case at path: %q, expected: %d, actual: %d\n", ifPath, exp1, act1)
		} else {
			fmt.Printf("PASSED (question 1) test case at path: %q\n", ifPath)
		}

		if act2 != exp2 {
			fails++
			if maxFails != -1 && fails >= maxFails {
				break
			}
			fmt.Printf("FAILED (question 2) test case at path: %q, expected: %d, actual: %d\n", ifPath, exp2, act2)
		} else {
			fmt.Printf("PASSED (question 2) test case at path: %q\n", ifPath)
		}
	}
	if fails > 0 {
		fmt.Printf("failed, %d/%d test cases\n", fails, total)
	} else {
		fmt.Println("success, passed all test cases!")
	}
}
