package interfaces

import (
	"context"
	"flex-frog-bot/dto"
)

type SearchService interface {
	GetImagesByChatIdByTags(ctx context.Context, chatID int64, tags []string) ([]*dto.ImageDTO, error)
	GetChatsByUserID(ctx context.Context, userID string) ([]*dto.ChatDTO, error)
}
