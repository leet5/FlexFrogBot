package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Chat struct {
	ID       int64  `db:"id"`
	ChatName string `db:"chat_name"`
	Watched  bool   `db:"watched"`
}

type ChatRepository struct {
	DB *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{DB: db}
}

// GetChatByID retrieves a chat by its ID.
func (r *ChatRepository) GetChatByID(ctx context.Context, chatID int64) (*Chat, error) {
	var chat Chat
	query := `
		SELECT id, chat_name, watched
		FROM chats
		WHERE id = $1
	`
	err := r.DB.QueryRowContext(ctx, query, chatID).Scan(&chat.ID, &chat.ChatName, &chat.Watched)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Chat not found
		}
		return nil, fmt.Errorf("failed to get chat: %v", err)
	}
	return &chat, nil
}

// InsertChat inserts a new chat into the database.
func (r *ChatRepository) InsertChat(ctx context.Context, chatID int64, chatName string, watched bool) error {
	var id int64
	query := `
		INSERT INTO chats (id, chat_name, watched)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	if err := r.DB.QueryRowContext(ctx, query, chatID, chatName, watched).Scan(&id); err != nil {
		return err
	}
	return nil
}

// CheckIfChatWatched checks if a chat is watched.
func (r *ChatRepository) CheckIfChatWatched(ctx context.Context, chatID int64) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 
			FROM chats 
			WHERE id = $1 AND watched = true
		)
	`
	err := r.DB.QueryRowContext(ctx, query, chatID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if chat is watched: %v", err)
	}
	return exists, nil
}

// UnwatchChat sets the watched status of a chat to false.
func (r *ChatRepository) UnwatchChat(ctx context.Context, chatID int64) error {
	query := `
		UPDATE chats 
		SET watched = false 
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, chatID)
	if err != nil {
		return fmt.Errorf("failed to unwatch chat: %v", err)
	}
	return nil
}

// WatchChat sets the watched status of a chat to true.
func (r *ChatRepository) WatchChat(ctx context.Context, chatID int64) error {
	query := `
		UPDATE chats 
		SET watched = true 
		WHERE id = $1
	`
	_, err := r.DB.ExecContext(ctx, query, chatID)
	if err != nil {
		return fmt.Errorf("failed to watch chat: %v", err)
	}
	return nil
}
