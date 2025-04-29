package repositories

import (
	"context"
	"database/sql"
	"errors"
	"flex-frog-bot/db/repositories/interfaces"
	"flex-frog-bot/dto"
	"github.com/lib/pq"
	"log"
)

type searchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) interfaces.SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) SearchImagesByTags(ctx context.Context, tags []string) ([]*dto.ImageDTO, error) {
	query := `
		SELECT
		    u.thumbnail  AS user_thumbnail,
		    u.name	     AS user_name,
		    i.message_id AS message_id,
		    i.created_at AS created_at,
		    i.thumbnail  AS image_thumbnail,
		    c.name       AS chat_name,
		    c.thumbnail  AS chat_thumbnail
		FROM images i
		INNER JOIN users u ON i.user_id = u.id
		INNER JOIN chats c ON i.chat_id = c.id
		INNER JOIN images_tags it ON i.id = it.image_id
		INNER JOIN tags t ON it.tag_id = t.id
		WHERE t.name = ANY ($1)
		GROUP BY
		    u.thumbnail,
		    u.name,
		    i.message_id,
		    i.created_at,
		    i.thumbnail,
		    c.name,
		    c.thumbnail
		HAVING array_agg(t.name) @> ARRAY[$1]
		ORDER BY i.created_at DESC
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

	var imageDTOs []*dto.ImageDTO
	for rows.Next() {
		imageDTO := &dto.ImageDTO{}
		err := rows.Scan(&imageDTO.UserThumbnail, &imageDTO.UserName, &imageDTO.MessageID, &imageDTO.CreatedAt, &imageDTO.ImageThumbnail, &imageDTO.ChatName, &imageDTO.ChatThumbnail)

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
