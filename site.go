package solaredge

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

// This file contains APIs from the Site Data API section of the SolarEdge API specifications.
// https://knowledge-center.solaredge.com/sites/kc/files/se_monitoring_api.pdf

// Site contains details for an installation registered under the supplied API Token
type Site struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	AccountID        int     `json:"accountId"`
	Status           string  `json:"status"`
	PeakPower        float64 `json:"peakPower"`
	LastUpdateTime   string  `json:"lastUpdateTime"`
	InstallationDate string  `json:"installationDate"`
	// TODO
	PtoDate  interface{} `json:"ptoDate"`
	Notes    string      `json:"notes"`
	Type     string      `json:"type"`
	Location struct {
		Address     string `json:"address"`
		Address2    string `json:"address2"`
		City        string `json:"city"`
		Zip         string `json:"zip"`
		Country     string `json:"country"`
		CountryCode string `json:"countryCode"`
		TimeZone    string `json:"timeZone"`
	} `json:"location"`
	PrimaryModule struct {
		ManufacturerName string  `json:"manufacturerName"`
		ModelName        string  `json:"modelName"`
		MaximumPower     float64 `json:"maximumPower"`
	} `json:"primaryModule"`
	PublicSettings struct {
		IsPublic bool `json:"isPublic"`
	} `json:"publicSettings"`
	Uris struct {
		SITEIMAGE  string `json:"SITE_IMAGE"`
		DATAPERIOD string `json:"DATA_PERIOD"`
		DETAILS    string `json:"DETAILS"`
		OVERVIEW   string `json:"OVERVIEW"`
	} `json:"uris"`
	client *Client
}

// GetID returns the site's ID
func (s *Site) GetID() int {
	return s.ID
}

// GetDataPeriod returns the energy production start and end date for the active site.
//
// Note: unlike the example in the specs, this only returns the date, not the time of day.
func (s *Site) GetDataPeriod(ctx context.Context) (time.Time, time.Time, error) {
	var output struct {
		DataPeriod struct {
			StartDate Date `json:"startDate"`
			EndDate   Date `json:"endDate"`
		} `json:"dataPeriod"`
	}

	err := s.client.call(ctx, "/site/"+strconv.Itoa(s.ID)+"/dataPeriod", url.Values{}, &output)
	return time.Time(output.DataPeriod.StartDate), time.Time(output.DataPeriod.EndDate), err
}

// Energy produced by the site in a time range
type Energy struct {
	MeasuredBy string  `json:"measuredBy"`
	TimeUnit   string  `json:"timeUnit"`
	Unit       string  `json:"unit"`
	Values     []Value `json:"values"`
}

// GetEnergy returns the energy produced by the site for a given time range, with a given time resolution (timeUnit)
//
// timeUnit must be one of the following values: QUARTER_OF_AN_HOUR, HOUR, DAY, WEEK, MONTH, YEAR.
//
// The server limits the time range depending on the chosen timeUnit:
//   - For QUARTER_OF_AN_HOUR and HOUR, the time range cannot exceed one month
//   - For DAY, the time range cannot exceed one year.
//
// If these conditions are not met, an APIError is returned.
func (s *Site) GetEnergy(ctx context.Context, timeUnit string, start, end time.Time) (Energy, error) {
	var output struct {
		Energy Energy `json:"energy"`
	}
	args, err := buildArgsFromTimeRange(start, end, "Date", "2006-01-02")
	if err == nil {
		args.Set("timeUnit", timeUnit)
	}
	err = s.client.call(ctx, "/site/"+strconv.Itoa(s.ID)+"/energy", args, &output)
	return output.Energy, err
}

// TimeFrameEnergy is the total energy produced for a given period.
type TimeFrameEnergy struct {
	Energy float64 `json:"energy"`
	Unit   string  `json:"unit"`
}

// GetTimeFrameEnergy returns the total energy produced for a given period.
//
// Note: the period between end and start must not exceed one year. If it does, an APIError is returned.
func (s *Site) GetTimeFrameEnergy(ctx context.Context, start, end time.Time) (TimeFrameEnergy, error) {
	var output struct {
		TimeFrameEnergy TimeFrameEnergy `json:"timeFrameEnergy"`
	}
	args, err := buildArgsFromTimeRange(start, end, "Date", "2006-01-02")
	if err == nil {
		err = s.client.call(ctx, "/site/"+strconv.Itoa(s.ID)+"/timeFrameEnergy", args, &output)
	}
	return output.TimeFrameEnergy, err
}

