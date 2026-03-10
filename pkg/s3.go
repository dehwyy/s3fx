package s3fx

import (
	"context"

	"github.com/dehwyy/s3fx/pkg/dto"
	"github.com/minio/minio-go/v7"
)

type ObjectStorage interface {
	Create(
		ctx context.Context,
		req dto.CreateRequest,
		opts ...minio.PutObjectOptions,
	) error
	Get(
		ctx context.Context,
		req dto.GetRequest,
	) (dto.GetResponse, error)
}
