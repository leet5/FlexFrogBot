package handlers

import (
	"context"
	"flex-frog-bot/services/interfaces"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func HandleNewChat(ctx context.Context, update *api.Update, chtSvc interfaces.ChatService) {
	chat, err := chtSvc.GetOrCreate(ctx, update)
	if err != nil {
		log.Printf("[bot][handle_new_chat] ❌ Failed to add new chat: %v", err)
		return
	}

	err = api.SendPayloadMessage(api.MessagePayload{
		ChatID: chat.Id,
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
