package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

// Gowatch's Configuration will be marshalled and unmarshalled into/from this datastructure. Note that we don't use
// pointers here, as their contents will not be compared in deep equals, and therefore they are harder to use in tests
// etc.
type GowatchConfig struct {
	Logfiles []LogfileConfig
	Summary  []SummaryConfig
}

type LogfileConfig struct {
	Filename string
	Tags     []string
	Parser   string
	Config   map[interface{}]interface{}
	Where    PredicateConfig
}

type SummaryConfig struct {
	Summarizer string
	Title      string
	Config     map[interface{}]interface{}
	Where      PredicateConfig
}

// This structure allows to express conditions on logentry.LogEntry in configuration files. It is not made for internal
// use, but solely for unmarshalling users configuration. Keys are either names of LogEntry fields, or the special
// values "not", "allof", "anyof" or "noneof".
type PredicateConfig map[string]interface{}

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
