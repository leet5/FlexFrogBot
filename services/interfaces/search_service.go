package interfaces

import (
	"context"
	"flex-frog-bot/dto"
)

type SearchService interface {
	GetImagesByTags(ctx context.Context, tags []string) ([]*dto.ImageDTO, error)
}
