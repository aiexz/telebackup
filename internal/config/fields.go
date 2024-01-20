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
	Target string `yaml:"target"`
	// Optional: Mapping of paths and topics
	Targets []string `yaml:"targets"`
}
