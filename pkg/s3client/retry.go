package s3client

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/minio/minio-go/v7"
)

const (
	retryMaxAttempts    = 5
	retryBaseBackoff    = 200 * time.Millisecond
	retryAttemptTimeout = 15 * time.Second
)

// retryNetwork runs fn with bounded retries and exponential backoff, giving each
// attempt its own timeout derived from ctx. It returns early on success, on a
// non-retryable error, or when the parent ctx is cancelled/expired. This makes
// object-store operations resilient to a lossy network path to the S3 endpoint.
func retryNetwork(ctx context.Context, fn func(context.Context) error) error {
	var lastErr error
	for attempt := 0; attempt < retryMaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		attemptCtx, cancel := context.WithTimeout(ctx, retryAttemptTimeout)
		lastErr = fn(attemptCtx)
		cancel()

		if lastErr == nil {
			return nil
		}
		if !isRetryable(lastErr) {
			return lastErr
		}

		timer := time.NewTimer(retryBaseBackoff << attempt)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}

	return lastErr
}

// isRetryable reports whether err is a transient failure worth retrying. Network
// failures and per-attempt timeouts are retried; permanent client errors (auth,
// validation, not-found) and parent-context cancellation are not.
func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.Canceled) {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}

	// A minio server response means the request reached S3 and got an HTTP status:
	// retry only transient 5xx / 429, never permanent 4xx.
	if response := minio.ToErrorResponse(err); response.StatusCode != 0 {
		return response.StatusCode >= 500 || response.StatusCode == 429
	}

	// No HTTP response at all (connection reset/refused before a reply) — transient.
	return true
}
