package solaredge

import (
	"context"
	"net/url"
	"strconv"
)

// This file contains all supported APIs from the Site Equipment API section of the SolarEdge API specifications.
// https://knowledge-center.solaredge.com/sites/kc/files/se_monitoring_api.pdf

// GetInverters returns all inverters at the active site.
func (s *Site) GetInverters(ctx context.Context) ([]Inverter, error) {
	var output struct {
		Reporters struct {
			Count int        `json:"count"`
			List  []Inverter `json:"list"`
		} `json:"reporters"`
	}
	err := s.client.call(ctx, "/equipment/"+strconv.Itoa(s.ID)+"/list", url.Values{}, &output)
	if err == nil {
		for index := range output.Reporters.List {
			output.Reporters.List[index].client = s.client
			output.Reporters.List[index].site = s
		}
	}
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

// GetInventory returns the full inventory of SolarEdge equipment at the active site.
func (s *Site) GetInventory(ctx context.Context) (Inventory, error) {
	var output struct {
		Inventory Inventory `json:"inventory"`
	}
	err := s.client.call(ctx, "/site/"+strconv.Itoa(s.ID)+"/inventory", url.Values{}, &output)
	if err == nil {
		for index := range output.Inventory.Inverters {
			output.Inventory.Inverters[index].client = s.client
			output.Inventory.Inverters[index].site = s

		}
	}

	return output.Inventory, err
}

// ChangeLogEntry contains an equipment component replacement.
type ChangeLogEntry struct {
	Date         Date   `json:"date"`
	PartNumber   string `json:"partNumber"`
	SerialNumber string `json:"serialNumber"`
}

// GetChangeLog returns a list of equipment component replacements, ordered by date. This method is applicable to inverters, optimizers, batteries and gateways.
func (s *Site) GetChangeLog(ctx context.Context, serialNr string) ([]ChangeLogEntry, error) {
	var output struct {
		ChangeLog struct {
			Count int              `json:"count"`
			List  []ChangeLogEntry `json:"list"`
		} `json:"ChangeLog"`
	}
	err := s.client.call(ctx, "/equipment/"+strconv.Itoa(s.ID)+"/"+serialNr+"/changeLog", url.Values{}, &output)
	return output.ChangeLog.List, err
}
