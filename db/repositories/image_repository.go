package repositories

import (
	"context"
	"database/sql"
	"errors"
	"flex-frog-bot/db/domain"
	"flex-frog-bot/db/repositories/interfaces"
	"fmt"
	"github.com/lib/pq"
	"log"
	"time"
)

type imageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) interfaces.ImageRepository {
	return &imageRepository{db: db}
}

// Create inserts a new image into the database.
func (r *imageRepository) Create(ctx context.Context, img *domain.Image) (int64, error) {
	var id int64
	query := `
		INSERT INTO images (data, thumbnail, message_id, user_id, chat_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, img.Data, img.Thumbnail, img.MessageId, img.UserId, img.ChatId).Scan(&id)
	return id, err
}

// GetByID retrieves an image by its ID from the database.
func (r *imageRepository) GetByID(ctx context.Context, id int64) (*domain.Image, error) {
	var img domain.Image
	query := `
		SELECT id, data, thumbnail, message_id, user_id, chat_id
		FROM images
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&img.Id, &img.Data, &img.Thumbnail, &img.MessageId, &img.UserId, &img.ChatId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get image by ID: %w", err)
	}
	return &img, nil
}

// GetImageIDsByTags retrieves image IDs associated with specific tags.
func (r *imageRepository) GetImageIDsByTags(ctx context.Context, tags []string) ([]int64, error) {
	query := `
        SELECT i.id, MAX(i.created_at) AS created_at
		FROM images i
		INNER JOIN images_tags it ON i.id = it.image_id
		INNER JOIN tags t ON it.tag_id = t.id
		WHERE t.name = ANY ($1)
		GROUP BY i.id
		HAVING array_agg(t.name) @> $1::text[]
		ORDER BY created_at DESC
    `

	rows, err := r.db.QueryContext(ctx, query, pq.Array(tags))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("[database] ⚠️ Error closing rows: %v", err)
		}
	}(rows)

	var imageIDs []int64
	for rows.Next() {
		var id int64
		var createdAt time.Time
		if err := rows.Scan(&id, &createdAt); err != nil {
			return nil, err
		}
		imageIDs = append(imageIDs, id)
	}
	return imageIDs, rows.Err()
}
