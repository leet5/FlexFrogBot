package handlers

import (
	"fmt"
	"log"
	api "mini-app-back/tg-bot-api"
)

const saveDir = "./images"

func HandleImage(update *api.Update, groups map[int64]bool) {
	chatID, err := GetChatID(update)
	if err != nil {
		log.Printf("[images] ❌ Failed to extract chat ID: %v", err)
		return
	}

	if !groups[chatID] {
		return
	}

	switch {
	case len(update.Message.Photo) > 0:
		photo := update.Message.Photo
		largest := photo[len(photo)-1]
		fileID := largest.FileID
		filename := fmt.Sprintf("%d_photo.jpg", update.Message.MessageID)
		api.DownloadFile(saveDir, fileID, filename)

	case update.Message.Document != nil && isImage(update.Message.Document.MimeType):
		doc := update.Message.Document
		fileID := doc.FileID
		filename := fmt.Sprintf("%d_%s", update.Message.MessageID, doc.FileName)
		api.DownloadFile(saveDir, fileID, filename)
	}
}
