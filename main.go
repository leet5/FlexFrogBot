package main

import (
	_ "fmt"
	"mini-app-back/bot"
	tgbotapi "mini-app-back/tg-bot-api"
)

func main() {
	updates := tgbotapi.GetUpdatesChan()
	bot.ProcessUpdates(updates)
}
