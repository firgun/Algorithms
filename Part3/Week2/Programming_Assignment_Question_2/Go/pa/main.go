package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	// "sort"
	"strconv"
	"strings"
	"bufio"
)

//
// The idea
//
// We need to find the largest k for which the max spacing between components is at least 
// 3.
//
// In other words, we need to put all 0, 1 and 2 hamming distance nodes in the same 
// component
//
// For 0, it's really easy because that just means the nodes have the same value. We 
// could sort or use maps with the bit pattern as the key and an array of nodes as the
// value.
//
// For 1, we could iterate through the individual bits, toggling them and seeing if the
// bit pattern with the ith bit toggled is in the map and just pick the first element of
// the value.
//
// For 2, things get a little bit more irritating, but I don't think this renders the
// solution undoable.
//
// We iterate over all pairs 1 <= i < j <= n and toggle those bits and check for matches
// in the hash map.
//
// There are at most 24 bits, there are n*(n-1)/2 = 24*23/2 = 276 distinct bit pairs that 
// we could try and toggle.
//
// There are 200'000 bit patterns so that's at most 200'000 * 276 <= 60'000'000
//
// If we were to naively check the hamming distance between each pair for this case
// we'd get...
//
// 200'000 * 200'000 = 4 * 10^10 = 40'000'000'000
//
// So while it's still a lot of work, it's several orders of magnitude better than the 
// naive version for our large test input at least.
//

func hammingDistance(a, b string) int {
	if len(a) != len(b) {
		panic("must be equal length strings")
	}
	d := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			d++
		}
	}
	return d
}

