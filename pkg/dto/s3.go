package dto

import "io"

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
