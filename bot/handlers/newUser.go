package handlers

import (
	"context"
	"flex-frog-bot/db/repository"
	"flex-frog-bot/img_tools"
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

	photoPath, err := api.GetUserProfilePhoto(saveDir, img_tools.GenerateUUID(), userID)
	if err != nil {
		log.Printf("[bot][handle_new_user] ❌ Error getting user profile photo: %v", err)
		photoPath = ""
	}

	thumbnail, err := img_tools.CreateThumbnailByPath(photoPath)
	if err != nil {
		log.Printf("[bot][handle_new_user] ⚠️ Error creating thumbnail: %v", err)
	}

	_, err = userRepo.InsertUser(ctx, userID, username, thumbnail)
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
