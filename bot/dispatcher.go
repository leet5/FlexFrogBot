package bot

import (
	"context"
	"flex-frog-bot/bot/handlers"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func processUpdates(ctx context.Context, updates <-chan *api.Update) {
	for update := range updates {
		switch {
		case isBotRemovedFromChat(update):
			handlers.HandleBotRemoved(ctx, update, ChatService)

		case isBotAddedToGroup(update):
			handlers.HandleNewChat(ctx, update, ChatService)

		case isCallbackQuery(update):
			handleCallback(ctx, update)

		case isCommand(update):
			handleCommand(ctx, update)

		case ImageService.HasImage(update):
			handlers.HandleImage(ctx, update, ImageService, ChatService, UserService)
		}

		checkNewUser(ctx, update)
	}
}

func checkNewUser(ctx context.Context, update *api.Update) {
	userID, err := UserService.GetUserID(update)
	if err != nil {
		return
	}

	user, err := UserService.GetByID(ctx, userID)
	if err != nil {
		log.Printf("[bot][check_new_user] ⚠️ Error getting user from database: %v", err)
	}

	if user == nil {
		handlers.HandleNewUser(ctx, update, UserService, ChatService)
	} else {
		chatID, err := ChatService.GetChatID(update)
		if err != nil {
			log.Printf("[bot][check_new_user] ⚠️ Error getting chat ID: %v", err)
			return
		}
		err = UserService.AssociateWithChat(ctx, userID, chatID)
		if err != nil {
			log.Printf("[bot][check_new_user] ⚠️ Error associating user with chat: %v", err)
			return
		}
	}
}
