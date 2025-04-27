package handlers

import (
	"context"
	"flex-frog-bot/db/repository"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

type Chat struct {
	ChatID    int64
	ChatName  string
	Thumbnail []byte
	Watched   bool
}

func HandleStart(ctx context.Context, update *api.Update, chatRepo *repository.ChatRepository) {
	chatID, err := GetChatID(update)
	if err != nil {
		log.Printf("[bot][handle_start] ❌ Failed to extract chat ID: %v", err)
		return
	}

	chat, err := chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		log.Printf("[bot][handle_start] ❌ Failed to get chat by ID: %v", err)
		return
	}

	addIfAbsent(ctx, chatID, update, chatRepo)

	if chat.Watched {
		log.Printf("[bot][handle_start] ❌ Bot is already active in chat '%s'.", chat.ChatName)
		err = api.SendTextMessage(chatID, "❌ Bot is already active in this chat.")
		if err != nil {
			log.Printf("[bot][handle_start] ❌ Failed to send message: %v", err)
		}
		return
	} else {
		log.Printf("[bot][handle_start] ✅ Bot is not active in chat '%s'. Activating now.", chat.ChatName)
		err = chatRepo.WatchChat(ctx, chatID)
		if err != nil {
			log.Printf("[bot][handle_start] ❌ Failed to update chat status: %v", err)
			return
		}
		log.Printf("[bot][handle_start] ✅ Chat '%s' status updated to active.", chat.ChatName)
		err = api.SendTextMessage(chatID, "✅ Bot is now active and ready to process images!")
		if err != nil {
			log.Printf("[bot][handle_start] ❌ Failed to send start confirmation: %v", err)
		}
	}
}
