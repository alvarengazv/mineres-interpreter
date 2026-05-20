package parser

type Stack []string

func (s *Stack) Push(v string) {
	*s = append(*s, v)
}

func (s *Stack) Pop() string {
	if len(*s) == 0 {
		return ""
	}
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element
}

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Top() string {
	if s.IsEmpty() {
		return ""
	}
	return (*s)[len(*s)-1]
}
