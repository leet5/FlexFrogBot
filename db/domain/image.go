package domain

import "time"

type Image struct {
	Id        int64
	Data      []byte
	Thumbnail []byte
	MessageId int64
	UserId    int64
	ChatId    int64
	CreatedAt time.Time
}
