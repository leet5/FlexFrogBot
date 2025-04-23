package handlers

import (
	"log"
	api "mini-app-back/tg-bot-api"
)

func HandleMenu(update *api.Update) {
	chatID, err := GetChatID(update)
	if err != nil {
		log.Printf("[bot] ❌ Failed to extract chat ID: %v", err)
		return
	}

	isPrivate := update.Message.Chat.Type == "private"

	buttons := [][]api.InlineKeyboardButton{
		{
			{Text: "▶ Start", CallbackData: "/start"},
			{Text: "⏹ Stop", CallbackData: "/stop"},
		},
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
