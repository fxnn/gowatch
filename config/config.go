package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

// Gowatch's Configuration will be marshalled and unmarshalled into/from this datastructure. Note that we don't use
// pointers here, as their contents will not be compared in deep equals, and therefore they are harder to use in tests
// etc.
type GowatchConfig struct {
	Logfiles []LogfileConfig
	Summary  []SummaryConfig
}

// Configures a logfile whose lines shall be parsed and processed by gowatch
type LogfileConfig struct {
	Filename string
	// Tags that will be appended to each logline
	Tags   []string
	Parser string
	// Configures the parser
	With map[interface{}]interface{}
	// Predicates that can filter the lines being processed by gowatch
	Where PredicateConfig
	// A layout providing the reference time, as described in Go's time.Parse function. Also see the predefined layouts
	// in the PredefinedTimeLayouts variable.
	TimeLayout string
}

// Configures one summary
type SummaryConfig struct {
	// Determines the type of summary
	Do    string
	Title string
	// Configures the summary
	With map[interface{}]interface{}
	// Prediactes that filter the lines to be considered in this summary
	Where PredicateConfig
}

// This structure allows to express conditions on logentry.LogEntry in configuration files. It is not made for internal
// use, but solely for unmarshalling users configuration. Keys are either names of LogEntry fields, or the special
// values "not", "allof", "anyof" or "noneof".
type PredicateConfig map[string]interface{}

var PredefinedTimeLayouts map[string]string = map[string]string{
	"ANSIC":       time.ANSIC,
	"UnixDate":    time.UnixDate,
	"RubyDate":    time.RubyDate,
	"RFC822":      time.RFC822,
	"RFC822Z":     time.RFC822Z,
	"RFC850":      time.RFC850,
	"RFC1123":     time.RFC1123,
	"RFC1123Z":    time.RFC1123Z,
	"RFC3339":     time.RFC3339,
	"RFC3339Nano": time.RFC3339Nano,
	"Kitchen":     time.Kitchen,
	"":            time.Stamp,
	"Stamp":       time.Stamp,
	"StampMilli":  time.StampMilli,
	"StampMicro":  time.StampMicro,
	"StampNano":   time.StampNano,
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
