package summary

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSummarizeWithTwoEchoes(t *testing.T) {

	// given
	echo1, echo2 := createEcho(), createEcho()
	sut := NewMultiplexer()
	sut.AddSummarizer(echo1)
	sut.AddSummarizer(echo2)

	// when
	entries := givenEntriesWithMessages("1", "2", "3")
	sut.SummarizeAsync(entries)

	// then
	require.Equal(t, 3, echo1.NumberOfLinesAfterSummarizeAsyncCompleted())
	require.Equal(t, 3, echo2.NumberOfLinesAfterSummarizeAsyncCompleted())

}

func createEcho() *Echo {
	return NewEcho(&logentry.AcceptAllPredicate{})
}
