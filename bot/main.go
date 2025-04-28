package bot

import (
	"context"
	"flex-frog-bot/db/repository"
	tgbotapi "flex-frog-bot/tg-bot-api"
)

const (
	botName = "img_srch_bot"
)

var (
	ImgRepo  *repository.ImageRepository
	ChatRepo *repository.ChatRepository
	UserRepo *repository.UserRepository
)

func RunBot(ctx context.Context, imgRepo *repository.ImageRepository, chatRepo *repository.ChatRepository, userRepo *repository.UserRepository) {
	ImgRepo = imgRepo
	ChatRepo = chatRepo
	UserRepo = userRepo

	updates := tgbotapi.GetUpdatesChan()
	processUpdates(ctx, updates)
}
