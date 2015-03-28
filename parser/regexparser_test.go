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

	parser := NewRegexpParser(linesource, regexp, submatchNameMap)
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, "abc", resultEntry.Message)
}

func TestSingleFieldAsTag(t *testing.T) {
	submatchNameMap := givenMapWithSingleSubmatch(t, "Tags")
	regexp := givenRegexpPattern(t, "^x([^x]*)x$")
	linesource := givenLineSource(t, "xabcx")

	parser := NewRegexpParser(linesource, regexp, submatchNameMap)
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, []string{"abc"}, resultEntry.Tags)
}

func TestSingleFieldAsLogLevel(t *testing.T) {
	submatchNameMap := givenMapWithSingleSubmatch(t, "Level")
	regexp := givenRegexpPattern(t, "^x([^x]*)x$")
	linesource := givenLineSource(t, "xDEBUGx")

	parser := NewRegexpParser(linesource, regexp, submatchNameMap)
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, logentry.DEBUG, resultEntry.Level)
}

func TestSingleFieldAsCustomEntry(t *testing.T) {
	submatchNameMap := givenMapWithSingleSubmatch(t, "MyCustomEntry")
	regexp := givenRegexpPattern(t, "^x([^x]*)x$")
	linesource := givenLineSource(t, "xabcx")

	parser := NewRegexpParser(linesource, regexp, submatchNameMap)
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, "abc", resultEntry.Custom["MyCustomEntry"])
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
