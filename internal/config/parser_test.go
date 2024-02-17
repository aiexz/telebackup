package config

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	data := []byte(`appId: 1
appHash: 2
botToken: 3
target: "@test"
targets:
    - /test
    - /test2
`)
	config, err := ParseConfig(data)
	if err != nil {
		t.Error(err)
	}
	if config.AppID != 1 {
		t.Error("ApiID not parsed correctly")
	}
	if config.AppHash != "2" {
		t.Error("ApiHash not parsed correctly")
	}
	if config.BotToken != "3" {
		t.Error("BotToken not parsed correctly")
	}
	if config.Target.Username != "@test" {
		t.Error("ChatID not parsed correctly")
	}
	if len(config.Targets) != 2 {
		t.Error("Targets not parsed correctly")
	}
	if config.Targets[0] != "/test" {
		t.Error("Topic not parsed correctly")
	}
	if config.Targets[1] != "/test2" {
		t.Error("Topic not parsed correctly")
	}
	data = []byte(`appId: 1
appHash: 2
botToken: 3
target: 123
targets:
    - /test
    - /test2
`)
	config, err = ParseConfig(data)
	if config.Target.ID != int64(123) {
		t.Error("ChatID not parsed correctly")
	}
}
