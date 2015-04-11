package config

import (
	"github.com/fxnn/gowatch/logentry"
	"github.com/fxnn/gowatch/parser"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateParserWithPredefinedTimeLayout(t *testing.T) {

	formattedTime := "2006-01-02T15:04:05-07:00"

	linesource := parser.NewSimpleLineSource()
	linesource.AddLine(formattedTime)
	predicate := logentry.AcceptAllPredicate{}

	config := LogfileConfig{Parser: "grok", TimeLayout: "RFC3339", Config: map[interface{}]interface{}{"pattern": "^%{DATA:Timestamp}$"}}
	parser := config.CreateParser(linesource, predicate)

	entries := parser.Parse()
	require.NotNil(t, entries)
	entry := <-entries
	require.Equal(t, formattedTime, entry.Timestamp.Format(time.RFC3339))

}

func TestParseTimeLayout(t *testing.T) {

	testcases := []struct {
		givenTimeLayout    string
		expectedTimeLayout string
	}{
		{"ANSIC", time.ANSIC},
		{"UnixDate", time.UnixDate},
		{"RubyDate", time.RubyDate},
		{"RFC822", time.RFC822},
		{"RFC822Z", time.RFC822Z},
		{"RFC850", time.RFC850},
		{"RFC1123", time.RFC1123},
		{"RFC1123Z", time.RFC1123Z},
		{"RFC3339", time.RFC3339},
		{"RFC3339Nano", time.RFC3339Nano},
		{"Kitchen", time.Kitchen},
		{"Stamp", time.Stamp},
		{"StampMilli", time.StampMilli},
		{"StampMicro", time.StampMicro},
		{"StampNano", time.StampNano},
	}

	for _, testcase := range testcases {
		actualTimeLayout := parseTimeLayout(testcase.givenTimeLayout)
		require.Equal(t, testcase.expectedTimeLayout, actualTimeLayout)
	}

}
