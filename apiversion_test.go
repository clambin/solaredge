package solaredge

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_GetCurrentAPIVersion(t *testing.T) {
	c, s, _ := makeTestServer(struct {
		Version struct {
			Release string `json:"release"`
		} `json:"version"`
	}{
		Version: struct {
			Release string `json:"release"`
		}{
			Release: "1.0.0",
		},
	})
	defer s.Close()

	version, err := c.GetCurrentAPIVersion(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "1.0.0", version)
}

func TestClient_GetSupportedAPIVersions(t *testing.T) {
	c, s, _ := makeTestServer(struct {
		Supported []struct {
			Release string `json:"release"`
		} `json:"supported"`
	}{
		Supported: []struct {
			Release string `json:"release"`
		}{
			{Release: "0.9.9"},
			{Release: "1.0.0"},
		},
	})
	defer s.Close()

	version, err := c.GetSupportedAPIVersions(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{"0.9.9", "1.0.0"}, version)
}
