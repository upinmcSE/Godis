package data_structure

type SimpleSet struct {
	key  string
	dict map[string]struct{}
}

func NewSimpleSet(key string) *SimpleSet {
	return &SimpleSet{
		key:  key,
		dict: make(map[string]struct{}),
	}
}

func (s *SimpleSet) Add(members ...string) int {
	return 0
}

func (s *SimpleSet) Rem(members ...string) int {
	return 0
}

func (s *SimpleSet) IsMember(member string) int {
	return 0
}

func (s *SimpleSet) Members() []string {
	return nil
}
