package main

import (
	"fmt"
	"sort"
)

// By following this same pattern of creating a custom type,
// implementing the three Interface methods on that type,
// and then calling sort.Sort on a collection of that custom type, we can sort Go slices by arbitrary functions.
type byLength []string

func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func main() {
	fruits := []string{"peach", "banana", "kiwi"}
	sort.Sort(byLength(fruits))
	fmt.Println(fruits)
}
