package solaredge_test

import (
	"context"
	"github.com/clambin/solaredge"
	"github.com/clambin/solaredge/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestInverter_GetTelemetry(t *testing.T) {
	c := solaredge.Client{Token: "1234"}
	c.HTTPClient = &http.Client{Transport: &testutil.Server{Token: "1234"}}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	inverters, err := sites[0].GetInverters(ctx)
	require.NoError(t, err)
	require.Len(t, inverters, 1)

	telemetries, err := inverters[0].GetTelemetry(ctx,
		time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.March, 2, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	assert.NotZero(t, telemetries)
}

func TestInverterEquipment_GetTelemetry(t *testing.T) {
	c := solaredge.Client{Token: "1234"}
	c.HTTPClient = &http.Client{Transport: &testutil.Server{Token: "1234"}}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	inventory, err := sites[0].GetInventory(ctx)
	require.NoError(t, err)
	require.Len(t, inventory.Inverters, 1)

	telemetries, err := inventory.Inverters[0].GetTelemetry(ctx,
		time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.March, 2, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	assert.NotZero(t, telemetries)
}
