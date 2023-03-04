package solaredge

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestClient_Site_E2E(t *testing.T) {
	apikey := os.Getenv("SOLAREDGE_APIKEY")
	if apikey == "" {
		t.Skip("SOLAREDGE_APIKEY not set. Skipping")
	}

	c := Client{Token: apikey}
	ctx := context.Background()

	_, err := c.GetSites(ctx)
	require.NoError(t, err)

	_, err = c.GetSiteDetails(ctx)
	require.NoError(t, err)

	start, end, err := c.GetSiteDataPeriod(ctx)
	require.NoError(t, err)
	assert.False(t, start.IsZero())
	assert.False(t, end.IsZero())

	energy, err := c.GetSiteEnergy(ctx, "YEAR", start, end)
	require.NoError(t, err)
	assert.NotEmpty(t, energy.Values)

	siteEnergy, err := c.GetSiteTimeFrameEnergy(ctx, end.Add(-365*24*time.Hour), end)
	require.NoError(t, err)
	assert.NotZero(t, siteEnergy.Energy)

	sitePower, err := c.GetPower(ctx, end.Add(-7*24*time.Hour), end)
	require.NoError(t, err)
	assert.Equal(t, "QUARTER_OF_AN_HOUR", sitePower.TimeUnit)
	assert.Equal(t, "W", sitePower.Unit)
	assert.NotEmpty(t, sitePower.Values)

	powerOverview, err := c.GetPowerOverview(ctx)
	require.NoError(t, err)
	assert.NotZero(t, powerOverview.LifeTimeData)
	assert.False(t, time.Time(powerOverview.LastUpdateTime).IsZero())

	powerDetails, err := c.GetPowerDetails(ctx, end.Add(7*24*time.Hour), end)
	require.NoError(t, err)
	assert.NotEmpty(t, powerDetails.Meters)

}

func TestClient_GetSiteIDs(t *testing.T) {
	c, s, _ := makeTestServer(nil)
	defer s.Close()
	siteIDs, err := c.GetSiteIDs(context.Background())

	require.NoError(t, err)
	require.Len(t, siteIDs, 1)
	assert.Equal(t, 1, siteIDs[0])
}

func TestClient_GetSites(t *testing.T) {
	c, s, _ := makeTestServer(nil)
	defer s.Close()
	sites, err := c.GetSites(context.Background())

	require.NoError(t, err)
	require.Len(t, sites, 1)
	assert.Equal(t, 1, sites[0].ID)
	assert.Equal(t, "home", sites[0].Name)
}

func TestClient_GetSiteDetails(t *testing.T) {
	details := Site{ID: 1, Name: "home"}
	c, s, _ := makeTestServer(struct {
		Details Site `json:"details"`
	}{
		Details: details,
	})
	defer s.Close()
	site, err := c.GetSiteDetails(context.Background())

	require.NoError(t, err)
	assert.Equal(t, details, site)
}

func TestClient_GetSiteDataPeriod(t *testing.T) {
	c, s, _ := makeTestServer(struct {
		DataPeriod struct {
			StartDate Date `json:"startDate"`
			EndDate   Date `json:"endDate"`
		} `json:"dataPeriod"`
	}{
		DataPeriod: struct {
			StartDate Date `json:"startDate"`
			EndDate   Date `json:"endDate"`
		}{
			StartDate: Date(time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)),
			EndDate:   Date(time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)),
		},
	},
	)
	defer s.Close()

	start, stop, err := c.GetSiteDataPeriod(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "2023-02-01", start.Format("2006-01-02"))
	assert.Equal(t, "2023-02-26", stop.Format("2006-01-02"))
}

func TestClient_GetSiteEnergy(t *testing.T) {
	info := Energy{
		MeasuredBy: "foo",
		TimeUnit:   "DAY",
		Unit:       "Wh",
		Values:     nil,
	}
	c, s, _ := makeTestServer(struct {
		Energy Energy `json:"energy"`
	}{
		Energy: info,
	})
	defer s.Close()
	start := time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)

	energy, err := c.GetSiteEnergy(context.Background(), "DAY", start, end)
	require.NoError(t, err)
	assert.Equal(t, info, energy)
}

func TestClient_GetSiteTimeFrameEnergy(t *testing.T) {
	info := TimeFrameEnergy{
		Energy: 4096,
		Unit:   "Wh",
	}
	c, s, _ := makeTestServer(struct {
		TimeFrameEnergy TimeFrameEnergy `json:"timeFrameEnergy"`
	}{
		TimeFrameEnergy: info,
	})
	defer s.Close()
	start := time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)

	energy, err := c.GetSiteTimeFrameEnergy(context.Background(), start, end)
	require.NoError(t, err)
	assert.Equal(t, info, energy)
}

func TestClient_GetPower(t *testing.T) {
	info := Power{
		TimeUnit: "DAY",
		Unit:     "W",
	}
	c, s, _ := makeTestServer(struct {
		Power Power `json:"power"`
	}{
		Power: info,
	})
	defer s.Close()
	start := time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)

	energy, err := c.GetPower(context.Background(), start, end)
	require.NoError(t, err)
	assert.Equal(t, info, energy)
}

