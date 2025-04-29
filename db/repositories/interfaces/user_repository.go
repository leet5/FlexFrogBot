package interfaces

import (
	"context"
	"flex-frog-bot/db/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetChatsByUserID(ctx context.Context, userID int64) (map[int64]struct{}, error)
	AssociateWithChat(ctx context.Context, userID int64, chatID int64) error
}
