package summary

import "github.com/fxnn/gowatch/logentry"

type Summarizer interface {
	Summarize(entries <-chan logentry.LogEntry)
	String() string
}
