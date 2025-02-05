package generic

import "slices"

type Stack[T any] struct {
	instance []T
}

func NewStack[T any](capacity int) Stack[T] {
	return Stack[T]{
		instance: make([]T, 0, capacity),
	}
}

func (stack *Stack[T]) Push(item T) {
	stack.instance = append(stack.instance, item)
}

func (stack *Stack[T]) Peek() (result T) {
	i := len(stack.instance) - 1
	if i == -1 {
		return
	}
	return stack.instance[i]
}

func (stack *Stack[T]) Pop() (result T) {
	i := len(stack.instance) - 1
	if i == -1 {
		return
	}
	result = stack.instance[i]
	stack.instance = slices.Delete(stack.instance, i, i+1)
	return
}
