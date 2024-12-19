package v2_test

import (
	"github.com/clambin/go-common/testutils"
	solaredge "github.com/clambin/solaredge/v2"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	testServer = NewTestServer()
	code := m.Run()
	testServer.Close()
	os.Exit(code)
}

func NewTestServer() *httptest.Server {
	paths := make(map[string]testutils.Path)
	for path, response := range responses {
		paths[path] = testutils.Path{Body: response}
	}
	s := testutils.TestServer{Paths: paths}
	return httptest.NewServer(&s)
}

var responses = map[string]any{
	"/sites/list": solaredge.GetSitesResponse{
		Sites: struct {
			Count int             `json:"count"`
			Site  solaredge.Sites `json:"site"`
		}{
			Count: 1,
			Site: solaredge.Sites{
				solaredge.SiteDetails{Id: 1, Name: "site1"},
			},
		},
	},
	"/site/1/details": solaredge.GetSiteDetailsResponse{
		Details: solaredge.SiteDetails{Id: 1, Name: "site1"},
	},
	"/site/1/dataPeriod": solaredge.GetSiteDataPeriodResponse{
		DataPeriod: solaredge.SiteDataPeriod{
			StartDate: solaredge.Date(time.Date(2020, time.December, 19, 0, 0, 0, 0, time.UTC)),
			EndDate:   solaredge.Date(time.Date(2024, time.December, 19, 0, 0, 0, 0, time.UTC)),
		},
	},
	"/site/1/energy": solaredge.GetSiteEnergyResponse{
		Energy: solaredge.SiteEnergy{
			Values: []solaredge.Value{{Value: 1000}},
		},
	},
	"/site/1/timeFrameEnergy": solaredge.GetSiteEnergyForTimeframeResponse{
		TimeFrameEnergy: solaredge.SiteEnergyForTimeframe{
			Energy:              10,
			StartLifetimeEnergy: solaredge.SiteLifetimeEnergy{Energy: 1},
			EndLifetimeEnergy:   solaredge.SiteLifetimeEnergy{Energy: 1000},
		},
	},
	"/site/1/power": solaredge.GetSitePowerResponse{
		Power: solaredge.SitePower{
			Values: []solaredge.Value{{Value: 1000}},
		},
	},
	"/site/1/overview": solaredge.GetSitePowerOverviewResponse{
		Overview: solaredge.PowerOverview{
			LifeTimeData:  solaredge.EnergyOverview{Energy: 1000},
			LastYearData:  solaredge.EnergyOverview{Energy: 100},
			LastMonthData: solaredge.EnergyOverview{Energy: 10},
			LastDayData:   solaredge.EnergyOverview{Energy: 1},
			CurrentPower:  solaredge.CurrentPower{Power: 200},
		},
	},
	"/site/1/powerDetails": solaredge.GetSitePowerDetailsResponse{
		PowerDetails: solaredge.PowerDetails{
			Meters: []solaredge.MeterReadings{{Values: []solaredge.Value{{Value: 1000}}}},
		},
	},
	"/site/1/energyDetails": solaredge.GetSiteEnergyDetailsResponse{
		EnergyDetails: solaredge.SiteEnergyDetails{
			Meters: []solaredge.MeterReadings{{Values: []solaredge.Value{{Value: 1000}}}},
		},
	},
	"/site/1/currentPowerFlow": solaredge.GetSitePowerFlowResponse{
		CurrentPowerFlow: solaredge.CurrentPowerFlow{
			Grid: solaredge.PowerFlowReading{CurrentPower: 100},
			Load: solaredge.PowerFlowReading{CurrentPower: 10},
			PV:   solaredge.PowerFlowReading{CurrentPower: 50},
		},
	},
	"/site/1/storageData": solaredge.GetSiteStorageDataResponse{
		StorageData: struct {
			BatteryCount int                 `json:"batteryCount"`
			Batteries    []solaredge.Battery `json:"batteries"`
		}{
			BatteryCount: 1,
			Batteries: []solaredge.Battery{{
				SerialNumber: "SN1",
			}},
		},
	},
	"/site/1/envBenefits": solaredge.SiteEnvBenefitsResponse{
		EnvBenefits: solaredge.EnvBenefits{
			LightBulbs:   1000,
			TreesPlanted: 20,
		},
	},
	"/equipment/1/list": solaredge.GetComponentsResponse{
		Reporters: struct {
			Count int                  `json:"count"`
			List  []solaredge.Inverter `json:"list"`
		}{
			Count: 1,
			List:  []solaredge.Inverter{{SerialNumber: "SN1"}},
		},
	},
	"/site/1/inventory": solaredge.GetInventoryResponse{
		Inventory: solaredge.Inventory{
			Inverters: []solaredge.InverterEquipment{{SN: "SN1"}},
		},
	},
	"/equipment/1/SN1/data": solaredge.GetInverterTechnicalDataResponse{
		Data: struct {
			Count       int                           `json:"count"`
			Telemetries []solaredge.InverterTelemetry `json:"telemetries"`
		}{
			Count: 1,
			Telemetries: []solaredge.InverterTelemetry{
				{
					DcVoltage:        380,
					TotalActivePower: 1000,
					TotalEnergy:      10,
				},
			},
		},
	},
	"/equipment/1/SN1/changeLog": solaredge.GetEquipmentChangeLogResponse{
		ChangeLog: struct {
			Count int                            `json:"count"`
			List  []solaredge.EquipmentChangeLog `json:"list"`
		}{
			Count: 1,
			List:  []solaredge.EquipmentChangeLog{{SerialNumber: "SN1"}},
		},
	},
	"/version/current": solaredge.GetCurrentAPIVersionResponse{
		Version: solaredge.APIRelease{Release: "1.0.0"},
	},
	"/version/supported": solaredge.GetSupportedAPIVersionsResponse{
		Supported: []solaredge.APIRelease{{Release: "1.0.0"}},
	},
}
