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
	matches := p.regexp.FindStringSubmatch(line)
	for matchNumber, matchedContent := range matches {
		if matchName, ok := p.submatchNameMap[matchNumber]; ok {
			p.matchToLogEntry(matchName, matchedContent, entry)
		}
	}
}

func (p *RegexpParser) matchToLogEntry(matchName string, matchedContent string, entry *logentry.LogEntry) {
	pointerToEntryValue := reflect.ValueOf(entry)
	entryValue := pointerToEntryValue.Elem()

	field := entryValue.FieldByName(matchName)
	switch {
	case isAnyStringField(field):
		field.SetString(matchedContent)
	case isLogLevelField(matchName, field):
		matchAsLogLevel, err := logentry.LevelFromString(matchedContent)
		if err == nil {
			field.SetInt(int64(matchAsLogLevel))
		} else {
			log.Print("Invalid log level [" + matchedContent + "]: ", err)
		}
	case isTagsField(matchName, field):
		field.Set(reflect.Append(field, reflect.ValueOf(matchedContent)))
	case isCustomField(field):
		entry.Custom[matchName] = matchedContent
	default:
		log.Print("Field not supported as Regexp target: ", matchName)
	}
}

func isLogLevelField(fieldName string, field reflect.Value) bool {
	return isAnyField(field) && field.Kind() == reflect.Int && fieldName == "Level"
}

func isTagsField(fieldName string, field reflect.Value) bool {
	return isAnyField(field) && field.Kind() == reflect.Slice && fieldName == "Tags"
}

func isAnyStringField(field reflect.Value) bool {
	return isAnyField(field) && field.Kind() == reflect.String
}

func isCustomField(field reflect.Value) bool {
	return !isAnyField(field)
}

func isAnyField(field reflect.Value) bool {
	return field.IsValid() && field.CanSet()
}