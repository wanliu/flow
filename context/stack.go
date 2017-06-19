package context

type Stack struct {
	Root     *ctxt
	Children []Context
}

func NewStack(root *ctxt) *Stack {
	return &Stack{
		Root:     root,
		Children: make([]Context, 0),
	}
}

func (s *Stack) Peek() Context {
	if len(s.Children) == 0 {
		return nil
	} else {
		return s.Children[len(s.Children)-1]
	}
}

func (s *Stack) Push(ctx Context) {
	s.Children = append(s.Children, ctx)
}

func (s *Stack) Pop() Context {
	if len(s.Children) == 0 {
		return nil
	} else {
		ctx := s.Peek()
		s.Children = s.Children[:len(s.Children)-1]
		return ctx
	}
}
