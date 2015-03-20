package parser

import (
	"github.com/fxnn/gowatch/logentry"
)

type Parser interface {
	Parse() <-chan logentry.LogEntry
}
