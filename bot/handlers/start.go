package handlers

import (
	"context"
	"flex-frog-bot/services/interfaces"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func HandleStart(ctx context.Context, update *api.Update, chtSvc interfaces.ChatService) {
	chatID, err := chtSvc.GetChatID(update)
	if err != nil {
		log.Printf("[bot][handle_start] ❌ Failed to extract chat ID: %v", err)
		return
	}

	chat, err := chtSvc.GetOrCreate(ctx, update)
	if err != nil {
		log.Printf("[bot][handle_new_chat] ❌ Failed to add new chat: %v", err)
		return
	}

	if chat.Watched {
		log.Printf("[bot][handle_start] ❌ Bot is already active in chat '%s'.", chat.Name)
		err = api.SendTextMessage(chatID, "❌ Bot is already active in this chat.")
		if err != nil {
			log.Printf("[bot][handle_start] ❌ Failed to send message: %v", err)
		}
		return
	} else {
		log.Printf("[bot][handle_start] ✅ Bot is not active in chat '%s'. Activating now.", chat.Name)
		err = chtSvc.Watch(ctx, chatID)
		if err != nil {
			log.Printf("[bot][handle_start] ❌ Failed to update chat status: %v", err)
			return
		}
		log.Printf("[bot][handle_start] ✅ Chat '%s' status updated to active.", chat.Name)
		err = api.SendTextMessage(chatID, "✅ Bot is now active and ready to process images!")
		if err != nil {
			log.Printf("[bot][handle_start] ❌ Failed to send start confirmation: %v", err)
		}
	}
}
