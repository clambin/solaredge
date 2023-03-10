// Package testutil is a supporting package for the solaredge package. It creates a test server for the solaredge library,
// to be used during testing of solaredge, or an application using the solaredge library.
//
// The package provides a stubbed SolarEdge server, implemented as an http.RoundTripper. As such, it can be used as the
// transport of a http.Client:
//
//	s := testutil.Server{Token: "1234"}
//	c := solaredge.Client{
//		Token:      "1234",
//		HTTPClient: &http.Client{Transport: &s},
//	}
//
// This creates a test server, using a Token of "1234" and default responses (which provides one site, with a set of responses
// for that site). The solaredge client uses the server as its transport layer, so all requests will be handled by the test Server.
//
// To create a server with multiple sites, use MakeResponses:
//
//	s := testutil.Server{Token: "1234", Responses: MakeResponses(1, 2, 3)}
//	c := solaredge.Client{
//		Token:      "1234",
//		HTTPClient: &http.Client{Transport: &s},
//	}
//
// The test Server will now provide 3 sites.
package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clambin/solaredge"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var _ http.RoundTripper = &Server{}

// Server implements a stubbed solarEdge server. The server implements the http.RoundTripper interface and
// can therefore be plugged in as the transport for the solaredge.Client HTTPClient:
//
//	  	c := solaredge.Client{
//			Token:      "1234",
//			HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
//		}
type Server struct {
	// Token is the API key clients should use. An incorrect token will result in a 403 - Forbidden
	Token string
	// Responses is a map of responses per URL path.  If left empty, DefaultResponses will be used.
	Responses map[string]any

	lock sync.Mutex
}

// RoundTrip implements the http.RoundTripper interface.
func (s *Server) RoundTrip(r *http.Request) (*http.Response, error) {
	//fmt.Println(r.URL.Path)
	if !s.validToken(r) {
		return &http.Response{StatusCode: http.StatusForbidden}, nil
	}

	responses := s.getResponses()
	response, ok := responses[r.URL.Path]
	if !ok {
		return &http.Response{StatusCode: http.StatusNotFound}, nil
	}

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(response)

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(&buf),
	}, nil
}

func (s *Server) validToken(r *http.Request) bool {
	values := r.URL.Query()
	value, ok := values["api_key"]
	return ok && len(value) > 0 && value[0] == s.Token
}

func (s *Server) getResponses() map[string]any {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.Responses == nil {
		s.Responses = MakeResponses(1)
	}
	return s.Responses
}

