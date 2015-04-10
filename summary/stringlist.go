package summary

import (
	"strings"
)

// taken from golang.org/x/text/collate/sort_test.go, implements collate's Lister interface
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

func (s stringList) Join(sep string) string {
	return strings.Join(s, sep)
}

func (s stringList) JoinBytes(sep string) []byte {
	return []byte(s.Join(sep))
}
