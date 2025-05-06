package repositories

import (
	"context"
	"database/sql"
	"errors"
	"flex-frog-bot/db/domain"
	"flex-frog-bot/db/repositories/interfaces"
	"fmt"
)

type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) interfaces.ChatRepository {
	return &chatRepository{db: db}
}

// Create inserts a new chat into the database.
func (r *chatRepository) Create(ctx context.Context, chat *domain.Chat) error {
	var id int64
	query := `
		INSERT INTO chats (id, name, thumbnail, watched, is_private)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	if err := r.db.QueryRowContext(ctx, query, chat.Id, chat.Name, chat.Thumbnail, chat.Watched, chat.IsPrivate).Scan(&id); err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a chat by its ID.
func (r *chatRepository) GetByID(ctx context.Context, chatID int64) (*domain.Chat, error) {
	var chat domain.Chat
	query := `
		SELECT id, name, thumbnail, watched, is_private
		FROM chats
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, chatID).Scan(&chat.Id, &chat.Name, &chat.Thumbnail, &chat.Watched, &chat.IsPrivate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Chat not found
		}
		return nil, fmt.Errorf("failed to get chat: %v", err)
	}
	return &chat, nil
}

// IsWatched checks if a chat is watched.
func (r *chatRepository) IsWatched(ctx context.Context, chatID int64) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM chats 
			WHERE id = $1 AND watched = true
		)
	`
	err := r.db.QueryRowContext(ctx, query, chatID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if chat is watched: %v", err)
	}
	return exists, nil
}

// UnwatchChat sets the watched status of a chat to false.
func (r *chatRepository) UnwatchChat(ctx context.Context, chatID int64) error {
	query := `
		UPDATE chats 
		SET watched = false 
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, chatID)
	if err != nil {
		return fmt.Errorf("failed to unwatch chat: %v", err)
	}
	return nil
}

// WatchChat sets the watched status of a chat to true.
func (r *chatRepository) WatchChat(ctx context.Context, chatID int64) error {
	query := `
		UPDATE chats 
		SET watched = true 
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, chatID)
	if err != nil {
		return fmt.Errorf("failed to watch chat: %v", err)
	}
	return nil
}

// Exists checks if a chat exists in the database.
func (r *chatRepository) Exists(ctx context.Context, chatID int64) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM chats 
			WHERE id = $1
		)
	`
	err := r.db.QueryRowContext(ctx, query, chatID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if chat exists: %v", err)
	}
	return exists, nil
}
