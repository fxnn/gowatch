package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type GowatchConfig struct {
	Logfiles []LogfileConfig
	Summary  []SummaryConfig
}

type LogfileConfig struct {
	Filename string
	Tags     []string
	Parser   string
	Config   map[interface{}]interface{}
}

type SummaryConfig struct {
	Summarizer string
	Title      string
	Config     map[interface{}]interface{}
}

// This structure allows to express conditions on logentry.LogEntry in configuration files. It is not made for internal
// use, but solely for unmarshalling users configuration. The "Field" fields must be set, and then at exactly one of the
// other fields. The latter define the condition on the value of the former.
type PredicateConfig struct {
	// Name of this field this condition applies to. Fields that are not in logentry.LogEntry will be treated as
	// logentry.LogEntry.Custom entry.
	Field string
	// If set, should be one of "empty", "not empty"
	Is string
	// If set, should be a string that is expected to be contained in the fields value.
	Contains string
	// If set, should be a regexp that is expected to match the fields value.
	Matches string
	// If set, all PredicateConfigs are expected to match.
	AllOf []PredicateConfig
	// If set, at least one of the PredicateConfigs is expected to match. "Field" must not be set when using this one.
	AnyOf []PredicateConfig
	// If set, none of the PredicateConfigs is expected to match. "Field" must not be set when using this one.
	NoneOf []PredicateConfig
	// If set, the PredicateConfig is expected not to match. "Field" must not be set when using this one.
	Not *PredicateConfig
}

func ReadConfigByFilename(filename string) *GowatchConfig {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return new(GowatchConfig)
	}

	return ReadConfigFromBytes(contents)
}

func ReadConfigFromBytes(contents []byte) *GowatchConfig {
	config := new(GowatchConfig)

	err := yaml.Unmarshal(contents, config)
	if err != nil {
		log.Fatal(err)
		return config
	}

	return config
}
