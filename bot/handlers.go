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
		handlers.HandleStart(ctx, update, ChatRepo)
	case "/stop":
		handlers.HandleStop(ctx, update, ChatRepo)
	default:
		log.Printf("[bot] Unknown callback: %s", update.Callback.Data)
	}
}

func handleCommand(ctx context.Context, update *api.Update) {
	cmd := update.Message.Text
	switch cmd {
	case "/menu":
		handlers.HandleMenu(ctx, update, ChatRepo)
	case "/start":
		handlers.HandleStart(ctx, update, ChatRepo)
	case "/stop":
		handlers.HandleStop(ctx, update, ChatRepo)
	default:
		log.Printf("Unknown command: %s", cmd)
		err := api.SendTextMessage(update.Message.Chat.ID, "❓ Unknown command. Try /menu")
		if err != nil {
			log.Printf("Error sending fallback message: %v", err)
		}
	}
}
