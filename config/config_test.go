package config

import (
    "testing"
    "github.com/stretchr/testify/require"
)

func TestSimpleLogfile(t *testing.T) {

    config := ReadConfigFromBytes([]byte(`
        logfiles:
        - filename: my-filename.log
          tags:
          - tag1
          - tag2
          parser: grok
          parserconfig: {
            pattern: my-pattern,
            other-config: abc
          }
    `))

    require.Len(t, config.Logfiles, 1)

    logfile := config.Logfiles[0]

    require.Equal(t, "my-filename.log", logfile.Filename)

    require.Len(t, logfile.Tags, 2)
    require.Contains(t, logfile.Tags, "tag1")
    require.Contains(t, logfile.Tags, "tag2")

    require.Equal(t, "grok", logfile.Parser)

    require.Len(t, logfile.ParserConfig, 2)
    require.Equal(t, "my-pattern", logfile.ParserConfig["pattern"])
    require.Equal(t, "abc", logfile.ParserConfig["other-config"])

}

func TestMultipleLogfiles(t *testing.T) {

    config := ReadConfigFromBytes([]byte(`
        logfiles:
        - filename: my-filename.log
          parser: grok
        - filename: otherfile.txt
          parser: regexp
    `))

    require.Len(t, config.Logfiles, 2)

    var logfile GowatchLogfile

    logfile = config.Logfiles[0]
    require.Equal(t, "my-filename.log", logfile.Filename)
    require.Equal(t, "grok", logfile.Parser)

    logfile = config.Logfiles[1]
    require.Equal(t, "otherfile.txt", logfile.Filename)
    require.Equal(t, "regexp", logfile.Parser)

}

