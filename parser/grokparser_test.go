package parser

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGrok_MultipleFields(t *testing.T) {
	linesource := givenLineSource(t, "WARNING This is the message")

	parser := grokParserWithLinesourceAndPattern(linesource, "%{LOGLEVEL:Level} %{DATA:Message}$")
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, logentry.WARNING, resultEntry.Level)
	require.Equal(t, "This is the message", resultEntry.Message)
}

func TestGrok_SingleMessage(t *testing.T) {
	linesource := givenLineSource(t, "abc")

	parser := grokParserWithLinesourceAndPattern(linesource, "^%{DATA:Message}$")
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, "abc", resultEntry.Message)
}

func TestGrok_SingleTag(t *testing.T) {
	linesource := givenLineSource(t, "abc")

	parser := grokParserWithLinesourceAndPattern(linesource, "^%{DATA:Tags}$")
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, []string{"abc"}, resultEntry.Tags)
}

func TestGrok_MultipleTags(t *testing.T) {
	linesource := givenLineSource(t, "abc def")

	parser := grokParserWithLinesourceAndPattern(linesource, "^%{DATA:Tags} %{DATA:Tags}$")
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, []string{"abc", "def"}, resultEntry.Tags)
}

func TestGrok_SingleLogLevel(t *testing.T) {
	linesource := givenLineSource(t, "DEBUG")

	parser := grokParserWithLinesourceAndPattern(linesource, "^%{LOGLEVEL:Level}$")
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, logentry.DEBUG, resultEntry.Level)
}

func TestGrok_SingleCustomEntry(t *testing.T) {
	linesource := givenLineSource(t, "abc")

	parser := grokParserWithLinesourceAndPattern(linesource, "^%{DATA:MyCustomEntry}$")
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, "abc", resultEntry.Custom["MyCustomEntry"])
}

func TestGrok_MultipleCustomEntries(t *testing.T) {
	linesource := givenLineSource(t, "28.03.2015 abc")

	parser := grokParserWithLinesourceAndPattern(linesource, "^%{DATE:CustomDate} %{USER:MyCustomEntry}$")
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, "28.03.2015", resultEntry.Custom["CustomDate"])
	require.Equal(t, "abc", resultEntry.Custom["MyCustomEntry"])
}

func TestGrok_Predicate(t *testing.T) {
	linesource := givenLineSource(t, "abc")

	parser := grokParserWithLinesourceAndPredicate(linesource, &logentry.ContainsPredicate{FieldName: "Message", ToBeContained: "xyz"})
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.Equal(t, logentry.LogEntry{}, resultEntry) // zero value --> no element in channel
}

func TestGrok_TimeLayout(t *testing.T) {
	linesource := givenLineSource(t, "Tue, 10 Nov 2009 23:00:00 +0000")
	expectedTime, _ := time.Parse(time.RFC1123Z, "Tue, 10 Nov 2009 23:00:00 +0000")

	parser := grokParserWithLinesourceAndTimeLayout(linesource, time.RFC1123Z)
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.Equal(t, expectedTime, resultEntry.Timestamp)
}

func grokParserWithLinesourceAndTimeLayout(linesource LineSource, timeLayout string) *GrokParser {
	return NewGrokParser(linesource, "^%{DATA:timestamp}$", timeLayout, acceptAllPredicate())
}

func grokParserWithLinesourceAndPredicate(linesource LineSource, predicate logentry.Predicate) *GrokParser {
	return NewGrokParser(linesource, "^%{DATA:Message}$", "", predicate)
}

func grokParserWithLinesourceAndPattern(linesource LineSource, pattern string) *GrokParser {
	return NewGrokParser(linesource, pattern, "", acceptAllPredicate())
}