// Power is the site power measurements for a site over a period of time.
type Power struct {
	TimeUnit string  `json:"timeUnit"`
	Unit     string  `json:"unit"`
	Values   []Value `json:"values"`
}

// GetPower returns the site power measurements in 15 minutes resolution.
//
// Note: This API is limited to a one-month period. If the provided time range exceeds one month, an APIError is returned.
func (s *Site) GetPower(ctx context.Context, start, end time.Time) (Power, error) {
	var powerStats struct {
		Power Power `json:"power"`
	}

	args, err := buildArgsFromTimeRange(start, end, "Time", "2006-01-02 15:04:05")
	if err == nil {
		err = s.client.call(ctx, "/site/"+strconv.Itoa(s.ID)+"/power", args, &powerStats)
	}
	return powerStats.Power, err
}

// PowerOverview contains the energy produced at a site for its entire lifetime, current year, month and day (in Wh) and current power (in W).
type PowerOverview struct {
	LastUpdateTime Time           `json:"lastUpdateTime"`
	LifeTimeData   EnergyOverview `json:"lifeTimeData"`
	LastYearData   EnergyOverview `json:"lastYearData"`
	LastMonthData  EnergyOverview `json:"lastMonthData"`
	LastDayData    EnergyOverview `json:"lastDayData"`
	CurrentPower   CurrentPower   `json:"currentPower"`
}

// EnergyOverview contains the energy produced in PowerOverview
type EnergyOverview struct {
	Energy  float64 `json:"energy"`
	Revenue float64 `json:"revenue"`
}

// CurrentPower contains the current power output, as contained in PowerOverview
type CurrentPower struct {
	Power float64 `json:"power"`
}

// GetPowerOverview returns the energy produced at the site for its entire lifetime, current year, month and day (in Wh) and current power (in W).
func (s *Site) GetPowerOverview(ctx context.Context) (PowerOverview, error) {
	var response struct {
		Overview PowerOverview `json:"overview"`
	}
	err := s.client.call(ctx, "/site/"+strconv.Itoa(s.ID)+"/overview", url.Values{}, &response)
	return response.Overview, err
}

// PowerDetails contains site power measurements from meters such as consumption, export (feed-in), import (purchase), etc.
type PowerDetails struct {
	Meters   []MeterReadings `json:"meters"`
	TimeUnit string          `json:"timeUnit"`
	Unit     string          `json:"unit"`
}

// MeterReadings contains power measurements for a type of meter.
type MeterReadings struct {
	Type   string  `json:"type"`
	Values []Value `json:"values"`
}

// GetPowerDetails returns site power measurements from meters such as consumption, export (feed-in), import (purchase), etc.
//
// Note: This API is limited to one-month period. If the provided time range exceeds one month, an APIError is returned.
func (s *Site) GetPowerDetails(ctx context.Context, start, end time.Time) (PowerDetails, error) {
	var response struct {
		PowerDetails PowerDetails `json:"powerDetails"`
	}
	args, err := buildArgsFromTimeRange(start, end, "Time", "2006-01-02 15:04:05")
	if err == nil {
		err = s.client.call(ctx, "/site/"+strconv.Itoa(s.ID)+"/powerDetails", args, &response)
	}
	return response.PowerDetails, err
}

// EnergyDetails contains site energy measurements from meters such as consumption, export (feed-in), import (purchase), etc.
type EnergyDetails struct {
	Meters   []MeterReadings `json:"meters"`
	TimeUnit string          `json:"timeUnit"`
	Unit     string          `json:"unit"`
}

// GetEnergyDetails returns site energy measurements from meters such as consumption, export (feed-in), import (purchase), etc.
// for a given time range and a specified time resolution.
//
// timeUnit must be one of the following values: QUARTER_OF_AN_HOUR, HOUR, DAY, WEEK, MONTH, YEAR.
//
// The server limits the time range depending on the chosen timeUnit:
//   - For QUARTER_OF_AN_HOUR and HOUR, the time range cannot exceed one month
//   - For DAY, the time range cannot exceed one year.
//
// If these conditions are not met, an APIError is returned.
func (s *Site) GetEnergyDetails(ctx context.Context, timeUnit string, start, end time.Time) (EnergyDetails, error) {
	var response struct {
		EnergyDetails EnergyDetails `json:"energyDetails"`
	}
	args, err := buildArgsFromTimeRange(start, end, "Time", "2006-01-02 15:04:05")
	if err == nil {
		args.Set("timeUnit", timeUnit)
		err = s.client.call(ctx, "/site/"+strconv.Itoa(s.ID)+"/energyDetails", args, &response)
	}
	return response.EnergyDetails, err
}

