package bot

import (
	"context"
	interfaces2 "flex-frog-bot/services/interfaces"
	tgbotapi "flex-frog-bot/tg-bot-api"
)

const (
	botName = "FlexFrogBot"
)

var (
	ImageService interfaces2.ImageService
	ChatService  interfaces2.ChatService
	UserService  interfaces2.UserService
)

func RunBot(ctx context.Context, imgService interfaces2.ImageService, chatService interfaces2.ChatService, userService interfaces2.UserService) {
	ImageService = imgService
	ChatService = chatService
	UserService = userService

	updates := tgbotapi.GetUpdatesChan()
	processUpdates(ctx, updates)
}
