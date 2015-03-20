package mapper

import (
	"github.com/fxnn/gowatch/logentry"
)

type Mapper interface {
	Map(entries <-chan logentry.LogEntry) <-chan logentry.LogEntry
}
