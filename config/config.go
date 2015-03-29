package config

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "log"
)

type GowatchConfig struct {
    Logfiles []LogfileConfig
    Summary []SummaryConfig
}

type LogfileConfig struct {
    Filename        string
    Tags            []string
    Parser            string
    Config    map[interface{}]interface{}
}

type SummaryConfig struct {
    Summarizer    string
    Title         string
    Config        map[interface{}]interface{}
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
