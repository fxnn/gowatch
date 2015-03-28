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

	for i := range config.Logfiles {
		logfile := config.Logfiles[i]
		linesource := parser.NewFileLineSource(logfile.Filename)
		parser := createParser(logfile, linesource)
		entries := parser.Parse()

		logfileMapper := mapper.NewConfigurationBasedMapper(logfile)
		mappedEntries := logfileMapper.Map(entries)

		summarizer := summary.NewTagCounter()
//		summarizer := summary.NewEcho()
		summarizer.Summarize(mappedEntries)

		fmt.Print(summarizer.String())
	}
}

func createParser(logfile config.GowatchLogfile, linesource parser.LineSource) parser.Parser {
	switch logfile.Parser {
	case "":
		return parser.NewSimpleParser(linesource)
	case "grok":
		if pattern, ok := logfile.ParserConfig["pattern"]; ok {
			return parser.NewGrokParser(linesource, pattern)
		} else {
			log.Fatal("Grok parser used without pattern on logfile '", logfile.Filename, "'")
		}
	case "regexp":
		if pattern, ok := logfile.ParserConfig["pattern"]; ok {
			// TODO: implement that map or remove the whole regexp parser
			if parser, err := parser.NewRegexpParser(linesource, pattern, make(map[int]string)); err == nil {
				return parser
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Regexp parser used without pattern on logfile '", logfile.Filename, "'")
		}
	default:
		log.Fatal("Unrecognized parser '", logfile.Parser, "' on logfile '", logfile.Filename, "'")
	}
	return nil // actually never reached
}