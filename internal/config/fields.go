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
	TelegramTarget TelegramTarget `yaml:"target"`
	// Optional: Mapping of paths and topics
	PathTarget []PathTarget `yaml:"targets"`
}

type TelegramTarget struct {
	Username string
	ID       int64
}

func (t *TelegramTarget) UnmarshalYAML(unmarshal func(interface{}) error) error {
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

type PathTarget struct {
	Path  string
	Forum Forum
}

func (t *PathTarget) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var path string
	var forum Forum
	if err := unmarshal(&path); err == nil {
		t.Path = path
		return nil
	}
	if err := unmarshal(&forum); err == nil {
		t.Forum = forum
		return nil
	}
	return nil

}

func (t *PathTarget) GetPath() string {
	if t.Path != "" {
		return t.Path
	} else if t.Forum.Path != "" {
		return t.Forum.Path
	}
	return ""
}

func (t *PathTarget) IsForum() bool {
	return t.Forum.Path != "" && t.Forum.Topic != 0
}

type Forum struct {
	Topic int32
	Path  string
}

func (f *Forum) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var forum struct {
		Topic int32
		Path  string
	}
	if err := unmarshal(&forum); err == nil {
		f.Topic = forum.Topic
		f.Path = forum.Path
		return nil
	}
	return nil
}
