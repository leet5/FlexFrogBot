package interfaces

import (
	"context"
	"flex-frog-bot/db/domain"
	api "flex-frog-bot/tg-bot-api"
)

type ImageService interface {
	Create(ctx context.Context, update *api.Update, image *domain.Image) error
	GetIdsByTags(ctx context.Context, tags []string) ([]int64, error)
	Download(update *api.Update) (string, error)
	HasImage(update *api.Update) bool
	GetByID(ctx context.Context, id int64) (*domain.Image, error)
}
