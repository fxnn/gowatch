package summary

import "github.com/fxnn/gowatch/logentry"

type Multiplexer struct {
    summarizers []Summarizer
}

func NewMultiplexer() *Multiplexer {
    return &Multiplexer{make([]Summarizer, 0)}
}

func (m *Multiplexer) AddSummarizer(s Summarizer) {
    m.summarizers = append(m.summarizers, s)
}

func (m *Multiplexer) Summarize(entries <-chan logentry.LogEntry) {
    channels := make([]chan logentry.LogEntry, len(m.summarizers))
    for i, s := range m.summarizers {
        ch := make(chan logentry.LogEntry)
        go s.Summarize(ch)
        channels[i] = ch
    }

    for entry := range entries {
        for _, channel := range channels {
            channel <- entry
        }
    }

    for _, channel := range channels {
        close(channel)
    }
}
