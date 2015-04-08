package parser

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGrok_MultipleFields(t *testing.T) {
	linesource := givenLineSource(t, "WARNING This is the message")

	parser := NewGrokParser(linesource, "%{LOGLEVEL:Level} %{DATA:Message}$", acceptAllPredicate())
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, logentry.WARNING, resultEntry.Level)
	require.Equal(t, "This is the message", resultEntry.Message)
}

func TestGrok_SingleMessage(t *testing.T) {
	linesource := givenLineSource(t, "abc")

	parser := NewGrokParser(linesource, "^%{DATA:Message}$", acceptAllPredicate())
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, "abc", resultEntry.Message)
}

func TestGrok_SingleTag(t *testing.T) {
	linesource := givenLineSource(t, "abc")

	parser := NewGrokParser(linesource, "^%{DATA:Tags}$", acceptAllPredicate())
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, []string{"abc"}, resultEntry.Tags)
}

func TestGrok_MultipleTags(t *testing.T) {
	linesource := givenLineSource(t, "abc def")

	parser := NewGrokParser(linesource, "^%{DATA:Tags} %{DATA:Tags}$", acceptAllPredicate())
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, []string{"abc", "def"}, resultEntry.Tags)
}

func TestGrok_SingleLogLevel(t *testing.T) {
	linesource := givenLineSource(t, "DEBUG")

	parser := NewGrokParser(linesource, "^%{LOGLEVEL:Level}$", acceptAllPredicate())
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, logentry.DEBUG, resultEntry.Level)
}

func TestGrok_SingleCustomEntry(t *testing.T) {
	linesource := givenLineSource(t, "abc")

	parser := NewGrokParser(linesource, "^%{DATA:MyCustomEntry}$", acceptAllPredicate())
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, "abc", resultEntry.Custom["MyCustomEntry"])
}

func TestGrok_MultipleCustomEntries(t *testing.T) {
	linesource := givenLineSource(t, "28.03.2015 abc")

	parser := NewGrokParser(linesource, "^%{DATE:CustomDate} %{USER:MyCustomEntry}$", acceptAllPredicate())
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, "28.03.2015", resultEntry.Custom["CustomDate"])
	require.Equal(t, "abc", resultEntry.Custom["MyCustomEntry"])
}

func TestGrok_Predicate(t *testing.T) {
	linesource := givenLineSource(t, "abc")

	parser := NewGrokParser(linesource, "^%{DATA:Message}$", &logentry.ContainsPredicate{FieldName: "Message", ToBeContained: "xyz"})
	result := parser.Parse()

	require.NotNil(t, result)

	resultEntry := <-result
	require.Equal(t, logentry.LogEntry{}, resultEntry) // zero value --> no element in channel
}

func acceptAllPredicate() logentry.Predicate {
	return &logentry.AcceptAllPredicate{}
}
