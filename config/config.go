package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const depCheckFileName = "dep-check-v2.yaml"

type depCheckConfig struct {
	Artifacts []string `yaml:"artifacts"`
}

var DepCheckConfig depCheckConfig

func init() {
	data, err := os.ReadFile(depCheckFileName)
	if err != nil {
		panic(fmt.Sprint("No", depCheckFileName, "found", err))
	}
	err = yaml.Unmarshal(data, &DepCheckConfig)
	if err != nil {
		panic(fmt.Sprint(depCheckFileName, " is not yaml parsable", err))
	}
}
