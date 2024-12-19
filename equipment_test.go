package solaredge

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestClient_GetComponents(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetComponents(context.Background(), 1)
	expect(t, resp, "/equipment/1/list", err)
}

func TestClient_GetInventory(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetInventory(context.Background(), 1)
	expect(t, resp, "/site/1/inventory", err)
}

func TestClient_GetInverterTechnicalData(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetInverterTechnicalData(context.Background(), 1, "SN1", time.Time{}, time.Time{})
	expect(t, resp, "/equipment/1/SN1/data", err)
}

func TestClient_GetEquipmentChangeLog(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetEquipmentChangeLog(context.Background(), 1, "SN1")
	expect(t, resp, "/equipment/1/SN1/changeLog", err)
}
