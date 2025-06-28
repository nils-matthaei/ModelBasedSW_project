mod stack;
use stack::Stack;

fn print_stack<T: std::fmt::Display>(stack: &Stack<T>) {
    // Push-style: for loop (uses IntoIterator Trait)
    for value in stack {
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

fn main() {
    let mut stack = Stack::new();
    for i in 0..5 {
        stack.push(i);
    }

    print_stack(&stack);

    print_pairwise(&stack);
}
