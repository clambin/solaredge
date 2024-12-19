package solaredge_test

import (
	"context"
	"github.com/clambin/solaredge/v2"
	"os"
	"testing"
	"time"
)

func TestClient_E2E(t *testing.T) {
	siteKey := os.Getenv("SOLAREDGE_TOKEN")
	if siteKey == "" {
		t.Skip("Skipping: environment variable SOLAREDGE_TOKEN is not set")
	}

	c := solaredge.Client{SiteKey: siteKey}
	ctx := context.Background()
	sites, err := c.GetSites(ctx)
	if err != nil {
		t.Fatalf("failed to get sites: %v", err)
	}
	if sites.Sites.Count != 1 {
		t.Fatalf("expected 1 site, got %d", sites.Sites.Count)
	}
	id := sites.Sites.Site[0].Id

	details, err := c.GetSiteDetails(ctx, id)
	if err != nil {
		t.Fatalf("failed to get site details: %v", err)
	}
	if details.Details.Id != id {
		t.Fatalf("expected id %q, got %q", id, details.Details.Id)
	}

	overview, err := c.GetPowerOverview(ctx, id)
	if err != nil {
		t.Fatalf("failed to get site overview: %v", err)
	}
	if overview.Overview.LifeTimeData.Energy == 0 {
		t.Fatalf("expected overview.Overview.LifeTimeData.Energy, got zero")
	}

	endTime := time.Now()
	startTime := endTime.AddDate(0, 0, -1)
	energy, err := c.GetEnergyDetails(ctx, id, solaredge.TimeUnitDay, startTime, endTime)
	if err != nil {
		t.Fatalf("failed to get site energy: %v", err)
	}
	if len(energy.EnergyDetails.Meters) == 0 {
		t.Fatalf("expected non-empty energyDetails.Meters")
	}

	components, err := c.GetComponents(ctx, id)
	if err != nil {
		t.Fatalf("failed to get site components: %v", err)
	}
	if n := components.Reporters.Count; n != 1 {
		t.Fatalf("expected 1 reporter, got %d", n)
	}

	data, err := c.GetInverterTechnicalData(ctx, id, components.Reporters.List[0].SerialNumber, startTime, endTime)
	if err != nil {
		t.Fatalf("failed to get site inverter technical data: %v", err)
	}
	if data.Data.Count == 0 {
		t.Fatalf("expected non-empty data.Count, got 0")
	}
}
