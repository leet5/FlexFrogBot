package bot

import (
	"context"
	"flex-frog-bot/bot/handlers"
	api "flex-frog-bot/tg-bot-api"
	"fmt"
	"log"
)

func handleCallback(ctx context.Context, update *api.Update) {
	switch update.Callback.Data {
	case "/start":
		handlers.HandleStart(ctx, update, ChatService, UserService)
	case "/stop":
		handlers.HandleStop(ctx, update, ChatService, UserService)
	default:
		log.Printf("[bot][handle_callback] ⚠️ Unknown callback: %s", update.Callback.Data)
	}
}

func handleCommand(ctx context.Context, update *api.Update) {
	cmd := update.Message.Text
	switch cmd {
	case fmt.Sprintf("/start@%s", botName):
		handlers.HandleStart(ctx, update, ChatService, UserService)
	case fmt.Sprintf("/stop@%s", botName):
		handlers.HandleStop(ctx, update, ChatService, UserService)
	default:
		log.Printf("[bot][handle_command] ⚠️ Unknown command: %s", cmd)
		err := api.SendTextMessage(update.Message.Chat.ID, "❓ Unknown command. Try /start /stop")
		if err != nil {
			log.Printf("[bot][handle_command] ⚠️ Error sending fallback message: %v", err)
		}
	}
}
