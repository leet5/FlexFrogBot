package handlers

import (
	"log"
	api "mini-app-back/tg-bot-api"
)

func HandleNewChat(update *api.Update) {
	chatID := update.Message.Chat.ID
	log.Printf("[bot] 🆕 New chat detected (chat_id=%d). Sending start button...", chatID)

	err := api.SendPayloadMessage(api.MessagePayload{
		ChatID: chatID,
		Text:   "Press 'Start' to activate the bot.",
		ReplyMarkup: api.InlineKeyboardMarkup{
			InlineKeyboard: [][]api.InlineKeyboardButton{
				{
					{Text: "▶ Start", CallbackData: "/start"},
				},
			},
		},
	})
	if err != nil {
		log.Println(err)
	}
}
