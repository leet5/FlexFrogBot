package repositories

import (
	"context"
	"database/sql"
	"errors"
	"flex-frog-bot/db/repositories/interfaces"
	"flex-frog-bot/dto"
	"fmt"
	"github.com/lib/pq"
	"log"
)

type searchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) interfaces.SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) SearchImagesByChatIdByTags(ctx context.Context, chatID int64, tags []string) ([]*dto.ImageDTO, error) {
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
		WHERE chat_id = $1 AND t.name = ANY ($2)
		GROUP BY
		    u.thumbnail,
		    u.name,
		    i.message_id,
		    i.created_at,
		    i.thumbnail,
		    c.name,
		    c.thumbnail
		HAVING array_agg(t.name) @> ARRAY[$2]
		ORDER BY i.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, chatID, pq.Array(tags))
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

func (r *searchRepository) SearchChatsByUserID(ctx context.Context, userID string) ([]*dto.ChatDTO, error) {
	query := `
		SELECT c.id, c.name, c.thumbnail
		FROM chats c
		JOIN users_chats uc ON uc.chat_id = c.id
		WHERE uc.user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("query chats by user_id: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("[search_repository][SearchChatsByUserID] ⚠️ Error closing rows: %v", err)
		}
	}()

	var chatDTOs []*dto.ChatDTO
	for rows.Next() {
		var chatDTO dto.ChatDTO
		if err := rows.Scan(&chatDTO.ID, &chatDTO.Name, &chatDTO.Thumbnail); err != nil {
			return nil, fmt.Errorf("scan chat row: %w", err)
		}
		chatDTOs = append(chatDTOs, &chatDTO)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return chatDTOs, nil
}
