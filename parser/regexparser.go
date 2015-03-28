package parser

import (
    "github.com/fxnn/gowatch/logentry"
    "log"
    "regexp"
)

type RegexpParser struct {
    linesource      LineSource
    regexp          *regexp.Regexp
    submatchNameMap map[int]string
}

func NewRegexpParser(linesource LineSource, pattern string, submatchNameMap map[int]string) (p *RegexpParser, err error) {
    regexp, err := regexp.Compile(pattern)
    return &RegexpParser{linesource, regexp, submatchNameMap}, err
}

func (p *RegexpParser) Parse() <-chan logentry.LogEntry {
    return parse(p.linesource, p.lineToLogEntry)
}

func (p *RegexpParser) lineToLogEntry(line string, entry *logentry.LogEntry) {
    matches := p.regexp.FindStringSubmatch(line)
    for matchNumber, matchedContent := range matches {
        if matchName, ok := p.submatchNameMap[matchNumber]; ok {
            p.matchToLogEntry(matchName, matchedContent, entry)
        }
    }
}

func (p *RegexpParser) matchToLogEntry(matchName string, matchedContent string, entry *logentry.LogEntry) {
    if err := entry.AssignValue(matchName, matchedContent); err != nil {
        log.Print(err)
    }
}