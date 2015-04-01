package summary

import (
	"bytes"
	"github.com/fxnn/gowatch/logentry"
	"github.com/gemsi/grok"
	"log"
	"strconv"
	"strings"
	"sync"
)

type GrokCounter struct {
	grok                *grok.Grok
	patternsByName      map[string]string
	countPerPatternName map[string]int
	waitGroup           sync.WaitGroup
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
	go func() {
		tc.Summarize(entries)
		tc.waitGroup.Done()
	}()
}

func (tc *GrokCounter) Summarize(entries <-chan logentry.LogEntry) {
	for entry := range entries {
		for name, pattern := range tc.patternsByName {
			matches, err := tc.grok.Parse(pattern, entry.Message)
			if err != nil {
				log.Print(err)
			} else if len(matches) > 0 {
				for matchName, matchContent := range matches {
					// TODO: Add escape possibility
					name = strings.Replace(name, "%{"+matchName+"}", matchContent, len(name))
				}
				tc.countPerPatternName[name] = tc.countPerPatternName[name] + 1
			}
		}
	}
}

func (tc *GrokCounter) StringAfterSummarizeAsyncCompleted() string {
	tc.waitGroup.Wait()
	return tc.String()
}

func (tc *GrokCounter) String() string {
	var buffer bytes.Buffer

	for patternName, count := range tc.countPerPatternName {
		buffer.WriteString(patternName + ": " + strconv.Itoa(count) + "\n")
	}

	return buffer.String()
}
