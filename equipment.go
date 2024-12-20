package solaredge

import (
	"context"
	"net/url"
	"time"
)

// This file implements the "Site Equipment API" section of the SolarEdge API specifications.
// https://knowledge-center.solaredge.com/sites/kc/files/se_monitoring_api.pdf

// GetComponents returns a list of inverters/SMIs in the specific site.
func (c *Client) GetComponents(ctx context.Context, id int) (GetComponentsResponse, error) {
	return call[GetComponentsResponse](ctx, c, makePath("/equipment/{siteId}/list", id), nil)
}

type GetComponentsResponse struct {
	Reporters struct {
		List  []Inverter `json:"list"`
		Count int        `json:"count"`
	} `json:"reporters"`
}

// Inverter contains an inverter's name, model, manufacturer and serial number, as returned by GetInverters.
type Inverter struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Name         string `json:"name"`
	SerialNumber string `json:"serialNumber"`
}

// GetInventory the inventory of SolarEdge equipment in the site, including inverters/SMIs, batteries, meters, gateways and sensors.
func (c *Client) GetInventory(ctx context.Context, id int) (GetInventoryResponse, error) {
	return call[GetInventoryResponse](ctx, c, makePath("/site/{siteId}/inventory", id), nil)
}

type GetInventoryResponse struct {
	Inventory Inventory `json:"inventory"`
}

// Inventory contains all batteries, gateways, inverters, meters and sensors at a site.
type Inventory struct {
	Batteries []BatteryEquipment  `json:"batteries"`
	Gateways  []GatewayEquipment  `json:"gateways"`
	Inverters []InverterEquipment `json:"inverters"`
	Meters    []MeterEquipment    `json:"meters"`
	Sensors   []SensorEquipment   `json:"sensors"`
}

// BatteryEquipment contains a battery's name, model, manufacturer, serial number, etc.
type BatteryEquipment struct {
	Name                string  `json:"name"`
	Manufacturer        string  `json:"manufacturer"`
	Model               string  `json:"model"`
	FirmwareVersion     string  `json:"firmwareVersion"`
	ConnectedInverterSn string  `json:"connectedInverterSn"`
	SN                  string  `json:"SN"`
	NameplateCapacity   float64 `json:"nameplateCapacity"`
}

// InverterEquipment contains an inverter's name, model, manufacturer, serial number, etc. as returned by GetEquipment.
type InverterEquipment struct {
	SN                  string `json:"SN"`
	CommunicationMethod string `json:"communicationMethod"`
	CPUVersion          string `json:"cpuVersion"`
	Manufacturer        string `json:"manufacturer"`
	Model               string `json:"model"`
	Name                string `json:"name"`
	ConnectedOptimizers int    `json:"connectedOptimizers"`
}

// GatewayEquipment contains a gateway's name, model, manufacturer, serial number, etc.
type GatewayEquipment struct {
	Name            string `json:"name"`
	FirmwareVersion string `json:"firmwareVersion"`
	SN              string `json:"SN"`
}

// MeterEquipment contains an meter's name, model, manufacturer, serial number, etc.
type MeterEquipment struct {
	Name                       string `json:"name"`
	Manufacturer               string `json:"manufacturer"`
	Model                      string `json:"model"`
	FirmwareVersion            string `json:"firmwareVersion"`
	ConnectedSolarEdgeDeviceSN string `json:"connectedSolaredgeDeviceSN"`
	Type                       string `json:"type"`
	Form                       string `json:"form"`
}

// SensorEquipment contains an sensor's name, model, manufacturer, serial number, etc.
type SensorEquipment struct {
	ConnectedSolarEdgeDeviceSN string `json:"connectedSolaredgeDeviceSN"`
	ID                         string `json:"id"`
	ConnectedTo                string `json:"connectedTo"`
	Category                   string `json:"category"`
	Type                       string `json:"type"`
}

// GetInverterTechnicalData returns specific inverter data for a given timeframe.
//
// Notes:
//   - This API is limited to a one-week period. If the time range exceeds one week, an error is returned.
//   - this may not be fully complete, as data returned for my account doesn't match the specifications.
func (c *Client) GetInverterTechnicalData(ctx context.Context, id int, serialNr string, startTime, endTime time.Time) (GetInverterTechnicalDataResponse, error) {
	args := url.Values{
		"startTime": []string{startTime.Format(timeFormat)},
		"endTime":   []string{endTime.Format(timeFormat)},
	}
	return call[GetInverterTechnicalDataResponse](ctx, c, makePath("/equipment/{siteId}/"+serialNr+"/data", id), args)
}

type GetInverterTechnicalDataResponse struct {
	Data struct {
		Telemetries []InverterTelemetry `json:"telemetries"`
		Count       int                 `json:"count"`
	} `json:"data"`
}

// InverterTelemetry contains technical data for an inverter.
type InverterTelemetry struct {
	Time                  Time                    `json:"date"`
	InverterMode          string                  `json:"inverterMode"`
	L1Data                InverterTelemetryL1Data `json:"L1Data"`
	DcVoltage             float64                 `json:"dcVoltage"`
	GroundFaultResistance float64                 `json:"groundFaultResistance,omitempty"`
	OperationMode         int                     `json:"operationMode"`
	PowerLimit            float64                 `json:"powerLimit"`
	Temperature           float64                 `json:"temperature"`
	TotalActivePower      float64                 `json:"totalActivePower"`
	TotalEnergy           float64                 `json:"totalEnergy"`
}

type InverterTelemetryL1Data struct {
	AcCurrent     float64 `json:"acCurrent"`
	AcFrequency   float64 `json:"acFrequency"`
	AcVoltage     float64 `json:"acVoltage"`
	ActivePower   float64 `json:"activePower"`
	ApparentPower float64 `json:"apparentPower"`
	CosPhi        float64 `json:"cosPhi"`
	ReactivePower float64 `json:"reactivePower"`
}

// GetEquipmentChangeLog returns a list of equipment component replacements ordered by date. This method is applicable to inverters, optimizers, batteries and gateways.
func (c *Client) GetEquipmentChangeLog(ctx context.Context, id int, serialNr string) (GetEquipmentChangeLogResponse, error) {
	return call[GetEquipmentChangeLogResponse](ctx, c, makePath("/equipment/{siteId}/"+serialNr+"/changeLog", id), nil)
}

type GetEquipmentChangeLogResponse struct {
	ChangeLog struct {
		List  []EquipmentChangeLog `json:"list"`
		Count int                  `json:"count"`
	} `json:"ChangeLog"`
}
type EquipmentChangeLog struct {
	SerialNumber string `json:"serialNumber"`
	PartNumber   string `json:"partNumber"`
	Date         string `json:"date"`
}
