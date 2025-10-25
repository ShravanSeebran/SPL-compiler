package analyser

type Binding struct {
	name   string
	nodeID int
}

type AuxillaryStack []Binding

func Empty() *AuxillaryStack {
	stack := AuxillaryStack{}
	return &stack
}

func (s *AuxillaryStack) bind(name string, nodeID int) {
	newBinding := Binding{name: name, nodeID: nodeID}
	*s = append(*s, newBinding)
}

func (s *AuxillaryStack) lookup(name string) (int, bool) {
	for i := len(*s) - 1; i >= 0; i-- {
		if (*s)[i].name == name {
			return (*s)[i].nodeID, true
		}
	}
	return -1, false
}

func (s *AuxillaryStack) enter(scopeID int) int {
	s.bind("$", scopeID)
	return scopeID
}

func (s *AuxillaryStack) exit() int {
	for len(*s) > 0 {
		top := (*s)[len(*s)-1]
		*s = (*s)[:len(*s)-1]
		if top.name == "$" {
			return top.nodeID
		}
	}
	return -1
}
