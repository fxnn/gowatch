package main

import (
	"code.google.com/p/getopt"
	"fmt"
	"github.com/fxnn/gowatch/config"
	"github.com/fxnn/gowatch/mapper"
	"github.com/fxnn/gowatch/parser"
	"github.com/fxnn/gowatch/summary"
	"log"
	"os"
)

func main() {
	configFilePath := getopt.StringLong("config", 'c', "", "Path to configuration file")
	getopt.Parse()

	if !getopt.Lookup("config").Seen() {
		log.Fatal("No configuration file specified")
		os.Exit(1)
	}

	config := config.ReadConfigByFilename(*configFilePath)

	for i := range config.Logfiles {
		logfile := &config.Logfiles[i]
		parser := parser.NewSimpleFileParser(logfile.Filename)
		entries := parser.Parse()

		logfileMapper := mapper.NewConfigurationBasedMapper(logfile)
		mappedEntries := logfileMapper.Map(entries)

		summarizer := summary.NewTagCounter()
		summarizer.Summarize(mappedEntries)

		fmt.Print(summarizer.String())
	}
}
