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

type chatService struct {
	chatRepo repos.ChatRepository
}

func NewChatService(chatRepository repos.ChatRepository) services.ChatService {
	return &chatService{
		chatRepo: chatRepository,
	}
}

func (svc *chatService) IsWatched(ctx context.Context, chatID int64) (bool, error) {
	watched, err := svc.chatRepo.IsWatched(ctx, chatID)
	if err != nil {
		return false, err
	}
	return watched, nil
}

func (svc *chatService) Watch(ctx context.Context, chatID int64) error {
	err := svc.chatRepo.WatchChat(ctx, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (svc *chatService) Unwatch(ctx context.Context, chatID int64) error {
	err := svc.chatRepo.UnwatchChat(ctx, chatID)
	if err != nil {
		return err
	}
	return nil
}

func (svc *chatService) GetByID(ctx context.Context, chatID int64) (*domain.Chat, error) {
	chat, err := svc.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (svc *chatService) GetOrCreate(ctx context.Context, update *api.Update) (*domain.Chat, error) {
	chatID, err := svc.GetChatID(update)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat ID: %v", err)
	}

	exists, err := svc.chatRepo.Exists(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if chat exists: %v", err)
	}

	if exists {
		chat, err := svc.chatRepo.GetByID(ctx, chatID)
		if err != nil {
			return nil, fmt.Errorf("failed to get chat by ID: %v", err)
		}
		return chat, nil
	}

	chatName, err := svc.GetChatName(update)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat name: %v", err)
	}

	photoPath, err := api.GetChatProfilePhoto(saveDir, img_tools.GenerateUUID(), chatID)
	if err != nil {
		log.Printf("[chat_service][add_if_absent] ⚠️ Error getting chat profile photo: %v", err)
	}

	thumbnail, err := img_tools.CreateThumbnailByPath(photoPath)
	if err != nil {
		log.Printf("[chat_service][add_if_absent] ⚠️ Error creating thumbnail: %v", err)
	}

	isPrivate, err := api.IsChatPrivate(chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to check whether chat is private: %v", err)
	}

	chat := &domain.Chat{
		Id:        chatID,
		Name:      chatName,
		Thumbnail: thumbnail,
		Watched:   false,
		IsPrivate: isPrivate,
	}

	err = svc.chatRepo.Create(ctx, chat)
	if err != nil {
		return nil, fmt.Errorf("failed to insert chat '%s' into DB: %v", chatName, err)
	}
	return chat, nil
}

func (svc *chatService) GetChatID(update *api.Update) (int64, error) {
	if update.Message != nil {
		return update.Message.Chat.ID, nil
	}
	if update.Callback != nil && update.Callback.Message != nil {
		return update.Callback.Message.Chat.ID, nil
	}
	if update.MyChatMember != nil {
		return update.MyChatMember.Chat.ID, nil
	}
	return 0, errors.New("chat ID not found in update")
}

func (svc *chatService) GetChatName(update *api.Update) (string, error) {
	if update.Message != nil {
		return update.Message.Chat.Title, nil
	}
	if update.Callback != nil && update.Callback.Message != nil {
		return update.Callback.Message.Chat.Title, nil
	}
	return "", errors.New("chat name not found in update")
}
