package solaredge

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

// PowerMeasurement contains a power measurement from the solar panels at a point in time
type PowerMeasurement struct {
	Time  time.Time
	Value float64
}

// TimeStamp represents a timestamp. Used to unmarshal the server response
type TimeStamp struct {
	TS time.Time
}

// UnmarshalJSON parses a timestamp in the server response
func (ts *TimeStamp) UnmarshalJSON(buf []byte) (err error) {
	var t time.Time
	t, err = time.Parse("\"2006-01-02 15:04:05\"", string(buf))

	if err == nil {
		ts.TS = t
	}
	return
}

// GetPower returns the PowerMeasurements for the specified site and timeframe
func (client *Client) GetPower(ctx context.Context, siteID int, startTime, endTime time.Time) (entries []PowerMeasurement, err error) {
	var powerStats struct {
		Power struct {
			TimeUnit   string
			Unit       string
			MeasuredBy string
			Values     []struct {
				Date  TimeStamp
				Value *float64
			}
		}
	}

	args := url.Values{}
	args.Set("startTime", startTime.Format("2006-01-02 15:04:05"))
	args.Set("endTime", endTime.Format("2006-01-02 15:04:05"))

	err = client.call(ctx, "/site/"+strconv.Itoa(siteID)+"/power", args, &powerStats)

	if err == nil {
		for _, entry := range powerStats.Power.Values {
			if entry.Value != nil {
				entries = append(entries, PowerMeasurement{
					Time:  entry.Date.TS,
					Value: *entry.Value,
				})
			}
		}
	}

	return
}

// GetPowerOverview returns the energy produced at the site for its entire lifetime, current year, month and day (in Wh) and current power (in W)
func (client *Client) GetPowerOverview(ctx context.Context, siteID int) (lifeTime, currentYear, currentMonth, currentDay, current float64, err error) {
	var overviewResponse struct {
		Overview struct {
			LastUpdateTime TimeStamp
			LifeTimeData   struct {
				Energy float64
			}
			LastYearData struct {
				Energy float64
			}
			LastMonthData struct {
				Energy float64
			}
			LastDayData struct {
				Energy float64
			}
			CurrentPower struct {
				Power float64
			}
			MeasuredBy string
		}
	}

	args := url.Values{}
	err = client.call(ctx, "/site/"+strconv.Itoa(siteID)+"/overview", args, &overviewResponse)

	if err == nil {
		lifeTime = overviewResponse.Overview.LifeTimeData.Energy
		currentYear = overviewResponse.Overview.LastYearData.Energy
		currentMonth = overviewResponse.Overview.LastMonthData.Energy
		currentDay = overviewResponse.Overview.LastDayData.Energy
		current = overviewResponse.Overview.CurrentPower.Power
	}
	return
}
