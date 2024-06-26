package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"telebackup/internal/compress"
	"telebackup/internal/config"
	"telebackup/internal/sender"
	"time"
)

func main() {
	configFile := flag.String("config", "config.yml", "config file")
	flag.Parse()
	reader, err := os.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}
	resultConfig, err := config.ParseConfig(reader)
	if err != nil {
		panic(err)
	}

	client, err := sender.NewSender(resultConfig.AppID, resultConfig.AppHash, resultConfig.BotToken)
	if err != nil {
		panic(err)
	}
	err = client.Start()
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	for _, target := range resultConfig.PathTarget {
		var thread int32
		var path string
		if target.IsForum() {
			thread, path = target.Forum.Topic, target.GetPath()
		} else {
			path = target.GetPath()
		}
		go func() {
			defer wg.Done()
			tempFile, err := os.CreateTemp("", "telebackup-*.tar.gz")
			if err != nil {
				log.Println("Error creating temp file", err)
				return
			}

			buf, _ := os.OpenFile(tempFile.Name(), os.O_CREATE|os.O_WRONLY, 0644)
			err = compress.CompressPath(path, buf)
			if err != nil {
				log.Println("Error compressing path", path, err)
				return
			}

			dirs := strings.Split(path, "/")
			lastDir := dirs[len(dirs)-1]
			err = client.SendMedia(resultConfig.TelegramTarget, tempFile.Name(), &sender.SendOptions{Caption: path, FileName: lastDir + fmt.Sprintf("-%d.tar.gz", time.Now().Unix()), Thread: thread})
			if err != nil {
				log.Println("Error sending file", path, err)
				return
			}

			log.Println(path, "sent")

		}()
		wg.Add(1)
	}
	wg.Wait()

}
