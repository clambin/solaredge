package solaredge

import (
	"context"
	"net/http"
	"testing"
)

func TestClient_GetCurrentAPIVersions(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetCurrentAPIVersion(context.Background())
	expect(t, resp, "/version/current", err)
}

func TestClient_GetSupportedAPIVersions(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetSupportedAPIVersions(context.Background())
	expect(t, resp, "/version/supported", err)
}
