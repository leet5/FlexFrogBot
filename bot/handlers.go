package bot

import (
	"log"
	"mini-app-back/bot/handlers"
	api "mini-app-back/tg-bot-api"
)

func handleCallback(update *api.Update) {
	switch update.Callback.Data {
	case "/start":
		handlers.HandleStart(update, chats)
	case "/stop":
		handlers.HandleStop(update, chats)
	case "/link":
		handlers.HandleLink(update, chats, userChats)
	case "/unlink":
		handlers.HandleUnlink(update, userChats)
	default:
		log.Printf("[bot] Unknown callback: %s", update.Callback.Data)
	}
}

func handleCommand(update *api.Update) {
	cmd := update.Message.Text
	switch cmd {
	case "/menu":
		handlers.HandleMenu(update, chats, userChats)
	case "/start":
		handlers.HandleStart(update, chats)
	case "/stop":
		handlers.HandleStop(update, chats)
	default:
		log.Printf("Unknown command: %s", cmd)
		err := api.SendTextMessage(update.Message.Chat.ID, "❓ Unknown command. Try /menu")
		if err != nil {
			log.Printf("Error sending fallback message: %v", err)
		}
	}
}
