package handlers

import (
	"errors"
	api "flex-frog-bot/tg-bot-api"
)

func GetMessageID(update *api.Update) (int64, error) {
	if update.Message != nil {
		return int64(update.Message.MessageID), nil
	}
	if update.Callback != nil && update.Callback.Message != nil {
		return int64(update.Callback.Message.MessageID), nil
	}
	return 0, errors.New("message ID not found in update")
}
