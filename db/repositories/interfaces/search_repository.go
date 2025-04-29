package interfaces

import (
	"context"
	"flex-frog-bot/dto"
)

type SearchRepository interface {
	SearchImagesByChatIdByTags(ctx context.Context, chatID int64, tags []string) ([]*dto.ImageDTO, error)
	SearchChatsByUserID(ctx context.Context, userID string) ([]*dto.ChatDTO, error)
}
