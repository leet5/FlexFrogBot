package bot

import (
	api "flex-frog-bot/tg-bot-api"
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
	if update.Message == nil {
		return false
	}
	if update.Message.Document != nil && isImage(update.Message.Document.MimeType) {
		return true
	}
	return len(update.Message.Photo) > 0
}

func isImage(mime string) bool {
	return mime == "image/jpeg" || mime == "image/png" || mime == "image/gif" || mime == "image/webp"
}
