package summary

import (
	"bytes"
	"github.com/fxnn/gowatch/logentry"
	"strconv"
	"sync"
	"github.com/gemsi/grok"
)

type GrokCounter struct {
	grok			*grok.Grok
	patternsByName        map[string]string
	countPerPatternName map[string]int
	waitGroup   	sync.WaitGroup
}

func NewGrokCounter(patternsByName map[string]string) (tc *GrokCounter) {
	tc = new(GrokCounter)
	tc.countPerPatternName = make(map[string]int)
	tc.patternsByName = patternsByName
	tc.grok = grok.New()

	return
}

func (tc *GrokCounter) SummarizeAsync(entries <-chan logentry.LogEntry) {
	tc.waitGroup.Add(1)
	go  tc.Summarize(entries)
}

func (tc *GrokCounter) Summarize(entries <-chan logentry.LogEntry) {
	for entry := range entries {
		for name, pattern := range tc.patternsByName {
			if ok, _ := tc.grok.Match(pattern, entry.Message); ok {
				// TODO: possibility to generate the name out of the parsed message
				tc.countPerPatternName[name] = tc.countPerPatternName[name] + 1
			}
		}
	}
	tc.waitGroup.Done()
}

func (tc *GrokCounter) String() string {
	var buffer bytes.Buffer

	tc.waitGroup.Wait()
	for patternName, count := range tc.countPerPatternName {
		buffer.WriteString(patternName + ": " + strconv.Itoa(count) + "\n")
	}

	return buffer.String()
}
