package config

import (
    "log"
    "github.com/fxnn/gowatch/summary"
    "fmt"
)

func (summaryConfig *SummaryConfig) CreateSummarizer() (summary.Summarizer, fmt.Stringer) {
    var title string
    if len(summaryConfig.Title) > 0 {
        title = summaryConfig.Title
    } else {
        title = summaryConfig.Summarizer
    }

    switch summaryConfig.Summarizer {
        case "echo":
            echo := summary.NewEcho()
            return echo, &summary.TitledStringer{title, echo}
        case "tagcounter":
            tagcounter := summary.NewTagCounter()
            return tagcounter, &summary.TitledStringer{title, tagcounter}
        case "grokcounter":
            patternsByName := make(map[string]string)
            for key, value := range summaryConfig.Config {
                if pattern, ok := value.(string); ok {
                    patternsByName[fmt.Sprint(key)] = pattern
                } else {
                    log.Fatalf("Unrecognized option for summarizer '%s': key '%s' with value '%s'", summaryConfig.Summarizer, key, value)
                }
            }

            grokcounter := summary.NewGrokCounter(patternsByName)
            return grokcounter, &summary.TitledStringer{title, grokcounter}
        default:
            log.Fatalf("Unrecognized parser '%s'", summaryConfig.Summarizer)
            return nil, nil // actually never reached
    }
}