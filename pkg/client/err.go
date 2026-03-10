package client

import "errors"

var (
	ErrMalformedURL = errors.New("invalid url: must be `scheme://host/bucket/file.ext`")
)
