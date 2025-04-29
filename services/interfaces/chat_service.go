package interfaces

import (
	"context"
	"flex-frog-bot/db/domain"
	api "flex-frog-bot/tg-bot-api"
)

type ChatService interface {
	IsWatched(ctx context.Context, chatID int64) (bool, error)
	Watch(ctx context.Context, chatID int64) error
	Unwatch(ctx context.Context, chatID int64) error
	GetByID(ctx context.Context, chatID int64) (*domain.Chat, error)
	GetOrCreate(ctx context.Context, update *api.Update) (*domain.Chat, error)
	GetChatID(update *api.Update) (int64, error)
	GetChatName(update *api.Update) (string, error)
}
