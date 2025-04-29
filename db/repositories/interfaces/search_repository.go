package interfaces

import (
	"context"
	"flex-frog-bot/dto"
)

type SearchRepository interface {
	SearchImagesByTags(ctx context.Context, tags []string) ([]*dto.ImageDTO, error)
}
