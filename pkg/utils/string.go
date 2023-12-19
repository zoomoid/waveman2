/*
Copyright 2022-2023 zoomoid.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

// empty is the map element, as we only want to check existence of keys
type empty struct{}

// stringSet is a map of Empty structs to be used as a Set
type stringSet map[string]empty

// newStringSet constructs a new StringSet from the given list of strings
func newStringSet(items ...string) stringSet {
	ss := stringSet{}
	ss.Insert(items...)
	return ss
}

// Insert adds all items to the set. Already existing items are replaced
func (s stringSet) Insert(items ...string) stringSet {
	for _, item := range items {
		s[item] = empty{}
	}
	return s
}

// Delete removes elements from the map
func (s stringSet) Delete(items ...string) stringSet {
	for _, item := range items {
		delete(s, item)
	}
	return s
}

// Has returns true if an item is already contained in the set
func (s stringSet) Has(item string) bool {
	_, contained := s[item]
	return contained
}

// Len is the canonic length of the underlying map
func (s stringSet) Len() int {
	return len(s)
}
