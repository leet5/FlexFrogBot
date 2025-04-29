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

func isBotRemovedFromChat(update *api.Update) bool {
	if update.MyChatMember == nil {
		return false
	}
	return update.MyChatMember.NewChatMember.User.Username == botName &&
		update.MyChatMember.NewChatMember.Status == "left"
}
