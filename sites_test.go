package solaredge_test

import (
	"context"
	"github.com/clambin/solaredge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetSiteIDs(t *testing.T) {
	server := &Server{token: "TESTTOKEN"}
	apiServer := httptest.NewServer(http.HandlerFunc(server.apiHandler))
	defer apiServer.Close()

	client := solaredge.Client{
		Token:      "TESTTOKEN",
		HTTPClient: &http.Client{},
		APIURL:     apiServer.URL,
	}

	siteIDs, err := client.GetSiteIDs(context.Background())

	require.NoError(t, err)
	require.Len(t, siteIDs, 1)
	assert.Equal(t, 1, siteIDs[0])
}
