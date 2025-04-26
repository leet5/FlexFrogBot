package handlers

import (
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func HandleNewChat(update *api.Update) {
	chatID := update.Message.Chat.ID
	log.Printf("[bot][handle_new_chat] 🆕 New chat detected (chat_id=%d). Sending start button...", chatID)

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
