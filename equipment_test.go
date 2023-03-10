package solaredge_test

import (
	"context"
	"github.com/clambin/solaredge"
	"github.com/clambin/solaredge/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestClient_Equipment_E2E(t *testing.T) {
	apikey := os.Getenv("SOLAREDGE_APIKEY")
	if apikey == "" {
		t.Skip("SOLAREDGE_APIKEY not set. Skipping")
	}

	c := solaredge.Client{Token: apikey}
	ctx := context.Background()

	sites, err := c.GetSites(ctx)
	require.NoError(t, err)

	for _, site := range sites {

		inverters, err := site.GetInverters(ctx)
		require.NoError(t, err)

		end := time.Now()
		start := end.Add(-24 * time.Hour)

		for _, inverter := range inverters {
			telemetry, err := inverter.GetTelemetry(ctx, start, end)
			require.NoError(t, err)
			assert.NotEmpty(t, telemetry)

			_, err = site.GetChangeLog(ctx, inverter.SerialNumber)
			require.NoError(t, err)
		}

		inventory, err := site.GetInventory(ctx)
		require.NoError(t, err)

		for _, inverter := range inventory.Inverters {
			telemetry, err := inverter.GetTelemetry(ctx, start, end)
			require.NoError(t, err)
			assert.NotEmpty(t, telemetry)
			_, err = site.GetChangeLog(ctx, inverter.SN)
			require.NoError(t, err)
		}
	}
}

func TestSite_GetInverters(t *testing.T) {
	c := solaredge.Client{Token: "1234"}
	c.HTTPClient = &http.Client{Transport: &testutil.Server{Token: "1234"}}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	inverters, err := sites[0].GetInverters(ctx)
	require.NoError(t, err)
	assert.NotZero(t, inverters)
}

func TestSite_GetInventory(t *testing.T) {
	c := solaredge.Client{Token: "1234"}
	c.HTTPClient = &http.Client{Transport: &testutil.Server{Token: "1234"}}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	inventory, err := sites[0].GetInventory(ctx)
	require.NoError(t, err)
	assert.NotZero(t, inventory)
}

func TestSite_GetChangeLog(t *testing.T) {
	c := solaredge.Client{Token: "1234"}
	c.HTTPClient = &http.Client{Transport: &testutil.Server{Token: "1234"}}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	changelog, err := sites[0].GetChangeLog(ctx, "SN1234")
	require.NoError(t, err)
	assert.NotZero(t, changelog)
}
