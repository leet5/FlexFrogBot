package handlers

import (
	api "flex-frog-bot/tg-bot-api"
	"log"
)

// HandleUnlink removes the specified group from the user's "watched" groups
func HandleUnlink(update *api.Update, userChats map[int64]map[int64]struct{}) {
	chatID, err := GetChatID(update)
	if err != nil {
		log.Printf("[bot] ❌ Failed to extract chat ID: %v", err)
		return
	}

	userID, err := GetUserID(update)
	if err != nil {
		log.Printf("❌ Could not get user ID: %v", err)
		return
	}

	// Check if the user has any linked groups
	if groups, ok := userChats[userID]; ok {
		// Remove the group if it exists
		if _, exists := groups[chatID]; exists {
			delete(groups, chatID)
			err := api.SendTextMessage(chatID, "⛔ You've removed this group from your watches!")
			if err != nil {
				log.Printf("[bot] ❌ Failed to send unlink confirmation: %v", err)
			}
			log.Printf("[bot] ✅ user ID=%d unlinked a chat with ID=%d.", userID, chatID)
		} else {
			// If the group wasn't found in the user's watch list
			log.Printf("[bot] ⚠️ user ID=%d tried to unlink a chat with ID=%d, but it was not linked.", userID, chatID)
			err := api.SendTextMessage(chatID, "⚠️ This group is not linked to your account.")
			if err != nil {
				log.Printf("[bot] ❌ Failed to send unlink error: %v", err)
			}
		}
	} else {
		// If the user has no linked groups
		log.Printf("[bot] ⚠️ user ID=%d has no linked groups.", userID)
		err := api.SendTextMessage(chatID, "⚠️ You have no groups linked to your account.")
		if err != nil {
			log.Printf("[bot] ❌ Failed to send unlink error: %v", err)
		}
	}
}
