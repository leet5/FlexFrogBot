package services

import (
	"context"
	"errors"
	"flex-frog-bot/db/domain"
	repos "flex-frog-bot/db/repositories/interfaces"
	"flex-frog-bot/img_tools"
	services "flex-frog-bot/services/interfaces"
	api "flex-frog-bot/tg-bot-api"
	"fmt"
	"log"
)

type userService struct {
	userRepo repos.UserRepository
}

func NewUserService(userRepository repos.UserRepository) services.UserService {
	return &userService{
		userRepo: userRepository,
	}
}

func (svc *userService) getThumbnail(userID int64) ([]byte, error) {
	photoPath, err := api.GetUserProfilePhoto(saveDir, img_tools.GenerateUUID(), userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user profile photo: %w", err)
	}

	thumbnail, err := img_tools.CreateThumbnailByPath(photoPath)
	if err != nil {
		return nil, fmt.Errorf("error creating thumbnail: %w", err)
	}

	return thumbnail, nil
}

func (svc *userService) Create(ctx context.Context, chatId int64, user *domain.User) error {
	thumbnail, err := svc.getThumbnail(user.Id)
	if err != nil {
		log.Printf("[user_service][create_user] ⚠️ Error getting thumbnail: %v", err)
	}
	user.Thumbnail = thumbnail

	_, err = svc.userRepo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	err = svc.userRepo.AssociateWithChat(ctx, user.Id, chatId)
	if err != nil {
		return fmt.Errorf("error associating user with chat: %w", err)
	}
	return nil
}

func (svc *userService) GetByID(ctx context.Context, userId int64) (*domain.User, error) {
	user, err := svc.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting user by ID: %w", err)
	}
	return user, nil
}

func (svc *userService) GetUserID(update *api.Update) (int64, error) {
	if update.Message != nil && update.Message.From != nil {
		return update.Message.From.ID, nil
	}
	if update.Callback != nil && update.Callback.From != nil {
		return update.Callback.From.ID, nil
	}
	return 0, errors.New("user ID not found in update")
}

func (svc *userService) GetUserName(update *api.Update) (string, error) {
	if update.Message != nil && update.Message.From != nil {
		user := update.Message.From
		if user.Username != "" {
			return user.Username, nil
		}
		return user.FirstName, nil // fallback if username is missing
	}
	if update.Callback != nil && update.Callback.From != nil {
		user := update.Callback.From
		if user.Username != "" {
			return user.Username, nil
		}
		return user.FirstName, nil
	}
	return "", errors.New("username not found in update")
}

func (svc *userService) AssociateWithChat(ctx context.Context, userId int64, chatId int64) error {
	err := svc.userRepo.AssociateWithChat(ctx, userId, chatId)
	if err != nil {
		return fmt.Errorf("error associating user with chat: %w", err)
	}
	return nil
}
