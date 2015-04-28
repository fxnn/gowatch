package summary

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/gemsi/grok"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
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
	predicate           logentry.Predicate
	collator            *collate.Collator
}

func NewGrokCounter(patternsByName map[string]string, predicate logentry.Predicate) (tc *GrokCounter) {
	return NewGrokCounterWithLocale("en-US", patternsByName, predicate)
}

func NewGrokCounterWithLocale(locale string, patternsByName map[string]string, predicate logentry.Predicate) (tc *GrokCounter) {
	// TODO #11 Due to a bug in golang.org/x/text, we need to hard-code the language...
	return NewGrokCounterWithLanguageTag(language.Und, patternsByName, predicate)
}

func NewGrokCounterWithLanguageTag(locale language.Tag, patternsByName map[string]string, predicate logentry.Predicate) (tc *GrokCounter) {
	tc = new(GrokCounter)
	tc.countPerPatternName = make(map[string]int)
	tc.patternsByName = patternsByName
	tc.grok = grok.New()
	tc.predicate = predicate
	tc.collator = collate.New(locale, collate.Numeric)

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
		if tc.predicate.Applies(&entry) {
			tc.summarizeEntry(&entry)
		}
	}
}

func (tc *GrokCounter) summarizeEntry(entry *logentry.LogEntry) {
	for name, pattern := range tc.patternsByName {
		matches, err := tc.grok.Parse(pattern, entry.Message)
		if err != nil {
			log.Printf("Error in GrokCounter.summarizeEntry: %s [entry=%s]", err, entry)
		} else if len(matches) > 0 {
			for matchName, matchContent := range matches {
				// TODO: Add escape possibility
				name = strings.Replace(name, "%{"+matchName+"}", matchContent, len(name))
			}
			tc.countPerPatternName[name] = tc.countPerPatternName[name] + 1
		}
	}
}

func (tc *GrokCounter) StringAfterSummarizeAsyncCompleted() string {
	tc.waitGroup.Wait()
	return tc.String()
}

func (tc *GrokCounter) String() string {
	var result stringList

	for patternName, count := range tc.countPerPatternName {
		result = result.Append(patternName + ": " + strconv.Itoa(count))
	}

	tc.collator.Sort(result)
	return result.Join("\n")
}
