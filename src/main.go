package main

import "iter"

func PrintStack[T any](s *Stack[T]) {
	// Iterate over stack using the All() function
	for v := range s.All() {
		println(v)
	}
}

func Pairwise[V any](seq iter.Seq[V]) iter.Seq2[V, V] {
	return func(yield func(V, V) bool) {
		next1, stop1 := iter.Pull(seq)
		defer stop1()
		next2, stop2 := iter.Pull(seq)
		defer stop2()

		// Advance the second iterator to start from the second element
		next2()

		for {
			v1, ok1 := next1()
			if !ok1 {
				return
			}
			v2, ok2 := next2()
			// Return if ok2 is false
			if !ok2 {
				return
			}
			// Yield the pair of values
			if !yield(v1, v2) {
				return
			}
		}
	}
}

func PrintPairs[T any](s *Stack[T]) {
	for v1, v2 := range Pairwise(s.All()) {
		println(v1, v2)
	}
}

func main() {
	// Build stack
	s := NewStack[int]()
	for i := 0; i < 5; i++ {
		s.Push(i)
	}

	PrintStack(s)
	PrintPairs(s)
}
