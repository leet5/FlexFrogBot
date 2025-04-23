package handlers

import (
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func HandleStop(update *api.Update, chats map[int64]bool) {
	chatID, err := GetChatID(update)
	if err != nil {
		log.Printf("[bot] ❌ Failed to extract chat ID: %v", err)
		return
	}

	chats[chatID] = false
	log.Printf("[bot] ⛔ Stop button pressed in chat_id=%d. Bot deactivated.", chatID)

	err = api.SendTextMessage(chatID, "⛔ Bot has been deactivated.")
	if err != nil {
		log.Printf("[bot] ❌ Failed to send stop confirmation: %v", err)
	}
}
