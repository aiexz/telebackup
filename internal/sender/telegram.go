package sender

import (
	"errors"
	"github.com/aiexz/telebackup/internal/config"
	"github.com/amarnathcjd/gogram/telegram"
	"log/slog"
	"strconv"
	"time"
)

type Sender struct {
	client   *telegram.Client
	botToken string
	cache    map[string]interface{} // cache for username resolution
}

func NewSender(AppID int32, AppHash string, BotToken string) (*Sender, error) {
	client, err := telegram.NewClient(telegram.ClientConfig{
		AppID:    AppID,
		AppHash:  AppHash,
		LogLevel: telegram.LogError,
	})
	if err != nil {
		return nil, err
	}
	return &Sender{client: client, botToken: BotToken, cache: make(map[string]interface{})}, nil
}

func (s *Sender) Start() error {
	// wait for 10 seconds to connect to telegram servers and try to reconnect if failed
	for i := 0; i < 10; i++ {
		errChan := make(chan error)
		go func() {
			errChan <- s.client.Connect()
		}()
		select {
		case err := <-errChan:
			if err != nil {
				return err
			}
			break
		case <-time.After(10 * time.Second):
			slog.Debug("telegram client is not connected, trying to reconnect")
		}
	}

	for i := 0; i < 10; i++ {
		// Authenticate the client using the bot token
		errChan := make(chan error)
		go func() {
			errChan <- s.client.LoginBot(s.botToken)
		}()
		select {
		case err := <-errChan:
			if err != nil {
				return err
			}
			break
		case <-time.After(10 * time.Second):
			slog.Debug("telegram client is not authorized, trying to authorize")
		}
	}

	return nil
}

type SendOptions struct {
	Caption  string
	FileName string
	Thread   int32
}

func (s *Sender) SendMedia(target config.TelegramTarget, path string, options *SendOptions) error {
	targetResolved, err := s.ResolveTarget(target)
	if err != nil {
		return err
	}

	_, err = s.client.SendMedia(targetResolved, path, &telegram.MediaOptions{Caption: options.Caption, FileName: options.FileName, ReplyID: options.Thread})
	return err
}

// ResolveTarget resolves the config.TelegramTarget to a telegram.ChatObj, telegram.Channel or telegram.UserObj
func (s *Sender) ResolveTarget(target config.TelegramTarget) (interface{}, error) {
	if target.Username != "" {
		if resolved, ok := s.cache[target.Username]; ok {
			return resolved, nil
		}
		// ResolveUsername is not cached by the library, so we have to do it ourselves
		resolved, err := s.client.ResolveUsername(target.Username)
		if err != nil {
			return nil, err
		}
		s.cache[target.Username] = resolved
		return resolved, nil
	}
	if target.ID != 0 {
		if resolved, ok := s.cache[strconv.FormatInt(target.ID, 10)]; ok {
			return resolved, nil
		}
		//	first we try to resolve it as a chat or channel
		chat, err := s.client.ChannelsGetChannels([]telegram.InputChannel{&telegram.InputChannelObj{ChannelID: target.ID, AccessHash: 0}})
		if err == nil && chat != nil {
			resolved := chat.(*telegram.MessagesChatsObj).Chats[0]
			s.cache[strconv.FormatInt(target.ID, 10)] = resolved
			return resolved, nil
		}

		// now we try to resolve it as a user
		user, err := s.client.UsersGetUsers([]telegram.InputUser{&telegram.InputUserObj{UserID: target.ID, AccessHash: 0}})
		if err != nil {
			return nil, err
		}
		if len(user) > 0 {
			resolved := user[0]
			s.cache[strconv.FormatInt(target.ID, 10)] = resolved
			return resolved, nil
		}
	}
	return nil, errors.New("no target found")
}
