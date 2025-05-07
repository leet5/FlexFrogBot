package domain

type Chat struct {
	Id        int64
	Title     string
	Username  string
	Thumbnail []byte
	Watched   bool
	IsPrivate bool
}
