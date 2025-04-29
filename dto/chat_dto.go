package dto

import (
	"encoding/base64"
	"encoding/json"
)

type ChatDTO struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Thumbnail []byte `json:"thumbnail,omitempty"`
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
