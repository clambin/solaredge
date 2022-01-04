package solaredge_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/clambin/solaredge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
	"net/http"
	"net/http/httptest"
	"net/url"
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
		assert.Equal(t, "failed to call server: 403 Forbidden", err.Error())
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
	start := time.Now()
	_, err := client.GetSiteIDs(ctx)
	require.Error(t, err)
	assert.Less(t, time.Since(start), 400*time.Millisecond)
}

func TestClient_Errors(t *testing.T) {
	server := &Server{token: "TESTTOKEN"}
	apiServer := httptest.NewServer(http.HandlerFunc(server.apiHandler))
	testURL, err := url.Parse(apiServer.URL)
	require.NoError(t, err)

	client := solaredge.Client{
		Token:      "BADTOKEN",
		HTTPClient: http.DefaultClient,
		APIURL:     apiServer.URL,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err = client.GetSiteIDs(ctx)
	require.Error(t, err)
	assert.Equal(t, "failed to call server: 403 Forbidden", err.Error())
	assert.NoError(t, errors.Unwrap(err))

	client.Token = "TESTTOKEN"
	server.garbage = true

	_, err = client.GetSiteIDs(ctx)
	require.Error(t, err)
	assert.Equal(t, `invalid server response: invalid character '=' after object key`, err.Error())
	assert.True(t, errors.Is(err, &solaredge.InvalidServerResponse{}))

	err = errors.Unwrap(err)
	require.Error(t, err)
	require.IsType(t, &json.SyntaxError{}, err)
	assert.Equal(t, int64(8), err.(*json.SyntaxError).Offset)

	apiServer.Close()
	_, err = client.GetSiteIDs(ctx)
	require.Error(t, err)
	assert.Equal(t, `failed to call server: Get "`+apiServer.URL+`/sites/list?api_key=TESTTOKEN": dial tcp 127.0.0.1:`+testURL.Port()+`: connect: connection refused`, err.Error())
	assert.True(t, errors.Is(err, unix.ECONNREFUSED))

	client.APIURL = "invalid url"
	_, err = client.GetSiteIDs(ctx)
	require.Error(t, err)
	assert.Equal(t, `failed to call server: Get "invalid%20url/sites/list?api_key=TESTTOKEN": unsupported protocol scheme ""`, err.Error())
	err = errors.Unwrap(err)
	require.Error(t, err)
	require.IsType(t, &url.Error{}, err)
	assert.Equal(t, "invalid%20url/sites/list?api_key=TESTTOKEN", err.(*url.Error).URL)
}

func TestClientEmpty(t *testing.T) {
	server := &Server{token: "TESTTOKEN", empty: true}
	apiServer := httptest.NewServer(http.HandlerFunc(server.apiHandler))

	client := solaredge.Client{
		Token:      "TESTTOKEN",
		HTTPClient: http.DefaultClient,
		APIURL:     apiServer.URL,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sites, err := client.GetSiteIDs(ctx)
	require.NoError(t, err)
	assert.Empty(t, sites)

	power, err := client.GetPower(ctx, 1, time.Now().Add(-time.Hour), time.Now())
	require.NoError(t, err)
	assert.Empty(t, power)

	lifeTime, lastYear, lastMonth, lastDay, current, err := client.GetPowerOverview(context.Background(), 1)
	require.NoError(t, err)
	assert.Zero(t, lifeTime)
	assert.Zero(t, lastYear)
	assert.Zero(t, lastMonth)
	assert.Zero(t, lastDay)
	assert.Zero(t, current)

}
