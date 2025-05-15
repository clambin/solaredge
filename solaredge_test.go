package solaredge

import (
	"codeberg.org/clambin/go-common/testutils"
	"net/http"
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
	if !reflect.DeepEqual(testResponses[path], response) {
		t.Errorf("expected response %v, got %v", testResponses[path], response)
	}
}

func NewTestServer() *httptest.Server {
	responses := make(testutils.Responses)
	for path, response := range testResponses {
		responses[path] = testutils.PathResponse{
			http.MethodGet: {Body: response},
		}
	}
	s := testutils.TestServer{Responses: responses}
	return httptest.NewServer(&s)
}

var testResponses = map[string]any{
	"/sites/list": GetSitesResponse{
		Sites: struct {
			Site  Sites `json:"site"`
			Count int   `json:"count"`
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
			Batteries    []Battery `json:"batteries"`
			BatteryCount int       `json:"batteryCount"`
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
			List  []Inverter `json:"list"`
			Count int        `json:"count"`
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
			Telemetries []InverterTelemetry `json:"telemetries"`
			Count       int                 `json:"count"`
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
			List  []EquipmentChangeLog `json:"list"`
			Count int                  `json:"count"`
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
