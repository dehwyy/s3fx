package s3client

import (
	"context"
	"fmt"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"
)

type MinioConfig struct {
	KeyEndpoint  string
	KeyAccessKey string
	KeySecret    string
}

type MinioStorageOpts struct {
	fx.In
	*MinioConfig

	SecretsProvider SecretsProvider
	Lifecycle       fx.Lifecycle
}

type MinioStorage struct {
	client *minio.Client

	mu             sync.RWMutex
	checkedBuckets map[string]struct{}
}

func NewMinioStorage(
	opts MinioStorageOpts,
) (*MinioStorage, error) {
	get := func(key string) string {
		return opts.SecretsProvider.MustGet(context.Background(), key).(string)
	}

	client, err := minio.New(
		get(opts.KeyEndpoint),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				get(opts.KeyAccessKey),
				get(opts.KeySecret),
				"", // ?
			),
		},
	)
	if err != nil {
		return nil, err
	}

	opts.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := client.ListBuckets(ctx)
			if err != nil {
				return fmt.Errorf("minio client could not access server: %w", err)
			}
			return nil
		},
	})

	return &MinioStorage{
		client:         client,
		checkedBuckets: make(map[string]struct{}),
	}, nil
}
