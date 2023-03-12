package solaredge

import (
	"context"
	"strconv"
	"time"
)

// Inverter contains an inverter's name, model, manufacturer and serial number, as returned by GetInverters.
type Inverter struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Name         string `json:"name"`
	SerialNumber string `json:"serialNumber"`
	client       *Client
	site         *Site
}

// GetSerialNumber returns the inverter's SerialNumber
func (i *Inverter) GetSerialNumber() string {
	return i.SerialNumber
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

// GetTelemetry returns the measurements of the inverter within the specified time range.
//
// This API is limited to a one-week period. If the time range exceeds one week, an APIError is returned.
//
// NOTE: this may not be fully complete, as data returned for my account doesn't match the specifications.
func (i *Inverter) GetTelemetry(ctx context.Context, from, to time.Time) ([]InverterTelemetry, error) {
	var output struct {
		Data struct {
			Count       int                 `json:"count"`
			Telemetries []InverterTelemetry `json:"telemetries"`
		} `json:"data"`
	}
	args, err := buildArgsFromTimeRange(from, to, "Time", "2006-01-02 15:04:05")
	if err == nil {
		err = i.client.call(ctx, "/equipment/"+strconv.Itoa(i.site.ID)+"/"+i.SerialNumber+"/data", args, &output)
	}
	return output.Data.Telemetries, err
}

// InverterEquipment contains an inverter's name, model, manufacturer, serial number, etc, as returned by GetEquipment.
type InverterEquipment struct {
	SN                  string `json:"SN"`
	CommunicationMethod string `json:"communicationMethod"`
	ConnectedOptimizers int    `json:"connectedOptimizers"`
	CPUVersion          string `json:"cpuVersion"`
	Manufacturer        string `json:"manufacturer"`
	Model               string `json:"model"`
	Name                string `json:"name"`
	client              *Client
	site                *Site
}

// GetSerialNumber returns the inverter's SerialNumber
func (i *InverterEquipment) GetSerialNumber() string {
	return i.SN
}

// GetTelemetry returns the measurements of the inverter within the specified time range.
//
// This API is limited to a one-week period. If the time range exceeds one week, an APIError is returned.
//
// NOTE: this may not be fully complete, as data returned for my account doesn't match the specifications.
func (i *InverterEquipment) GetTelemetry(ctx context.Context, from, to time.Time) ([]InverterTelemetry, error) {
	var output struct {
		Data struct {
			Count       int                 `json:"count"`
			Telemetries []InverterTelemetry `json:"telemetries"`
		} `json:"data"`
	}
	args, err := buildArgsFromTimeRange(from, to, "Time", "2006-01-02 15:04:05")
	if err == nil {
		err = i.client.call(ctx, "/equipment/"+strconv.Itoa(i.site.ID)+"/"+i.SN+"/data", args, &output)
	}
	return output.Data.Telemetries, err
}
