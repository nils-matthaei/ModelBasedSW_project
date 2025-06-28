pub struct Stack<T> {
    data: Vec<T>,
}

impl<T> Stack<T> {
    pub fn new() -> Self {
        Stack { data: Vec::new() }
    }
    pub fn push(&mut self, value: T) {
        self.data.push(value);
    }
    pub fn pop(&mut self) -> Option<T> {
        self.data.pop()
    }
}

// Push-style iterator over Stack by value (consuming)
impl<T> IntoIterator for Stack<T> {
    type Item = T;
    type IntoIter = std::vec::IntoIter<T>;
    fn into_iter(self) -> Self::IntoIter {
        self.data.into_iter()
    }
}

// Push-style iterator over Stack by reference (non-consuming)
impl<'a, T> IntoIterator for &'a Stack<T> {
    type Item = &'a T;
    type IntoIter = std::slice::Iter<'a, T>;
    fn into_iter(self) -> Self::IntoIter {
        self.data.iter()
    }
}

// Pull-style iterator for Stack
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

impl<T> Stack<T> {
    pub fn iter(&self) -> StackIter<T> {
        StackIter {
            stack: self,
            index: 0,
        }
    }
}
