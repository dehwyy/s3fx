package client

import "context"

type SecretsProvider interface {
	MustGet(ctx context.Context, key string) any
}
