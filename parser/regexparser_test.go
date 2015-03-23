package parser

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestSingleFieldAsMessage(t *testing.T) {
	submatchNameMap := givenMapWithSingleSubmatch(t, "Message")
	regexp := givenRegexpPattern(t, "^x([^x]*)x$")
	linesource := givenLineSource(t, "xabcx")

	result := NewRegexpParser(linesource, regexp, submatchNameMap).Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, "abc", resultEntry.Message)
}

func TestSingleFieldAsTag(t *testing.T) {
	submatchNameMap := givenMapWithSingleSubmatch(t, "Tags")
	regexp := givenRegexpPattern(t, "^x([^x]*)x$")
	linesource := givenLineSource(t, "xabcx")

	result := NewRegexpParser(linesource, regexp, submatchNameMap).Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, []string{"abc"}, resultEntry.Tags)
}

func TestSingleFieldAsLogLevel(t *testing.T) {
	submatchNameMap := givenMapWithSingleSubmatch(t, "Level")
	regexp := givenRegexpPattern(t, "^x([^x]*)x$")
	linesource := givenLineSource(t, "xDEBUGx")

	result := NewRegexpParser(linesource, regexp, submatchNameMap).Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, logentry.DEBUG, resultEntry.Level)
}

func TestSingleFieldAsCustomEntry(t *testing.T) {
	submatchNameMap := givenMapWithSingleSubmatch(t, "MyCustomEntry")
	regexp := givenRegexpPattern(t, "^x([^x]*)x$")
	linesource := givenLineSource(t, "xabcx")

	result := NewRegexpParser(linesource, regexp, submatchNameMap).Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, "abc", resultEntry.Custom["MyCustomEntry"])
}

func givenLineSource(t *testing.T, lines ...string) LineSource {
	linesource := NewSimpleLineSource()
	for _, line := range lines {
		linesource.AddLine(line)
	}
	return linesource
}

func givenRegexpPattern(t *testing.T, pattern string) *regexp.Regexp {
	regexp, err := regexp.Compile(pattern)
	if err != nil {
		t.Error(err)
	}
	return regexp
}

func givenMapWithSingleSubmatch(t *testing.T, submatchName string) map[int]string {
	submatchNameMap := make(map[int]string)
	submatchNameMap[1] = submatchName
	return submatchNameMap
}
