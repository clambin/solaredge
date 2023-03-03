package solaredge

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestClient_Equipment_E2E(t *testing.T) {
	apikey := os.Getenv("SOLAREDGE_APIKEY")
	if apikey == "" {
		t.Skip("SOLAREDGE_APIKEY not set. Skipping")
	}

	c := Client{Token: apikey}
	ctx := context.Background()

	inverters, err := c.GetInverters(ctx)
	require.NoError(t, err)

	end := time.Now()
	start := end.Add(-24 * time.Hour)

	for _, inverter := range inverters {
		telemetry, err := c.GetInverterTelemetry(ctx, inverter.SerialNumber, start, end)
		require.NoError(t, err)
		assert.NotEmpty(t, telemetry)

		_, err = c.GetChangeLog(ctx, inverter.SerialNumber)
		require.NoError(t, err)
	}

	inventory, err := c.GetInventory(ctx)
	require.NoError(t, err)

	for _, inverter := range inventory.Inverters {
		telemetry, err := c.GetInverterTelemetry(ctx, inverter.SN, start, end)
		require.NoError(t, err)
		assert.NotEmpty(t, telemetry)

		_, err = c.GetChangeLog(ctx, inverter.SN)
		require.NoError(t, err)
	}
}

func TestClient_GetInverters(t *testing.T) {
	info := []Inverter{{
		Manufacturer: "foo",
		Model:        "bar",
		Name:         "snafu",
		SerialNumber: "SN1234",
	}}
	c, s, _ := makeTestServer(struct {
		Reporters struct {
			Count int        `json:"count"`
			List  []Inverter `json:"list"`
		} `json:"reporters"`
	}{
		Reporters: struct {
			Count int        `json:"count"`
			List  []Inverter `json:"list"`
		}(struct {
			Count int
			List  []Inverter
		}{Count: 1, List: info}),
	})
	defer s.Close()

	inverters, err := c.GetInverters(context.Background())
	require.NoError(t, err)
	assert.Equal(t, info, inverters)
}

func TestClient_GetInventory(t *testing.T) {
	info := Inventory{
		Inverters: []InverterEquipment{
			{
				SN:                  "SN1234",
				CommunicationMethod: "ETHERNET",
				ConnectedOptimizers: 20,
				CpuVersion:          "1.0.0",
				Manufacturer:        "foo",
				Model:               "bar",
				Name:                "snafu",
			},
		},
	}
	c, s, _ := makeTestServer(struct {
		Inventory Inventory `json:"inventory"`
	}{Inventory: info})
	defer s.Close()

	inventory, err := c.GetInventory(context.Background())
	require.NoError(t, err)
	assert.Equal(t, info, inventory)
}

func TestClient_GetInverterTelemetry(t *testing.T) {
	info := []InverterTelemetry{
		{
			L1Data: struct {
				AcCurrent     float64 `json:"acCurrent"`
				AcFrequency   float64 `json:"acFrequency"`
				AcVoltage     float64 `json:"acVoltage"`
				ActivePower   float64 `json:"activePower"`
				ApparentPower float64 `json:"apparentPower"`
				CosPhi        float64 `json:"cosPhi"`
				ReactivePower float64 `json:"reactivePower"`
			}{},
			Time:                  Time{},
			DcVoltage:             400,
			GroundFaultResistance: 0,
			InverterMode:          "",
			OperationMode:         0,
			PowerLimit:            0,
			Temperature:           0,
			TotalActivePower:      0,
			TotalEnergy:           0,
		},
	}
	c, s, _ := makeTestServer(struct {
		Data struct {
			Count       int                 `json:"count"`
			Telemetries []InverterTelemetry `json:"telemetries"`
		} `json:"data"`
	}{
		Data: struct {
			Count       int                 `json:"count"`
			Telemetries []InverterTelemetry `json:"telemetries"`
		}{
			Count:       1,
			Telemetries: info,
		},
	})
	defer s.Close()

	telemetries, err := c.GetInverterTelemetry(context.Background(),
		"SN1234",
		time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.March, 2, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	assert.Equal(t, info, telemetries)
}

func TestClient_GetChangeLog(t *testing.T) {
	info := []ChangeLogEntry{
		{
			Date:         Date(time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC)),
			PartNumber:   "foo",
			SerialNumber: "SN1234",
		},
	}
	c, s, _ := makeTestServer(struct {
		ChangeLog struct {
			Count int              `json:"count"`
			List  []ChangeLogEntry `json:"list"`
		} `json:"ChangeLog"`
	}{
		ChangeLog: struct {
			Count int              `json:"count"`
			List  []ChangeLogEntry `json:"list"`
		}{
			Count: 1,
			List:  info,
		},
	})
	defer s.Close()

	telemetries, err := c.GetChangeLog(context.Background(), "SN1234")
	require.NoError(t, err)
	assert.Equal(t, info, telemetries)
}
