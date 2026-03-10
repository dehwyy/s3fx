package client

import (
	"errors"
	"net/url"
	"strings"
)

func parseURL(
	URL string,
) (
	bucket string,
	path string,
	err error,
) {
	parsedURL, err := url.Parse(URL)
	if err != nil {
		err = errors.Join(ErrMalformedURL, err)
		return
	}

	pathParts := strings.SplitN(strings.TrimPrefix(parsedURL.Path, "/"), "/", 2)
	if len(pathParts) < 2 {
		err = ErrMalformedURL
		return
	}

	bucket = pathParts[0]
	path = pathParts[1]

	return
}
