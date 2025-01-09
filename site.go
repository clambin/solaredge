package solaredge

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

// This file implements the "Site Data API" section of the SolarEdge API specifications.
// https://knowledge-center.solaredge.com/sites/kc/files/se_monitoring_api.pdf
//
// Omitted APIs:
// - Site Image (/site/{siteId}/siteImage/{name})
// - Installer Logo Image (/site/{siteId}/installerImage/{name})

// GetSites returns a list of sites related to the given SiteKey, which is the account api_key.
//
// Note: the API accepts parameters for convenient search, sort and pagination, but this is currently not implemented.
// This means the result is limited to the first 100 sites registered under the SiteKey.
func (c *Client) GetSites(ctx context.Context) (GetSitesResponse, error) {
	return call[GetSitesResponse](ctx, c, "/sites/list", nil)
}

type GetSitesResponse struct {
	Sites struct {
		Site  Sites `json:"site"`
		Count int   `json:"count"`
	} `json:"sites"`
}

type Sites []SiteDetails

// FindByID returns the SiteDetails with the specified ID. Returns false if the SiteDetails could not be found.
func (s Sites) FindByID(id int) (SiteDetails, bool) {
	for _, site := range s {
		if site.Id == id {
			return site, true
		}
	}
	return SiteDetails{}, false
}

// FindByName returns the SiteDetails with the specified name. Returns false if the SiteDetails could not be found.
func (s Sites) FindByName(name string) (SiteDetails, bool) {
	for _, site := range s {
		if site.Name == name {
			return site, true
		}
	}
	return SiteDetails{}, false
}

// GetSiteDetails returns the site details, such as name, location, status, etc.
func (c *Client) GetSiteDetails(ctx context.Context, id int) (GetSiteDetailsResponse, error) {
	return call[GetSiteDetailsResponse](ctx, c, "/site/"+strconv.Itoa(id)+"/details", nil)
}

type GetSiteDetailsResponse struct {
	Details SiteDetails `json:"details"`
}

