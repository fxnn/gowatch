package main

import (
	"code.google.com/p/getopt"
	"fmt"
	"github.com/fxnn/gowatch/config"
	"github.com/fxnn/gowatch/mapper"
	"github.com/fxnn/gowatch/parser"
	"github.com/fxnn/gowatch/summary"
	"log"
	"strings"
)

func main() {
	configFilePath := getopt.StringLong("config", 'c', "", "Path to configuration file", "/path/to/config.yml")
	getopt.Parse()

	if !getopt.Lookup("config").Seen() {
		log.Fatal("No configuration file given. Specify one using `-c /path/to/config.yml`")
	}

	config := config.ReadConfigByFilename(*configFilePath)

	multiplexer := summary.NewMultiplexer()
	summarizerTitles := make([]string, len(config.Summary))
	for i, summaryConfig := range config.Summary {
		multiplexer.AddSummarizer(summaryConfig.CreateSummarizer())
		if summaryConfig.Title != "" {
			summarizerTitles[i] = summaryConfig.Title
		} else {
			summarizerTitles[i] = summaryConfig.Summarizer
		}
	}

	for _, logfile := range config.Logfiles {
		linesource := parser.NewFileLineSource(logfile.Filename)
		parser := logfile.CreateParser(linesource, logfile.Where.CreatePredicate())
		entries := parser.Parse()

		logfileMapper := mapper.NewConfigurationBasedMapper(logfile)
		mappedEntries := logfileMapper.Map(entries)

		multiplexer.SummarizeAsync(mappedEntries)
	}

	for i, summarizer := range multiplexer.Summarizers {
		title := summarizerTitles[i]
		fmt.Printf("%s\n", title)
		fmt.Printf("%s\n", strings.Repeat("=", len(title)))
		fmt.Printf("%s\n\n", summarizer.StringAfterSummarizeAsyncCompleted())
	}
}
