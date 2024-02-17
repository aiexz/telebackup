package config

// Config for the whole application
type Config struct {
	// Telegram API ID credentials
	AppID int32 `yaml:"appId"`
	// Telegram API Hash credentials
	AppHash string `yaml:"appHash"`
	// Telegram Bot Token
	BotToken string `yaml:"botToken"`
	// Telegram user to send messages to
	Target Target `yaml:"target"`
	// Optional: Mapping of paths and topics
	Targets []string `yaml:"targets"`
}

type Target struct {
	Username string
	ID       int64
}

func (t *Target) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var username string
	var ID int64
	if err := unmarshal(&ID); err == nil {
		t.ID = ID
		return nil
	}
	if err := unmarshal(&username); err == nil {
		t.Username = username
		return nil
	}
	return nil
}