func TestClient_GetPowerOverview(t *testing.T) {
	info := PowerOverview{
		LastUpdateTime: Time(time.Date(2023, time.March, 4, 0, 0, 0, 0, time.UTC)),
		LifeTimeData: struct {
			Energy  float64 `json:"energy"`
			Revenue float64 `json:"revenue"`
		}{
			Energy:  1000,
			Revenue: 0,
		},
		LastYearData: struct {
			Energy  float64 `json:"energy"`
			Revenue float64 `json:"revenue"`
		}{
			Energy:  100,
			Revenue: 0,
		},
		LastMonthData: struct {
			Energy  float64 `json:"energy"`
			Revenue float64 `json:"revenue"`
		}{
			Energy:  10,
			Revenue: 0,
		},
		LastDayData: struct {
			Energy  float64 `json:"energy"`
			Revenue float64 `json:"revenue"`
		}{
			Energy:  1,
			Revenue: 0,
		},
		CurrentPower: struct {
			Power float64 `json:"power"`
		}{
			Power: 1,
		},
	}
	c, s, _ := makeTestServer(struct {
		Overview PowerOverview `json:"overview"`
	}{
		Overview: info,
	})
	defer s.Close()

	power, err := c.GetPowerOverview(context.Background())
	require.NoError(t, err)
	assert.Equal(t, info, power)
}

func TestClient_GetPowerDetails(t *testing.T) {
	info := PowerDetails{
		TimeUnit: "QUARTER_OF_AN_HOUR",
		Unit:     "W",
		Meters: []MeterReadings{
			{
				Type: "Consumption",
				Values: []Value{
					{
						Date:  Time(time.Date(2023, time.February, 26, 22, 12, 0, 0, time.UTC)),
						Value: 10,
					},
				},
			},
		},
	}
	c, s, _ := makeTestServer(struct {
		PowerDetails PowerDetails `json:"powerDetails"`
	}{
		PowerDetails: info,
	})
	defer s.Close()

	end := time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)
	start := end.Add(-7 * 24 * time.Hour)
	power, err := c.GetPowerDetails(context.Background(), start, end)
	require.NoError(t, err)
	assert.Equal(t, info, power)
}

func TestClient_GetEnergyDetails(t *testing.T) {
	info := EnergyDetails{
		Meters: []MeterReadings{
			{
				Type: "Production",
				Values: []Value{
					{
						Date:  Time(time.Date(2023, time.March, 4, 0, 0, 0, 0, time.UTC)),
						Value: 1000,
					},
				},
			},
		},
		TimeUnit: "DAY",
		Unit:     "Wh",
	}
	c, s, _ := makeTestServer(struct {
		EnergyDetails EnergyDetails `json:"energyDetails"`
	}{
		EnergyDetails: info,
	})
	defer s.Close()

	ctx := context.Background()
	end := time.Date(2023, time.March, 4, 0, 0, 0, 0, time.UTC)
	start := end.Add(-7 * 24 * time.Hour)
	details, err := c.GetEnergyDetails(ctx, "DAY", start, end)
	require.NoError(t, err)
	assert.Equal(t, info, details)
}

func TestClient_GetPowerFlow(t *testing.T) {
	info := CurrentPowerFlow{
		Unit: "Wh",
		Connections: []struct {
			From string `json:"from"`
			To   string `json:"to"`
		}{{From: "GRID", To: "Load"}},
	}
	c, s, _ := makeTestServer(struct {
		SiteCurrentPowerFlow CurrentPowerFlow `json:"siteCurrentPowerFlow"`
	}{
		SiteCurrentPowerFlow: info,
	})
	defer s.Close()

	flow, err := c.GetPowerFlow(context.Background())
	require.NoError(t, err)
	assert.Equal(t, info, flow)
}

func TestClient_GetStorageData(t *testing.T) {
	info := []Battery{
		{
			Nameplate:    1,
			SerialNumber: "foo",
			ModelNumber:  "bar",
		},
	}
	c, s, _ := makeTestServer(struct {
		StorageData struct {
			BatteryCount int       `json:"batteryCount"`
			Batteries    []Battery `json:"batteries"`
		} `json:"storageData"`
	}{
		StorageData: struct {
			BatteryCount int       `json:"batteryCount"`
			Batteries    []Battery `json:"batteries"`
		}{
			BatteryCount: 1,
			Batteries:    info,
		},
	})
	defer s.Close()

	end := time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)
	start := end.Add(-7 * 24 * time.Hour)

	data, err := c.GetStorageData(context.Background(), start, end)
	require.NoError(t, err)
	assert.Equal(t, info, data)
}

func TestClient_GetEnvBenefits(t *testing.T) {
	info := EnvBenefits{
		TreesPlanted: 1024,
	}
	c, s, _ := makeTestServer(struct {
		EnvBenefits EnvBenefits `json:"envBenefits"`
	}{EnvBenefits: info})
	defer s.Close()
	response, err := c.GetEnvBenefits(context.Background())
	require.NoError(t, err)
	assert.Equal(t, info, response)
}
