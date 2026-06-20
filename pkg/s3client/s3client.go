package s3client

import (
	"context"
	"log"
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
			attemptCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			if _, err := client.ListBuckets(attemptCtx); err != nil {
				log.Printf("s3fx: initial connectivity check to %s failed, starting anyway (runtime ops retry): %v", opts.Endpoint, err)
			}
			return nil
		},
	})

	return &MinioStorage{
		client:         client,
		checkedBuckets: make(map[string]struct{}),
	}, nil
}
