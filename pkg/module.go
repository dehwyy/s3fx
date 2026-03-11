package s3fx

import (
	"github.com/dehwyy/s3fx/pkg/s3client"
	"go.uber.org/fx"
)

type Opts = s3client.MinioStorageOpts

var (
	_ ObjectStorage = &s3client.MinioStorage{}
)

func Module(config s3client.MinioConfig) fx.Option {
	return fx.Module(
		"s3",
		fx.Provide(
			func() *s3client.MinioConfig {
				return &config
			},
			fx.Annotate(
				s3client.NewMinioStorage,
				fx.As(new(ObjectStorage)),
			),
		),
	)
}
