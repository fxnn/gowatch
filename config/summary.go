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
	default:
		log.Fatalf("Unrecognized parser '%s'", summaryConfig.Do)
		return nil // actually never reached
	}
}
