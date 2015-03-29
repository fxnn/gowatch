package mapper

import (
	"github.com/fxnn/gowatch/config"
	"github.com/fxnn/gowatch/logentry"
)

type ConfigurationBasedMapper struct {
	config.LogfileConfig
}

func NewConfigurationBasedMapper(config config.LogfileConfig) (m *ConfigurationBasedMapper) {
	m = new(ConfigurationBasedMapper)
	m.LogfileConfig = config
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
