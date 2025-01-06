package config

import (
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {
	data := []byte(`appId: 1
appHash: 2
botToken: 3
target: "@test"
targets:
    - /test
    - topic: 2
      path: /test2
`)
	config, err := ParseFileConfig(data)
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
	if config.TelegramTarget.Username != "@test" {
		t.Error("ChatID not parsed correctly")
	}
	if len(config.PathTarget) != 2 {
		t.Error("PathTarget not parsed correctly")
	}
	if config.PathTarget[0].GetPath() != "/test" {
		t.Error("Topic not parsed correctly")
	}
	if !config.PathTarget[1].IsForum() {
		t.Error("Forum not parsed correctly")
	}
	if config.PathTarget[1].GetPath() != "/test2" || config.PathTarget[1].Forum.Topic != 2 {
		t.Error("Forum not parsed correctly")
	}
	data = []byte(`appId: 1
appHash: 2
botToken: 3
target: 123
targets:
    - /test
    - /test2
`)
	config, err = ParseFileConfig(data)
	if config.TelegramTarget.ID != int64(123) {
		t.Error("ChatID not parsed correctly")
	}
}

func TestParseEnvConfig(t *testing.T) {
	os.Setenv("APP_ID", "1")
	os.Setenv("APP_HASH", "2")
	os.Setenv("BOT_TOKEN", "3")
	os.Setenv("TARGET", "123")
	os.Setenv("TARGETS", "/test\n/test2")
	config, err := ParseEnvConfig()
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
	if config.TelegramTarget.ID != int64(123) {
		t.Error("ChatID not parsed correctly")
	}
	if len(config.PathTarget) != 2 {
		t.Error("PathTarget not parsed correctly")
	}
	if config.PathTarget[0].GetPath() != "/test" {
		t.Error("Topic not parsed correctly")
	}
	if config.PathTarget[1].GetPath() != "/test2" {
		t.Error("Forum not parsed correctly")
	}
}
