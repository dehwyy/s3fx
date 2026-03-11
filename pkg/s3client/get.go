package s3client

import (
	"context"

	"github.com/dehwyy/s3fx/pkg/dto"
	"github.com/minio/minio-go/v7"
)

func (storage *MinioStorage) Get(
	ctx context.Context,
	req dto.GetRequest,
) (dto.GetResponse, error) {
	bucket, objectPath, err := parseURL(req.URL)
	if err != nil {
		return dto.GetResponse{}, err
	}

	object, err := storage.client.GetObject(
		ctx,
		bucket,
		objectPath,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return dto.GetResponse{}, err
	}

	return dto.GetResponse{
		Object: object, // *minio.Object inherently implements io.ReadCloser
	}, nil
}
