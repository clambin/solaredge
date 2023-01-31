package solaredge

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestClient_BuildURL(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		path     string
		args     url.Values
		pass     bool
		expected string
	}{
		{
			name:     "default url",
			pass:     true,
			expected: apiURL,
		},
		{
			name:     "path",
			target:   "https://example.com",
			path:     "/foo",
			pass:     true,
			expected: "https://example.com/foo",
		},
		{
			name:     "args",
			target:   "https://example.com",
			path:     "/foo",
			args:     url.Values{"foo": []string{"1"}, "bar": []string{"2"}},
			pass:     true,
			expected: "https://example.com/foo?bar=2&foo=1",
		},
		{
			name:     "encoded args",
			target:   "https://example.com",
			path:     "/foo",
			args:     url.Values{"foo": []string{"1", "\"2\""}},
			pass:     true,
			expected: "https://example.com/foo?foo=1&foo=%222%22",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{APIURL: tt.target}
			output, err := c.buildURL(tt.path, tt.args)
			if !tt.pass {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, output)
		})
	}
}
