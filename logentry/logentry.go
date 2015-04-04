package logentry

import (
	"errors"
	"fmt"
	"reflect"
	"time"
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
	return entryValue.FieldByName(fieldName), nil
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
