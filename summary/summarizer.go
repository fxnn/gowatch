package summary

import "github.com/fxnn/gowatch/logentry"

type Summarizer interface {
	SummarizeAsync(entries <-chan logentry.LogEntry)
	StringAfterSummarizeAsyncCompleted() string
}
