package summary

import (
	"github.com/fxnn/gowatch/logentry"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
	"strconv"
	"sync"
)

type TagCounter struct {
	countPerTag map[string]int
	waitGroup   sync.WaitGroup
	predicate   logentry.Predicate
	collator    *collate.Collator
}

func NewTagCounter(predicate logentry.Predicate) (tc *TagCounter) {
	return NewTagCounterWithLocale("en_US", predicate)
}

func NewTagCounterWithLocale(locale string, predicate logentry.Predicate) (tc *TagCounter) {
	return NewTagCounterWithLanguageTag(language.Make(locale), predicate)
}

func NewTagCounterWithLanguageTag(locale language.Tag, predicate logentry.Predicate) (tc *TagCounter) {
	tc = new(TagCounter)
	tc.countPerTag = make(map[string]int)
	tc.predicate = predicate
	tc.collator = collate.New(locale, collate.Numeric)

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
		if tc.predicate.Applies(&entry) {
			for _, tag := range entry.Tags {
				count := tc.countPerTag[tag]
				count++
				tc.countPerTag[tag] = count
			}
		}
	}
}

func (tc *TagCounter) StringAfterSummarizeAsyncCompleted() string {
	tc.waitGroup.Wait()
	return tc.String()
}

func (tc *TagCounter) String() string {
	var result stringList

	for tag, count := range tc.countPerTag {
		result = result.Append(tag + ": " + strconv.Itoa(count))
	}

	tc.collator.Sort(result)
	return result.Join("\n")
}
