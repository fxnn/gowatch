package config

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/fxnn/gowatch/parser"
	"testing"
)

func acceptAllPredicate() logentry.AcceptAllPredicate {
	return logentry.AcceptAllPredicate{}
}

func givenLineSource(t *testing.T, lines ...string) parser.LineSource {
	linesource := parser.NewSimpleLineSource()
	for _, line := range lines {
		linesource.AddLine(line)
	}
	return linesource
}
