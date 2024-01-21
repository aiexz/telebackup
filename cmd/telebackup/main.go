package main

import (
	"flag"
	"fmt"
	"github.com/amarnathcjd/gogram/telegram"
	"log"
	"os"
	"strings"
	"sync"
	"telebackup/internal/compress"
	"telebackup/internal/config"
	"time"
)

func main() {
	configFile := flag.String("config", "config.yaml", "config file")
	flag.Parse()
	reader, err := os.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}
	resultConfig, err := config.ParseConfig(reader)
	if err != nil {
		panic(err)
	}

	client, _ := telegram.NewClient(telegram.ClientConfig{
		AppID:    resultConfig.AppID,
		AppHash:  resultConfig.AppHash,
		LogLevel: telegram.LogWarn,
	})

	if err := client.Connect(); err != nil {
		panic(err)
	}

	// Authenticate the client using the bot token
	if err := client.LoginBot(resultConfig.BotToken); err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	for _, path := range resultConfig.Targets {
		path := path
		go func() {
			tempFile, err := os.CreateTemp("", "telebackup-*.tar.gz")
			if err != nil {
				log.Println("Error creating temp file", err)
				return
			}
			buf, _ := os.OpenFile(tempFile.Name(), os.O_CREATE|os.O_WRONLY, 0644)
			err = compress.CompressPath(path, buf)
			if err != nil {
				log.Println("Error compressing path", err)
				return
			}
			dirs := strings.Split(path, "/")
			lastDir := dirs[len(dirs)-1]
			file, err := client.UploadFile(tempFile.Name(), &telegram.UploadOptions{FileName: lastDir + fmt.Sprintf("-%d.tar.gz", time.Now().Unix())})
			if err != nil {
				log.Println("Error uploading file", err)
				return
			}
			_, err = client.SendMedia(resultConfig.Target, file, &telegram.MediaOptions{Caption: path})
			if err != nil {
				log.Println("Error sending file", err)
				return
			}
			log.Println(path, "sent")
			wg.Done()
		}()
		wg.Add(1)
	}
	wg.Wait()

}
