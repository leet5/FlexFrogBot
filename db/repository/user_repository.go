package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type User struct {
	ID       int64
	Username string
}

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// InsertUser inserts a new user into the database.
func (r *UserRepository) InsertUser(ctx context.Context, userID int64, username string, thumbnail []byte) (int64, error) {
	var id int64
	query := `
		INSERT INTO users (id, username, thumbnail)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.DB.QueryRowContext(ctx, query, userID, username, thumbnail).Scan(&id)
	return id, err
}

// AssociateUserWithChat associates a user with a chat in the database.
func (r *UserRepository) AssociateUserWithChat(ctx context.Context, userID int64, chatID int64) error {
	query := `
		INSERT INTO users_chats (user_id, chat_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, chat_id) DO NOTHING
	`
	_, err := r.DB.ExecContext(ctx, query, userID, chatID)
	return err
}

// GetUserChats retrieves all chats a user is part of and returns them as a map.
func (r *UserRepository) GetUserChats(ctx context.Context, userID int64) (map[int64]struct{}, error) {
	query := `SELECT chat_id FROM user_chats WHERE user_id = ?`

	rows, err := r.DB.QueryContext(ctx, query, userID)
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

// GetUserByID retrieves a user by their ID.
func (r *UserRepository) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	query := `SELECT id, username FROM users WHERE id = $1`
	var user User
	err := r.DB.QueryRowContext(ctx, query, userID).Scan(&user.ID, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to query user by ID: %w", err)
	}
	return &user, nil
}
