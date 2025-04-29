package interfaces

import (
	"context"
	"flex-frog-bot/db/domain"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *domain.Chat) error
	GetByID(ctx context.Context, chatID int64) (*domain.Chat, error)
	IsWatched(ctx context.Context, chatID int64) (bool, error)
	UnwatchChat(ctx context.Context, chatID int64) error
	WatchChat(ctx context.Context, chatID int64) error
	Exists(ctx context.Context, chatID int64) (bool, error)
}
