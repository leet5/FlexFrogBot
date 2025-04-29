package handlers

import (
	interfaces2 "flex-frog-bot/services/interfaces"
	api "flex-frog-bot/tg-bot-api"
	"fmt"
	"log"
)

func HandleWebApp(update *api.Update, chtSvc interfaces2.ChatService, usrSvc interfaces2.UserService) {
	chatID, err := chtSvc.GetChatID(update)
	if err != nil {
		log.Printf("[bot][handle_webapp] ⚠️ Error getting chat ID: %v", err)
		return
	}
	userID, err := usrSvc.GetUserID(update)
	if err != nil {
		log.Printf("[bot][handle_webapp] ⚠️ Error getting user ID: %v", err)
		return
	}

	button := api.InlineKeyboardButton{
		Text: "🧩 Open Mini App",
		WebApp: &api.WebApp{
			URL: fmt.Sprintf("https://flexfrog.ddns.net?user_id=%d", userID),
		},
	}

	payload := api.MessagePayload{
		ChatID: chatID,
		Text:   "📋 Choose an action:",
		ReplyMarkup: &api.InlineKeyboardMarkup{
			InlineKeyboard: [][]api.InlineKeyboardButton{
				{
					button,
				},
			},
		},
	}

	err = api.SendPayloadMessage(payload)
	if err != nil {
		log.Printf("[bot][handle_webapp] ❌ Failed to send web app button: %v", err)
		return
	}
}
