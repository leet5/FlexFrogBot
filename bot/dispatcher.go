package bot

import (
	"context"
	"flex-frog-bot/bot/handlers"
	"flex-frog-bot/db/repository"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

const (
	botName = "img_srch_bot"
)

var (
	ImageRepo *repository.ImageRepository
	ChatRepo  *repository.ChatRepository
	UserRepo  *repository.UserRepository
)

func ProcessUpdates(ctx context.Context, updates <-chan *api.Update) {
	for update := range updates {
		switch {
		case isBotRemovedFromChat(update):
			handlers.HandleBotRemoved(ctx, update, ChatRepo)

		case isBotAddedToGroup(update):
			handlers.HandleNewChat(update)

		case isCallbackQuery(update):
			handleCallback(ctx, update)

		case isCommand(update):
			handleCommand(ctx, update)

		case hasImage(update):
			handlers.HandleImage(ctx, update, ImageRepo, ChatRepo)
		}

		checkNewUser(ctx, update)
	}
}

func checkNewUser(ctx context.Context, update *api.Update) {
	userID, err := handlers.GetUserID(update)
	if err != nil {
		log.Printf("[bot][check_new_user] ⚠️ Error getting user ID: %v", err)
		return
	}

	user, err := UserRepo.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("[bot][check_new_user] ⚠️ Error getting user from database: %v", err)
	}

	if user == nil {
		handlers.HandleNewUser(ctx, update, UserRepo)
	}
}
