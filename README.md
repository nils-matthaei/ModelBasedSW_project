# Range Over Function Types
Range Over Function Types in Go - This repository contains my final project for the lecture: Model-Based Software Development.

## 1 About the project
Range over function types is a relatively new feature in Go. It was added in version 1.23 on August 13th, 2024. This project aims to answer the following questions:
- What problem does this feature solve?
  - What are its typical use cases?
- Is the provided solution complete and satisfactory?
- Is this an important addition to the Go language?

## 2 Motivation and use cases

### 2.1 Motivation
Go version 1.18 added the ability for users to define their own generic container types such as Sets, Queues, Trees, etc. However, what was still lacking was a standardized way to iterate over all elements in a user-defined collection. In many other programming languages like C++, Python, Java, or Rust, just to name a few, this is solved via _Iterators_. Iterators serve as a way to access all elements in a container, or aggregate object, without directly exposing the underlying data structure.

In Go version 1.23, a standardized form of Iterators was added. To achieve this, certain function types can act as Iterators and thus the **_ranging_ over these function types** was added. Finally, these features can be leveraged to range over user-defined containers, but also offer some additional benefits.
