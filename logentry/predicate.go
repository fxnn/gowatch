package logentry

import (
	"github.com/fxnn/grok"
	"reflect"
	"strings"
)

type Predicate interface {
	Applies(*LogEntry) bool
}

type AcceptAllPredicate struct{}

func (this *AcceptAllPredicate) Applies(*LogEntry) bool {
	return true
}

type ContainsPredicate struct {
	FieldName     string
	ToBeContained string
}

func (this *ContainsPredicate) Applies(logEntry *LogEntry) bool {
	stringValue, err := logEntry.FieldAsString(this.FieldName)
	if err == nil {
		return strings.Contains(stringValue, this.ToBeContained)
	}
	return false // in case of error, let's say it doesn't contain
}

type MatchesPredicate struct {
	FieldName   string
	GrokPattern string
	grok        *grok.Grok
}

func (this *MatchesPredicate) Applies(logEntry *LogEntry) bool {
	stringValue, err := logEntry.FieldAsString(this.FieldName)
	if err == nil {
		g := this.grok
		if g == nil {
			g = grok.New()
		}

		result, err := g.Match(this.GrokPattern, stringValue)
		if err == nil {
			return result
		}
	}
	return false // in case of error, let's say it doesn't contain
}

type IsEmptyPredicate struct{ FieldName string }

func (this *IsEmptyPredicate) Applies(logEntry *LogEntry) bool {
	fieldValue, err := logEntry.FieldValue(this.FieldName)
	if err == nil {
		if fieldValue.IsValid() {
			return isZero(fieldValue)
		}
		return logEntry.Custom[this.FieldName] == ""
	}

	return true // in case of error, let's say it's empty
}

type IsNotEmptyPredicate struct{ IsEmptyPredicate }

func (this *IsNotEmptyPredicate) Applies(logEntry *LogEntry) bool {
	return !this.IsEmptyPredicate.Applies(logEntry)
}

type NotPredicate struct{ SubPredicate Predicate }

func (this *NotPredicate) Applies(logEntry *LogEntry) bool {
	return !this.SubPredicate.Applies(logEntry)
}

type AllOfPredicate struct{ SubPredicates []Predicate }

func (this *AllOfPredicate) Applies(logEntry *LogEntry) bool {
	for _, p := range this.SubPredicates {
		if !p.Applies(logEntry) {
			return false
		}
	}

	return true
}

type AnyOfPredicate struct{ SubPredicates []Predicate }

func (this *AnyOfPredicate) Applies(logEntry *LogEntry) bool {
	for _, p := range this.SubPredicates {
		if p.Applies(logEntry) {
			return true
		}
	}

	return false
}

type NoneOfPredicate struct{ SubPredicates []Predicate }

func (this *NoneOfPredicate) Applies(logEntry *LogEntry) bool {
	for _, p := range this.SubPredicates {
		if p.Applies(logEntry) {
			return false
		}
	}

	return true
}

func isZero(v reflect.Value) bool {
	if v.Kind() == reflect.Slice {
		return v.Len() == 0
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
