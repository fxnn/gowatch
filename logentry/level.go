package logentry

import "errors"

type Level int

const (
	TRACE Level = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

func LevelFromString(name string) (level Level, err error) {
	switch name {
	case "TRACE":
		level = TRACE
	case "DEBUG":
		level = DEBUG
	case "INFO":
		level = INFO
	case "WARNING":
		level = WARNING
	case "ERROR":
		level = ERROR
	case "FATAL":
		level = FATAL
	default:
		err = errors.New("Invalid log level: " + name)
	}

	return
}
