package dto

import (
	"encoding/base64"
	"encoding/json"
)

type ChatDTO struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Username  string `json:"username,omitempty"`
	Thumbnail []byte `json:"thumbnail,omitempty"`
	IsPrivate bool   `json:"is_private"`
}

func (dto *ChatDTO) MarshalJSON() ([]byte, error) {
	type Alias ChatDTO
	return json.Marshal(&struct {
		Thumbnail string `json:"thumbnail"`
		*Alias
	}{
		Thumbnail: base64.StdEncoding.EncodeToString(dto.Thumbnail),
		Alias:     (*Alias)(dto),
	})
}
