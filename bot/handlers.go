package bot

import (
	"context"
	"flex-frog-bot/bot/handlers"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func handleCallback(ctx context.Context, update *api.Update) {
	switch update.Callback.Data {
	case "/start":
		handlers.HandleStart(ctx, update, ChatService)
	case "/stop":
		handlers.HandleStop(ctx, update, ChatService)
	case "/webapp":
		handlers.HandleWebApp(update, ChatService, UserService)
	default:
		log.Printf("[bot][handle_callback] ⚠️ Unknown callback: %s", update.Callback.Data)
	}
}

func handleCommand(ctx context.Context, update *api.Update) {
	cmd := update.Message.Text
	switch cmd {
	case "/start":
		handlers.HandleStart(ctx, update, ChatService)
	case "/stop":
		handlers.HandleStop(ctx, update, ChatService)
	case "/webapp":
		handlers.HandleWebApp(update, ChatService, UserService)
	default:
		log.Printf("[bot][handle_command] ⚠️ Unknown command: %s", cmd)
		err := api.SendTextMessage(update.Message.Chat.ID, "❓ Unknown command. Try /start /stop")
		if err != nil {
			log.Printf("[bot][handle_command] ⚠️ Error sending fallback message: %v", err)
		}
	}
}
