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

func (tc *TagCounter) SummarizeAsync(entries <-chan logentry.LogEntry) {
	tc.waitGroup.Add(1)
	go func() {
		tc.Summarize(entries)
		tc.waitGroup.Done()
	}()
}

func (tc *TagCounter) Summarize(entries <-chan logentry.LogEntry) {
	for entry := range entries {
		for _, tag := range entry.Tags {
			count := tc.countPerTag[tag]
			count++
			tc.countPerTag[tag] = count
		}
	}
}

func (tc *TagCounter) StringAfterSummarizeAsyncCompleted() string {
	tc.waitGroup.Wait()
	return tc.String()
}

func (tc *TagCounter) String() string {
	var buffer bytes.Buffer

	for tag, count := range tc.countPerTag {
		buffer.WriteString(tag + ": " + strconv.Itoa(count) + "\n")
	}

	return buffer.String()
}
