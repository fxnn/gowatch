package parser

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/gemsi/grok"
	"log"
)

type GrokParser struct {
	linesource LineSource
	grok       *grok.Grok
	pattern    string
	predicate  logentry.Predicate
}

func NewGrokParser(linesource LineSource, pattern string, predicate logentry.Predicate) (p *GrokParser) {
	return &GrokParser{linesource: linesource, grok: grok.New(), pattern: pattern, predicate: predicate}
}

func (p *GrokParser) AddPattern(name string, pattern string) {
	p.grok.AddPattern(name, pattern)
}

func (p *GrokParser) Parse() <-chan logentry.LogEntry {
	return parse(p.linesource, p.predicate, p.lineToLogEntry)
}

func (p *GrokParser) lineToLogEntry(line string, entry *logentry.LogEntry) {
	matches, err := p.grok.ParseToMultiMap(p.pattern, line)
	if err != nil {
		log.Print(err)
		return
	}

	for field, values := range matches {
		for _, value := range values {
			if err := entry.AssignValue(field, value); err != nil {
				log.Print(err)
			}
		}
	}
}
