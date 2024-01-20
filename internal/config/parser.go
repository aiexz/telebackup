package config

import (
	"gopkg.in/yaml.v3"
)

func ParseConfig(data []byte) (*Config, error) {
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}
