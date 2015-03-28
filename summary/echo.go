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
	e.waitGroup.Add(1)
	return
}

func (tc *Echo) Summarize(entries <-chan logentry.LogEntry) {
	for entry := range entries {
		tc.outputLines = append(tc.outputLines, fmt.Sprint(entry))
	}
	tc.waitGroup.Done()
}

func (tc *Echo) String() string {
	var buffer bytes.Buffer

	tc.waitGroup.Wait()
	for _, line := range tc.outputLines {
		buffer.WriteString(line + "\n")
	}

	return buffer.String()
}
