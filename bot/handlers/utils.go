package handlers

import (
	"errors"
	api "mini-app-back/tg-bot-api"
)

func GetChatID(update *api.Update) (int64, error) {
	if update.Message != nil {
		return update.Message.Chat.ID, nil
	}
	if update.Callback != nil && update.Callback.Message != nil {
		return update.Callback.Message.Chat.ID, nil
	}
	return 0, errors.New("chat ID not found in update")
}
