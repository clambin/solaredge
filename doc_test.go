package solaredge_test

import (
	"context"
	"fmt"
	"github.com/clambin/solaredge/v2"
	"os"
	"time"
)

func ExampleClient_GetSites() {
	ctx := context.Background()
	c := solaredge.Client{SiteKey: os.Getenv("SOLAREDGE_APIKEY")}

	resp, err := c.GetSites(ctx)
	if err != nil {
		panic(err)
	}
	for _, site := range resp.Sites.Site {
		fmt.Printf("Site '%s' (%s), Peak Power: %.1f\n", site.Name, site.Status, site.PeakPower)
	}
}

func ExampleClient_GetPowerOverview() {
	ctx := context.Background()
	c := solaredge.Client{SiteKey: os.Getenv("SOLAREDGE_APIKEY")}

	resp, err := c.GetSites(ctx)
	if err != nil {
		panic(err)
	}

	for _, site := range resp.Sites.Site {
		resp, err := c.GetPowerOverview(ctx, site.Id)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Site: %s\nLast update: %s\nCurrent power: %.1fW", site.Name, time.Time(resp.Overview.LastUpdateTime).String(), resp.Overview.CurrentPower.Power)
	}
}

func ExampleClient_GetInverterTechnicalData() {
	c := solaredge.Client{SiteKey: os.Getenv("SOLAREDGE_APIKEY")}

	ctx := context.Background()
	end := time.Now()
	start := end.Add(-1 * 24 * time.Hour)

	resp, err := c.GetSites(ctx)
	if err != nil {
		panic(err)
	}

	for _, site := range resp.Sites.Site {
		inventory, err := c.GetComponents(ctx, site.Id)
		if err != nil {
			panic(err)
		}

		for _, inverter := range inventory.Reporters.List {
			telemetry, err := c.GetInverterTechnicalData(ctx, site.Id, inverter.SerialNumber, start, end)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			for _, entry := range telemetry.Data.Telemetries {
				fmt.Printf("%s - %s - %5.1f V - %4.1f ºC - %6.1f\n", inverter.Name, time.Time(entry.Time).String(), entry.DcVoltage, entry.Temperature, entry.TotalActivePower)
			}
		}
	}
}
