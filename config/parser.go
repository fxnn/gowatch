package config

import (
    "github.com/fxnn/gowatch/parser"
    "fmt"
    "log"
)

func (logfile *LogfileConfig) CreateParser(linesource parser.LineSource) parser.Parser {
    switch logfile.Parser {
        case "":
            return parser.NewSimpleParser(linesource)
        case "grok":
            if pattern, ok := logfile.Config["pattern"]; ok {
                return parser.NewGrokParser(linesource, fmt.Sprint(pattern))
            }
            log.Fatal("Grok parser used without pattern on logfile '", logfile.Filename, "'")
            return nil // actually never reached
        default:
            log.Fatal("Unrecognized parser '", logfile.Parser, "' on logfile '", logfile.Filename, "'")
            return nil // actually never reached
    }
}