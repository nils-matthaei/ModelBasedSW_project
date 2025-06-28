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
// Iterator
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

// IntoIterator
impl<'a, T> IntoIterator for &'a Stack<T> {
    type Item = &'a T;
    type IntoIter = StackIter<'a, T>;
    fn into_iter(self) -> Self::IntoIter {
        self.iter()
    }
}
