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
			"logfiles: [{filename: my-filename.log, where: {a: {matches: b}} }]",
			GowatchConfig{
				Logfiles: []LogfileConfig{
					LogfileConfig{
						Filename: "my-filename.log",
						Where:    PredicateConfig{"a": map[interface{}]interface{}{"matches": "b"}},
					},
				},
			},
		},
		// parser
		{
			"logfiles: [{parser: grok, with: {pattern: my-pattern, patterns: {a: b, c: d}}}]",
			GowatchConfig{
				Logfiles: []LogfileConfig{LogfileConfig{
					Parser: "grok",
					With: map[interface{}]interface{}{
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
			"summary: [{summarizer: echo, title: Title, with: {a: {b: c}}}]",
			GowatchConfig{Summary: []SummaryConfig{
				SummaryConfig{
					Summarizer: "echo",
					Title:      "Title",
					With: map[interface{}]interface{}{
						"a": map[interface{}]interface{}{"b": "c"},
					},
				},
			}},
		}, {
			"summary: [{title: Title, where: {a: {matches: b}} }]",
			GowatchConfig{
				Summary: []SummaryConfig{
					SummaryConfig{
						Title: "Title",
						Where: PredicateConfig{"a": map[interface{}]interface{}{"matches": "b"}},
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
