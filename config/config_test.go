package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLogParsing(t *testing.T) {

	testcases := []struct {
		input    string
		expected GowatchConfig
	}{
		// filename and tags
		{
			"logfiles:\n- filename: my-filename.log\n  tags:\n  - tag1\n  - tag2",
			GowatchConfig{
				Logfiles: []LogfileConfig{
					LogfileConfig{
						Filename: "my-filename.log",
						Tags:     []string{"tag1", "tag2"},
					},
				},
			},
		}, {
			"logfiles: [{filename: my-filename.log, tags: [tag1, tag2]}]",
			GowatchConfig{
				Logfiles: []LogfileConfig{
					LogfileConfig{
						Filename: "my-filename.log",
						Tags:     []string{"tag1", "tag2"},
					},
				},
			},
		}, {
			"logfiles: [{filename: my-filename.log}, {filename: another-log.txt}]",
			GowatchConfig{
				Logfiles: []LogfileConfig{
					LogfileConfig{
						Filename: "my-filename.log",
					},
					LogfileConfig{
						Filename: "another-log.txt",
					},
				},
			},
		}, {
			"logfiles: [{filename: my-filename.log, where: {field: a, matches: b}}]",
			GowatchConfig{
				Logfiles: []LogfileConfig{
					LogfileConfig{
						Filename: "my-filename.log",
						Where:    PredicateConfig{Field: "a", Matches: "b"},
					},
				},
			},
		},
		// parser
		{
			"logfiles: [{parser: grok, config: {pattern: my-pattern, patterns: {a: b, c: d}}}]",
			GowatchConfig{
				Logfiles: []LogfileConfig{LogfileConfig{
					Parser: "grok",
					Config: map[interface{}]interface{}{
						"pattern": "my-pattern",
						"patterns": map[interface{}]interface{}{
							"a": "b",
							"c": "d",
						},
					},
				}},
			},
		},
		// summary
		{
			"summary: [{summarizer: echo, title: Title, config: {a: {b: c}}}]",
			GowatchConfig{Summary: []SummaryConfig{
				SummaryConfig{
					Summarizer: "echo",
					Title:      "Title",
					Config: map[interface{}]interface{}{
						"a": map[interface{}]interface{}{"b": "c"},
					},
				},
			}},
		}, {
			"summary: [{title: Title, where: {field: a, matches: b}}]",
			GowatchConfig{
				Summary: []SummaryConfig{
					SummaryConfig{
						Title: "Title",
						Where: PredicateConfig{Field: "a", Matches: "b"},
					},
				},
			},
		},
	}

	for _, testcase := range testcases {
		actual := ReadConfigFromBytes([]byte(testcase.input))
		require.Equal(t, &testcase.expected, actual)
	}

}
