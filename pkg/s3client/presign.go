package s3client

import (
	"context"
	"fmt"
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

	queryParams := url.Values{}
	if req.Filename != "" {
		queryParams.Set(
			"response-content-disposition",
			fmt.Sprintf("attachment; filename=\"%s\"", req.Filename),
		)
	}

	presignedURL, err := storage.client.PresignedGetObject(
		ctx,
		bucket,
		objectPath,
		req.Expiry,
		queryParams,
	)
	if err != nil {
		return dto.CreatePresignedURLResponse{}, err
	}

	return dto.CreatePresignedURLResponse{
		URL: presignedURL.String(),
	}, nil
}
