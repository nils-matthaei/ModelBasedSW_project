package main

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

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

func CollectMapValues[K comparable, V any](m map[K]V) []V {
	return slices.Collect(maps.Values(m))
}

func CollectMapsValuesSorted[K comparable, V cmp.Ordered](m map[K]V) []V {
	return slices.Sorted(maps.Values(m))
}

func CollectMapKeys[K comparable, V any](m map[K]V) []K {
	return slices.Collect(maps.Keys(m))
}

func CollectMapKeysSorted[K cmp.Ordered, V any](m map[K]V) []K {
	return slices.Sorted(maps.Keys(m))
}

func STLfunctions() {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	s_v := CollectMapValues(m)
	for _, v := range s_v {
		println(v)
	}

	s_k := CollectMapKeys(m)
	for _, k := range s_k {
		println(k)
	}

	s_v_sorted := CollectMapsValuesSorted(m)
	for _, v := range s_v_sorted {
		println(v)
	}

	// Print values in reverse order using slices.Backward
	for _, v := range slices.Backward(s_v_sorted) {
		println(v)
	}

	s_k_sorted := CollectMapKeysSorted(m)
	for _, k := range s_k_sorted {
		println(k)
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
	STLfunctions()
}
