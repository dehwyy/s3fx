package s3fx

import (
	"github.com/dehwyy/s3fx/pkg/client"
	"go.uber.org/fx"
)

type Opts = client.MinioStorageOpts

var (
	_ ObjectStorage = &client.MinioStorage{}
)

func Module(opts Opts) fx.Option {
	return fx.Module(
		"s3",
		fx.Provide(
			fx.Annotate(
				client.NewMinioStorage,
				fx.As(new(ObjectStorage)),
			),
		),
	)
}
