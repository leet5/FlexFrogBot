package main

import (
	_ "fmt"
	"mini-app-back/bot"
	tgbotapi "mini-app-back/tg-bot-api"
	"mini-app-back/watcher"
	"sync"
)

func main() {
	updates := tgbotapi.GetUpdatesChan()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		bot.ProcessUpdates(updates)
	}()
	go func() {
		defer wg.Done()
		watcher.WatchForNewImages()
	}()

	wg.Wait()
}
