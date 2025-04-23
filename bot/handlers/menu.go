package handlers

import (
	"log"
	api "mini-app-back/tg-bot-api"
)

func HandleMenu(update *api.Update, groups map[int64]bool, userGroups map[int64]map[int64]struct{}) {
	chatID, err := GetChatID(update)
	if err != nil {
		log.Printf("[bot] ❌ Failed to get chat ID: %v", err)
		return
	}
	userID, err := GetUserID(update)
	if err != nil {
		log.Printf("[bot] ❌ Failed to get user ID: %v", err)
		return
	}

	isPrivate := update.Message.Chat.Type == "private"
	var buttons [][]api.InlineKeyboardButton

	isAdmin, err := api.IsUserAdmin(chatID, userID)
	if err != nil {
		log.Printf("[bot] Failed to check user admin rights: %v", err)
		return
	}

	if isAdmin {
		if groups[chatID] {
			buttons = append(buttons, []api.InlineKeyboardButton{{Text: "⏹️ Stop", CallbackData: "/stop"}})
		} else {
			buttons = append(buttons, []api.InlineKeyboardButton{{Text: "▶️ Start", CallbackData: "/start"}})
		}
	}

	if userGroups[userID] == nil {
		userGroups[userID] = make(map[int64]struct{})
	}

	if _, ok := userGroups[userID][chatID]; !ok {
		buttons = append(buttons, []api.InlineKeyboardButton{{Text: "⏺️ Add to watches", CallbackData: "/link"}})
	} else {
		buttons = append(buttons, []api.InlineKeyboardButton{{Text: "↩️ Remove from watches", CallbackData: "/unlink"}})
	}

	if isPrivate {
		buttons = append(buttons, []api.InlineKeyboardButton{
			{
				Text: "🧩 Open Mini App",
				WebApp: &api.WebApp{
					URL: "https://your-mini-app-url.com",
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
		log.Printf("[bot] ❌ Failed to send menu: %v", err)
	} else {
		log.Printf("[bot] 📋 Sent menu to chat_id=%d", chatID)
	}
}
