package interfaces

import (
	"context"
	"flex-frog-bot/db/domain"
	api "flex-frog-bot/tg-bot-api"
)

type UserService interface {
	Create(ctx context.Context, chatId int64, user *domain.User) error
	GetByID(ctx context.Context, userId int64) (*domain.User, error)
	GetUserID(update *api.Update) (int64, error)
	GetUserName(update *api.Update) (string, error)
	AssociateWithChat(ctx context.Context, userId int64, chatId int64) error
}
