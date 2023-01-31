package solaredge_test

import (
	"context"
	"github.com/clambin/solaredge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_GetPower(t *testing.T) {
	server := &Server{token: "TESTTOKEN"}
	apiServer := httptest.NewServer(http.HandlerFunc(server.apiHandler))
	defer apiServer.Close()

	client := solaredge.Client{
		Token:  "TESTTOKEN",
		APIURL: apiServer.URL,
	}

	siteIDs, err := client.GetSiteIDs(context.Background())
	require.NoError(t, err)
	require.Len(t, siteIDs, 1)
	assert.Equal(t, 1, siteIDs[0])

	entries, err := client.GetPower(context.Background(), siteIDs[0], time.Now().Add(-1*time.Hour), time.Now())

	require.NoError(t, err)
	require.Len(t, entries, 2)
	assert.Equal(t, 12.0, entries[0].Value)
	assert.Equal(t, 24.0, entries[1].Value)
}

func TestClient_GetPowerOverview(t *testing.T) {
	server := &Server{token: "TESTTOKEN"}
	apiServer := httptest.NewServer(http.HandlerFunc(server.apiHandler))
	defer apiServer.Close()

	client := solaredge.Client{
		Token:  "TESTTOKEN",
		APIURL: apiServer.URL,
	}

	siteIDs, err := client.GetSiteIDs(context.Background())
	require.NoError(t, err)
	require.Len(t, siteIDs, 1)
	assert.Equal(t, 1, siteIDs[0])

	lifeTime, currentYear, currentMonth, currentDay, current, err := client.GetPowerOverview(context.Background(), siteIDs[0])
	require.NoError(t, err)
	assert.NotZero(t, lifeTime)
	assert.NotZero(t, currentYear)
	assert.NotZero(t, currentMonth)
	assert.NotZero(t, currentDay)
	assert.NotZero(t, current)
}
