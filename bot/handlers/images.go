package handlers

import (
	"context"
	"flex-frog-bot/db/repository"
	api "flex-frog-bot/tg-bot-api"
	"fmt"
	"log"
	"os"
)

const saveDir = "./images"

// HandleImage processes an incoming image or image-document from a Telegram update,
// saves it to disk, inserts it into the DB, and deletes the file afterward.
func HandleImage(ctx context.Context, update *api.Update, imageRepo *repository.ImageRepository, chatRepo *repository.ChatRepository) {
	userID, err := GetUserID(update)
	if err != nil {
		log.Printf("[bot][images] ❌ Unable to extract user ID: %v", err)
		return
	}

	chatID, err := GetChatID(update)
	if err != nil {
		log.Printf("[bot][images] ❌ Unable to extract chat ID: %v", err)
		return
	}
	watched, err := chatRepo.CheckIfChatWatched(ctx, chatID)
	if err != nil {
		log.Printf("[bot][images] ❌ Unable to check if chat is watched: %v", err)
		return
	}
	if !watched {
		return
	}

	messageID, err := GetMessageID(update)
	if err != nil {
		log.Printf("[bot][images] ❌ Unable to extract message ID: %v", err)
		return
	}

	dest, err := downloadImageFromUpdate(update)
	if err != nil {
		log.Printf("[bot][images] ❌ %v", err)
		return
	}

	if err := storeImage(ctx, dest, chatID, userID, messageID, imageRepo); err != nil {
		log.Printf("[bot][images] ❌ %v", err)
	} else {
		log.Printf("[bot][images] ✅ Stored image from message %d", messageID)
	}
}

func downloadImageFromUpdate(update *api.Update) (string, error) {
	if update.Message == nil {
		return "", fmt.Errorf("no message content found")
	}

	switch {
	case len(update.Message.Photo) > 0:
		largest := update.Message.Photo[len(update.Message.Photo)-1]
		filename := fmt.Sprintf("%d_photo.jpg", update.Message.MessageID)
		return api.DownloadFile(saveDir, largest.FileID, filename)

	case update.Message.Document != nil && isImage(update.Message.Document.MimeType):
		doc := update.Message.Document
		filename := fmt.Sprintf("%d_%s", update.Message.MessageID, doc.FileName)
		return api.DownloadFile(saveDir, doc.FileID, filename)

	default:
		return "", fmt.Errorf("no valid image found in message")
	}
}

func storeImage(ctx context.Context, filepath string, chatID, userID, messageID int64, imageRepo *repository.ImageRepository) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read image file %s: %w", filepath, err)
	}

	_, err = imageRepo.InsertImage(ctx, &repository.Image{
		ChatID:    chatID,
		UserID:    userID,
		MessageID: messageID,
		ImageData: data,
	})
	if err != nil {
		return fmt.Errorf("failed to insert image into DB: %w", err)
	}

	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("failed to delete image file %s: %w", filepath, err)
	}
	log.Printf("[bot][images] 🗑 Deleted image file %s", filepath)

	return nil
}
