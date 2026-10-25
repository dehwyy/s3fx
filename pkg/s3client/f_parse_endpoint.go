package s3client

import "strings"

func parseEndpoint(
	endpoint string,
	secure bool,
) (
	host string,
	isSecure bool,
) {
	switch {
	case strings.HasPrefix(endpoint, "https://"):
		return strings.TrimSuffix(strings.TrimPrefix(endpoint, "https://"), "/"), true
	case strings.HasPrefix(endpoint, "http://"):
		return strings.TrimSuffix(strings.TrimPrefix(endpoint, "http://"), "/"), false
	default:
		return endpoint, secure
	}
}
