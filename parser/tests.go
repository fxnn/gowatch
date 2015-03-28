package parser

import "testing"

func givenLineSource(t *testing.T, lines ...string) LineSource {
    linesource := NewSimpleLineSource()
    for _, line := range lines {
        linesource.AddLine(line)
    }
    return linesource
}

