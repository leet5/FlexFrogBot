package repository

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Image struct {
	ID        int64
	ImageName string
	ImageData []byte
	MessageID int64
	CreatedAt time.Time
}

type ImageRepository struct {
	DB *sql.DB
}

func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{DB: db}
}

// InsertImage inserts a new image into the database.
func (r *ImageRepository) InsertImage(ctx context.Context, image *Image) (int64, error) {
	var id int64
	query := `
		INSERT INTO images (image_name, image_data, message_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.DB.QueryRowContext(ctx, query, image.ImageName, image.ImageData, image.MessageID).Scan(&id)
	return id, err
}

// GetImageTags retrieves tag names associated with an image.
func (r *ImageRepository) GetImageTags(ctx context.Context, imageID int64) ([]string, error) {
	query := `
		SELECT t.name
		FROM tags t
		JOIN images_tags it ON t.id = it.tag_id
		WHERE it.image_id = $1
	`
	rows, err := r.DB.QueryContext(ctx, query, imageID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("[database] ⚠️ Error closing rows: %v", err)
		}
	}(rows)

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}
