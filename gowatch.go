package main

import (
	"code.google.com/p/getopt"
	"fmt"
	"github.com/fxnn/gowatch/config"
	"github.com/fxnn/gowatch/mapper"
	"github.com/fxnn/gowatch/parser"
	"github.com/fxnn/gowatch/summary"
	"log"
	"bytes"
	"strings"
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
		summarizer, stringer := createSummarizer(summaryConfig)
		summarizers.AddSummarizer(summarizer)
		summaryStringers[i] = stringer
	}

	for _, logfile := range config.Logfiles {
		linesource := parser.NewFileLineSource(logfile.Filename)
		parser := createParser(logfile, linesource)
		entries := parser.Parse()

		logfileMapper := mapper.NewConfigurationBasedMapper(logfile)
		mappedEntries := logfileMapper.Map(entries)

		summarizers.Summarize(mappedEntries)
	}

	for _, summaryStringer := range summaryStringers {
		fmt.Print(summaryStringer.String())
	}
}

func createParser(logfile config.LogfileConfig, linesource parser.LineSource) parser.Parser {
	switch logfile.Parser {
	case "":
		return parser.NewSimpleParser(linesource)
	case "grok":
		if pattern, ok := logfile.Config["pattern"]; ok {
			return parser.NewGrokParser(linesource, fmt.Sprint(pattern))
		} else {
			log.Fatal("Grok parser used without pattern on logfile '", logfile.Filename, "'")
		}
	case "regexp":
		if pattern, ok := logfile.Config["pattern"]; ok {
			// TODO: implement that map or remove the whole regexp parser
			if parser, err := parser.NewRegexpParser(linesource, fmt.Sprint(pattern), make(map[int]string)); err == nil {
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

func createSummarizer(summaryConfig config.SummaryConfig) (summary.Summarizer, fmt.Stringer) {
	var title string
	if len(summaryConfig.Title) > 0 {
		title = summaryConfig.Title
	} else {
		title = summaryConfig.Summarizer
	}

	switch summaryConfig.Summarizer {
		case "echo":
			echo := summary.NewEcho()
			return echo, &TitledStringer{title, echo}
		case "tagcounter":
			tagcounter := summary.NewTagCounter()
			return tagcounter, &TitledStringer{title, tagcounter}
		default:
			log.Fatal("Unrecognized parser '" + summaryConfig.Summarizer)
	}

	return nil, nil // actually never reached
}

type TitledStringer struct{
	title		string
	stringer	fmt.Stringer
}

func (s *TitledStringer) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(s.title)
	buffer.WriteString("\n")
	buffer.WriteString(strings.Repeat("=", len(s.title)))
	buffer.WriteString("\n")

	buffer.WriteString(s.stringer.String())
	buffer.WriteString("\n\n")

	return buffer.String()
}