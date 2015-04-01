package config

import (
	"fmt"
	"github.com/fxnn/gowatch/summary"
	"log"
)

func (summaryConfig *SummaryConfig) CreateSummarizer() summary.Summarizer {
	switch summaryConfig.Summarizer {
	case "echo":
		return summary.NewEcho()
	case "tagcounter":
		return summary.NewTagCounter()
	case "grokcounter":
		patternsByName := make(map[string]string)
		for key, value := range summaryConfig.Config {
			if pattern, ok := value.(string); ok {
				patternsByName[fmt.Sprint(key)] = pattern
			} else {
				log.Fatalf("Unrecognized option for summarizer '%s': key '%s' with value '%s'", summaryConfig.Summarizer, key, value)
			}
		}

		return summary.NewGrokCounter(patternsByName)
	default:
		log.Fatalf("Unrecognized parser '%s'", summaryConfig.Summarizer)
		return nil // actually never reached
	}
}
