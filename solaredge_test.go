package solaredge_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/clambin/solaredge"
	"github.com/clambin/solaredge/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestClient_Authentication(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}

	goodToken := c.Token
	c.Token = "BADTOKEN"

	_, err := c.GetSites(context.Background())

	require.Error(t, err)
	require.ErrorIs(t, err, &solaredge.APIError{})

	c.Token = goodToken
	_, err = c.GetSites(context.Background())

	assert.NoError(t, err)
}

func TestClient_Server_Fail(t *testing.T) {
	c := solaredge.Client{
		Token: "1234",
		HTTPClient: &http.Client{
			Transport: &stubbedServer{
				roundTripper: func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Status:     "500 Internal Server Error",
						Body:       io.NopCloser(bytes.NewBufferString(`{"error": "software no workie"}`)),
					}, nil
				}}},
	}
	ctx := context.Background()

	_, err := c.GetSites(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, &solaredge.HTTPError{})
	assert.Equal(t, "500 Internal Server Error", err.Error())

	var err2 *solaredge.HTTPError
	require.ErrorAs(t, err, &err2)
	assert.Equal(t, `{"error": "software no workie"}`, string(err2.Body))
}

func TestClient_Client_Fail(t *testing.T) {
	c := solaredge.Client{
		Token: "1234",
		HTTPClient: &http.Client{
			Transport: &stubbedServer{
				roundTripper: func(r *http.Request) (*http.Response, error) {
					return nil, fmt.Errorf("fail")
				}}},
	}
	ctx := context.Background()

	_, err := c.GetSites(ctx)
	require.Error(t, err)
	assert.Equal(t, "Get \"https://monitoringapi.solaredge.com/sites/list?api_key=<REDACTED>&version=1.0.0\": fail", err.Error())
}

func TestClient_ParseError(t *testing.T) {
	c := solaredge.Client{
		Token: "1234",
		HTTPClient: &http.Client{
			Transport: &stubbedServer{
				roundTripper: func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString("definitely not a json object")),
					}, nil
				}}},
	}
	ctx := context.Background()

	_, err := c.GetSites(ctx)
	require.Error(t, err)
	assert.Equal(t, `json parse error: invalid character 'd' looking for beginning of value`, err.Error())
	assert.ErrorIs(t, err, &solaredge.ParseError{})
}

type stubbedServer struct {
	roundTripper func(r *http.Request) (*http.Response, error)
}

func (f stubbedServer) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.roundTripper(req)
}
