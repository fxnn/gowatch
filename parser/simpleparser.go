package parser

import (
	"github.com/fxnn/gowatch/logentry"
)

type SimpleParser struct {
	linesource LineSource
}

func NewSimpleParser(linesource LineSource) (p *SimpleParser) {
	p = new(SimpleParser)

	p.linesource = linesource

	return
}

func (p *SimpleParser) Parse() <-chan logentry.LogEntry {
	return parse(p.linesource, p.lineToLogEntry)
}

func (p *SimpleParser) lineToLogEntry(line string, entry *logentry.LogEntry) {
	entry.Message = line
}
