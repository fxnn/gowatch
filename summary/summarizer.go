package summary

import "github.com/fxnn/gowatch/logentry"

// TODO: Currently, we're only specifying the input function. Later on, we could demand to return an output object.
type Summarizer interface {
	// The deal is that before it provides any summary, the Summarizer waits for the most recently given channel being
	// closed.
	SummarizeAsync(entries <-chan logentry.LogEntry)
}
