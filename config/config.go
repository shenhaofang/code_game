package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App AppCfg `json:"app" yaml:"app"`
	Log LogCfg `json:"log" yaml:"log"`
}

func (c *Config) Load(cfgFile string) error {
	file, err := os.OpenFile(cfgFile, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	decoder := yaml.NewDecoder(file)
	return decoder.Decode(c)
}

var stdConfig = new(Config)

func Load(cfgFile string) error {
	return stdConfig.Load(cfgFile)
}

func App() AppCfg {
	return stdConfig.App
}

func Log() LogCfg {
	return stdConfig.Log
}
