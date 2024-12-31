package util

import (
	"testing"
)

func TestStack(t *testing.T) {
	stack := &Stack[int]{}

	// Test IsEmpty on a new stack
	if !stack.IsEmpty() {
		t.Errorf("expected stack to be empty")
	}

	// Test Push and Peek
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	if top, ok := stack.Peek(); !ok || top != 3 {
		t.Errorf("expected top element to be 3, got %v", top)
	}

	// Test Pop
	if popped, ok := stack.Pop(); !ok || popped != 3 {
		t.Errorf("expected popped element to be 3, got %v", popped)
	}

	if popped, ok := stack.Pop(); !ok || popped != 2 {
		t.Errorf("expected popped element to be 2, got %v", popped)
	}

	if popped, ok := stack.Pop(); !ok || popped != 1 {
		t.Errorf("expected popped element to be 1, got %v", popped)
	}

	// Test Pop on an empty stack
	if _, ok := stack.Pop(); ok {
		t.Errorf("expected Pop on empty stack to return false")
	}

	// Test Peek on an empty stack
	if _, ok := stack.Peek(); ok {
		t.Errorf("expected Peek on empty stack to return false")
	}

	// Test IsEmpty after popping all elements
	if !stack.IsEmpty() {
		t.Errorf("expected stack to be empty after popping all elements")
	}
}

func TestGenericStack(t *testing.T) {
	stack := &Stack[string]{}

	stack.Push("A")
	stack.Push("B")
	stack.Push("C")

	if top, ok := stack.Peek(); !ok || top != "C" {
		t.Errorf("expected top element to be 'C', got %v", top)
	}

	if popped, ok := stack.Pop(); !ok || popped != "C" {
		t.Errorf("expected popped element to be 'C', got %v", popped)
	}

	if popped, ok := stack.Pop(); !ok || popped != "B" {
		t.Errorf("expected popped element to be 'B', got %v", popped)
	}

	if popped, ok := stack.Pop(); !ok || popped != "A" {
		t.Errorf("expected popped element to be 'A', got %v", popped)
	}

	if !stack.IsEmpty() {
		t.Errorf("expected stack to be empty")
	}
}

func TestNewStack(t *testing.T) {
	// Test NewStack with no initial items
	stack := NewStack[int]()
	if !stack.IsEmpty() {
		t.Errorf("expected new stack to be empty")
	}

	// Test NewStack with initial items
	stackWithItems := NewStack(1, 2, 3)
	if stackWithItems.IsEmpty() {
		t.Errorf("expected stack with items to not be empty")
	}

	if top, ok := stackWithItems.Peek(); !ok || top != 3 {
		t.Errorf("expected top element to be 3, got %v", top)
	}
}
