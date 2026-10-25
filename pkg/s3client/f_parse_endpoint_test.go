package s3client

import "testing"

func TestParseEndpoint(t *testing.T) {
	cases := []struct {
		endpoint string
		secure   bool
		host     string
		isSecure bool
	}{
		{"s3.twcstorage.ru", false, "s3.twcstorage.ru", false},
		{"s3.twcstorage.ru", true, "s3.twcstorage.ru", true},
		{"https://hel1.your-objectstorage.com", false, "hel1.your-objectstorage.com", true},
		{"https://hel1.your-objectstorage.com/", false, "hel1.your-objectstorage.com", true},
		{"http://minio:9000", true, "minio:9000", false},
	}

	for _, c := range cases {
		host, isSecure := parseEndpoint(
			c.endpoint,
			c.secure,
		)
		if host != c.host || isSecure != c.isSecure {
			t.Errorf("parseEndpoint(%q, %v) = (%q, %v), want (%q, %v)", c.endpoint, c.secure, host, isSecure, c.host, c.isSecure)
		}
	}
}
