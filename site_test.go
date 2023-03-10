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

func TestSite_GetID(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}
	ctx := context.Background()

	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)
	assert.Equal(t, 1, sites[0].GetID())
}

func TestSite_GetDataPeriod(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}
	ctx := context.Background()

	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	start, stop, err := sites[0].GetDataPeriod(ctx)
	require.NoError(t, err)
	assert.Equal(t, "2023-02-01", start.Format("2006-01-02"))
	assert.Equal(t, "2023-02-26", stop.Format("2006-01-02"))
}

func TestSite_GetEnergy(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}

	ctx := context.Background()
	start := time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)

	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	energy, err := sites[0].GetEnergy(ctx, "DAY", start, end)
	require.NoError(t, err)
	assert.NotZero(t, energy.Values)
}

func TestSite_GetTimeFrameEnergy(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	start := time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)

	energy, err := sites[0].GetTimeFrameEnergy(ctx, start, end)
	require.NoError(t, err)
	assert.NotZero(t, energy)
}

func TestSite_GetPower(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	start := time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)

	power, err := sites[0].GetPower(ctx, start, end)
	require.NoError(t, err)
	assert.NotZero(t, power)
}

func TestSite_GetPowerOverview(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	power, err := sites[0].GetPowerOverview(ctx)
	require.NoError(t, err)
	assert.NotZero(t, power)
}

func TestSite_GetPowerDetails(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	end := time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)
	start := end.Add(-7 * 24 * time.Hour)

	power, err := sites[0].GetPowerDetails(ctx, start, end)
	require.NoError(t, err)
	assert.NotZero(t, power)
}

func TestSite_GetEnergyDetails(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	end := time.Date(2023, time.March, 4, 0, 0, 0, 0, time.UTC)
	start := end.Add(-7 * 24 * time.Hour)

	details, err := sites[0].GetEnergyDetails(ctx, "DAY", start, end)
	require.NoError(t, err)
	assert.NotZero(t, details)
}

func TestSize_GetPowerFlow(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	flow, err := sites[0].GetPowerFlow(ctx)
	require.NoError(t, err)
	assert.NotZero(t, flow)
}

func TestSite_GetStorageData(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	end := time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)
	start := end.Add(-7 * 24 * time.Hour)

	data, err := sites[0].GetStorageData(ctx, start, end)
	require.NoError(t, err)
	assert.NotZero(t, data)
}

func TestSite_GetEnvBenefits(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}

	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	require.NoError(t, err)
	require.Len(t, sites, 1)

	response, err := sites[0].GetEnvBenefits(ctx)
	require.NoError(t, err)
	assert.NotZero(t, response)
}
