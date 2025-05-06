package domain

type Chat struct {
	Id        int64
	Name      string
	Thumbnail []byte
	Watched   bool
	IsPrivate bool
}
