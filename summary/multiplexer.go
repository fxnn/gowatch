package summary

import (
	"bytes"
	"github.com/fxnn/gowatch/logentry"
)

type Multiplexer struct {
	Summarizers []Summarizer
}

func NewMultiplexer() *Multiplexer {
	return &Multiplexer{make([]Summarizer, 0)}
}

func (m *Multiplexer) AddSummarizer(s Summarizer) {
	m.Summarizers = append(m.Summarizers, s)
}

func (m *Multiplexer) SummarizeAsync(entries <-chan logentry.LogEntry) {
	channels := make([]chan logentry.LogEntry, len(m.Summarizers))
	for i, s := range m.Summarizers {
		ch := make(chan logentry.LogEntry)
		s.SummarizeAsync(ch)
		channels[i] = ch
	}

	go m.Summarize(channels, entries)
}

func (m *Multiplexer) Summarize(channels []chan logentry.LogEntry, entries <-chan logentry.LogEntry) {
	for entry := range entries {
		for _, channel := range channels {
			channel <- entry
		}
	}

	for _, channel := range channels {
		close(channel)
	}
}

func (m *Multiplexer) StringAfterSummarizeAsyncCompleted() string {
	result := new(bytes.Buffer)
	for _, summarizer := range m.Summarizers {
		result.WriteString(summarizer.StringAfterSummarizeAsyncCompleted())
	}
	return result.String()
}
