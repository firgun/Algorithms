package main

import (
	"fmt"
	"strconv"
)

func reversed(s []int) []int {
	r := make([]int, len(s))
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = s[j], s[i]
	}
	if len(s) % 2 != 0 {
		r[len(s)/2] = s[len(s)/2]
	}
	return r
}

// mul1 computes the product of two n-digit integers using long-multiplication
func mul1(a, b []int) []int {
	if len(a) != len(b) || len(a) <= 0 {
		panic("invalid args")
	}
	product := make([]int, 2*len(a)+1)
	for i, x := range a {
		par := make([]int, len(a)+1)
		var j, c int
		for j, c = 0, 0; j < len(b); j++ {
			p := x * b[j] + c
			par[j], c = p % 10, p / 10
		}
		par[j] = c;
		for j, c = i, 0; j < len(product); j++ {
			x := 0
			if j-i < len(par) {
				x = par[j-i]	
			} 
			s := product[j] + x + c
			product[j], c = s % 10, s / 10
		}
	}
	return product
}

func add(a, b []int, shift int) []int {
	c := 0
	for i := shift; i < len(a); i++ {
		n := 0
		if i-shift < len(b) {
			n = b[i-shift]
		}
		s := a[i] + n + c
		a[i], c = s % 10, s / 10
	}
	return a
}

// mul2 computes the product of two n-digit integers using recursive multiplication
func mul2(a, b []int) []int {
	if len(a) != len(b) || len(a) <= 0 {
		panic("invalid args")
	}
	if len(a) == 1 {
		p := a[0] * b[0]
		return []int{p%10, p/10}
	}
	ta := make([]int, len(a)+(len(a)%2))
	copy(ta, a)
	a = ta
	tb := make([]int, len(b)+(len(b)%2))
	copy(tb, b)
	b = tb
	p := make([]int, len(a)*2+1)
	pb, pa := a[:len(a)/2], a[len(a)/2:]
	pd, pc := b[:len(a)/2], b[len(a)/2:]
	ac, bc, ad, bd := mul2(pa, pc), mul2(pb, pc), mul2(pa, pd), mul2(pb, pd)
	p = add(p, ac, 2*(len(a)/2))
	p = add(p, bc, len(a)/2)
	p = add(p, ad, len(a)/2)
	p = add(p, bd, 0)
	return p 
}

// sub: subtracts the integer b from a, a must be greater than b -- we don't
// need negative numbers.
func sub(a []int, b []int) []int {
	c, i := 0, 0
	for ; i < len(b); i++ {
		a[i] -= b[i] + c
		if a[i] < 0 {
			c = 1
			a[i] += 10
		} else {
			c = 0
		}
	}
	if c != 0 {
		if i < len(a) {
			a[i] -= c
		} else {
			panic("invalid args")
		}
	}
	return a
}


// mul3 computes the product of two n-digit integers using karatsuba multiplication
func mul3(a, b []int) []int {
	if len(a) != len(b) || len(a) <= 0 {
		panic("invalid args")
	}
	if len(a) < 4 {
		return mul1(a, b)
	}
	ta := make([]int, len(a)+(len(a)%2))
	copy(ta, a)
	a = ta
	tb := make([]int, len(b)+(len(b)%2))
	copy(tb, b)
	b = tb
	p := make([]int, len(a)*2+1)
	pb, pa := a[:len(a)/2], a[len(a)/2:]
	pd, pc := b[:len(a)/2], b[len(a)/2:]
	ac, bd := mul3(pa, pc), mul3(pb, pd)
	paCopy := make([]int, len(pa)+1); copy(paCopy, pa)
	pcCopy := make([]int, len(pc)+1); copy(pcCopy, pc)
	sumAB := add(paCopy, pb, 0);
	sumCD := add(pcCopy, pd, 0);
	g := mul3(sumAB, sumCD)
	g = sub(g, ac)
	g = sub(g, bd)
	p = add(p, ac, 2*(len(a)/2))
	p = add(p, g, len(a)/2)
	p = add(p, bd, 0)
	return p 
}

func stobi(s string) []int {
	r := make([]int, len(s))
	for i, c := range s {
		d, err := strconv.Atoi(string(c))
		if err != nil {
			panic("invalid args")
		}
		r[len(r)-i-1] = d
	}
	return r
}

func bitos(n []int) string {
	s := ""
	for i := len(n)-1; i >= 0; i-- {
		s += strconv.Itoa(n[i])
	}
	return s
}

func main() {
	fmt.Println(reversed(mul1([]int{8,7,6,5}, []int{4,3,2,1})))
	fmt.Println(reversed(mul1([]int{9,9,9}, []int{9,9,9})))

	fmt.Println(reversed(mul2([]int{8,7,6,5}, []int{4,3,2,1})))
	fmt.Println(reversed(mul2([]int{9,9,9}, []int{9,9,9})))

	fmt.Println(reversed(mul3([]int{8,7,6,5}, []int{4,3,2,1})))
	fmt.Println(reversed(mul3([]int{9,9,9}, []int{9,9,9})))

	x := stobi("3141592653589793238462643383279502884197169399375105820974944592")
	y := stobi("2718281828459045235360287471352662497757247093699959574966967627")

	fmt.Println(reversed(x))
	fmt.Println(reversed(y))

	fmt.Println(bitos(mul1(x, y)))
	fmt.Println(bitos(mul2(x, y)))
	fmt.Println(bitos(mul3(x, y)))

	// fmt.Println(reversed(add([]int{3,2,1,0}, []int{6,5,4}, 0)))
	// fmt.Println(reversed(add([]int{1,1,1,0}, []int{9,1,1}, 0)))

	// fmt.Println(reversed(sub(reversed([]int{5,6,7}), reversed([]int{1,2,3}))))
	// fmt.Println(reversed(sub(reversed([]int{1,1,1}), reversed([]int{0,9,9}))))
}


