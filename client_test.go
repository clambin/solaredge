package solaredge_test

import (
	"context"
	"github.com/clambin/solaredge"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Authentication(t *testing.T) {
	server := &Server{token: "TESTTOKEN"}
	apiServer := httptest.NewServer(http.HandlerFunc(server.apiHandler))
	defer apiServer.Close()

	client := solaredge.NewClient("BADTOKEN", nil)
	client.APIURL = apiServer.URL

	_, err := client.GetSiteIDs(context.Background())

	if assert.Error(t, err) {
		assert.Equal(t, "403 Forbidden", err.Error())
	}

	client.Token = "TESTTOKEN"
	_, err = client.GetSiteIDs(context.Background())

	assert.NoError(t, err)
}

func TestClient_Timeout(t *testing.T) {
	server := &Server{token: "TESTTOKEN", slow: true}
	apiServer := httptest.NewServer(http.HandlerFunc(server.apiHandler))
	defer apiServer.Close()

	client := solaredge.NewClient("TESTTOKEN", &http.Client{Timeout: 100 * time.Millisecond})
	client.APIURL = apiServer.URL

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// this should finish after 100 ms (http.Client timeout)
	_, err := client.GetSiteIDs(ctx)

	assert.Error(t, err)
}
