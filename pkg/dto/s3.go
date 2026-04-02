package dto

import (
	"io"
	"time"
)

type CreateRequest struct {
	URL  string
	Data io.Reader
	Size int64
}

type GetRequest struct {
	URL string
}

type GetResponse struct {
	Object io.ReadCloser
}

type CreatePresignedURLRequest struct {
	URL    string
	Expiry time.Duration
}

type CreatePresignedURLResponse struct {
	URL string
}
