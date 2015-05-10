package config

import (
	"fmt"
	"github.com/fxnn/gowatch/logentry"
	"log"
	"strings"
	"time"
)

var PredicateTimeLayout string = time.RFC3339

func (config PredicateConfig) CreatePredicate() logentry.Predicate {
	predicates := make([]logentry.Predicate, 0, 3)

	for key, value := range config {
		switch strings.ToLower(key) {
		case "not":
			subPredicate := createPredicateOrFail(value)
			predicates = append(predicates, logentry.NotPredicate{subPredicate})
		case "allof":
			subPredicates := createSubPredicatesOrFail(value)
			predicates = append(predicates, logentry.AllOfPredicate{subPredicates})
		case "anyof":
			subPredicates := createSubPredicatesOrFail(value)
			predicates = append(predicates, logentry.AnyOfPredicate{subPredicates})
		case "noneof":
			subPredicates := createSubPredicatesOrFail(value)
			predicates = append(predicates, logentry.NoneOfPredicate{subPredicates})
		default:
			subPredicate := createPredicateForField(key, value)
			predicates = append(predicates, subPredicate)
		}
	}

	switch len(predicates) {
	case 0:
		return logentry.AcceptAllPredicate{}
	case 1:
		return predicates[0]
	default:
		return logentry.AllOfPredicate{predicates}
	}
}

func createPredicateForField(field string, predicateValue interface{}) logentry.Predicate {
	predicates := make([]logentry.Predicate, 0, 3)

	if predicateConfig, err := tryDecodeAsPredicateConfig(predicateValue); err != nil {
		log.Fatalf("No valid predicate configuration \"%s\" for field \"%s\", expected a map", predicateValue, field)
		return logentry.AcceptNothingPredicate{} // actually never executed
	} else {
		for keyValue, value := range predicateConfig {
			key := fmt.Sprint(keyValue)
			if stringValue, ok := value.(string); !ok {
				log.Fatalf("No valid value \"%s\" for predicate \"%s\" on field \"%s\", expected a string", value, key, field)
				return logentry.AcceptNothingPredicate{} // actually never executed
			} else {
				switch strings.ToLower(key) {
				case "is":
					switch strings.ToLower(stringValue) {
					case "empty":
						predicates = append(predicates, logentry.IsEmptyPredicate{field})
					case "not empty":
						predicates = append(predicates, logentry.IsNotEmptyPredicate{logentry.IsEmptyPredicate{field}})
					default:
						log.Fatalf("No valid predicate \"%s:%s\" for field \"%s\", expected \"empty\" or \"not empty\"", key, stringValue, field)
						return logentry.AcceptNothingPredicate{} // actually never executed
					}

				case "contains":
					predicates = append(predicates, logentry.ContainsPredicate{FieldName: field, ToBeContained: stringValue})
				case "matches":
					predicates = append(predicates, logentry.MatchesPredicate{FieldName: field, GrokPattern: stringValue})

				case "after":
					predicates = append(predicates, createAfterPredicate(field, PredicateTimeLayout, stringValue))
				case "before":
					predicates = append(predicates, createBeforePredicate(field, PredicateTimeLayout, stringValue))

				case "younger than":
					predicates = append(predicates, createYoungerThanPredicate(field, stringValue))

				default:
					log.Fatalf("No valid predicate \"%s\" for field \"%s\", expected \"is\", \"contains\", \"matches\" or \"after\"", key, field)
					return logentry.AcceptNothingPredicate{} // actually never executed
				}
			}
		}
	}

	switch len(predicates) {
	case 0:
		return logentry.AcceptAllPredicate{}
	case 1:
		return predicates[0]
	default:
		return logentry.AllOfPredicate{predicates}
	}
}

func createYoungerThanPredicate(fieldName string, value string) logentry.Predicate {
	duration, err := time.ParseDuration("-" + value)
	if err != nil {
		log.Fatalf("No valid duration \"%s\" to compare with field \"%s\": %s", value, fieldName, err.Error())
		return logentry.AcceptNothingPredicate{} // actually never executed
	}
	return logentry.AfterPredicate{FieldName: fieldName, EarlierTimestamp: time.Now().Add(duration)}
}

func createAfterPredicate(fieldName string, timeLayout string, value string) logentry.Predicate {
	operand, err := time.Parse(timeLayout, value)
	if err != nil {
		log.Fatalf("No valid timestamp \"%s\" to compare with field \"%s\", expected format\"%s\"", value, fieldName, timeLayout)
		return logentry.AcceptNothingPredicate{} // actually never executed
	}
	return logentry.AfterPredicate{FieldName: fieldName, EarlierTimestamp: operand}
}

func createBeforePredicate(fieldName string, timeLayout string, value string) logentry.Predicate {
	operand, err := time.Parse(timeLayout, value)
	if err != nil {
		log.Fatalf("No valid timestamp \"%s\" to compare with field \"%s\", expected format\"%s\"", value, fieldName, timeLayout)
		return logentry.AcceptNothingPredicate{} // actually never executed
	}
	return logentry.BeforePredicate{FieldName: fieldName, LaterTimestamp: operand}
}

func createSubPredicatesOrFail(value interface{}) []logentry.Predicate {
	if predicateConfigSlice, ok := value.([]PredicateConfig); ok {
		results := make([]logentry.Predicate, len(predicateConfigSlice))
		for i := 0; i < len(predicateConfigSlice); i++ {
			results[i] = predicateConfigSlice[i].CreatePredicate()
		}
		return results
	}

	if predicateConfig, err := tryDecodeAsPredicateConfig(value); err == nil {
		results := make([]logentry.Predicate, len(predicateConfig))
		i := 0
		// NOTE: Return map as single predicates to be able to implement allof, anyof, noneof logic
		for key, value := range predicateConfig {
			singlePredicateConfig := PredicateConfig(map[string]interface{}{key: value})
			results[i] = singlePredicateConfig.CreatePredicate()
			i++
		}
		return results
	} else {
		log.Fatal(err)
		return []logentry.Predicate{logentry.AcceptNothingPredicate{}} // actually never executed
	}
}

func createPredicateOrFail(value interface{}) logentry.Predicate {
	if predicateConfig, err := tryDecodeAsPredicateConfig(value); err == nil {
		return predicateConfig.CreatePredicate()
	} else {
		log.Fatal(err)
		return logentry.AcceptNothingPredicate{} // actually never executed
	}
}

func tryDecodeAsPredicateConfig(value interface{}) (PredicateConfig, error) {
	// TODO: Is it correct that we have to go this strange way? I had situations where the YAML parser returned me all three of them
	if predicateConfig, ok := value.(PredicateConfig); ok {
		return predicateConfig, nil
	}
	if predicateMap, ok := value.(map[interface{}]interface{}); ok {
		predicateConfig := make(PredicateConfig)
		for key, value := range predicateMap {
			predicateConfig[fmt.Sprint(key)] = value
		}
		return predicateConfig, nil
	}
	if predicateMap, ok := value.(map[string]interface{}); ok {
		predicateConfig := make(PredicateConfig)
		for key, value := range predicateMap {
			predicateConfig[key] = value
		}
		return predicateConfig, nil
	}

	return nil, fmt.Errorf("Invalid predicate configuration \"%s\", expected a map", value)
}
