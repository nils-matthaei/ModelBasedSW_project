# Range Over Function Types
Range Over Function Types in Go - This repository contains my final project for the lecture: Model-Based Software Development.

## 1 About the project
Range over function types is a relatively new feature in Go. It was added in version 1.23 on August 13th, 2024. This project aims to answer the following questions:
1. What problem does this feature solve? What are its typical use cases?
2. Is the provided solution complete and satisfactory?
3. Is this an important addition to the Go language?

## 2 Motivation
Go version 1.18 added the ability for users to define their own generic container types such as Sets, Queues, Trees, etc. However, what was still lacking was a standardized way to iterate over all elements in a user-defined collection. In many other programming languages like C++, Python, Java, or Rust, just to name a few, this is solved via _Iterators_. Iterators serve as a way to access all elements in a container, or aggregate object, without directly exposing the underlying data structure.

In Go version 1.23, a standardized form of Iterators was added. To achieve this, certain function types can act as Iterators and thus the **_ranging_ over these function types** was added. Finally, these features can be leveraged to range over user-defined containers, but also offer some additional benefits.

## 3 Iterators
The new `for/range` statement does not directly support arbitrary function types. As of Go 1.23 it supports ranging over specific functions that:
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
In Go, a function of one of these types, is called a _push iterator_. This is the default kind of iterator and therefore also often just called iterator. There is, however, also another kind of iterator called the _pull iterator_.

### 3.1 Push Iterators
A push iterator is a function that "pushes" values to a provided yield function, one at a time. When ranging over a push iterator, the Go runtime repeatedly calls the iterator, which in turn invokes the yield function for each element in the collection. This allows the iterator to control the flow of elements, sending each value to the yield function until there are no more elements or the yield function signals to stop.

To implement iterators, Go 1.23 introduced a new standard library package called _iter_. This package defines the types `Seq` and `Seq2`. These are names for two of the iterator function types:

```Go
package iter

type Seq[V any] func(yield func(V) bool)

type Seq2[K, V any] func(yield func(K, V) bool)
```
So far, there is no `Seq0` that would define `func(yield func() bool)`, for some reason.

Using these, it's quite simple to define an iterator over a generic container type.
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
By convention, the method to return an iterator over a container is called `All()`. For our Stack, this method could look like this:
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
With the `All()` method implemented, it is possible to iterate over the Stack using the iterator it provides.
```Go
func PrintStack[T any](s *Stack[T]) {
	for v := range s.All() {
		println(v)
	}
}
```
The returned iterator function takes the _yield_ function as an argument. In the case of the `for/range` statement, that _yield_ function is automatically generated at compile time. This is relatively complex to do efficiently, but luckily the details of this are not important for using iterators. The way iterators like this work is:
- For a sequence of values, the iterator calls a yield function with each value in the sequence.
- If the yield function returns false, no more values are needed.
   - The iterator can now return and do any needed cleanup.
- If the yield function never returns false: return after calling yield with all values.

### 3.2 Pull Iterators

Pull iterators can be used to access one value at a time, instead of ranging over the entire sequence as is the case with push iterators. To get a pull iterator, the `iter` package provides two functions: `Pull` and `Pull2`. They convert a push iterator into a pull iterator by returning two functions, `next` and `stop`:
- `next` returns the next value in the sequence and a boolean indicating whether the value is valid.
- `stop` ends the iteration.
  - Must be called when the caller no longer wants to iterate, but `next` has not yet indicated that the sequence is over.
  - Typically, callers should `defer stop()`.

For example, this could be used to create a pairwise iterator:

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
Together with iterators came some new standard library functions in the `slices` and `maps` packages that work with iterators. These provide some useful new features when working with slices and maps.

Here are some examples:

- `All`: this function takes a slice or a map and returns a `Seq2` iterator over the container
- `Values`: takes a slice or map and returns a `Seq` iterator over the container's values
- `Keys`: takes a map and returns a `Seq` iterator over its keys
- `Collect`:
  - for slices: takes a `Seq` iterator and collects its values in a slice
  - for maps: takes a `Seq2` iterator and collects its key-value pairs in a map
- `Backward`: takes a slice and returns an iterator that iterates over the values of the slice backwards
- `Sorted`: takes an iterator over ordered values and collects the values into a slice, then sorts the slice and returns it.

These new functions provide some extra convenience when handling slices and maps.

For example:

```Go
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
```

## 5 Comparison to Rust
Now we quickly want to look at how another modern programming language, in our case _Rust_, solves the issue of iteration over user-defined container types. In Rust this is solved via the `Iterator` and `IntoIterator` _Traits_ (a _Trait_ in Rust is like an _Interface_ in other programming languages). For this purpose, the earlier Go examples were translated into Rust, the entire source code for this can be found [here](./src_rs).

### 5.1 Iterator
The `Iterator` Trait requires one method: `next`. As usual, this method advances the iterator and returns the next value. Creating an iterator in Rust requires two steps:

1. Create a struct holding the iterator state.
2. Implement the `Iterator` trait by implementing `next`.

