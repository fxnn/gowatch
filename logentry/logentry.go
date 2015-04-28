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

func (l *LogEntry) FieldValue(fieldName string) (fieldValue reflect.Value, actualFieldName string) {
	exactFieldValue := l.exactFieldValue(fieldName)

	if !exactFieldValue.IsValid() && !isFormattedLikeGoStructPublicField(fieldName) {
		alternativeFieldName := formatLikeGoStructPublicField(fieldName)
		alternativeFieldValue := l.exactFieldValue(alternativeFieldName)
		if alternativeFieldValue.IsValid() {
			return alternativeFieldValue, alternativeFieldName
		}
	}

	return exactFieldValue, fieldName
}

func (l *LogEntry) exactFieldValue(fieldName string) reflect.Value {
	entryValuePointer := reflect.ValueOf(l)
	entryValue := entryValuePointer.Elem()
	return entryValue.FieldByName(fieldName)
}

func formatLikeGoStructPublicField(s string) string {
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

func isFormattedLikeGoStructPublicField(s string) bool {
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

func (l *LogEntry) FieldAsTime(fieldName string) (time.Time, error) {
	if l.IsTimestamp(fieldName) {
		return l.Timestamp, nil
	}

	return time.Time{}, errors.New("No timestamp field: " + fieldName)
}

func (l *LogEntry) FieldAsString(fieldName string) (string, error) {
	field, actualFieldName := l.FieldValue(fieldName)

	switch {
	case l.isAnyStringField(field):
		return field.String(), nil
	case l.isLogLevelField(field, actualFieldName):
		return l.Level.String(), nil
	case l.isTagsField(field, actualFieldName):
		return fmt.Sprint(field.Interface()), nil
	case l.isCustomField(field):
		// HINT: Don't use actualFieldName here, but the one the user intended
		return l.Custom[fieldName], nil
	default:
		return "", errors.New("Cannot be read as logentry field: " + fieldName)
	}
}

func (l *LogEntry) AssignValue(fieldName string, value string) error {
	field, actualFieldName := l.FieldValue(fieldName)

	switch {
	case l.isAnyStringField(field):
		field.SetString(value)
		return nil
	case l.isLogLevelField(field, actualFieldName):
		// TODO: Support wider range of Log Level names (cf. names defined in Grok)
		matchAsLogLevel, err := LevelFromString(value)
		if err != nil {
			return err
		}
		field.SetInt(int64(matchAsLogLevel))
		return nil
	case l.isTagsField(field, actualFieldName):
		field.Set(reflect.Append(field, reflect.ValueOf(value)))
		return nil
	case l.isCustomField(field):
		// HINT: Don't use actualFieldName here, but the one the user intended
		l.Custom[fieldName] = value
		return nil
	default:
		return fmt.Errorf("Cannot be set as logentry field: \"%s\" with value \"%s\"", fieldName, value)
	}
}

func (l *LogEntry) IsTimestamp(fieldName string) bool {
	field, actualFieldName := l.FieldValue(fieldName)
	return l.isTimestampField(field, actualFieldName)
}

func (l *LogEntry) isTimestampField(field reflect.Value, fieldName string) bool {
	return l.isAnyField(field) && fieldName == "Timestamp"
}

func (l *LogEntry) isLogLevelField(field reflect.Value, fieldName string) bool {
	return l.isAnyField(field) && fieldName == "Level"
}

func (l *LogEntry) IsTags(fieldName string) bool {
	field, actualFieldName := l.FieldValue(fieldName)
	return l.isTagsField(field, actualFieldName)
}

func (l *LogEntry) isTagsField(field reflect.Value, fieldName string) bool {
	return l.isAnyField(field) && fieldName == "Tags"
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
