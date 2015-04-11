package logentry

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type LogEntry struct {
	Timestamp time.Time
	Level     Level
	Tags      []string
	Message   string
	Host      string
	User      string
	Thread    string
	Process   string
	Custom    map[string]string
}

func New() (entry *LogEntry) {
	result := new(LogEntry)
	result.Tags = make([]string, 0, 10)
	result.Custom = make(map[string]string)
	return result
}

func (l *LogEntry) FieldValue(fieldName string) (reflect.Value, error) {
	entryValuePointer := reflect.ValueOf(l)
	entryValue := entryValuePointer.Elem()
	fieldValue := entryValue.FieldByName(fieldName)
	if fieldValue.IsValid() {
		return fieldValue, nil
	}
	if !isOnlyFirstLetterUpperCase(fieldName) {
		return l.FieldValue(onlyFirstLetterToUpper(fieldName))
	}
	return reflect.Value{}, errors.New("No valid field: " + fieldName)
}

func onlyFirstLetterToUpper(s string) string {
	result := ""
	if len(s) > 0 {
		// NOTE https://blog.golang.org/strings: s[0] would be wrong, as it gives bytes, not runes
		firstRune, width := utf8.DecodeRuneInString(s)
		result += string(unicode.ToUpper(firstRune))

		if len(s)-width > 0 {
			result += strings.ToLower(s[width:])
		}
	}
	return result
}

func isOnlyFirstLetterUpperCase(s string) bool {
	for i, w := 0, 0; i < len(s); i += w {
		// NOTE https://blog.golang.org/strings: s[i] would be wrong, as it gives bytes, not runes
		rune, width := utf8.DecodeRuneInString(s[i:])

		if i == 0 && !unicode.IsUpper(rune) {
			return false
		}
		if i > 0 && unicode.IsUpper(rune) {
			return false
		}

		w = width
	}
	return true
}

func (l *LogEntry) FieldAsString(fieldName string) (string, error) {
	field, err := l.FieldValue(fieldName)
	if err != nil {
		return "", err
	}

	switch {
	case l.isAnyStringField(field):
		return field.String(), err
	case l.isLogLevelField(field, fieldName):
		return l.Level.String(), err
	case l.isTagsField(field, fieldName):
		return fmt.Sprint(field.Interface()), err
	case l.isCustomField(field):
		return l.Custom[fieldName], err
	default:
		return "", errors.New("Cannot be read as logentry field: " + fieldName)
	}
}

func (l *LogEntry) AssignValue(fieldName string, value string) error {
	entryValuePointer := reflect.ValueOf(l)
	entryValue := entryValuePointer.Elem()
	field := entryValue.FieldByName(fieldName)

	switch {
	case l.isAnyStringField(field):
		field.SetString(value)
		return nil
	case l.isLogLevelField(field, fieldName):
		// TODO: Support wider range of Log Level names (cf. names defined in Grok)
		matchAsLogLevel, err := LevelFromString(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(matchAsLogLevel))
		return nil
	case l.isTagsField(field, fieldName):
		field.Set(reflect.Append(field, reflect.ValueOf(value)))
		return nil
	case l.isCustomField(field):
		l.Custom[fieldName] = value
		return nil
	default:
		return errors.New("Cannot be set as logentry field: " + fieldName)
	}
}

func (l *LogEntry) isLogLevelField(field reflect.Value, fieldName string) bool {
	return l.isAnyField(field) && field.Kind() == reflect.Int && fieldName == "Level"
}

func (l *LogEntry) isTagsField(field reflect.Value, fieldName string) bool {
	return l.isAnyField(field) && field.Kind() == reflect.Slice && fieldName == "Tags"
}

func (l *LogEntry) isAnyStringField(field reflect.Value) bool {
	return l.isAnyField(field) && field.Kind() == reflect.String
}

func (l *LogEntry) isCustomField(field reflect.Value) bool {
	return !l.isAnyField(field)
}

func (l *LogEntry) isAnyField(field reflect.Value) bool {
	return field.IsValid() && field.CanSet()
}
