package handlers

import (
	"context"
	"flex-frog-bot/services/interfaces"
	api "flex-frog-bot/tg-bot-api"
	"log"
)

func HandleStop(ctx context.Context, update *api.Update, chtSvc interfaces.ChatService, usrSvc interfaces.UserService) {
	chatID, err := chtSvc.GetChatID(update)
	if err != nil {
		log.Printf("[bot][handle_stop] ❌ Failed to extract chat ID: %v", err)
		return
	}

	userID, err := usrSvc.GetUserID(update)
	if err != nil {
		log.Printf("[bot][handle_stop] ❌ Failed to extract user ID: %v", err)
		return
	}

	isAdmin, err := api.IsUserAdmin(chatID, userID)
	if err != nil {
		log.Printf("[bot][handle_stop] ❌ Failed to check if user is admin: %v", err)
		return
	}

	if !isAdmin {
		log.Printf("[bot][handle_stop] ❌ User is not an admin in chat %d", chatID)
		err = api.SendTextMessage(chatID, "❌ Only admins can activate the bot.")
		if err != nil {
			log.Printf("[bot][handle_stop] ❌ Failed to send message: %v", err)
		}
		return
	}

	watched, err := chtSvc.IsWatched(ctx, chatID)
	if err != nil {
		log.Printf("[bot][handle_stop] ❌ Failed to check if chat is already watched: %v", err)
		return
	}

	if watched {
		if err := chtSvc.Unwatch(ctx, chatID); err != nil {
			log.Printf("[bot][handle_stop] ❌ Failed to unwatch chat '%d': %v", chatID, err)
			return
		}
		log.Printf("[bot][handle_stop] ⛔ Stop button pressed in chat_id=%d. Bot deactivated.", chatID)
		err = api.SendTextMessage(chatID, "⛔ Bot has been deactivated.")
		if err != nil {
			log.Printf("[bot][handle_stop] ❌ Failed to send stop confirmation: %v", err)
		}
	} else {
		log.Printf("[bot][handle_stop] ❌ Chat '%d' is not watched. Bot is already inactive.", chatID)
		err = api.SendTextMessage(chatID, "❌ Bot is already inactive in this chat.")
		if err != nil {
			log.Printf("[bot][handle_stop] ❌ Failed to send message: %v", err)
		}
		return
	}
}
