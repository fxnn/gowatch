package logentry

import (
	"github.com/fxnn/grok"
	"reflect"
	"strings"
	"time"
)

type Predicate interface {
	Applies(*LogEntry) bool
}

type AcceptNothingPredicate struct{}

func (this AcceptNothingPredicate) Applies(*LogEntry) bool {
	return false
}

type AcceptAllPredicate struct{}

func (this AcceptAllPredicate) Applies(*LogEntry) bool {
	return true
}

type ContainsPredicate struct {
	FieldName     string
	ToBeContained string
}

func (this ContainsPredicate) Applies(logEntry *LogEntry) bool {
	if logEntry.IsTags(this.FieldName) {
		for _, tag := range logEntry.Tags {
			if tag == this.ToBeContained {
				return true
			}
		}
	} else {
		stringValue, err := logEntry.FieldAsString(this.FieldName)
		if err == nil {
			return strings.Contains(stringValue, this.ToBeContained)
		}
		// in case of error, let's say it doesn't contain
	}
	return false
}

type MatchesPredicate struct {
	FieldName   string
	GrokPattern string
	grok        *grok.Grok
}

func (this MatchesPredicate) Applies(logEntry *LogEntry) bool {
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

type AfterPredicate struct {
	FieldName        string
	EarlierTimestamp time.Time
}

func (this AfterPredicate) Applies(logEntry *LogEntry) bool {
	time, err := logEntry.FieldAsTime(this.FieldName)
	if err == nil {
		return time.After(this.EarlierTimestamp)
	}
	return false // in case of error, let's say it doesn't contain
}

type BeforePredicate struct {
	FieldName      string
	LaterTimestamp time.Time
}

func (this BeforePredicate) Applies(logEntry *LogEntry) bool {
	time, err := logEntry.FieldAsTime(this.FieldName)
	if err == nil {
		return time.Before(this.LaterTimestamp)
	}
	return false // in case of error, let's say it doesn't contain
}

type IsEmptyPredicate struct{ FieldName string }

func (this IsEmptyPredicate) Applies(logEntry *LogEntry) bool {
	fieldValue, _ := logEntry.FieldValue(this.FieldName)
	if fieldValue.IsValid() {
		return isZero(fieldValue)
	}

	return logEntry.Custom[this.FieldName] == ""
}

type IsNotEmptyPredicate struct{ IsEmptyPredicate }

func (this IsNotEmptyPredicate) Applies(logEntry *LogEntry) bool {
	return !this.IsEmptyPredicate.Applies(logEntry)
}

type NotPredicate struct{ SubPredicate Predicate }

func (this NotPredicate) Applies(logEntry *LogEntry) bool {
	return !this.SubPredicate.Applies(logEntry)
}

type AllOfPredicate struct{ SubPredicates []Predicate }

func (this AllOfPredicate) Applies(logEntry *LogEntry) bool {
	for _, p := range this.SubPredicates {
		if !p.Applies(logEntry) {
			return false
		}
	}

	return true
}

type AnyOfPredicate struct{ SubPredicates []Predicate }

func (this AnyOfPredicate) Applies(logEntry *LogEntry) bool {
	for _, p := range this.SubPredicates {
		if p.Applies(logEntry) {
			return true
		}
	}

	return false
}

type NoneOfPredicate struct{ SubPredicates []Predicate }

func (this NoneOfPredicate) Applies(logEntry *LogEntry) bool {
	for _, p := range this.SubPredicates {
		if p.Applies(logEntry) {
			return false
		}
	}

	return true
}

func isZero(v reflect.Value) bool {
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Map {
		return v.Len() == 0
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
