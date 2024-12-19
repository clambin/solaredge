package v2_test

import (
	"context"
	v2 "github.com/clambin/solaredge/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestClient_GetCurrentAPIVersions(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetCurrentAPIVersion(ctx)
	require.NoError(t, err)
	assert.Equal(t, responses["/version/current"], resp)
}

func TestClient_GetSupportedAPIVersions(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSupportedAPIVersions(ctx)
	require.NoError(t, err)
	assert.Equal(t, responses["/version/supported"], resp)
}
