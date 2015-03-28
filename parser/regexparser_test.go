package parser

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSingleFieldAsMessage(t *testing.T) {
	submatchNameMap := givenMapWithSingleSubmatch(t, "Message")
	linesource := givenLineSource(t, "xabcx")

	parser, err := NewRegexpParser(linesource, "^x([^x]*)x$", submatchNameMap)
	result := parser.Parse()

	require.Nil(t, err)
	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, "abc", resultEntry.Message)
}

func TestSingleFieldAsTag(t *testing.T) {
	submatchNameMap := givenMapWithSingleSubmatch(t, "Tags")
	linesource := givenLineSource(t, "xabcx")

	parser, err := NewRegexpParser(linesource, "^x([^x]*)x$", submatchNameMap)
	result := parser.Parse()

	require.Nil(t, err)
	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, []string{"abc"}, resultEntry.Tags)
}

func TestSingleFieldAsLogLevel(t *testing.T) {
	submatchNameMap := givenMapWithSingleSubmatch(t, "Level")
	linesource := givenLineSource(t, "xDEBUGx")

	parser, err := NewRegexpParser(linesource, "^x([^x]*)x$", submatchNameMap)
	result := parser.Parse()

	require.Nil(t, err)
	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, logentry.DEBUG, resultEntry.Level)
}

func TestSingleFieldAsCustomEntry(t *testing.T) {
	submatchNameMap := givenMapWithSingleSubmatch(t, "MyCustomEntry")
	linesource := givenLineSource(t, "xabcx")

	parser, err := NewRegexpParser(linesource, "^x([^x]*)x$", submatchNameMap)
	result := parser.Parse()

	require.Nil(t, err)
	require.NotNil(t, result)

	resultEntry := <-result
	require.NotNil(t, resultEntry)
	require.Equal(t, "abc", resultEntry.Custom["MyCustomEntry"])
}

func givenMapWithSingleSubmatch(t *testing.T, submatchName string) map[int]string {
	submatchNameMap := make(map[int]string)
	submatchNameMap[1] = submatchName
	return submatchNameMap
}
