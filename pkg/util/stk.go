package util

// Stack represents a simple generic stack data structure.
type Stack[T any] struct {
	items []T
}

// NewStack creates a new stack with optional initial items.
func NewStack[T any](items ...T) *Stack[T] {
	return &Stack[T]{items: items}
}

// Push adds an item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item of the stack.
// It returns the zero value of T if the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if s.IsEmpty() {
		return zero, false
	}
	topIndex := len(s.items) - 1
	topItem := s.items[topIndex]
	s.items = s.items[:topIndex]
	return topItem, true
}

// Peek returns the top item of the stack without removing it.
// It returns the zero value of T if the stack is empty.
func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if s.IsEmpty() {
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// IsEmpty checks if the stack is empty.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}
