package dto

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type ImageDTO struct {
	UserThumbnail  []byte    `json:"user_thumbnail"`
	ImageThumbnail []byte    `json:"image_thumbnail"`
	ChatThumbnail  []byte    `json:"chat_thumbnail"`
	MessageID      int64     `json:"message_id"`
	UserName       string    `json:"user_name"`
	ChatName       string    `json:"chat_name"`
	CreatedAt      time.Time `json:"created_at"`
}

func (dto *ImageDTO) MarshalJSON() ([]byte, error) {
	type Alias ImageDTO
	return json.Marshal(&struct {
		UserThumbnail  string `json:"user_thumbnail"`
		ImageThumbnail string `json:"image_thumbnail"`
		ChatThumbnail  string `json:"chat_thumbnail"`
		*Alias
	}{
		UserThumbnail:  base64.StdEncoding.EncodeToString(dto.UserThumbnail),
		ImageThumbnail: base64.StdEncoding.EncodeToString(dto.ImageThumbnail),
		ChatThumbnail:  base64.StdEncoding.EncodeToString(dto.ChatThumbnail),
		Alias:          (*Alias)(dto),
	})
}
