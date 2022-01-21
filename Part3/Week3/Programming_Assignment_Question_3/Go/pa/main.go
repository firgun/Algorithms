package main

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"io/ioutil"
)

// solve, computes the max-weight independent set of a path graph.
func solve(arr []int) string {
	m := make([]int, len(arr)+1)
	m[0] = 0
	m[1] = arr[0]
	for i := 2; i < len(m); i++ {
		if m[i-1] > m[i-2] + arr[i-1] {
			m[i] = m[i-1]
		} else {
			m[i] = m[i-2] + arr[i-1]
		}
	}
	s := make([]bool, len(arr))
	i := len(m)-1
	for ; i >= 2; {
		if m[i-1] >= m[i-2] + arr[i-1] {
			i -= 1
		} else {
			s[i-1] = true
			i -= 2
		}
	}
	if i == 1 {
		s[0] = true
	}
	b := make([]byte, 8)
	for i, n := range []int{1, 2, 3, 4, 17, 117, 517, 997} {
		if n-1 < len(s) && s[n-1] {
			b[i] = '1'
		} else {
			b[i] = '0'
		}
	}
	return string(b)
}

func main() {
	algoPath, ok := os.LookupEnv("ALGO_PATH")
	if !ok {
		fmt.Println("error: ALGO_PATH environment variable not set, cannot locate generated tests, exiting.")
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Println("incorrect usage: supply [<filename>|-g]")
		os.Exit(1)
	}

	if os.Args[1] != "-g" {
		p := path.Join(algoPath, os.Args[1])
		g, err := readInput(p)
		if err != nil {
			fmt.Printf("error: failed to read input from %q, error: %v\n", p, err)
			os.Exit(1)
		}
		fmt.Println(solve(g))
		os.Exit(0)
	}

	basePath := path.Join(algoPath, "Tests/Part3/Week3/Question_2/Generated")
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

		arr, err := readInput(ifPath)
		if err != nil {
			fmt.Printf("error: failed to read input for path: %q, error: %v\n", ifPath, err)
			os.Exit(1)
		}

		exp, err := readOutput(ofPath)
		if err != nil {
			fmt.Printf("error: failed to read output for path: %q, error: %v\n", ifPath, err)
			os.Exit(1)
		}

		act := solve(arr)

		total++
		msg := "pass"
		if exp != act {
			fails++
			if maxFails != -1 && fails >= maxFails {
				break
			}
			msg = "FAIL"
		}

		fmt.Printf("want %s, got %s  %s %s\n", exp, act, msg, strings.Replace(ifName, ".txt", "", -1))
	}

	if fails > 0 {
		fmt.Printf("failed, %d/%d test cases\n", fails, total)
	} else {
		fmt.Println("success, passed all test cases!")
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
// I/O
//

func readInput(path string) ([]int, error) {
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
	return arr[1:], nil
}

func readOutput(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("readOutput: failed: %v", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("readOutput: failed: %v", err)
	}
	toks := strings.Split(string(b), "\n")
	if len(toks) < 2 {
		return "", fmt.Errorf("readOutput: want 2, got %d tokens", len(toks))
	}
	return toks[0], nil
}