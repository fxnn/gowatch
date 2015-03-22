package parser

import (
	"github.com/fxnn/gowatch/logentry"
)

type Parser interface {
	Parse() <-chan logentry.LogEntry
}

// used by this package for simplyfing the technical problem of "converting a channel of strings into a channel of
// LogEntries" down to "filling information from one string into one LogEntry".
func parse(linesource LineSource, lineToLogEntry func(line string, entry *logentry.LogEntry)) <-chan logentry.LogEntry {
	out := make(chan logentry.LogEntry)

	go func() {
		lines := linesource.Lines()
		for line := range lines {
			entry := logentry.New()
			lineToLogEntry(line, entry)
			out <- *entry
		}
		close(out)
	}()

	return out
}
