package mapper

import (
	"github.com/fxnn/gowatch/config"
	"github.com/fxnn/gowatch/logentry"
)

type LogfileConfigMapper struct {
	config.GowatchLogfile
}

func NewLogfileConfigMapper(config *config.GowatchLogfile) (m *LogfileConfigMapper) {
	m = new(LogfileConfigMapper)
	m.GowatchLogfile = *config
	return
}

func (m *LogfileConfigMapper) Map(entries <-chan logentry.LogEntry) <-chan logentry.LogEntry {
	out := make(chan logentry.LogEntry)

	go func() {
		for entry := range entries {
			var mappedEntry = entry
			mappedEntry.Tags = m.Tags
			out <- mappedEntry
		}
		close(out)
	}()

	return out
}
