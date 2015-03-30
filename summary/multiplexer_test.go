package summary

import (
    "testing"
    "github.com/stretchr/testify/require"
    "github.com/fxnn/gowatch/logentry"
)

func TestSummarizeWithTwoEchoes(t *testing.T) {

    // given
    echo1, echo2 := NewEcho(), NewEcho()
    sut := NewMultiplexer()
    sut.AddSummarizer(echo1)
    sut.AddSummarizer(echo2)

    // when
    entries := make(chan logentry.LogEntry, 3)
    sut.SummarizeAsync(entries)

    entries <- *logentry.New()
    entries <- *logentry.New()
    entries <- *logentry.New()
    close(entries)

    // then
    require.Equal(t, 3, echo1.NumberOfLines())
    require.Equal(t, 3, echo2.NumberOfLines())

}