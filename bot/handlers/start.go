package handlers

import (
	"context"
	"flex-frog-bot/db/repository"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

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
	if chat == nil {
		chatName, err := GetChatName(update)
		if err != nil {
			log.Printf("[bot][handle_start] ❌ Failed to extract chat name: %v", err)
			return
		}
		log.Printf("[bot][handle_start] ✅ Chat '%s' not found in DB. Adding it now.", chatName)

		err = chatRepo.InsertChat(ctx, chatID, chatName, true)
		if err != nil {
			log.Printf("[bot][handle_start] ❌ Failed to insert chat '%s' into DB: %v", chatName, err)
			return
		}
		log.Printf("[bot][handle_start] ✅ Chat '%s' added to DB.", chatName)

		err = api.SendTextMessage(chatID, "✅ Bot is now active and ready to process images!")
		if err != nil {
			log.Printf("[bot][handle_start] ❌ Failed to send start confirmation: %v", err)
		}
	} else {
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
}