func solve(arr []string) int {
	leaders := make([]int, len(arr))
	for i := range leaders {
		leaders[i] = i
	}
	
	ccs := make([][]int, len(arr))
	for i := range ccs {
		ccs[i] = []int{i}
	}
	
	k := len(arr)
	
	m := make(map[string][]int)
	for index, item := range arr {
		m[item] = append(m[item], index)
	}

	var hdist0_cont, hdist1_cont, hdist2_cont int
		
	for {
		var a, b int
		
		minDist := -1
		
		/*
		
		naive solution
		
		for i := range arr {
			// naive solution
			for j := i+1; j < len(arr); j++ {
				d := hammingDistance(arr[i], arr[j])
				if d <= 2 && (minDist == -1 || d < minDist) && leaders[i] != leaders[j] {
					a, b = i, j
					minDist = d
				}
			}
		}
		*/
		
		// Hamming distance of 0
		for index := hdist0_cont; index < len(arr); index++ {
			item := arr[index]
			if otherIndexes, ok := m[item]; ok {
				for _, otherIndex := range otherIndexes {
					if index != otherIndex && leaders[index] != leaders[otherIndex] {
						a, b = index, otherIndex
						minDist = 0
						goto done0
					}
				}
			}
			hdist0_cont = index
		}
		done0:
		
		// Hamming distance of 1
		if minDist == -1 {
			for index := hdist1_cont; index < len(arr); index++ {
				item := arr[index]
				buf := []byte(item)
				for i := range buf {
					cur := buf[i]
					if cur == '0' {
						buf[i] = '1'
					} else if cur == '1' {
						buf[i] = '0'
					}
					if otherIndexes, ok := m[string(buf)]; ok {
						for _, otherIndex := range otherIndexes {
							if index != otherIndex && leaders[index] != leaders[otherIndex] {
								a, b = index, otherIndex
								minDist = 1
								goto done1
							}
						}
					}
					buf[i] = cur
					hdist1_cont = index
				}
			}
		}
		done1:

		// Hamming distance of 2		
		if minDist == -1 {
			for index := hdist2_cont; index < len(arr); index++ {
				item := arr[index]
				buf := []byte(item)
				for i := range buf {
					for j := i+1; j < len(buf); j++ {
						curi, curj := buf[i], buf[j]
						
						if curi == '0' {
							buf[i] = '1'
						} else if curi == '1' {
							buf[i] = '0'
						} else {
							panic("invalid!!!")
						}
						
						if curj == '0' {
							buf[j] = '1'
						} else if curj == '1' {
							buf[j] = '0'
						} else {
							panic("invalid!!!")
						}
						
						if otherIndexes, ok := m[string(buf)]; ok {
							for _, otherIndex := range otherIndexes {
								if index != otherIndex && leaders[index] != leaders[otherIndex] {
									a, b = index, otherIndex
									minDist = 2
									goto done2
								}
							}
						}
						buf[i] = curi
						buf[j] = curj
					}
				}
				hdist2_cont = index
			}
		}
		done2:
		
		if minDist == -1 { 
			return k
		}
		
		if leaders[a] != leaders[b] {
			lesser, greater := leaders[a], leaders[b]
			if len(ccs[lesser]) > len(ccs[greater]) {
				lesser, greater = greater, lesser
			}
			
			for _, index := range ccs[lesser] {
				leaders[index] = greater
				ccs[greater] = append(ccs[greater], index)
			}
			
			ccs[lesser] = nil
			k--
		}
	}
	
	return -1
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

	basePath := path.Join(algoPath, "Tests/Part3/Week2/Question_2/Generated")
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
			os.	Exit(1)
		}

		exp, err := readOutput(ofPath)
		if err != nil {
			fmt.Printf("error: failed to read output for path: %q, error: %v\n", ifPath, err)
			os.Exit(1)
		}

		act := solve(arr)

		total++
		msg := "pass"
		if act != exp {
			fails++
			if maxFails != -1 && fails >= maxFails {
				break
			}
			msg = "FAIL"
		}

		fmt.Printf("want %5d    got %5d    %s %s\n", exp, act, msg, strings.Replace(ifName, ".txt", "", -1))
	}

	if fails > 0 {
		fmt.Printf("failed, %d/%d test cases\n", fails, total)
	} else {
		fmt.Println("success, passed all test cases!")
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
// IO utils
//

func readInput(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	var numStrings, numBits int
	if s.Scan() {
		tokens := strings.Split(s.Text(), " ")
		if len(tokens) != 2 {
			return nil, fmt.Errorf("readInput: failed to scan header line")
		}
		numStrings, err = strconv.Atoi(strings.Trim(tokens[0], " \n\t\r"))
		if err == nil {
			numBits, err = strconv.Atoi(strings.Trim(tokens[1], " \n\t\r"))
		}
		if err != nil {
			return nil, fmt.Errorf("readInput: error parsing header token: %v", err)
		}
	} else if s.Err() != nil {	
		return nil, fmt.Errorf("readInput: scanner error: %v", s.Err())
	} else {
		return nil, fmt.Errorf("readInput: unexpected EOF")
	}
	bitStrings := make([]string, numStrings)
	for i := 0; s.Scan(); i++ {
		tokens := strings.Split(s.Text(), " ")
		if len(tokens) != numBits {
			return nil, fmt.Errorf("readInput: invalid bit-string length, want: %d, got: %d", numBits, len(tokens))
		}
		bitArray := make([]string, numBits)		
		for i, tok := range tokens {
			tok = strings.Trim(tok, " \t\r\n")
			if tok != "0" && tok != "1" {
				return nil, fmt.Errorf("readInput: expected 0 or 1 got %v", tok)
			}
			bitArray[i] = tok
		}
		if i >= numStrings {
			return nil, fmt.Errorf("readInput: expected %d lines found one more", numStrings)
		}
		bitStrings[i] = strings.Join(bitArray, "")
	}
	if s.Err() != nil {
		return nil, fmt.Errorf("readInput: scanner error: %v", s.Err())
	}
	return bitStrings, nil
}

func readOutput(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return -1, fmt.Errorf("readOutput: failed: %v", err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return -1, fmt.Errorf("readOutput: failed: %v", err)
	}
	return strconv.Atoi(strings.Trim(string(b), " \t\r\n"))
}