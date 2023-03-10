package solaredge_test

import (
	"context"
	"fmt"
	"github.com/clambin/solaredge"
	"os"
	"time"
)

func ExampleClient_GetSites() {
	ctx := context.Background()
	c := solaredge.Client{Token: os.Getenv("SOLAREDGE_APIKEY")}

	sites, err := c.GetSites(ctx)
	if err != nil {
		panic(err)
	}
	for _, site := range sites {
		fmt.Printf("Site '%s' (%s), Peak Power: %.1f\n", site.Name, site.Status, site.PeakPower)
	}
}

func ExampleSite_GetPowerOverview() {
	ctx := context.Background()
	c := solaredge.Client{Token: os.Getenv("SOLAREDGE_APIKEY")}

	sites, err := c.GetSites(ctx)
	if err != nil {
		panic(err)
	}

	for _, site := range sites {
		overview, err := site.GetPowerOverview(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Site: %s\nLast update: %s\nCurrent power: %.1fW", site.Name, overview.LastUpdateTime, overview.CurrentPower.Power)
	}
}

func ExampleInverterEquipment_GetTelemetry() {
	c := solaredge.Client{Token: os.Getenv("SOLAREDGE_APIKEY")}

	ctx := context.Background()
	end := time.Now()
	start := end.Add(-1 * 24 * time.Hour)

	sites, err := c.GetSites(ctx)
	if err != nil {
		panic(err)
	}

	for _, site := range sites {
		inventory, err := site.GetInventory(ctx)
		if err != nil {
			panic(err)
		}

		for _, inverter := range inventory.Inverters {
			telemetry, err := inverter.GetTelemetry(ctx, start, end)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			for _, entry := range telemetry {
				fmt.Printf("%s - %s - %5.1f V - %4.1f ºC - %6.1f\n", inverter.Name, entry.Time, entry.DcVoltage, entry.Temperature, entry.TotalActivePower)
			}
		}
	}
}
