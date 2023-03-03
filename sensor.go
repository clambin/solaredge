package solaredge

/*
type Sensor struct {
	ConnectedTo string `json:"connectedTo"`
	Count       int    `json:"count"`
	Sensors     []struct {
		Name        string `json:"name"`
		Measurement string `json:"measurement"`
		Type        string `json:"type"`
	} `json:"sensors"`
}

// TODO: doesn't return any data?

func (c *Client) GetSensors(ctx context.Context) ([]Sensor, error) {
	var output struct {
		SiteSensors struct {
			Count int      `json:"count"`
			List  []Sensor `json:"list"`
		} `json:"SiteSensors"`
	}
	err := c.call(ctx, "/equipment/%d/sensors", url.Values{}, &output)
	return output.SiteSensors.List, err
}

type SensorData struct {
	ConnectedTo string `json:"connectedTo"`
	Count       int    `json:"count"`
	Telemetries []struct {
		// TODO: this list isn't complete in the API documentation
		Date               Time    `json:"date"`
		AmbientTemperature float64 `json:"ambientTemperature"`
		ModuleTemperature  float64 `json:"moduleTemperature"`
		WindSpeed          float64 `json:"windSpeed"`
	} `json:"telemetries"`
}

// TODO: returns 404?

func (c *Client) GetSensorData(ctx context.Context, from, to time.Time) ([]SensorData, error) {
	args, err := buildArgsFromTimeRange(from, to, "Date", "2006-01-02 03:04:05")
	if err != nil {
		return nil, err
	}
	var output struct {
		SiteSensors struct {
			Data []SensorData `json:"data"`
		}
	}
	err = c.call(ctx, "/site/%d/sensors", args, &output)
	return output.SiteSensors.Data, err
}


*/
