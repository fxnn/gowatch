package summary

import (
	"bytes"
	"github.com/fxnn/gowatch/logentry"
	"sync"
	"fmt"
)

type Echo struct {
	outputLines []string
	waitGroup sync.WaitGroup
}

func NewEcho() (e *Echo) {
	e = new(Echo)
	return
}

func (e *Echo) SummarizeAsync(entries <-chan logentry.LogEntry) {
	e.waitGroup.Add(1)
	go e.Summarize(entries)
}

func (e *Echo) Summarize(entries <-chan logentry.LogEntry) {
	for entry := range entries {
		e.outputLines = append(e.outputLines, fmt.Sprint(entry))
	}
	e.waitGroup.Done()
}

func (e *Echo) NumberOfLines() int {
	e.waitGroup.Wait()
	return len(e.outputLines)
}

func (e *Echo) String() string {
	var buffer bytes.Buffer

	e.waitGroup.Wait()
	for _, line := range e.outputLines {
		buffer.WriteString(line + "\n")
	}

	return buffer.String()
}
