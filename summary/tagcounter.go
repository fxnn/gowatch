package summary

import (
	"bytes"
	"github.com/fxnn/gowatch/logentry"
	"strconv"
	"sync"
)

type TagCounter struct {
	countPerTag map[string]int
	waitGroup   sync.WaitGroup
}

func NewTagCounter() (tc *TagCounter) {
	tc = new(TagCounter)
	tc.countPerTag = make(map[string]int)

	return
}

func (tc *TagCounter) Summarize(entries <-chan logentry.LogEntry) {
	tc.waitGroup.Add(1)
	for entry := range entries {
		for _, tag := range entry.Tags {
			count := tc.countPerTag[tag]
			count++
			tc.countPerTag[tag] = count
		}
	}
	tc.waitGroup.Done()
}

func (tc *TagCounter) String() string {
	var buffer bytes.Buffer

	tc.waitGroup.Wait()
	for tag, count := range tc.countPerTag {
		buffer.WriteString(tag + ": " + strconv.Itoa(count) + "\n")
	}

	return buffer.String()
}
