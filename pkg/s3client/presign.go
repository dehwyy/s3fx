package s3client

import (
	"context"
	"net/url"

	"github.com/dehwyy/s3fx/pkg/dto"
)

func (storage *MinioStorage) CreatePresignedURL(
	ctx context.Context,
	req dto.CreatePresignedURLRequest,
) (dto.CreatePresignedURLResponse, error) {
	bucket, objectPath, err := parseURL(req.URL)
	if err != nil {
		return dto.CreatePresignedURLResponse{}, err
	}

	presignedURL, err := storage.client.PresignedGetObject(
		ctx,
		bucket,
		objectPath,
		req.Expiry,
		url.Values{},
	)
	if err != nil {
		return dto.CreatePresignedURLResponse{}, err
	}

	return dto.CreatePresignedURLResponse{
		URL: presignedURL.String(),
	}, nil
}
