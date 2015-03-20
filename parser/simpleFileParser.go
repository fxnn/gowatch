package parser

import (
	"github.com/fxnn/gowatch/logentry"
)

type SimpleFileParser struct {
	FileParser
}

func NewSimpleFileParser(filename string) (p *SimpleFileParser) {
	p = new(SimpleFileParser)

	p.filename = filename
	p.logTextToEntryFunction = func(line string) (entry logentry.LogEntry) {
		entry = *logentry.New()
		entry.Message = line
		return
	}

	return
}
