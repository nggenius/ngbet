package bet365

import (
	"container/list"
)

type Stack struct {
	l *list.List
}

func NewStack() *Stack {
	s := new(Stack)
	s.l = list.New()
	return s
}

func (s *Stack) Push(e interface{}) {
	s.l.PushFront(e)
}

func (s *Stack) Pop() {
	if s.l.Front() != nil {
		s.l.Remove(s.l.Front())
	}
}

func (s *Stack) Top() interface{} {
	if s.l.Front() == nil {
		return nil
	}
	return s.l.Front().Value
}
