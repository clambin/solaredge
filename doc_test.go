package solaredge_test

import (
	"context"
	"fmt"
	"github.com/clambin/solaredge"
	"os"
	"time"
)

func ExampleClient_GetSiteDetails() {
	ctx := context.Background()
	c := solaredge.Client{Token: os.Getenv("SOLAREDGE_APIKEY")}

	site, err := c.GetSiteDetails(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Site '%s' (%s), Peak Power: %.1f\n", site.Name, site.Status, site.PeakPower)
}

func ExampleClient_GetPowerOverview() {
	ctx := context.Background()
	c := solaredge.Client{Token: os.Getenv("SOLAREDGE_APIKEY")}

	site, err := c.GetSiteDetails(ctx)
	if err != nil {
		panic(err)
	}
	overview, err := c.GetPowerOverview(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Site: %s\nLast update: %s\nCurrent power: %.1fW", site.Name, overview.LastUpdateTime, overview.CurrentPower.Power)
}

func ExampleClient_GetInverterTelemetry() {
	c := solaredge.Client{Token: os.Getenv("SOLAREDGE_APIKEY")}

	ctx := context.Background()
	end := time.Now()
	start := end.Add(-1 * 24 * time.Hour)

	inventory, err := c.GetInventory(ctx)
	if err != nil {
		panic(err)
	}

	for _, inverter := range inventory.Inverters {
		telemetry, err := c.GetInverterTelemetry(ctx, inventory.Inverters[0].SN, start, end)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, entry := range telemetry {
			fmt.Printf("%s - %s - %5.1f V - %4.1f ºC - %6.1f\n", inverter.Name, entry.Time, entry.DcVoltage, entry.Temperature, entry.TotalActivePower)
		}
	}
}
