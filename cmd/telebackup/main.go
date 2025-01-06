package main

import (
	"flag"
	"fmt"
	"github.com/aiexz/telebackup/internal/compress"
	"github.com/aiexz/telebackup/internal/config"
	"github.com/aiexz/telebackup/internal/sender"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	levelDebug = "debug"
	levelInfo  = "info"
	levelWarn  = "error"
	leveNone   = "none"
)

func setupLogger(level string) {
	var handler slog.Handler

	switch strings.ToLower(level) {
	case levelDebug:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case levelInfo:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	case levelWarn:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError})
	case leveNone:
		handler = slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError})
	default:
		fmt.Printf("Unknown log level: %s\n", level)
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func main() {
	logLevel := flag.String("log", "info", "log level")
	configFile := flag.String("config", "config.yml", "config file")
	flag.Parse()
	setupLogger(*logLevel)
	slog.Debug("Starting telebackup", "config", *configFile)
	resultConfig := &config.Config{}
	if os.Getenv("APP_ID") != "" {
		// maybe use another way to check if env should be used
		slog.Debug("Using env variables for config")
		var err error
		resultConfig, err = config.ParseEnvConfig()
		if err != nil {
			panic(err)
		}
	} else {
		slog.Debug("Using config file for config")
		workingDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		configFilePath := filepath.Join(workingDir, *configFile)
		slog.Debug("config file path", "path", configFilePath)
		reader, err := os.ReadFile(configFilePath)

		if err != nil {
			panic(err)
		}
		resultConfig, err = config.ParseFileConfig(reader)
		if err != nil {
			panic(err)
		}
	}

	client, err := sender.NewSender(resultConfig.AppID, resultConfig.AppHash, resultConfig.BotToken)
	if err != nil {
		panic(err)
	}
	err = client.Start()
	if err != nil {
		panic(err)
	}
	slog.Info("Telegram client started")

	slog.Debug("Starting workers for sending files", "count", len(resultConfig.PathTarget))
	wg := &sync.WaitGroup{}
	for _, target := range resultConfig.PathTarget {
		slog.Debug("Starting worker for sending file", "path", target.GetPath())
		var thread int32
		var path string
		if target.IsForum() {
			thread, path = target.Forum.Topic, target.GetPath()
		} else {
			path = target.GetPath()
		}
		slog.Debug("Sending file", "path", path)
		go func() {
			defer wg.Done()
			tempFile, err := os.CreateTemp("", "telebackup-*.tar.gz")
			if err != nil {
				slog.Error("Error creating temp file", "error", err)
				return
			}
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					slog.Debug("Error removing temp file", "error", err)
					//	it is debug because it is not critical
				}
			}(tempFile.Name())

			buf, err := os.OpenFile(tempFile.Name(), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				slog.Error("Error opening temp file", "error", err)
				return
			}
			err = compress.CompressPath(path, buf)
			if err != nil {
				slog.Error("Error compressing path", "error", err)
				return
			}

			dirs := strings.Split(path, "/")
			lastDir := dirs[len(dirs)-1]
			err = client.SendMedia(resultConfig.TelegramTarget, tempFile.Name(), &sender.SendOptions{Caption: path, FileName: lastDir + fmt.Sprintf("-%d.tar.gz", time.Now().Unix()), Thread: thread})
			if err != nil {
				slog.Error("Error sending file", "path", path, "error", err)
				return
			}
			slog.Info("File sent", "path", path)

		}()
		wg.Add(1)
	}
	wg.Wait()

}
