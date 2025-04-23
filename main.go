package main

import (
	"flex-frog-bot/bot"
	tgbotapi "flex-frog-bot/tg-bot-api"
	_ "fmt"
)

func main() {
	updates := tgbotapi.GetUpdatesChan()
	bot.ProcessUpdates(updates)
}
