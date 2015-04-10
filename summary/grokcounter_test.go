package summary

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGrokcounterWithSimpleName(t *testing.T) {

	// given
	patternsByName := make(map[string]string)
	patternsByName["ip"] = "%{IP}"
	sut := NewGrokCounter(patternsByName, &logentry.AcceptAllPredicate{})
	entries := givenEntriesWithMessages("127.0.0.1", "prefix192.168.0.0.1 ", "this aint no ip")

	// when
	sut.Summarize(entries)
	result := sut.StringAfterSummarizeAsyncCompleted()

	// then
	require.Equal(t, "ip: 2", result)

}

func TestGrokcounterWithReplacementInName(t *testing.T) {

	// given
	patternsByName := make(map[string]string)
	patternsByName["prefix %{IP} suffix"] = "%{IP}"
	sut := NewGrokCounter(patternsByName, &logentry.AcceptAllPredicate{})
	entries := givenEntriesWithMessages("127.0.0.1", "prefix 192.168.0.1 ", "this aint no ip")

	// when
	sut.Summarize(entries)
	result := sut.StringAfterSummarizeAsyncCompleted()

	// then
	require.Contains(t, result, "prefix 127.0.0.1 suffix: 1")
	require.Contains(t, result, "prefix 192.168.0.1 suffix: 1")

}

func TestGrokcounterWithPredicate(t *testing.T) {

	// given
	patternsByName := make(map[string]string)
	patternsByName["ip"] = "%{IP}"
	sut := NewGrokCounter(patternsByName, &logentry.ContainsPredicate{FieldName: "Message", ToBeContained: "127"})
	entries := givenEntriesWithMessages("127.0.0.1", "prefix192.168.0.0.1 ", "this aint no ip")

	// when
	sut.Summarize(entries)
	result := sut.StringAfterSummarizeAsyncCompleted()

	// then
	require.Equal(t, "ip: 1", result)

}

func TestGrokcounterSortsOutput(t *testing.T) {

	// given
	patternsByName := make(map[string]string)
	patternsByName["%{GREEDYDATA}"] = "%{GREEDYDATA}"
	sut := NewGrokCounter(patternsByName, &logentry.AcceptAllPredicate{})
	entries := givenEntriesWithMessages("ba", "ab", "aa")

	// when
	sut.Summarize(entries)
	result := sut.StringAfterSummarizeAsyncCompleted()

	// then
	require.Equal(t, "aa: 1\nab: 1\nba: 1", result)

}

func TestGrokcounterSortsWithNumbersInMind(t *testing.T) {

	// given
	patternsByName := make(map[string]string)
	patternsByName["%{GREEDYDATA}"] = "%{GREEDYDATA}"
	sut := NewGrokCounter(patternsByName, &logentry.AcceptAllPredicate{})
	entries := givenEntriesWithMessages("1", "10", "2")

	// when
	sut.Summarize(entries)
	result := sut.StringAfterSummarizeAsyncCompleted()

	// then
	require.Equal(t, "1: 1\n2: 1\n10: 1", result)

}
