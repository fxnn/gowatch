package config

import (
	"github.com/fxnn/gowatch/logentry"
	"log"
	"reflect"
	"strings"
)

func (config *PredicateConfig) CreatePredicate() logentry.Predicate {
	switch {
	case config.exactlySet("Not"):
		return &logentry.NotPredicate{config.Not.CreatePredicate()}

	case config.exactlySet("AllOf"):
		predicates := make([]logentry.Predicate, len(config.AllOf))
		for i, subConfig := range config.AllOf {
			predicates[i] = subConfig.CreatePredicate()
		}
		return &logentry.AllOfPredicate{predicates}

	case config.exactlySet("AnyOf"):
		predicates := make([]logentry.Predicate, len(config.AnyOf))
		for i, subConfig := range config.AnyOf {
			predicates[i] = subConfig.CreatePredicate()
		}
		return &logentry.AnyOfPredicate{predicates}

	case config.exactlySet("NoneOf"):
		predicates := make([]logentry.Predicate, len(config.NoneOf))
		for i, subConfig := range config.NoneOf {
			predicates[i] = subConfig.CreatePredicate()
		}
		return &logentry.NoneOfPredicate{predicates}

	case config.exactlySet("Field", "Is") && strings.ToLower(config.Is) == "empty":
		return &logentry.IsEmptyPredicate{config.Field}

	case config.exactlySet("Field", "Is") && strings.ToLower(config.Is) == "not empty":
		return &logentry.IsNotEmptyPredicate{logentry.IsEmptyPredicate{config.Field}}

	case config.exactlySet("Field", "Contains"):
		return nil // TODO
	case config.exactlySet("Field", "Matches"):
		return nil // TODO
	default:
		log.Fatalf("Predicate configuration not allowed: %s", config)
		return nil // actually never executed
	}
}

func (config *PredicateConfig) exactlySet(fieldNames ...string) bool {
	for _, fieldName := range fieldNames {
		if !config.fieldSet(fieldName) {
			return false
		}
	}

	return config.noFieldsSetExcept(fieldNames...)
}

func (config *PredicateConfig) noFieldsSetExcept(allowedFieldNames ...string) bool {
	structValue := reflect.ValueOf(config).Elem()

	allowedFieldsSet := 0
	for _, allowedFieldName := range allowedFieldNames {
		if config.fieldSet(allowedFieldName) {
			allowedFieldsSet++
		}
	}

	fieldsSet := 0
	for i := 0; i < structValue.NumField(); i++ {
		if !isZero(structValue.Field(i)) {
			fieldsSet++
		}
	}

	return allowedFieldsSet == fieldsSet
}

func (config *PredicateConfig) fieldSet(fieldName string) bool {
	fieldValue := reflect.ValueOf(config).Elem().FieldByName(fieldName)
	return !isZero(fieldValue) // means: is not the zero value for that type
}

func isZero(v reflect.Value) bool {
	if v.Kind() == reflect.Slice {
		return v.Len() == 0
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
