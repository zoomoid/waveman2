package cmdutils

// Empty is the map element, as we only want to check existence of keys
type Empty struct{}

// StringSet is a map of Empty structs to be used as a Set
type StringSet map[string]Empty

// NewStringSet constructs a new StringSet from the given list of strings
func NewStringSet(items ...string) StringSet {
	ss := StringSet{}
	ss.Insert(items...)
	return ss
}

// Insert adds all items to the set. Already existing items are replaced
func (s StringSet) Insert(items ...string) StringSet {
	for _, item := range items {
		s[item] = Empty{}
	}
	return s
}

// Delete removes elements from the map
func (s StringSet) Delete(items ...string) StringSet {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Has returns true if an item is already contained in the set
func (s StringSet) Has(item string) bool {
	_, contained := s[item]
	return contained
}

// Len is the canonic length of the underlying map
func (s StringSet) Len() int {
	return len(s)
}
