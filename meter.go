package solaredge

/*
type MeterEnergyDetails struct {
	Meters   []Meter `json:"meters"`
	TimeUnit string  `json:"timeUnit"`
	Unit     string  `json:"unit"`
}

type Meter struct {
	MeterSerialNumber          string  `json:"meterSerialNumber"`
	ConnectedSolarEdgeDeviceSN string  `json:"connectedSolaredgeDeviceSN"`
	Model                      string  `json:"model"`
	MeterType                  string  `json:"meterType"`
	Values                     []Value `json:"values"`
}

// TODO: doesn't return any data?

func (c *Client) GetMeters(ctx context.Context, timeUnit string, start, end time.Time) (MeterEnergyDetails, error) {
	args, err := buildArgsFromTimeRange(start, end, "Time", "2006-01-02 03:04:05")
	if err != nil {
		return MeterEnergyDetails{}, err
	}
	// TODO: valid timeUnit vs time range
	args.Set("timeUnit", timeUnit)
	var output struct {
		MeterEnergyDetails MeterEnergyDetails `json:"meterEnergyDetails"`
	}

	err = c.call(ctx, "/site/%d/meters", args, &output)
	return output.MeterEnergyDetails, err
}
*/
