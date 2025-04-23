package handlers

import (
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func HandleStart(update *api.Update, chats map[int64]bool) {
	chatID, err := GetChatID(update)
	if err != nil {
		log.Printf("[bot] ❌ Failed to extract chat ID: %v", err)
		return
	}

	chats[chatID] = true
	log.Printf("[bot] ✅ Start button pressed in chat_id=%d. Bot activated.", chatID)

	err = api.SendTextMessage(chatID, "✅ Bot is now active and ready to process images!")
	if err != nil {
		log.Printf("[bot] ❌ Failed to send start confirmation: %v", err)
	}
}
