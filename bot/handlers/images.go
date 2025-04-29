package handlers

import (
	"context"
	"flex-frog-bot/db/domain"
	"flex-frog-bot/services/interfaces"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

// HandleImage processes an incoming image or image-document from a Telegram update,
// saves it to disk, inserts it into the DB, and deletes the file afterward.
func HandleImage(ctx context.Context, update *api.Update, imgSvc interfaces.ImageService, chtSvc interfaces.ChatService, usrSvc interfaces.UserService) {
	chatID, err := chtSvc.GetChatID(update)
	if err != nil {
		log.Printf("[bot][images] ❌ Unable to extract chat ID: %v", err)
		return
	}
	watched, err := chtSvc.IsWatched(ctx, chatID)
	if err != nil {
		log.Printf("[bot][images] ❌ Unable to check if chat is watched: %v", err)
		return
	}
	if !watched {
		return
	}

	userID, err := usrSvc.GetUserID(update)
	if err != nil {
		log.Printf("[bot][images] ❌ Unable to extract user ID: %v", err)
		return
	}

	messageID, err := GetMessageID(update)
	if err != nil {
		log.Printf("[bot][images] ❌ Unable to extract message ID: %v", err)
		return
	}

	if err := imgSvc.Create(ctx, update, &domain.Image{
		MessageId: messageID,
		UserId:    userID,
		ChatId:    chatID,
	}); err != nil {
		log.Printf("[bot][images] ❌ Error saving image: %v", err)
	} else {
		log.Printf("[bot][images] ✅ Stored image from message %d", messageID)
	}
}
