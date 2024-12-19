package solaredge

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestSites_FindByID(t *testing.T) {
	tests := []struct {
		name  string
		s     Sites
		id    int
		want  SiteDetails
		want1 bool
	}{
		{
			name:  "match",
			s:     Sites{{Id: 1, Name: "foo"}, {Id: 2, Name: "bar"}},
			id:    1,
			want:  SiteDetails{Id: 1, Name: "foo"},
			want1: true,
		},
		{
			name:  "no match",
			s:     Sites{{Id: 1, Name: "foo"}, {Id: 2, Name: "bar"}},
			id:    3,
			want1: false,
		},
		{
			name:  "empty",
			s:     Sites{},
			id:    1,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.FindByID(tt.id)
			if got != tt.want {
				t.Errorf("FindByID() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindByID() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSites_FindByName(t *testing.T) {
	tests := []struct {
		name  string
		s     Sites
		desc  string
		want  SiteDetails
		want1 bool
	}{
		{
			name:  "match",
			s:     Sites{{Id: 1, Name: "foo"}, {Id: 2, Name: "bar"}},
			desc:  "foo",
			want:  SiteDetails{Id: 1, Name: "foo"},
			want1: true,
		},
		{
			name:  "no match",
			s:     Sites{{Id: 1, Name: "foo"}, {Id: 2, Name: "bar"}},
			desc:  "snafu",
			want1: false,
		},
		{
			name: "empty",
			s:    Sites{},
			desc: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.FindByName(tt.desc)
			if got != tt.want {
				t.Errorf("FindByName() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FindByName() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestClient_GetSites(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetSites(context.Background())
	expect(t, resp, "/sites/list", err)
}

func TestClient_GetSiteDetails(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetSiteDetails(context.Background(), 1)
	expect(t, resp, "/site/1/details", err)
}

func TestClient_GetDataPeriod(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetDataPeriod(context.Background(), 1)
	expect(t, resp, "/site/1/dataPeriod", err)
}

func TestClient_GetEnergyMeasurements(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetEnergyMeasurements(context.Background(), 1, TimeUnitDay, time.Time{}, time.Time{})
	expect(t, resp, "/site/1/energy", err)
}

func TestClient_GetEnergyForTimeFrame(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetEnergyForTimeFrame(context.Background(), 1, time.Time{}, time.Time{})
	expect(t, resp, "/site/1/timeFrameEnergy", err)
}

func TestClient_GetPowerMeasurements(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetPowerMeasurements(context.Background(), 1, time.Time{}, time.Time{})
	expect(t, resp, "/site/1/power", err)
}

func TestClient_GetPowerOverview(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetPowerOverview(context.Background(), 1)
	expect(t, resp, "/site/1/overview", err)
}

func TestClient_GetPowerDetails(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetPowerDetails(context.Background(), 1, time.Time{}, time.Time{})
	expect(t, resp, "/site/1/powerDetails", err)
}

func TestClient_GetEnergyDetails(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetEnergyDetails(context.Background(), 1, TimeUnitQuarter, time.Time{}, time.Time{})
	expect(t, resp, "/site/1/energyDetails", err)
}

func TestClient_GetPowerFlow(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetPowerFlow(context.Background(), 1)
	expect(t, resp, "/site/1/currentPowerFlow", err)
}

func TestClient_GetStorageData(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetStorageData(context.Background(), 1, time.Time{}, time.Time{})
	expect(t, resp, "/site/1/storageData", err)
}

func TestClient_GetEnvBenefits(t *testing.T) {
	c := Client{baseURL: testServer.URL, HTTPClient: http.DefaultClient}
	resp, err := c.GetEnvBenefits(context.Background(), 1)
	expect(t, resp, "/site/1/envBenefits", err)
}
