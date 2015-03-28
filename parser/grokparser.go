package parser

import (
    "github.com/gemsi/grok"
    "github.com/fxnn/gowatch/logentry"
    "log"
)

type GrokParser struct {
    linesource      LineSource
    grok            *grok.Grok
    pattern         string
}

func NewGrokParser(linesource LineSource, pattern string) (p *GrokParser) {
    return &GrokParser{linesource:linesource, grok:grok.New(), pattern:pattern}
}

func (p *GrokParser) AddPattern(name string, pattern string) {
    p.grok.AddPattern(name, pattern)
}

func (p *GrokParser) Parse() <-chan logentry.LogEntry {
    return parse(p.linesource, p.lineToLogEntry)
}

func (p *GrokParser) lineToLogEntry(line string, entry *logentry.LogEntry) {
    matches, err := p.grok.Parse(p.pattern, line)
    if err != nil {
        log.Print(err)
        return
    }

    for field, value := range matches {
        if err := entry.AssignValue(field, value); err != nil {
            log.Print(err)
        }
    }
}
