package summary

import (
	"strings"
)

// A list of strings that might be swapped, accessed by index, and joined alltogether.
// Taken from golang.org/x/text/collate/sort_test.go, implements collate's Lister interface
type stringList []string

func (s stringList) Append(str string) stringList {
	return append(s, str)
}

func (s stringList) Len() int {
	return len(s)
}

func (s stringList) Swap(i, j int) {
	s[j], s[i] = s[i], s[j]
}

func (s stringList) Bytes(i int) []byte {
	return []byte(s[i])
}

// Joins all the strings in this stringList, separating each by the given separator.
func (s stringList) Join(sep string) string {
	return strings.Join(s, sep)
}

// Returns the result of Join(sep) as []byte.
func (s stringList) JoinBytes(sep string) []byte {
	return []byte(s.Join(sep))
}
