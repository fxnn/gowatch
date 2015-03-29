package parser

import (
    "github.com/stretchr/testify/require"
    "github.com/fxnn/gowatch/logentry"
    "testing"
)

func TestGrok_MultipleFields(t *testing.T) {
    linesource := givenLineSource(t, "WARNING This is the message")

    parser := NewGrokParser(linesource, "%{LOGLEVEL:Level} %{DATA:Message}$")
    result := parser.Parse()

    require.NotNil(t, result)

    resultEntry := <-result
    require.NotNil(t, resultEntry)
    require.Equal(t, logentry.WARNING, resultEntry.Level)
    require.Equal(t, "This is the message", resultEntry.Message)
}

func TestGrok_SingleMessage(t *testing.T) {
    linesource := givenLineSource(t, "abc")

    parser := NewGrokParser(linesource, "^%{DATA:Message}$")
    result := parser.Parse()

    require.NotNil(t, result)

    resultEntry := <-result
    require.NotNil(t, resultEntry)
    require.Equal(t, "abc", resultEntry.Message)
}

func TestGrok_SingleTag(t *testing.T) {
    linesource := givenLineSource(t, "abc")

    parser := NewGrokParser(linesource, "^%{DATA:Tags}$")
    result := parser.Parse()

    require.NotNil(t, result)

    resultEntry := <-result
    require.NotNil(t, resultEntry)
    require.Equal(t, []string{"abc"}, resultEntry.Tags)
}

func TestGrok_MultipleTags(t *testing.T) {
    linesource := givenLineSource(t, "abc def")

    parser := NewGrokParser(linesource, "^%{DATA:Tags} %{DATA:Tags}$")
    result := parser.Parse()

    require.NotNil(t, result)

    resultEntry := <-result
    require.NotNil(t, resultEntry)
    require.Equal(t, []string{"abc", "def"}, resultEntry.Tags)
}

func TestGrok_SingleLogLevel(t *testing.T) {
    linesource := givenLineSource(t, "DEBUG")

    parser := NewGrokParser(linesource, "^%{LOGLEVEL:Level}$")
    result := parser.Parse()

    require.NotNil(t, result)

    resultEntry := <-result
    require.NotNil(t, resultEntry)
    require.Equal(t, logentry.DEBUG, resultEntry.Level)
}

func TestGrok_SingleCustomEntry(t *testing.T) {
    linesource := givenLineSource(t, "abc")

    parser := NewGrokParser(linesource, "^%{DATA:MyCustomEntry}$")
    result := parser.Parse()

    require.NotNil(t, result)

    resultEntry := <-result
    require.NotNil(t, resultEntry)
    require.Equal(t, "abc", resultEntry.Custom["MyCustomEntry"])
}

func TestGrok_MultipleCustomEntries(t *testing.T) {
    linesource := givenLineSource(t, "28.03.2015 abc")

    parser := NewGrokParser(linesource, "^%{DATE:CustomDate} %{USER:MyCustomEntry}$")
    result := parser.Parse()

    require.NotNil(t, result)

    resultEntry := <-result
    require.NotNil(t, resultEntry)
    require.Equal(t, "28.03.2015", resultEntry.Custom["CustomDate"])
    require.Equal(t, "abc", resultEntry.Custom["MyCustomEntry"])
}

