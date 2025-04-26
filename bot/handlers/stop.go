package handlers

import (
	"context"
	"flex-frog-bot/db/repository"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func HandleStop(ctx context.Context, update *api.Update, chatRepo *repository.ChatRepository) {
	chatID, err := GetChatID(update)
	if err != nil {
		log.Printf("[bot][handle_stop] ❌ Failed to extract chat ID: %v", err)
		return
	}

	watched, err := chatRepo.CheckIfChatWatched(ctx, chatID)
	if err != nil {
		log.Printf("[bot][handle_stop] ❌ Failed to check if chat is already watched: %v", err)
		return
	}

	if watched {
		if err := chatRepo.UnwatchChat(ctx, chatID); err != nil {
			log.Printf("[bot][handle_stop] ❌ Failed to unwatch chat '%d': %v", chatID, err)
			return
		}
		log.Printf("[bot][handle_stop] ⛔ Stop button pressed in chat_id=%d. Bot deactivated.", chatID)
		err = api.SendTextMessage(chatID, "⛔ Bot has been deactivated.")
		if err != nil {
			log.Printf("[bot][handle_stop] ❌ Failed to send stop confirmation: %v", err)
		}
	} else {
		log.Printf("[bot][handle_stop] ❌ Chat '%d' is not watched. Bot is already inactive.", chatID)
		err = api.SendTextMessage(chatID, "❌ Bot is already inactive in this chat.")
		if err != nil {
			log.Printf("[bot][handle_stop] ❌ Failed to send message: %v", err)
		}
		return
	}
}
