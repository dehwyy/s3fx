package s3client

import (
	"context"

	"github.com/dehwyy/s3fx/pkg/dto"
	"github.com/minio/minio-go/v7"
)

func (storage *MinioStorage) Create(
	ctx context.Context,
	req dto.CreateRequest,
	opts ...minio.PutObjectOptions,
) error {
	bucketName, objectPath, err := parseURL(req.URL)
	if err != nil {
		return err
	}

	storage.mu.RLock()
	_, checked := storage.checkedBuckets[bucketName]
	storage.mu.RUnlock()

	if !checked {
		exists, err := storage.client.BucketExists(ctx, bucketName)
		if err != nil {
			return err
		}

		if !exists {
			err = storage.client.MakeBucket(
				ctx,
				bucketName,
				minio.MakeBucketOptions{},
			)
			if err != nil {
				return err
			}
		}

		storage.mu.Lock()
		storage.checkedBuckets[bucketName] = struct{}{}
		storage.mu.Unlock()
	}

	putOpts := minio.PutObjectOptions{}
	if len(opts) != 0 {
		putOpts = opts[0]
	}

	if _, err := storage.client.PutObject(
		ctx,
		bucketName,
		objectPath,
		req.Data,
		req.Size,
		putOpts,
	); err != nil {
		return err
	}

	return nil
}
