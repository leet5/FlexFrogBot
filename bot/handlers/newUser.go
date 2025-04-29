package handlers

import (
	"context"
	"flex-frog-bot/db/domain"
	interfaces2 "flex-frog-bot/services/interfaces"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func HandleNewUser(ctx context.Context, update *api.Update, usrSvc interfaces2.UserService, chtSvc interfaces2.ChatService) {
	username, err := usrSvc.GetUserName(update)
	if err != nil {
		log.Printf("[bot][handle_new_user] ❌ Error getting user name: %v", err)
		return
	}
	userID, err := usrSvc.GetUserID(update)
	if err != nil {
		log.Printf("[bot][handle_new_user] ❌ Error getting user ID: %v", err)
		return
	}
	chatID, err := chtSvc.GetChatID(update)
	if err != nil {
		log.Printf("[bot][handle_new_user] ❌ Error getting chat ID: %v", err)
		return
	}

	_, err = chtSvc.GetOrCreate(ctx, update)
	if err != nil {
		log.Printf("[bot][handle_new_user] ❌ Error checking or adding new chat: %v", err)
		return
	}

	err = usrSvc.Create(ctx, chatID, &domain.User{
		Id:   userID,
		Name: username,
	})
	if err != nil {
		log.Printf("[bot][handle_new_user] ❌ Error creating user: %v", err)
		return
	}
	log.Printf("[bot][handle_new_user] ✅ User %d (%s) linked to chat %d", userID, username, chatID)
}
