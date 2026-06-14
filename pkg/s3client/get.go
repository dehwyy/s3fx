package s3client

import (
	"bytes"
	"context"
	"io"

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

	var data []byte

	if err := retryNetwork(ctx, func(ctx context.Context) error {
		object, err := storage.client.GetObject(
			ctx,
			bucket,
			objectPath,
			minio.GetObjectOptions{},
		)
		if err != nil {
			return err
		}
		defer object.Close()

		// GetObject is lazy — the fetch happens here on read. Buffering inside the
		// retry makes a transient read failure recoverable and detaches the returned
		// body from the per-attempt context that would otherwise cancel it.
		buf, err := io.ReadAll(object)
		if err != nil {
			return err
		}

		data = buf
		return nil
	}); err != nil {
		return dto.GetResponse{}, err
	}

	return dto.GetResponse{
		Object: io.NopCloser(bytes.NewReader(data)),
	}, nil
}
