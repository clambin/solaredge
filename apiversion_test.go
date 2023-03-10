package solaredge_test

import (
	"context"
	"github.com/clambin/solaredge"
	"github.com/clambin/solaredge/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestClient_GetCurrentAPIVersion(t *testing.T) {
	c := solaredge.Client{Token: "1234"}
	c.HTTPClient = &http.Client{Transport: &testutil.Server{Token: "1234"}}

	version, err := c.GetCurrentAPIVersion(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "1.0.0", version)
}

func TestClient_GetSupportedAPIVersions(t *testing.T) {
	c := solaredge.Client{Token: "1234"}
	c.HTTPClient = &http.Client{Transport: &testutil.Server{Token: "1234"}}

	version, err := c.GetSupportedAPIVersions(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []string{"0.9.9", "1.0.0"}, version)
}
