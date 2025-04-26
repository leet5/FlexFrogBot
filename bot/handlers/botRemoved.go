package handlers

import (
	"context"
	"flex-frog-bot/db/repository"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func HandleBotRemoved(ctx context.Context, update *api.Update, chatRepo *repository.ChatRepository) {
	chatID, err := GetChatID(update)
	if err != nil {
		log.Printf("[bot][handle_bot_removed] ❌ Error getting chat ID: %v", err)
		return
	}

	log.Printf("[bot][handle_bot_removed] ❌ Bot was removed from chat %d", chatID)

	err = chatRepo.UnwatchChat(ctx, chatID)
	if err != nil {
		log.Printf("[bot][handle_bot_removed] ⚠️ Error unwatching chat: %v", err)
		return
	}
}
