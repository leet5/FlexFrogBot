package handlers

import (
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

	photo := update.Message.Photo
	largest := photo[len(photo)-1]
	fileID := largest.FileID

	api.DownloadFile(update, saveDir, fileID)
}
