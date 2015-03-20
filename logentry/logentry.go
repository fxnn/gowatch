package logentry

import "time"

type LogEntry struct {
	Timestamp time.Time
	Level     Level
	Tags      []string
	Message   string
	Host      string
	User      string
	Thread    string
	Process   string
}

func New() (entry *LogEntry) {
	return new(LogEntry)
}
