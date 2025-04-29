package services

import (
	"context"
	"flex-frog-bot/db/domain"
	repositories "flex-frog-bot/db/repositories/interfaces"
	"flex-frog-bot/img_tools"
	services "flex-frog-bot/services/interfaces"
	api "flex-frog-bot/tg-bot-api"
	"fmt"
	"log"
	"os"
)

const (
	saveDir = "./images"
)

type imageService struct {
	imgRepo repositories.ImageRepository
}

func NewImageService(imgRepo repositories.ImageRepository) services.ImageService {
	return &imageService{
		imgRepo: imgRepo,
	}
}

func (svc *imageService) Create(ctx context.Context, update *api.Update, image *domain.Image) error {
	filepath, err := svc.Download(update)
	if err != nil {
		return fmt.Errorf("download image from update error: %v", err)
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read image file %s: %w", filepath, err)
	}

	thumbnail, err := img_tools.CreateThumbnail(data)
	if err != nil {
		return fmt.Errorf("failed to create thumbnail: %w", err)
	}

	_, err = svc.imgRepo.Create(ctx, &domain.Image{
		Data:      data,
		Thumbnail: thumbnail,
		MessageId: image.MessageId,
		UserId:    image.UserId,
		ChatId:    image.ChatId,
	})
	if err != nil {
		return fmt.Errorf("failed to insert image into DB: %w", err)
	}

	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("failed to delete image file %s: %w", filepath, err)
	}
	log.Printf("[image_service][store_image] 🗑 Deleted image file %s", filepath)

	return nil
}

func (svc *imageService) Download(update *api.Update) (string, error) {
	if update.Message == nil {
		return "", fmt.Errorf("no message content found")
	}

	switch {
	case len(update.Message.Photo) > 0:
		largest := update.Message.Photo[len(update.Message.Photo)-1]
		filename := fmt.Sprintf("%d_photo.jpg", update.Message.MessageID)
		return api.DownloadFile(saveDir, largest.FileID, filename)

	case update.Message.Document != nil && isImage(update.Message.Document.MimeType):
		doc := update.Message.Document
		filename := fmt.Sprintf("%d_%s", update.Message.MessageID, doc.FileName)
		return api.DownloadFile(saveDir, doc.FileID, filename)

	default:
		return "", fmt.Errorf("no valid image found in message")
	}
}

func isImage(mime string) bool {
	return mime == "image/jpeg" || mime == "image/png" || mime == "image/gif" || mime == "image/webp"
}

func (svc *imageService) HasImage(update *api.Update) bool {
	if update.Message == nil {
		return false
	}
	if update.Message.Document != nil && isImage(update.Message.Document.MimeType) {
		return true
	}
	return len(update.Message.Photo) > 0
}

func (svc *imageService) GetIdsByTags(ctx context.Context, tags []string) ([]int64, error) {
	ids, err := svc.imgRepo.GetImageIDsByTags(ctx, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to get image IDs by tags: %w", err)
	}
	return ids, nil
}

func (svc *imageService) GetByID(ctx context.Context, id int64) (*domain.Image, error) {
	image, err := svc.imgRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get image by ID: %w", err)
	}
	return image, nil
}