type SiteDetails struct {
	LastUpdateTime   Date  `json:"lastUpdateTime"`
	InstallationDate Date  `json:"installationDate"`
	PtoDate          *Date `json:"ptoDate"`
	Location         struct {
		Country     string `json:"country"`
		City        string `json:"city"`
		Address     string `json:"address"`
		Address2    string `json:"address2"`
		Zip         string `json:"zip"`
		TimeZone    string `json:"timeZone"`
		CountryCode string `json:"countryCode"`
	} `json:"location"`
	Uris struct {
		SITEIMAGE  string `json:"SITE_IMAGE"`
		DATAPERIOD string `json:"DATA_PERIOD"`
		DETAILS    string `json:"DETAILS"`
		OVERVIEW   string `json:"OVERVIEW"`
	} `json:"uris"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	Notes         string `json:"notes"`
	Type          string `json:"type"`
	PrimaryModule struct {
		ManufacturerName string  `json:"manufacturerName"`
		ModelName        string  `json:"modelName"`
		MaximumPower     float64 `json:"maximumPower"`
	} `json:"primaryModule"`
	Id             int     `json:"id"`
	AccountId      int     `json:"accountId"`
	PeakPower      float64 `json:"peakPower"`
	PublicSettings struct {
		IsPublic bool `json:"isPublic"`
	} `json:"publicSettings"`
}

// GetDataPeriod returns the energy production start and end dates of the site.
//
// Note: unlike the example in the specs, this only returns the date, not the time of day.
func (c *Client) GetDataPeriod(ctx context.Context, id int) (GetDataPeriodResponse, error) {
	return call[GetDataPeriodResponse](ctx, c, makePath("/site/{siteId}/dataPeriod", id), nil)
}

type GetDataPeriodResponse struct {
	DataPeriod DataPeriod `json:"dataPeriod"`
}

type DataPeriod struct {
	StartDate Date `json:"startDate"`
	EndDate   Date `json:"endDate"`
}

// GetEnergyMeasurements return the site energy measurements.
//
// timeUnit must be one of the following values: QUARTER_OF_AN_HOUR, HOUR, DAY, WEEK, MONTH, YEAR.
//
// The server limits the time range depending on the chosen timeUnit:
//   - For QUARTER_OF_AN_HOUR and HOUR, the time range cannot exceed one month
//   - For DAY, the time range cannot exceed one year.
//
// If these conditions are not met, the call returns an error.
func (c *Client) GetEnergyMeasurements(ctx context.Context, id int, timeUnit TimeUnit, startDate time.Time, endDate time.Time) (GetEnergyMeasurementsResponse, error) {
	args := url.Values{
		"startDate": []string{startDate.Format(time.DateOnly)},
		"endDate":   []string{endDate.Format(time.DateOnly)},
		"timeUnit":  []string{string(timeUnit)},
	}
	return call[GetEnergyMeasurementsResponse](ctx, c, makePath("/site/{siteId}/energy", id), args)
}

type GetEnergyMeasurementsResponse struct {
	Energy EnergyMeasurements `json:"energy"`
}

type EnergyMeasurements struct {
	TimeUnit   string   `json:"timeUnit"`
	Unit       TimeUnit `json:"unit"`
	MeasuredBy string   `json:"measuredBy"`
	Values     []Value  `json:"values"`
}

// GetEnergyForTimeFrame returns the site total energy produced for a given period.
//
// Notes:
//   - This API only returns on-grid energy for the requested period. In sites with storage/backup, this may mean that results can differ from what appears in the Site Dashboard. Use the regular Site EnergyMeasurements API to obtain results that match the Site Dashboard calculation.
//   - The period between end and start must not exceed one year.
func (c *Client) GetEnergyForTimeFrame(ctx context.Context, id int, startDate, endDate time.Time) (GetEnergyForTimeframeResponse, error) {
	args := url.Values{
		"startDate": []string{startDate.Format(time.DateOnly)},
		"endDate":   []string{endDate.Format(time.DateOnly)},
	}
	return call[GetEnergyForTimeframeResponse](ctx, c, makePath("/site/{siteId}/timeFrameEnergy", id), args)
}

type GetEnergyForTimeframeResponse struct {
	TimeFrameEnergy SiteEnergyForTimeframe `json:"timeFrameEnergy"`
}

type SiteEnergyForTimeframe struct {
	Unit                string         `json:"unit"`
	MeasuredBy          string         `json:"measuredBy"`
	StartLifetimeEnergy LifetimeEnergy `json:"startLifetimeEnergy"`
	EndLifetimeEnergy   LifetimeEnergy `json:"endLifetimeEnergy"`
	Energy              float64        `json:"energy"`
}

type LifetimeEnergy struct {
	Date   string  `json:"date"`
	Unit   string  `json:"unit"`
	Energy float64 `json:"energy"`
}

const timeFormat = "2006-01-02 15:04:05"

// GetPowerMeasurements returns the site power measurements in 15 minute intervals.
//
// Note: This API is limited to a one-month period. If the provided time range exceeds one month, an APIError is returned.
func (c *Client) GetPowerMeasurements(ctx context.Context, id int, startTime, endTime time.Time) (GetPowerMeasurementsResponse, error) {
	args := url.Values{
		"startTime": []string{startTime.Format(timeFormat)},
		"endTime":   []string{endTime.Format(timeFormat)},
	}
	return call[GetPowerMeasurementsResponse](ctx, c, makePath("/site/{siteId}/power", id), args)
}

type GetPowerMeasurementsResponse struct {
	Power PowerMeasurements `json:"power"`
}

type PowerMeasurements struct {
	TimeUnit   string  `json:"timeUnit"`
	Unit       string  `json:"unit"`
	MeasuredBy string  `json:"measuredBy"`
	Values     []Value `json:"values"`
}

// GetPowerOverview returns the energy produced at the site for its entire lifetime, current year, month and day (in Wh) and current power (in W).
func (c *Client) GetPowerOverview(ctx context.Context, id int) (GetPowerOverviewResponse, error) {
	return call[GetPowerOverviewResponse](ctx, c, makePath("/site/{siteId}/overview", id), nil)
}

type GetPowerOverviewResponse struct {
	Overview PowerOverview `json:"overview"`
}

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

// GetPowerDetails returns site power measurements from meters such as consumption, export (feed-in), import (purchase), etc.
//
// Note: This API is limited to one-month period. If the provided time range exceeds one month, an APIError is returned.
func (c *Client) GetPowerDetails(ctx context.Context, id int, start, end time.Time) (GetPowerDetailsResponse, error) {
	args := url.Values{
		"startTime": []string{start.Format(timeFormat)},
		"endTime":   []string{end.Format(timeFormat)},
	}
	return call[GetPowerDetailsResponse](ctx, c, makePath("/site/{siteId}/powerDetails", id), args)
}

type GetPowerDetailsResponse struct {
	PowerDetails PowerDetails `json:"powerDetails"`
}

type PowerDetails struct {
	TimeUnit TimeUnit        `json:"timeUnit"`
	Unit     string          `json:"unit"`
	Meters   []MeterReadings `json:"meters"`
}

// MeterReadings contains power measurements for a type of meter.
type MeterReadings struct {
	Type   string  `json:"type"`
	Values []Value `json:"values"`
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
func (c *Client) GetEnergyDetails(ctx context.Context, id int, timeUnit TimeUnit, startTime, endTime time.Time) (GetEnergyDetailsResponse, error) {
	args := url.Values{
		"startTime": []string{startTime.Format(timeFormat)},
		"endTime":   []string{endTime.Format(timeFormat)},
		"timeUnit":  []string{string(timeUnit)},
	}
	return call[GetEnergyDetailsResponse](ctx, c, makePath("/site/{siteId}/energyDetails", id), args)
}

type GetEnergyDetailsResponse struct {
	EnergyDetails EnergyDetails `json:"energyDetails"`
}

// EnergyDetails contains site energy measurements from meters such as consumption, export (feed-in), import (purchase), etc.
type EnergyDetails struct {
	TimeUnit TimeUnit        `json:"timeUnit"`
	Unit     string          `json:"unit"`
	Meters   []MeterReadings `json:"meters"`
}

// GetPowerFlow returns the current power flow between all elements of the site including PV array, storage (battery), loads (consumption) and grid.
//
// TODO: test this during daytime
func (c *Client) GetPowerFlow(ctx context.Context, id int) (GetPowerFlowResponse, error) {
	return call[GetPowerFlowResponse](ctx, c, makePath("/site/{siteId}/currentPowerFlow", id), nil)
}

type GetPowerFlowResponse struct {
	CurrentPowerFlow PowerFlow `json:"siteCurrentPowerFlow"`
}

// PowerFlow contains current power flow between all elements of the site including PV array, storage (battery), loads (consumption) and grid.
type PowerFlow struct {
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

// GetStorageData returns detailed information from batteries installed at the active site.
//
// This API is limited to a one-week period.
func (c *Client) GetStorageData(ctx context.Context, id int, startTime, endTime time.Time) (GetStorageDataResponse, error) {
	args := url.Values{
		"startTime": []string{startTime.Format(timeFormat)},
		"endTime":   []string{endTime.Format(timeFormat)},
	}
	return call[GetStorageDataResponse](ctx, c, makePath("/site/{siteId}/storageData", id), args)
}

type GetStorageDataResponse struct {
	StorageData StorageData `json:"storageData"`
}

type StorageData struct {
	Batteries    []Battery `json:"batteries"`
	BatteryCount int       `json:"batteryCount"`
}

// Battery contains detailed storage information from batteries: the state of energy, power and lifetime energy.
type Battery struct {
	SerialNumber   string             `json:"serialNumber"`
	ModelNumber    string             `json:"modelNumber"`
	Telemetries    []BatteryTelemetry `json:"telemetries"`
	Nameplate      int                `json:"nameplate"`
	TelemetryCount int                `json:"telemetryCount"`
}

type BatteryTelemetry struct {
	TimeStamp                Time    `json:"timeStamp"`
	Power                    float64 `json:"power"`
	BatteryState             int     `json:"batteryState"`
	LifeTimeEnergyCharged    float64 `json:"lifeTimeEnergyCharged"`
	LifeTimeEnergyDischarged float64 `json:"lifeTimeEnergyDischarged"`
	FullPackEnergyAvailable  float64 `json:"fullPackEnergyAvailable"`
	InternalTemp             float64 `json:"internalTemp"`
	ACGridCharging           float64 `json:"ACGridCharging"`
}

// GetEnvBenefits returns all environmental benefits based on site energy production: gas emissions saved, equivalent trees planted and light bulbs powered for a day.
func (c *Client) GetEnvBenefits(ctx context.Context, id int) (GetEnvBenefitsResponse, error) {
	return call[GetEnvBenefitsResponse](ctx, c, makePath("/site/{siteId}/envBenefits", id), nil)
}

type GetEnvBenefitsResponse struct {
	EnvBenefits EnvBenefits `json:"envBenefits"`
}

type EnvBenefits struct {
	GasEmissionSaved struct {
		Units string  `json:"units"`
		Co2   float64 `json:"co2"`
		Nox   float64 `json:"nox"`
		So2   float64 `json:"so2"`
	} `json:"gasEmissionSaved"`
	LightBulbs   float64 `json:"lightBulbs"`
	TreesPlanted float64 `json:"treesPlanted"`
}
