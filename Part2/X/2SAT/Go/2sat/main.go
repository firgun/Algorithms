// Problem
// ===
// In the 2SAT problem, you are given a set of clauses, where each clause is the
// disjunction of two literals (a literal is a Boolean variable or the negation
// of a Boolean variable). You are looking for a way to assign a value "true" or
// "false" to each of the variables so that all clauses are satisfied --- that
// is, there is at least one true literal in each clause. For this problem,
// design an algorithm that determines whether or not a given 2SAT instance has
// a satisfying assignment. (Your algorithm does not need to exhibit a satisfying
// assignment, just decide whether or not one exists.) Your algorithm should run
// in O(m + n) time, where m and n are the number of clauses and variables,
// respectively. [Hint: strongly connected components.]

// Ideas
// ===
// Fix the order in which you make assignments.
// 
// You can build an undirected graph, the will allow you to find.
//
// You can build a directed graph based on the the arbitrary fixed ordering.
//
// You use the undirected graph as a guide to remove edges from the directed one.

package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")
}
