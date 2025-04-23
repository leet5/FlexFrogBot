package bot

import (
	"mini-app-back/bot/handlers"
	api "mini-app-back/tg-bot-api"
)

const (
	botName = "img_srch_bot"
)

var (
	chats     = make(map[int64]bool)
	userChats = make(map[int64]map[int64]struct{})
)

func ProcessUpdates(updates <-chan *api.Update) {
	for update := range updates {
		switch {
		case isBotAddedToGroup(update):
			handlers.HandleNewChat(update)

		case isCallbackQuery(update):
			handleCallback(update)

		case isCommand(update):
			handleCommand(update)

		case hasImage(update):
			handlers.HandleImage(update, chats)
		}
	}

}
