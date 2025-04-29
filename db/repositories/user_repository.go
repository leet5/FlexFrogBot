package repositories

import (
	"context"
	"database/sql"
	"errors"
	"flex-frog-bot/db/domain"
	"flex-frog-bot/db/repositories/interfaces"
	"fmt"
	"log"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) interfaces.UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database.
func (r *userRepository) Create(ctx context.Context, user *domain.User) (int64, error) {
	var id int64
	query := `
		INSERT INTO users (id, name, thumbnail)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, user.Id, user.Name, user.Thumbnail).Scan(&id)
	return id, err
}

// GetByID retrieves a user by their ID.
func (r *userRepository) GetByID(ctx context.Context, userID int64) (*domain.User, error) {
	query := `SELECT id, name, thumbnail FROM users WHERE id = $1`
	var user domain.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&user.Id, &user.Name, &user.Thumbnail)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to query user by ID: %w", err)
	}
	return &user, nil
}

// GetChatsByUserID retrieves all chats a user is part of and returns them as a map.
func (r *userRepository) GetChatsByUserID(ctx context.Context, userID int64) (map[int64]struct{}, error) {
	query := `SELECT chat_id FROM users_chats WHERE user_id = ?`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("[database] Error closing rows: %v", err)
		}
	}(rows)

	chats := make(map[int64]struct{})
	for rows.Next() {
		var chatID int64
		if err := rows.Scan(&chatID); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		chats[chatID] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return chats, nil
}

// AssociateWithChat associates a user with a chat in the database.
func (r *userRepository) AssociateWithChat(ctx context.Context, userID int64, chatID int64) error {
	query := `
		INSERT INTO users_chats (user_id, chat_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, chat_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, userID, chatID)
	return err
}
