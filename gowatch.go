package main

import (
	"code.google.com/p/getopt"
	"fmt"
	"github.com/fxnn/gowatch/config"
	"github.com/fxnn/gowatch/mapper"
	"github.com/fxnn/gowatch/parser"
	"github.com/fxnn/gowatch/summary"
	"log"
)

func main() {
	configFilePath := getopt.StringLong("config", 'c', "", "Path to configuration file")
	getopt.Parse()

	if !getopt.Lookup("config").Seen() {
		log.Fatal("No configuration file specified")
	}

	config := config.ReadConfigByFilename(*configFilePath)

	summarizers := summary.NewMultiplexer()
	summaryStringers := make([]fmt.Stringer, len(config.Summary))
	for i, summaryConfig := range config.Summary {
		summarizer, stringer := summaryConfig.CreateSummarizer()
		summarizers.AddSummarizer(summarizer)
		summaryStringers[i] = stringer
	}

	for _, logfile := range config.Logfiles {
		linesource := parser.NewFileLineSource(logfile.Filename)
		parser := logfile.CreateParser(linesource)
		entries := parser.Parse()

		logfileMapper := mapper.NewConfigurationBasedMapper(logfile)
		mappedEntries := logfileMapper.Map(entries)

		summarizers.Summarize(mappedEntries)
	}

	for _, summaryStringer := range summaryStringers {
		fmt.Print(summaryStringer.String())
	}
}


