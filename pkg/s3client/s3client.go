package s3client

import (
	"context"
	"fmt"
	"sync"
	"time"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"
)

type MinioConfig struct {
	Endpoint  string
	AccessKey string
	Secret    string
}

type MinioStorageOpts struct {
	fx.In
	*MinioConfig

	Lifecycle fx.Lifecycle
}

type MinioStorage struct {
	client *minio.Client

	mu             sync.RWMutex
	checkedBuckets map[string]struct{}
}

func NewMinioStorage(
	opts MinioStorageOpts,
) (*MinioStorage, error) {
	client, err := minio.New(
		opts.Endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				opts.AccessKey,
				opts.Secret,
				"",
			),
		},
	)
	if err != nil {
		return nil, err
	}

	opts.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			const attempts = 5
			var lastErr error
			for attempt := 1; attempt <= attempts; attempt++ {
				attemptCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
				_, lastErr = client.ListBuckets(attemptCtx)
				cancel()
				if lastErr == nil {
					return nil
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(500 * time.Millisecond):
				}
			}
			return fmt.Errorf("minio client could not access server after %d attempts: %w", attempts, lastErr)
		},
	})

	return &MinioStorage{
		client:         client,
		checkedBuckets: make(map[string]struct{}),
	}, nil
}
