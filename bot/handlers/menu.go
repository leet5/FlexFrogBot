package handlers

import (
	"context"
	"flex-frog-bot/db/repository"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func HandleMenu(ctx context.Context, update *api.Update, chatRepo *repository.ChatRepository) {
	chatID, err := GetChatID(update)
	if err != nil {
		log.Printf("[bot][handle_menu] ❌ Failed to get chat ID: %v", err)
		return
	}
	userID, err := GetUserID(update)
	if err != nil {
		log.Printf("[bot][handle_menu] ❌ Failed to get user ID: %v", err)
		return
	}

	isPrivate := update.Message.Chat.Type == "private"
	var buttons [][]api.InlineKeyboardButton

	isAdmin, err := api.IsUserAdmin(chatID, userID)
	if err != nil {
		log.Printf("[bot][handle_menu] Failed to check user admin rights: %v", err)
		return
	}

	if isAdmin {
		watched, err := chatRepo.CheckIfChatWatched(ctx, chatID)
		if err != nil {
			log.Printf("[bot][handle_menu] ❌ Failed to check if chat is watched: %v", err)
			return
		}
		if watched {
			buttons = append(buttons, []api.InlineKeyboardButton{{Text: "⏹️ Stop", CallbackData: "/stop"}})
		} else {
			buttons = append(buttons, []api.InlineKeyboardButton{{Text: "▶️ Start", CallbackData: "/start"}})
		}
	}

	if isPrivate {
		buttons = append(buttons, []api.InlineKeyboardButton{
			{
				Text: "🧩 Open Mini App",
				WebApp: &api.WebApp{
					URL: "https://flex-frog.ddns.net",
				},
			},
		})
	}

	err = api.SendPayloadMessage(api.MessagePayload{
		ChatID: chatID,
		Text:   "📋 Choose an action:",
		ReplyMarkup: api.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		},
	})

	if err != nil {
		log.Printf("[bot][handle_menu] ❌ Failed to send menu: %v", err)
	} else {
		log.Printf("[bot][handle_menu] 📋 Sent menu to chat_id=%d", chatID)
	}
}
