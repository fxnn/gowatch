package config

import (
	"fmt"
	"github.com/fxnn/gowatch/summary"
	"log"
)

func (summaryConfig *SummaryConfig) CreateSummarizer() summary.Summarizer {
	predicate := summaryConfig.Where.CreatePredicate()

	switch summaryConfig.Summarizer {
	case "echo":
		return summary.NewEcho(predicate)
	case "tagcounter":
		return summary.NewTagCounterWithLocale("en_US", predicate)
	case "count", "counter", "grokcounter":
		patternsByName := make(map[string]string)
		for key, value := range summaryConfig.Config {
			if pattern, ok := value.(string); ok {
				patternsByName[fmt.Sprint(key)] = pattern
			} else {
				log.Fatalf("Unrecognized option for summarizer '%s': key '%s' with value '%s'", summaryConfig.Summarizer, key, value)
			}
		}

		return summary.NewGrokCounterWithLocale("en_US", patternsByName, predicate)
	default:
		log.Fatalf("Unrecognized parser '%s'", summaryConfig.Summarizer)
		return nil // actually never reached
	}
}
