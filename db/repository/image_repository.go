package repository

import (
	"context"
	"database/sql"
	"errors"
	"flex-frog-bot/dto"
	"flex-frog-bot/img_tools"
	"github.com/lib/pq"
	"log"
	"time"
)

type Image struct {
	ID        int64
	ImageData []byte
	UserID    int64
	ChatID    int64
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
	thumbnail, err := img_tools.CreateThumbnail(image.ImageData)
	if err != nil {
		return 0, err
	}

	var id int64
	query := `
		INSERT INTO images (image_data, thumbnail, message_id, user_id, chat_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err = r.DB.QueryRowContext(ctx, query, image.ImageData, thumbnail, image.MessageID, image.UserID, image.ChatID).Scan(&id)
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

func (r *ImageRepository) GetImagesByTags(ctx context.Context, tags []string) ([]*dto.ImageDTO, error) {
	query := `
		SELECT
		    u.thumbnail AS user_thumbnail,
		    u.username  AS user_name,
		    i.message_id,
		    i.created_at,
		    i.thumbnail AS image_thumbnail,
		    c.chat_name,
		    c.thumbnail AS chat_thumbnail
		FROM images i
		INNER JOIN users u ON i.user_id = u.id
		INNER JOIN chats c ON i.chat_id = c.id
		INNER JOIN images_tags it ON i.id = it.image_id
		INNER JOIN tags t ON it.tag_id = t.id
		WHERE t.name = ANY ($1)
		GROUP BY
		    u.thumbnail,
		    u.username,
		    i.message_id,
		    i.created_at,
		    i.thumbnail,
		    c.chat_name,
		    c.thumbnail
		HAVING array_agg(t.name) @> ARRAY[$1]
		ORDER BY i.created_at DESC
	`

	rows, err := r.DB.QueryContext(ctx, query, pq.Array(tags))
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("[database] ⚠️ Error closing rows: %v", err)
		}
	}(rows)

	var imageDTOs []*dto.ImageDTO
	for rows.Next() {
		imageDTO := &dto.ImageDTO{}
		err := rows.Scan(
			&imageDTO.UserThumbnail,
			&imageDTO.UserName,
			&imageDTO.MessageID,
			&imageDTO.CreatedAt,
			&imageDTO.ImageThumbnail,
			&imageDTO.ChatName,
			&imageDTO.ChatThumbnail,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Printf("[database] ⚠️ No images found for tags: %v", tags)
				return nil, nil
			} else {
				log.Printf("[database] ⚠️ Error scanning row: %v", err)
				return nil, err
			}
		}
		imageDTOs = append(imageDTOs, imageDTO)
	}
	return imageDTOs, rows.Err()
}
