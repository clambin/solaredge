package solaredge

import (
	"github.com/clambin/go-common/testutils"
	"net/http/httptest"
	"os"
	"reflect"
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

func expect(t *testing.T, response any, path string, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(responses[path], response) {
		t.Errorf("expected response %v, got %v", responses[path], response)
	}
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
	"/sites/list": GetSitesResponse{
		Sites: struct {
			Count int   `json:"count"`
			Site  Sites `json:"site"`
		}{
			Count: 1,
			Site: Sites{
				SiteDetails{Id: 1, Name: "site1"},
			},
		},
	},
	"/site/1/details": GetSiteDetailsResponse{
		Details: SiteDetails{Id: 1, Name: "site1"},
	},
	"/site/1/dataPeriod": GetDataPeriodResponse{
		DataPeriod: DataPeriod{
			StartDate: Date(time.Date(2020, time.December, 19, 0, 0, 0, 0, time.UTC)),
			EndDate:   Date(time.Date(2024, time.December, 19, 0, 0, 0, 0, time.UTC)),
		},
	},
	"/site/1/energy": GetEnergyMeasurementsResponse{
		Energy: EnergyMeasurements{
			Values: []Value{{Value: 1000}},
		},
	},
	"/site/1/timeFrameEnergy": GetEnergyForTimeframeResponse{
		TimeFrameEnergy: SiteEnergyForTimeframe{
			Energy:              10,
			StartLifetimeEnergy: LifetimeEnergy{Energy: 1},
			EndLifetimeEnergy:   LifetimeEnergy{Energy: 1000},
		},
	},
	"/site/1/power": GetPowerMeasurementsResponse{
		Power: PowerMeasurements{
			Values: []Value{{Value: 1000}},
		},
	},
	"/site/1/overview": GetPowerOverviewResponse{
		Overview: PowerOverview{
			LifeTimeData:  EnergyOverview{Energy: 1000},
			LastYearData:  EnergyOverview{Energy: 100},
			LastMonthData: EnergyOverview{Energy: 10},
			LastDayData:   EnergyOverview{Energy: 1},
			CurrentPower:  CurrentPower{Power: 200},
		},
	},
	"/site/1/powerDetails": GetPowerDetailsResponse{
		PowerDetails: PowerDetails{
			Meters: []MeterReadings{{Values: []Value{{Value: 1000}}}},
		},
	},
	"/site/1/energyDetails": GetEnergyDetailsResponse{
		EnergyDetails: EnergyDetails{
			Meters: []MeterReadings{{Values: []Value{{Value: 1000}}}},
		},
	},
	"/site/1/currentPowerFlow": GetPowerFlowResponse{
		CurrentPowerFlow: PowerFlow{
			Grid: PowerFlowReading{CurrentPower: 100},
			Load: PowerFlowReading{CurrentPower: 10},
			PV:   PowerFlowReading{CurrentPower: 50},
		},
	},
	"/site/1/storageData": GetStorageDataResponse{
		StorageData: struct {
			BatteryCount int       `json:"batteryCount"`
			Batteries    []Battery `json:"batteries"`
		}{
			BatteryCount: 1,
			Batteries: []Battery{{
				SerialNumber: "SN1",
			}},
		},
	},
	"/site/1/envBenefits": GetEnvBenefitsResponse{
		EnvBenefits: EnvBenefits{
			LightBulbs:   1000,
			TreesPlanted: 20,
		},
	},
	"/equipment/1/list": GetComponentsResponse{
		Reporters: struct {
			Count int        `json:"count"`
			List  []Inverter `json:"list"`
		}{
			Count: 1,
			List:  []Inverter{{SerialNumber: "SN1"}},
		},
	},
	"/site/1/inventory": GetInventoryResponse{
		Inventory: Inventory{
			Inverters: []InverterEquipment{{SN: "SN1"}},
		},
	},
	"/equipment/1/SN1/data": GetInverterTechnicalDataResponse{
		Data: struct {
			Count       int                 `json:"count"`
			Telemetries []InverterTelemetry `json:"telemetries"`
		}{
			Count: 1,
			Telemetries: []InverterTelemetry{
				{
					DcVoltage:        380,
					TotalActivePower: 1000,
					TotalEnergy:      10,
				},
			},
		},
	},
	"/equipment/1/SN1/changeLog": GetEquipmentChangeLogResponse{
		ChangeLog: struct {
			Count int                  `json:"count"`
			List  []EquipmentChangeLog `json:"list"`
		}{
			Count: 1,
			List:  []EquipmentChangeLog{{SerialNumber: "SN1"}},
		},
	},
	"/version/current": GetCurrentAPIVersionResponse{
		Version: APIRelease{Release: "1.0.0"},
	},
	"/version/supported": GetSupportedAPIVersionsResponse{
		Supported: []APIRelease{{Release: "1.0.0"}},
	},
}
