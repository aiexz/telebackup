package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

func ParseFileConfig(data []byte) (*Config, error) {
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}
	return config, nil
}

func ParseEnvConfig() (*Config, error) {
	config := &Config{}
	appID, err := strconv.ParseInt(os.Getenv("APP_ID"), 10, 64)
	if err != nil {
		return nil, err
	}
	config.AppID = int32(appID)
	config.AppHash = os.Getenv("APP_HASH")
	config.BotToken = os.Getenv("BOT_TOKEN")
	target, err := strconv.ParseInt(os.Getenv("TARGET"), 10, 64)
	if err != nil {
		// TODO: support for telegram username
		return nil, errors.New("TARGET is not a valid chat id")
	}
	config.TelegramTarget = TelegramTarget{
		Username: "",
		ID:       target,
	}
	targets := os.Getenv("TARGETS")
	if strings.Contains(targets, "\n") {
		targets = strings.ReplaceAll(targets, "\n", ",")
	}
	for _, target := range strings.Split(targets, ",") {
		target = strings.TrimSpace(target)
		if target == "" {
			continue
		}
		config.PathTarget = append(config.PathTarget, PathTarget{
			Path:  target,
			Forum: Forum{},
		})
	}
	return config, nil
}