// MakeResponses creates a set of responses for the provided list of siteIDs.
func MakeResponses(siteIDs ...int) map[string]any {
	responses := make(map[string]any)
	responses["/sites/list"] = makeSitesResponse(siteIDs...)
	responses["/version/current"] = struct {
		Version struct {
			Release string `json:"release"`
		} `json:"version"`
	}{
		Version: struct {
			Release string `json:"release"`
		}{Release: "1.0.0"},
	}
	responses["/version/supported"] = struct {
		Supported []struct {
			Release string `json:"release"`
		} `json:"supported"`
	}{
		Supported: []struct {
			Release string `json:"release"`
		}{
			{Release: "0.9.9"},
			{Release: "1.0.0"},
		},
	}

	for _, siteID := range siteIDs {
		responses["/site/"+strconv.Itoa(siteID)+"/overview"] = struct {
			Overview solaredge.PowerOverview `json:"overview"`
		}{
			Overview: solaredge.PowerOverview{
				LastUpdateTime: solaredge.Time(time.Now()),
				LifeTimeData: struct {
					Energy  float64 `json:"energy"`
					Revenue float64 `json:"revenue"`
				}{
					Energy: 1000,
				},
				LastYearData: struct {
					Energy  float64 `json:"energy"`
					Revenue float64 `json:"revenue"`
				}{
					Energy: 1000,
				},
				LastMonthData: struct {
					Energy  float64 `json:"energy"`
					Revenue float64 `json:"revenue"`
				}{
					Energy: 100,
				},
				LastDayData: struct {
					Energy  float64 `json:"energy"`
					Revenue float64 `json:"revenue"`
				}{
					Energy: 10,
				},
				CurrentPower: struct {
					Power float64 `json:"power"`
				}{
					Power: 3400,
				},
			},
		}
		responses["/site/"+strconv.Itoa(siteID)+"/dataPeriod"] = struct {
			DataPeriod struct {
				StartDate solaredge.Date `json:"startDate"`
				EndDate   solaredge.Date `json:"endDate"`
			} `json:"dataPeriod"`
		}{
			DataPeriod: struct {
				StartDate solaredge.Date `json:"startDate"`
				EndDate   solaredge.Date `json:"endDate"`
			}{
				StartDate: solaredge.Date(time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC)),
				EndDate:   solaredge.Date(time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)),
			},
		}
		responses["/site/"+strconv.Itoa(siteID)+"/energy"] = struct {
			Energy solaredge.Energy `json:"energy"`
		}{
			Energy: solaredge.Energy{
				MeasuredBy: "foo",
				TimeUnit:   "DAY",
				Unit:       "Wh",
				Values: []solaredge.Value{
					{
						Date:  solaredge.Time(time.Date(2023, time.February, 26, 0, 0, 0, 0, time.UTC)),
						Value: 10,
					},
				},
			},
		}
		responses["/site/"+strconv.Itoa(siteID)+"/timeFrameEnergy"] = struct {
			TimeFrameEnergy solaredge.TimeFrameEnergy `json:"timeFrameEnergy"`
		}{
			TimeFrameEnergy: solaredge.TimeFrameEnergy{
				Energy: 4096,
				Unit:   "Wh",
			},
		}
		responses["/site/"+strconv.Itoa(siteID)+"/power"] = struct {
			Power solaredge.Power `json:"power"`
		}{
			Power: solaredge.Power{
				TimeUnit: "DAY",
				Unit:     "W",
			},
		}
		responses["/site/"+strconv.Itoa(siteID)+"/powerDetails"] = struct {
			PowerDetails solaredge.PowerDetails `json:"powerDetails"`
		}{
			PowerDetails: solaredge.PowerDetails{
				TimeUnit: "QUARTER_OF_AN_HOUR",
				Unit:     "W",
				Meters: []solaredge.MeterReadings{
					{
						Type: "Consumption",
						Values: []solaredge.Value{
							{
								Date:  solaredge.Time(time.Date(2023, time.February, 26, 22, 12, 0, 0, time.UTC)),
								Value: 10,
							},
						},
					},
				},
			},
		}
		responses["/site/"+strconv.Itoa(siteID)+"/energyDetails"] = struct {
			EnergyDetails solaredge.EnergyDetails `json:"energyDetails"`
		}{
			EnergyDetails: solaredge.EnergyDetails{
				Meters: []solaredge.MeterReadings{
					{
						Type: "Production",
						Values: []solaredge.Value{
							{
								Date:  solaredge.Time(time.Date(2023, time.March, 4, 0, 0, 0, 0, time.UTC)),
								Value: 1000,
							},
						},
					},
				},
				TimeUnit: "DAY",
				Unit:     "Wh",
			},
		}
		responses["/site/"+strconv.Itoa(siteID)+"/currentPowerFlow"] = struct {
			SiteCurrentPowerFlow solaredge.CurrentPowerFlow `json:"siteCurrentPowerFlow"`
		}{
			SiteCurrentPowerFlow: solaredge.CurrentPowerFlow{
				Unit: "Wh",
				Connections: []struct {
					From string `json:"from"`
					To   string `json:"to"`
				}{{From: "GRID", To: "Load"}},
			},
		}
		responses["/site/"+strconv.Itoa(siteID)+"/storageData"] = struct {
			StorageData struct {
				BatteryCount int                 `json:"batteryCount"`
				Batteries    []solaredge.Battery `json:"batteries"`
			} `json:"storageData"`
		}{
			StorageData: struct {
				BatteryCount int                 `json:"batteryCount"`
				Batteries    []solaredge.Battery `json:"batteries"`
			}{
				BatteryCount: 1,
				Batteries: []solaredge.Battery{
					{
						Nameplate:    1,
						SerialNumber: "foo",
						ModelNumber:  "bar",
					},
				},
			},
		}
		responses["/site/"+strconv.Itoa(siteID)+"/envBenefits"] = struct {
			EnvBenefits solaredge.EnvBenefits `json:"envBenefits"`
		}{
			EnvBenefits: solaredge.EnvBenefits{
				TreesPlanted: 1024,
			},
		}
		responses["/equipment/"+strconv.Itoa(siteID)+"/list"] = struct {
			Reporters struct {
				Count int                  `json:"count"`
				List  []solaredge.Inverter `json:"list"`
			} `json:"reporters"`
		}{
			Reporters: struct {
				Count int                  `json:"count"`
				List  []solaredge.Inverter `json:"list"`
			}{
				Count: 1,
				List: []solaredge.Inverter{
					{
						Manufacturer: "foo",
						Model:        "bar",
						Name:         "snafu",
						SerialNumber: "SN1234",
					},
				},
			},
		}
		responses["/site/"+strconv.Itoa(siteID)+"/inventory"] = struct {
			Inventory solaredge.Inventory `json:"inventory"`
		}{
			Inventory: solaredge.Inventory{
				Inverters: []solaredge.InverterEquipment{
					{
						SN:                  "SN1234",
						CommunicationMethod: "ETHERNET",
						ConnectedOptimizers: 20,
						CPUVersion:          "1.0.0",
						Manufacturer:        "foo",
						Model:               "bar",
						Name:                "snafu",
					},
				},
			},
		}
		responses["/equipment/"+strconv.Itoa(siteID)+"/SN1234/changeLog"] = struct {
			ChangeLog struct {
				Count int                        `json:"count"`
				List  []solaredge.ChangeLogEntry `json:"list"`
			} `json:"ChangeLog"`
		}{
			ChangeLog: struct {
				Count int                        `json:"count"`
				List  []solaredge.ChangeLogEntry `json:"list"`
			}{
				Count: 1,
				List: []solaredge.ChangeLogEntry{
					{
						Date:         solaredge.Date(time.Date(2023, time.March, 1, 0, 0, 0, 0, time.UTC)),
						PartNumber:   "foo",
						SerialNumber: "SN1234",
					},
				},
			},
		}
		responses["/equipment/"+strconv.Itoa(siteID)+"/SN1234/data"] = struct {
			Data struct {
				Count       int                           `json:"count"`
				Telemetries []solaredge.InverterTelemetry `json:"telemetries"`
			} `json:"data"`
		}{
			Data: struct {
				Count       int                           `json:"count"`
				Telemetries []solaredge.InverterTelemetry `json:"telemetries"`
			}{
				Count: 1,
				Telemetries: []solaredge.InverterTelemetry{
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
						Time:                  solaredge.Time{},
						DcVoltage:             400,
						GroundFaultResistance: 0,
						InverterMode:          "",
						OperationMode:         0,
						PowerLimit:            0,
						Temperature:           0,
						TotalActivePower:      0,
						TotalEnergy:           0,
					},
				},
			},
		}
	}

	return responses
}

func makeSitesResponse(ids ...int) any {
	var sites solaredge.Sites

	for i, id := range ids {
		sites = append(sites, solaredge.Site{ID: id, Name: fmt.Sprintf("site %d", 1+i)})
	}

	return struct {
		Sites struct {
			Count int              `json:"count"`
			Site  []solaredge.Site `json:"site"`
		} `json:"sites"`
	}{
		Sites: struct {
			Count int              `json:"count"`
			Site  []solaredge.Site `json:"site"`
		}{
			Count: len(sites),
			Site:  sites,
		},
	}
}
