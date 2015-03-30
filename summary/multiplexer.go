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

func (m *Multiplexer) SummarizeAsync(entries <-chan logentry.LogEntry) {
    channels := make([]chan logentry.LogEntry, len(m.summarizers))
    for i, s := range m.summarizers {
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