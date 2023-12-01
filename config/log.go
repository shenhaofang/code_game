package config

type LogCfg struct {
	Path   string `json:"path" yaml:"path"`
	Format string `json:"format" yaml:"format"`
	Level  string `json:"level" yaml:"level"`
}
