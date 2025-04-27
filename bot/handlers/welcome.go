package handlers

import (
	"context"
	"flex-frog-bot/db/repository"
	"flex-frog-bot/img_tools"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func HandleNewChat(ctx context.Context, update *api.Update, chatRepo *repository.ChatRepository) {
	chatID := update.Message.Chat.ID
	log.Printf("[bot][handle_new_chat] 🆕 New chat detected (chat_id=%d). Sending start button...", chatID)

	addIfAbsent(ctx, chatID, update, chatRepo)

	err := api.SendPayloadMessage(api.MessagePayload{
		ChatID: chatID,
		Text:   "Press 'Start' to activate the bot.",
		ReplyMarkup: &api.InlineKeyboardMarkup{
			InlineKeyboard: [][]api.InlineKeyboardButton{
				{
					{Text: "▶ Start", CallbackData: "/start"},
				},
			},
		},
	})
	if err != nil {
		log.Printf("[bot][handle_new_chat] ❌ Failed to send start button: %v", err)
	}
}

func addIfAbsent(ctx context.Context, chatID int64, update *api.Update, chatRepo *repository.ChatRepository) {
	chat, err := chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		log.Printf("[bot][handle_new_chat] ❌ Failed to get chat by ID: %v", err)
		return
	}

	if chat == nil {
		chatName, err := GetChatName(update)
		if err != nil {
			log.Printf("[bot][handle_new_chat] ❌ Failed to extract chat name: %v", err)
			return
		}
		log.Printf("[bot][handle_new_chat] ✅ Chat '%s' not found in DB. Adding it now.", chatName)

		photoPath, err := api.GetChatProfilePhoto(saveDir, img_tools.GenerateUUID(), chatID)
		if err != nil {
			log.Printf("[bot][handle_new_chat] ❌ Error getting chat profile photo: %v", err)
			photoPath = ""
		}

		thumbnail, err := img_tools.CreateThumbnailByPath(photoPath)
		if err != nil {
			log.Printf("[bot][handle_new_user] ⚠️ Error creating thumbnail: %v", err)
		}

		err = chatRepo.InsertChat(ctx, chatID, chatName, thumbnail, false)
		if err != nil {
			log.Printf("[bot][handle_new_chat] ❌ Failed to insert chat '%s' into DB: %v", chatName, err)
			return
		}
		log.Printf("[bot][handle_new_chat] ✅ Chat '%s' added to DB.", chatName)
	}
}
