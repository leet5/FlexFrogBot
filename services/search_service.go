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

func (svc *searchService) GetImagesByTags(ctx context.Context, tags []string) ([]*dto.ImageDTO, error) {
	return svc.searchRepo.SearchImagesByTags(ctx, tags)
}
