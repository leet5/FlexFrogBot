package handlers

import (
	"fmt"
	"log"
	api "mini-app-back/tg-bot-api"
)

// HandleLink appends the specified group to the user's "watched" groups
func HandleLink(update *api.Update, userGroups map[int64]map[int64]struct{}) {
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

	if _, ok := userGroups[userID]; !ok {
		userGroups[userID] = make(map[int64]struct{})
	}

	userGroups[userID][chatID] = struct{}{}

	username, err := GetUserName(update)
	if err != nil {
		log.Printf("[bot] ❌ Could not get username: %v", err)
		username = ""
	}

	err = api.SendTextMessage(chatID, fmt.Sprintf("✅ %s, you've added this group to your watches!", username))
	if err != nil {
		log.Printf("[bot] ❌ Failed to send link confirmation: %v", err)
	}
	log.Printf("[bot] ✅ user ID=%d linked a chat with ID=%d.", userID, chatID)
}
