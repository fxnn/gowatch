package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type GowatchLogfile struct {
	Filename string
	Tags     []string
}

type GowatchConfig struct {
	Logfiles []GowatchLogfile
}

func ReadConfigByFilename(filename string) GowatchConfig {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		return *new(GowatchConfig)
	}

	return ReadConfigFromBytes(contents)
}

func ReadConfigFromBytes(contents []byte) GowatchConfig {
	config := new(GowatchConfig)

	err := yaml.Unmarshal(contents, config)
	if err != nil {
		log.Fatal(err)
		return *config
	}

	return *config
}