// CurrentPowerFlow contains current power flow between all elements of the site including PV array, storage (battery), loads (consumption) and grid.
type CurrentPowerFlow struct {
	Unit        string `json:"unit"`
	Connections []struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"connections"`
	Grid    PowerFlowReading `json:"GRID"`
	Load    PowerFlowReading `json:"LOAD"`
	PV      PowerFlowReading `json:"PV"`
	Storage struct {
		Status       string  `json:"status"`
		CurrentPower float64 `json:"currentPower"`
		ChargeLevel  float64 `json:"chargeLevel"`
		Critical     bool    `json:"critical"`
	} `json:"STORAGE"`
}

// PowerFlowReading contains the current power flow for one element of the site.
type PowerFlowReading struct {
	Status       string  `json:"status"`
	CurrentPower float64 `json:"currentPower"`
}

// GetPowerFlow returns the current power flow between all elements of the site including PV array, storage (battery), loads (consumption) and grid.
func (s *Site) GetPowerFlow(ctx context.Context) (CurrentPowerFlow, error) {
	var response struct {
		SiteCurrentPowerFlow CurrentPowerFlow `json:"siteCurrentPowerFlow"`
	}
	err := s.client.call(ctx, "/site/"+strconv.Itoa(s.ID)+"/currentPowerFlow", url.Values{}, &response)
	return response.SiteCurrentPowerFlow, err
}

// Battery contains detailed storage information from batteries: the state of energy, power and lifetime energy.
type Battery struct {
	Nameplate      int    `json:"nameplate"`
	SerialNumber   string `json:"serialNumber"`
	ModelNumber    string `json:"modelNumber"`
	TelemetryCount int    `json:"telemetryCount"`
	Telemetries    []struct {
		TimeStamp                Time    `json:"timeStamp"`
		Power                    float64 `json:"power"`
		BatteryState             int     `json:"batteryState"`
		LifeTimeEnergyCharged    float64 `json:"lifeTimeEnergyCharged"`
		LifeTimeEnergyDischarged float64 `json:"lifeTimeEnergyDischarged"`
		FullPackEnergyAvailable  float64 `json:"fullPackEnergyAvailable"`
		InternalTemp             float64 `json:"internalTemp"`
		ACGridCharging           float64 `json:"ACGridCharging"`
	} `json:"telemetries"`
}

// GetStorageData returns detailed information from batteries installed at the active site.
//
// This API is limited to a one-week period. If the time range exceeds one week, an APIError is returned.
func (s *Site) GetStorageData(ctx context.Context, start, end time.Time) ([]Battery, error) {
	var response struct {
		StorageData struct {
			BatteryCount int       `json:"batteryCount"`
			Batteries    []Battery `json:"batteries"`
		} `json:"storageData"`
	}
	args, err := buildArgsFromTimeRange(start, end, "Time", "2006-01-02 15:04:05")
	if err == nil {
		err = s.client.call(ctx, "/site/"+strconv.Itoa(s.ID)+"/storageData", args, &response)
	}
	return response.StorageData.Batteries, err
}

// TODO
// func (c *Client) GetSiteImage(ctx context.Context)

// TODO: Installer Logo Image

// EnvBenefits lists all environmental benefits based on site energy production: gas emissions saved, equivalent trees planted and light bulbs powered for a day.
type EnvBenefits struct {
	GasEmissionSaved struct {
		Co2   float64 `json:"co2"`
		Nox   float64 `json:"nox"`
		So2   float64 `json:"so2"`
		Units string  `json:"units"`
	} `json:"gasEmissionSaved"`
	LightBulbs   float64 `json:"lightBulbs"`
	TreesPlanted float64 `json:"treesPlanted"`
}

// GetEnvBenefits returns all environmental benefits based on site energy production: gas emissions saved, equivalent trees planted and light bulbs powered for a day.
func (s *Site) GetEnvBenefits(ctx context.Context) (EnvBenefits, error) {
	var response struct {
		EnvBenefits EnvBenefits `json:"envBenefits"`
	}
	err := s.client.call(ctx, "/site/"+strconv.Itoa(s.ID)+"/envBenefits", url.Values{}, &response)
	return response.EnvBenefits, err
}
