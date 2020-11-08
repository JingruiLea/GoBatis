package utils

type Stack []string

// IsEmpty: check if stack is empty
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Length() int {
	return len(*s)
}

// Push a new value onto the stack
func (s *Stack) Push(str string) {
	*s = append(*s, str) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len(*s) - 1   // Get the index of the top most element.
		element := (*s)[index] // Index into the slice and obtain the element.
		*s = (*s)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

type Queue []string

// IsEmpty: check if stack is empty
func (s *Queue) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Queue) Length() int {
	return len(*s)
}

// Push a new value onto the stack
func (s *Queue) Append(str string) {
	*s = append(*s, str) // Simply append the new value to the end of the stack
}

// Remove and return top element of stack. Return false if stack is empty.
func (s *Queue) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		element := (*s)[0] // Index into the slice and obtain the element.
		*s = (*s)[1:]
		return element, true
	}
}
