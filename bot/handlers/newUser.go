package handlers

import (
	"context"
	"flex-frog-bot/db/repository"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func HandleNewUser(ctx context.Context, update *api.Update, userRepo *repository.UserRepository) {
	username, err := GetUserName(update)
	if err != nil {
		log.Printf("[bot][handle_new_user] ❌ Error getting user name: %v", err)
		return
	}
	userID, err := GetUserID(update)
	if err != nil {
		log.Printf("[bot][handle_new_user] ❌ Error getting user ID: %v", err)
		return
	}
	chatID, err := GetChatID(update)
	if err != nil {
		log.Printf("[bot][handle_new_user] ❌ Error getting chat ID: %v", err)
		return
	}

	_, err = userRepo.InsertUser(ctx, userID, username)
	if err != nil {
		log.Printf("[bot][handle_new_user] ⚠️ Error inserting user: %v", err)
	}

	err = userRepo.AssociateUserWithChat(ctx, userID, chatID)
	if err != nil {
		log.Printf("[bot][handle_new_user] ❌ Error associating user with chat: %v", err)
		return
	}

	log.Printf("[bot][handle_new_user] ✅ User %d (%s) linked to chat %d", userID, username, chatID)
}
