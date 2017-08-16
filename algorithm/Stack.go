package algorithm

import (
	"github.com/pkg/errors"
)

type Stack struct {
	array []interface{}
}

func (s *Stack) Size() int {
	return len(s.array)
}

func (s *Stack) Empty() bool {
	return len(s.array) == 0
}

func (s *Stack) PushStack(value interface{}) {
	s.array = append(s.array, value)
}
func (s *Stack) PopStack() (interface{}, error) {
	stackLen := len(s.array)
	if stackLen == 0 {
		return nil, errors.New("stack is empty")
	}
	e := s.array[stackLen-1]
	s.array = s.array[:stackLen-1]
	if e != nil {
		return e, nil
	}
	return nil, errors.New("value is nil")

}
func (s *Stack) Last() interface{} {
	if len(s.array) == 0 {
		return nil
	}
	return s.array[len(s.array)-1]
}
func (s *Stack) StackTop(i int) (interface{}, error) {
	stackLen := len(s.array)
	if stackLen+i > stackLen-1 {
		return nil, errors.Errorf("the index exceeds the boundary :%d", stackLen+i)

	}
	return s.array[stackLen+i], nil
}

func NewStack() *Stack {
	return &Stack{}
}
