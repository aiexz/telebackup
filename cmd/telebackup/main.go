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
		var configFilePath string
		if *configFile != "config.yml" {
			configFilePath = *configFile
		} else {
			configFilePath = filepath.Join(workingDir, "config.yml")
		}
		//TODO: add test for config file detection
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
			
			// Create a split writer that will handle splitting into 2GB parts
			splitWriter, err := compress.NewSplitWriter(tempFile.Name())
			if err != nil {
				slog.Error("Error creating split writer", "error", err)
				os.Remove(tempFile.Name())
				return
			}
			
			// Compress to the split writer
			err = compress.CompressPath(path, splitWriter)
			if err != nil {
				slog.Error("Error compressing path", "error", err)
				splitWriter.Close()
				// Clean up all parts
				for _, part := range splitWriter.Parts() {
					os.Remove(part)
				}
				return
			}
			
			// Close the split writer
			if err := splitWriter.Close(); err != nil {
				slog.Error("Error closing split writer", "error", err)
				return
			}
			
			// Get all parts
			parts := splitWriter.Parts()
			
			// Clean up function for all parts
			defer func() {
				for _, part := range parts {
					err := os.Remove(part)
					if err != nil {
						slog.Debug("Error removing temp file", "error", err, "file", part)
					}
				}
			}()

			dirs := strings.Split(path, "/")
			lastDir := dirs[len(dirs)-1]
			timestamp := time.Now().Unix()
			
			// Send all parts
			for i, partPath := range parts {
				var fileName string
				if len(parts) == 1 {
					fileName = lastDir + fmt.Sprintf("-%d.tar.gz", timestamp)
				} else {
					fileName = lastDir + fmt.Sprintf("-%d.tar.gz.part%d", timestamp, i+1)
				}
				
				caption := fmt.Sprintf("%s (part %d/%d)", path, i+1, len(parts))
				if len(parts) == 1 {
					caption = path
				}
				
				slog.Info("Sending file part", "path", path, "part", i+1, "total", len(parts))
				err = client.SendMedia(resultConfig.TelegramTarget, partPath, &sender.SendOptions{
					Caption:  caption,
					FileName: fileName,
					Thread:   thread,
				})
				if err != nil {
					slog.Error("Error sending file part", "path", path, "part", i+1, "error", err)
					return
				}
			}
			
			slog.Info("File sent successfully", "path", path, "parts", len(parts))

		}()
		wg.Add(1)
	}
	wg.Wait()

}
