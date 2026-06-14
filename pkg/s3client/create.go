package s3client

import (
	"bytes"
	"context"
	"io"

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

	// The upload may be retried on transient object-store failures, so the body
	// must be re-readable from the start on each attempt. A seekable reader is
	// rewound in place; a one-shot reader is buffered once (receipts are small).
	rewind, err := makeRewindable(req.Data)
	if err != nil {
		return err
	}

	putOpts := minio.PutObjectOptions{}
	if len(opts) != 0 {
		putOpts = opts[0]
	}

	return retryNetwork(ctx, func(ctx context.Context) error {
		if err := storage.ensureBucket(ctx, bucketName); err != nil {
			return err
		}

		body, err := rewind()
		if err != nil {
			return err
		}

		if _, err := storage.client.PutObject(
			ctx,
			bucketName,
			objectPath,
			body,
			req.Size,
			putOpts,
		); err != nil {
			return err
		}

		return nil
	})
}

// ensureBucket creates the bucket if it does not exist, caching the result so the
// existence check is performed at most once per bucket per process.
func (storage *MinioStorage) ensureBucket(
	ctx context.Context,
	bucketName string,
) error {
	storage.mu.RLock()
	_, checked := storage.checkedBuckets[bucketName]
	storage.mu.RUnlock()

	if checked {
		return nil
	}

	exists, err := storage.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		if err := storage.client.MakeBucket(
			ctx,
			bucketName,
			minio.MakeBucketOptions{},
		); err != nil {
			return err
		}
	}

	storage.mu.Lock()
	storage.checkedBuckets[bucketName] = struct{}{}
	storage.mu.Unlock()

	return nil
}

// makeRewindable returns a function that yields a reader positioned at the start
// of data on every call. Seekable readers are rewound in place; non-seekable
// readers are buffered into memory once.
func makeRewindable(data io.Reader) (func() (io.Reader, error), error) {
	if seeker, ok := data.(io.Seeker); ok {
		return func() (io.Reader, error) {
			if _, err := seeker.Seek(0, io.SeekStart); err != nil {
				return nil, err
			}
			return data, nil
		}, nil
	}

	buf, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}

	return func() (io.Reader, error) {
		return bytes.NewReader(buf), nil
	}, nil
}
