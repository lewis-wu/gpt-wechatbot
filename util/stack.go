package util

type Stack struct {
	items    []interface{}
	top      int
	capacity int
}

func NewStack(capacity int) *Stack {
	if capacity < 0 {
		panic("capacity should more than zero")
	}
	return &Stack{
		items:    make([]interface{}, capacity),
		top:      -1,
		capacity: capacity,
	}
}

func (s *Stack) Push(item interface{}) {
	if s.top == s.capacity-1 {
		s.capacity *= 2
		newItems := make([]interface{}, s.capacity)
		copy(newItems, s.items)
		s.items = newItems
	}

	s.top++
	s.items[s.top] = item
}

func (s *Stack) PushMany(items ...interface{}) {
	for _, item := range items {
		s.Push(item)
	}
}

func (s *Stack) Pop() (interface{}, bool) {
	if s.top < 0 {
		return nil, false // 栈已空
	}

	item := s.items[s.top]
	s.items[s.top] = nil
	s.top--
	return item, true
}

func (s *Stack) PopAll() []interface{} {
	result := make([]interface{}, 0, s.capacity)
	for {
		item, ok := s.Pop()
		if !ok {
			break
		}
		result = append(result, item)
	}
	return result
}

func (s *Stack) Peek() (interface{}, bool) {
	if s.top < 0 {
		return nil, false // 栈已空
	}

	item := s.items[s.top]
	return item, true
}