An iterator implementation for a custom _Stack_ type could look like this.
```Rust
pub struct StackIter<'a, T> {
    stack: &'a Stack<T>,
    index: usize,
}

impl<'a, T> Iterator for StackIter<'a, T> {
    type Item = &'a T;
    fn next(&mut self) -> Option<Self::Item> {
        if self.index >= self.stack.data.len() {
            None
        } else {
            let item = &self.stack.data[self.index];
            self.index += 1;
            Some(item)
        }
    }
}
```
By convention, the `iter()` method is used to return an iterator on a container:
```Rust
impl<T> Stack<T> {
    pub fn iter(&self) -> StackIter<T> {
        StackIter {
            stack: self,
            index: 0,
        }
    }
}
```
The returned iterator can then be used both as a push iterator, to iterate over the container using `for/in` (like for/range in Go); and as a pull iterator using the `next` method:
```Rust
fn print_stack_iter<T: std::fmt::Display>(stack: &Stack<T>) {
    // Push-style: for loop using an Iterator
    for value in stack.iter() {
        println!("{value}");
    }
}

fn print_pairwise<T: std::fmt::Display>(stack: &Stack<T>){
    let mut iter1 = stack.iter();
    let mut iter2 = stack.iter();
    // Advance iter2 by one Element
    iter2.next();

    // Pull-style: manual iteration
    while let Some(value1) = iter1.next() && let Some(value2) = iter2.next() {
        println!("{value1} {value2}")
    }
}
```

## 5.2 IntoIterator
The `IntoIterator` defines how values of a type can be _converted into_ an iterator. One major benefit of this is that it adds some _syntactic sugar_ when working with `for` loops.

`IntoIterator` requires the `into_iter()` method to be implemented for a type.

For our Stack example, this could look like this:
```Rust
impl<'a, T> IntoIterator for &'a Stack<T> {
    type Item = &'a T;
    type IntoIter = StackIter<'a, T>;
    fn into_iter(self) -> Self::IntoIter {
        self.iter()
    }
}
```
Having implemented `IntoIterator` for a type means that `for/in` can be used directly on a value of that type, without needing to manually call a function returning an iterator:

```Rust
fn print_stack<T: std::fmt::Display>(stack: &Stack<T>) {
    // Push-style: for loop using IntoIterator Trait
    for value in stack {
        println!("{value}");
    }
}
```

## 5.3 Iterator provided functions
A major difference between Traits in Rust and Interfaces in other programming languages is that Traits can also contain so-called "provided methods", which are default method implementations that types adopting the trait automatically inherit unless they choose to override them. Implementing `Iterator` provides a type with **75** such methods. These are similar to the standard library functions described in [chapter 4](#4-standard-library-functions). The main difference is that the provided methods in Rust can be used for _any type_ implementing `Iterator`, instead of just the standard library types.

Those methods include:

- `last`
- `count`
- `collect`, `collect_into`
- `map`
- `filter`
- `find`
- `max`, `min`
- `rev` (reverse)
- `sum`

just to name a few useful ones.

## 6 Summary
The goal of this project was to answer the following three questions:

> What problem does this feature solve? What are its typical use cases?

This was answered extensively in chapters 2 and 3. In short: the ability to range over (certain) function types was added to provide a standardized way to iterate over user-defined container types. This was supposed to be what _Iterators_ are in other programming languages.

> Is this an important addition to the Go language?

The answer to this is entirely subjective. In my opinion the answer is: no... and yes.
- No because: This feature is not actually a new feature, Go did not get more powerful by adding it, it just introduced a standardized way of doing something that was already possible.
- Yes because: Personally I think standardization is important because it helps make code more readable and thus better maintainable. Especially in cases as common as iterating over a user-defined container.

> Is the provided solution complete and satisfactory?

_Is the solution complete?_ Mostly yes. It allows for push-style and pull-style iteration over containers, and the `iter` package provides a standard way to implement this. However, the goal of standardization still somewhat fails in my eyes, since iterating over user-defined containers and standard library containers is still different; the former need to be iterated over using an iterator _conventionally_ provided by the `All` method, while the latter don't need this and can just be iterated over using `for/range` directly.

_Is the solution satisfactory?_ Again, this is very subjective. Personally, I prefer the way that Rust chose with an interface that requires the `next` function. When thinking of an interface, this is what I would expect and find intuitive, and thus I found the way chosen in Go very unintuitive at first. However, once I understood it and got used to it, it was not really any harder to use or implement than I was used to from other languages. But in the end, I'm still left with some questions:
- Why was this unconventional way of implementing iterators chosen?
- Why the difference between iterating over user-defined containers and standard library containers?
- Why isn't there an interface that requires implementing the `All` method?

In conclusion, I say that I think providing a standardized iterator implementation was a useful addition to the Go language; however, the way this was done sadly feels more like an added-on hack than a smoothly integrated feature.

## 7 Sources
- https://go.dev/blog/range-functions
- https://pkg.go.dev/iter
- https://pkg.go.dev/slices
- https://pkg.go.dev/maps
- https://doc.rust-lang.org/std/iter/index.html
- https://doc.rust-lang.org/std/iter/trait.Iterator.html
- https://doc.rust-lang.org/std/iter/trait.IntoIterator.html
