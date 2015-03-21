package mapper

import (
	"github.com/fxnn/gowatch/config"
	"github.com/fxnn/gowatch/logentry"
)

type ConfigurationBasedMapper struct {
	config.GowatchLogfile
}

func NewConfigurationBasedMapper(config *config.GowatchLogfile) (m *ConfigurationBasedMapper) {
	m = new(ConfigurationBasedMapper)
	m.GowatchLogfile = *config
	return
}

func (m *ConfigurationBasedMapper) Map(entries <-chan logentry.LogEntry) <-chan logentry.LogEntry {
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
