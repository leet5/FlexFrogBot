package interfaces

import (
	"context"
	"flex-frog-bot/db/domain"
)

type ImageRepository interface {
	Create(ctx context.Context, image *domain.Image) (int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Image, error)
	GetImageIDsByTags(ctx context.Context, tags []string) ([]int64, error)
}
