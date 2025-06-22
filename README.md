# Range Over Function Types
Range Over Function Types in Go - This repository contains my final project for the lecture: Model-Based Software Development.

## 1 About the project
Range over function types is a relatively new feature in Go. It was added in version 1.23 on August 13th, 2024. This project aims to answer the following questions:
- What problem does this feature solve?
  - What are its typical use cases?
- Is the provided solution complete and satisfactory?
- Is this an important addition to the Go language?

## 2 Motivation
Go version 1.18 added the ability for users to define their own generic container types such as Sets, Queues, Trees, etc. However, what was still lacking was a standardized way to iterate over all elements in a user-defined collection. In many other programming languages like C++, Python, Java, or Rust, just to name a few, this is solved via _Iterators_. Iterators serve as a way to access all elements in a container, or aggregate object, without directly exposing the underlying data structure.

In Go version 1.23, a standardized form of Iterators was added. To achieve this, certain function types can act as Iterators and thus the **_ranging_ over these function types** was added. Finally, these features can be leveraged to range over user-defined containers, but also offer some additional benefits.

## 3 Iterators
The new for/range statement does not directly support arbitrary function types. As of Go 1.23 it supports ranging over specific functions that:
- take a single argument
- the argument itself must be a function that:
  - takes zero to two arguments
  - returns a bool

This, by convention, is called the _yield function_.

```Go
func(yield func() bool)

func(yield func(V) bool)

func(yield func(K, V) bool)
```
In Go, a function of one of these types, is called a _push iterator_. This is the default kind of iterator and therefore also often just called iterator. There is, howerver, also another kind of iterator called the _pull iterator_.

### 3.1 Push Iterators
A push iterator is a function that "pushes" values to a provided yield function, one at a time. When ranging over a push iterator, the Go runtime repeatedly calls the iterator, which in turn invokes the yield function for each element in the collection. This allows the iterator to control the flow of elements, sending each value to the yield function until there are no more elements or the yield function signals to stop.

To implement iterators, Go 1.23 introduced a new standard library package called _iter_. This packages defines the types `Seq` and `Seq2`. These are names for two of the iterator function types:

```Go
package iter

type Seq[V any] func(yield func(V) bool)

type Seq2[K, V any] func(yield func(K, V) bool)
```
So far there is no `Seq0` that would define `func(yield func() bool)`, for some reason.

Using these it's quite simple to define an iterator over a generic container type.
As an example, let's look at this simple implementation of a Stack:
```Go
type Stack[T any] struct {
	data []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(value T) {
	s.data = append(s.data, value)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	index := len(s.data) - 1
	val := s.data[index]
	s.data = s.data[:index]
	return val, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.data) == 0 {
		var zero T
		return zero, false
	}
	return s.data[len(s.data)-1], true
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}

func (s *Stack[T]) IsEmpty() bool {
	return len(s.data) == 0
}
```
By convention, the method to return an iterator over a contaier is called `All()`. For our Stack, this method could look like this:
```Go
import "iter"
/*...*/
func (s *Stack[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s.data {
			if !yield(v) {
				return
			}
		}
	}
}
```
With the `All()` method implemented it is possible to iterate over the Stack using the iterator it provides.
```Go
func PrintStack[T any](s *Stack[T]) {
	for v := range s.All() {
		println(v)
	}
}
```
The returned iterator function takes the _yield_ function as an argument. In the case of the for/range statement that _yield_ fuction is automatically generated at compile time. This is relatively complex to do efficiently, but luckily the details of this are not important for using iterators. The way iteratiors like work is:
- For a sequence of values, the iterator calls a yield function with each value in the sequence
- If the yield function returns false, no more values are needed
   - The iterator can now return and do any needed cleanup
- If the yield function never returns false: return after calling yield with all values

### 3.2 Pull Iterators

Pull iterators can be used to access one value at a time, instead of ranging over the entire sequence as is the case with Push iterators. To get a pull iterator the `iter` package provides two functions `Pull` and `Pull2`. They convert a push iterator into a pull iterator, by returning two functions `next` and `stop`
- `next` returns the next value in the sequence and a boolean indicating whether the value is valid
- `stop` ends the iteration
  - must be called when caller no longer wants to iterate, but `next` has not yet indicated that the sequence is over
  - typically callers should `defer stop()`

For example this could be used to create a pairwise iterator

```Go
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
```
The `Pairwise` function takes a standard `Seq` iterator and returns a `Seq2` iterator over pairs of elements. Calling the `PrintPairs` function would create an output like this:

```sh
0 1
1 2
2 3
3 4
```

## 4 Standard Library Functions
TODO
### 4.1 Slices
TODO
### 4.2 Maps
TODO

## 5 Summary
TODO

## 6 Sources
- https://go.dev/blog/range-functions
- https://pkg.go.dev/iter
