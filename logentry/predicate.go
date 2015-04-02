package logentry

import "reflect"

type Predicate interface {
	Applies(*LogEntry) bool
}

type IsEmptyPredicate struct{ FieldName string }

func (this *IsEmptyPredicate) Applies(logEntry *LogEntry) bool {
	fieldValue := reflect.ValueOf(logEntry).Elem().FieldByName(this.FieldName)
	if fieldValue.IsValid() {
		return isZero(fieldValue)
	}
	return logEntry.Custom[this.FieldName] == ""
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
