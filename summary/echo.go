package summary

import (
	"fmt"
	"github.com/fxnn/gowatch/logentry"
	"sync"
)

type Echo struct {
	outputLines []string
	waitGroup   sync.WaitGroup
	predicate   logentry.Predicate
}

func NewEcho(predicate logentry.Predicate) (e *Echo) {
	return &Echo{predicate: predicate}
}

func (e *Echo) SummarizeAsync(entries <-chan logentry.LogEntry) {
	e.waitGroup.Add(1)
	go func() {
		e.Summarize(entries)
		e.waitGroup.Done()
	}()
}

func (e *Echo) Summarize(entries <-chan logentry.LogEntry) {
	for entry := range entries {
		if e.predicate.Applies(&entry) {
			e.outputLines = append(e.outputLines, fmt.Sprint(entry))
		}
	}
}

func (e *Echo) NumberOfLinesAfterSummarizeAsyncCompleted() int {
	e.waitGroup.Wait()
	return len(e.outputLines)
}

func (e *Echo) StringAfterSummarizeAsyncCompleted() string {
	e.waitGroup.Wait()
	return e.String()
}

func (e *Echo) String() string {
	var result stringList

	for _, line := range e.outputLines {
		result = result.Append(line)
	}

	return result.Join("\n")
}
