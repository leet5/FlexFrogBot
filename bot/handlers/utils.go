package handlers

import (
	"errors"
	api "flex-frog-bot/tg-bot-api"
)

func GetChatID(update *api.Update) (int64, error) {
	if update.Message != nil {
		return update.Message.Chat.ID, nil
	}
	if update.Callback != nil && update.Callback.Message != nil {
		return update.Callback.Message.Chat.ID, nil
	}
	if update.MyChatMember != nil {
		return update.MyChatMember.Chat.ID, nil
	}
	return 0, errors.New("chat ID not found in update")
}

func GetChatName(update *api.Update) (string, error) {
	if update.Message != nil {
		return update.Message.Chat.Title, nil
	}
	if update.Callback != nil && update.Callback.Message != nil {
		return update.Callback.Message.Chat.Title, nil
	}
	return "", errors.New("chat name not found in update")
}

func GetMessageID(update *api.Update) (int64, error) {
	if update.Message != nil {
		return int64(update.Message.MessageID), nil
	}
	if update.Callback != nil && update.Callback.Message != nil {
		return int64(update.Callback.Message.MessageID), nil
	}
	return 0, errors.New("message ID not found in update")
}

func GetUserID(update *api.Update) (int64, error) {
	if update.Message != nil && update.Message.From != nil {
		return update.Message.From.ID, nil
	}
	if update.Callback != nil && update.Callback.From != nil {
		return update.Callback.From.ID, nil
	}
	return 0, errors.New("user ID not found in update")
}

func GetUserName(update *api.Update) (string, error) {
	if update.Message != nil && update.Message.From != nil {
		user := update.Message.From
		if user.Username != "" {
			return user.Username, nil
		}
		return user.FirstName, nil // fallback if username is missing
	}
	if update.Callback != nil && update.Callback.From != nil {
		user := update.Callback.From
		if user.Username != "" {
			return user.Username, nil
		}
		return user.FirstName, nil
	}
	return "", errors.New("username not found in update")
}

func isImage(mime string) bool {
	return mime == "image/jpeg" || mime == "image/png" || mime == "image/gif" || mime == "image/webp"
}
