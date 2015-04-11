package config

import (
	"fmt"
	"github.com/fxnn/gowatch/logentry"
	"github.com/fxnn/gowatch/parser"
	"log"
)

func (logfile *LogfileConfig) CreateParser(linesource parser.LineSource, predicate logentry.Predicate) parser.Parser {
	timeLayout := parseTimeLayout(logfile.TimeLayout)

	switch logfile.Parser {
	case "":
		return parser.NewSimpleParser(linesource, predicate)
	case "grok":
		if pattern, ok := logfile.Config["pattern"]; ok {
			return parser.NewGrokParser(linesource, fmt.Sprint(pattern), timeLayout, predicate)
		}
		log.Fatal("Grok parser used without pattern on logfile '", logfile.Filename, "'")
		return nil // actually never reached
	default:
		log.Fatal("Unrecognized parser '", logfile.Parser, "' on logfile '", logfile.Filename, "'")
		return nil // actually never reached
	}
}

func parseTimeLayout(givenTimeLayout string) string {
	if interpretedTimeLayout, ok := PredefinedTimeLayouts[givenTimeLayout]; ok {
		return interpretedTimeLayout
	}
	return givenTimeLayout
}
