package logentry

import (
	"errors"
	"strconv"
)

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

func (l Level) String() string {
	switch l {
	case TRACE:
		return "TRACE"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "Level=" + strconv.Itoa(int(l))
	}
}
