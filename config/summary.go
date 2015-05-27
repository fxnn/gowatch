package config

import (
	"fmt"
	"github.com/fxnn/gowatch/summary"
	"log"
)

func (summaryConfig *SummaryConfig) CreateSummarizer() summary.Summarizer {
	predicate := summaryConfig.Where.CreatePredicate()

	switch summaryConfig.Do {
	case "echo":
		return summary.NewEcho(predicate)

	case "tagcounter", "counttags", "count tags":
		return summary.NewTagCounterWithLocale("en_US", predicate)

	case "count", "counter", "grokcounter":
		patternsByName := make(map[string]string)
		for key, value := range summaryConfig.With {
			if pattern, ok := value.(string); ok {
				patternsByName[fmt.Sprint(key)] = pattern
			} else {
				log.Fatalf("Unrecognized option for summarizer '%s': key '%s' with value '%s'", summaryConfig.Do, key, value)
			}
		}

		return summary.NewGrokCounterWithLocale("en_US", patternsByName, predicate)

	case "list", "lister", "groklister":
		patternsByListValueByName := make(map[string]map[string]string)
		for key, value := range summaryConfig.With {
			if patternsByListValues, ok := value.(map[interface{}]interface{}); ok {
				patternsByListValueByName[fmt.Sprint(key)] = make(map[string]string)
				for listValue, patternName := range patternsByListValues {
					patternsByListValueByName[fmt.Sprint(key)][fmt.Sprint(listValue)] = fmt.Sprint(patternName)
				}
			} else {
				log.Fatalf("Unrecognized option for summarizer '%s': key '%s' with value '%s'", summaryConfig.Do, key, value)
			}
		}

		return summary.NewGrokListerWithLocale("en_US", patternsByListValueByName, predicate)

	default:
		log.Fatalf("Unrecognized parser '%s'", summaryConfig.Do)
		return nil // actually never reached
	}
}
