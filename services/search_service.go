package services

import (
	"context"
	"flex-frog-bot/db/repositories/interfaces"
	"flex-frog-bot/dto"
	interfaces2 "flex-frog-bot/services/interfaces"
)

type searchService struct {
	searchRepo interfaces.SearchRepository
}

func NewSearchService(searchRepo interfaces.SearchRepository) interfaces2.SearchService {
	return &searchService{
		searchRepo: searchRepo,
	}
}

func (svc *searchService) GetImagesByChatIdByTags(ctx context.Context, chatID int64, tags []string) ([]*dto.ImageDTO, error) {
	return svc.searchRepo.SearchImagesByChatIdByTags(ctx, chatID, tags)
}

func (svc *searchService) GetChatsByUserID(ctx context.Context, userID string) ([]*dto.ChatDTO, error) {
	return svc.searchRepo.SearchChatsByUserID(ctx, userID)
}
