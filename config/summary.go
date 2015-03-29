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
        default:
            log.Fatal("Unrecognized parser '" + summaryConfig.Summarizer)
            return nil, nil // actually never reached
    }
}