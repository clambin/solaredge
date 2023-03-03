package solaredge

import (
	"context"
	"net/url"
	"time"
)

// This file contains all supported APIs from the Site Equipment API section of the SolarEdge API specifications.
// https://knowledge-center.solaredge.com/sites/kc/files/se_monitoring_api.pdf

// Inverter contains an inverter's name, model, manufacturer and serial number.
type Inverter struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Name         string `json:"name"`
	SerialNumber string `json:"serialNumber"`
}

// GetInverters returns all inverters at the active site.
func (c *Client) GetInverters(ctx context.Context) ([]Inverter, error) {
	var output struct {
		Reporters struct {
			Count int        `json:"count"`
			List  []Inverter `json:"list"`
		} `json:"reporters"`
	}
	err := c.call(ctx, "/equipment/%d/list", url.Values{}, &output)
	return output.Reporters.List, err
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
	NameplateCapacity   float64 `json:"nameplateCapacity"`
	SN                  string  `json:"SN"`
}

// GatewayEquipment contains a gateway's name, model, manufacturer, serial number, etc.
type GatewayEquipment struct {
	Name            string `json:"name"`
	FirmwareVersion string `json:"firmwareVersion"`
	SN              string `json:"SN"`
}

// InverterEquipment contains an inverter's name, model, manufacturer, serial number, etc.
type InverterEquipment struct {
	SN                  string `json:"SN"`
	CommunicationMethod string `json:"communicationMethod"`
	ConnectedOptimizers int    `json:"connectedOptimizers"`
	CpuVersion          string `json:"cpuVersion"`
	Manufacturer        string `json:"manufacturer"`
	Model               string `json:"model"`
	Name                string `json:"name"`
}

// MeterEquipment contains an meter's name, model, manufacturer, serial number, etc.
type MeterEquipment struct {
	Name                       string `json:"name"`
	Manufacturer               string `json:"manufacturer"`
	Model                      string `json:"model"`
	FirmwareVersion            string `json:"firmwareVersion"`
	ConnectedSolaredgeDeviceSN string `json:"connectedSolaredgeDeviceSN"`
	Type                       string `json:"type"`
	Form                       string `json:"form"`
}

// SensorEquipment contains an sensor's name, model, manufacturer, serial number, etc.
type SensorEquipment struct {
	ConnectedSolaredgeDeviceSN string `json:"connectedSolaredgeDeviceSN"`
	Id                         string `json:"id"`
	ConnectedTo                string `json:"connectedTo"`
	Category                   string `json:"category"`
	Type                       string `json:"type"`
}

// GetInventory returns the full inventory of SolarEdge equipment at the active site.
func (c *Client) GetInventory(ctx context.Context) (Inventory, error) {
	var output struct {
		Inventory Inventory `json:"inventory"`
	}
	err := c.call(ctx, "/site/%d/inventory", url.Values{}, &output)
	return output.Inventory, err
}

// InverterTelemetry contains technical data for an inverter.
type InverterTelemetry struct {
	L1Data struct {
		AcCurrent     float64 `json:"acCurrent"`
		AcFrequency   float64 `json:"acFrequency"`
		AcVoltage     float64 `json:"acVoltage"`
		ActivePower   float64 `json:"activePower"`
		ApparentPower float64 `json:"apparentPower"`
		CosPhi        float64 `json:"cosPhi"`
		ReactivePower float64 `json:"reactivePower"`
	} `json:"L1Data"`
	Time                  Time    `json:"date"`
	DcVoltage             float64 `json:"dcVoltage"`
	GroundFaultResistance float64 `json:"groundFaultResistance,omitempty"`
	InverterMode          string  `json:"inverterMode"`
	OperationMode         int     `json:"operationMode"`
	PowerLimit            float64 `json:"powerLimit"`
	Temperature           float64 `json:"temperature"`
	TotalActivePower      float64 `json:"totalActivePower"`
	TotalEnergy           float64 `json:"totalEnergy"`
}

// GetInverterTelemetry returns technical data for the inverter with the provided serialNr for a given timeframe.
// Telemetry data is returned in 5-minute intervals.
//
// This API is limited to a one-week period. If the time range exceeds one week, an APIError is returned.
//
// NOTE: this may not be fully complete, as data returned for my account doesn't match the specifications.
func (c *Client) GetInverterTelemetry(ctx context.Context, serialNr string, from, to time.Time) ([]InverterTelemetry, error) {
	var output struct {
		Data struct {
			Count       int                 `json:"count"`
			Telemetries []InverterTelemetry `json:"telemetries"`
		} `json:"data"`
	}
	args, err := buildArgsFromTimeRange(from, to, "Time", "2006-01-02 03:04:05")
	if err == nil {
		err = c.call(ctx, "/equipment/%d/"+serialNr+"/data", args, &output)
	}
	return output.Data.Telemetries, err
}

// ChangeLogEntry contains an equipment component replacement.
type ChangeLogEntry struct {
	Date         Date   `json:"date"`
	PartNumber   string `json:"partNumber"`
	SerialNumber string `json:"serialNumber"`
}

// GetChangeLog returns a list of equipment component replacements, ordered by date. This method is applicable to inverters, optimizers, batteries and gateways.
func (c *Client) GetChangeLog(ctx context.Context, serialNr string) ([]ChangeLogEntry, error) {
	var output struct {
		ChangeLog struct {
			Count int              `json:"count"`
			List  []ChangeLogEntry `json:"list"`
		} `json:"ChangeLog"`
	}
	err := c.call(ctx, "/equipment/%d/"+serialNr+"/changeLog", url.Values{}, &output)
	return output.ChangeLog.List, err
}
