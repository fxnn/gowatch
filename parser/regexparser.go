package parser

import (
	"github.com/fxnn/gowatch/logentry"
	"log"
	"reflect"
	"regexp"
)

type RegexpParser struct {
	linesource      LineSource
	regexp          *regexp.Regexp
	submatchNameMap map[int]string
}

func NewRegexpParser(linesource LineSource, regexp *regexp.Regexp, submatchNameMap map[int]string) (p *RegexpParser) {
	return &RegexpParser{linesource, regexp, submatchNameMap}
}

func (p *RegexpParser) Parse() <-chan logentry.LogEntry {
	return parse(p.linesource, p.lineToLogEntry)
}

func (p *RegexpParser) lineToLogEntry(line string, entry *logentry.LogEntry) {
	pointerToEntryValue := reflect.ValueOf(entry)
	entryValue := pointerToEntryValue.Elem()

	matches := p.regexp.FindStringSubmatch(line)
	for matchNumber, match := range matches {
		if name, ok := p.submatchNameMap[matchNumber]; ok {
			field := entryValue.FieldByName(name)
			if field.IsValid() && field.CanSet() {
				if field.Kind() == reflect.String {
					field.SetString(match)
				} else if field.Kind() == reflect.Int && name == "Level" {
					matchAsLogLevel, err := logentry.LevelFromString(match)
					if err != nil {
						log.Print(err)
					} else {
						field.SetInt(int64(matchAsLogLevel))
					}
				} else if field.Kind() == reflect.Slice { // esp. for Tags
					field.Set(reflect.Append(field, reflect.ValueOf(match)))
				}
			} else {
				entry.Custom[name] = match
			}
		}
	}
}
