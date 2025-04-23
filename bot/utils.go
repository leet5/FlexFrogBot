package bot

import (
	api "mini-app-back/tg-bot-api"
	"strings"
)

func isCommand(update *api.Update) bool {
	return update.Message != nil && strings.HasPrefix(update.Message.Text, "/")
}

func isCallbackQuery(update *api.Update) bool {
	return update.Callback != nil && update.Callback.Data != ""
}

func isBotAddedToGroup(update *api.Update) bool {
	if update.Message == nil || update.Message.NewChatMembers == nil {
		return false
	}
	for _, member := range *update.Message.NewChatMembers {
		if member.Username == botName {
			return true
		}
	}
	return false
}

func hasImage(update *api.Update) bool {
	return update.Message != nil && len(update.Message.Photo) > 0
}
