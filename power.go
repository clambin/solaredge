package solaredge

import (
	"net/url"
	"strconv"
	"time"
)

type PowerMeasurement struct {
	Time  time.Time
	Value float64
}

type TimeStamp struct {
	TS time.Time
}

func (ts *TimeStamp) UnmarshalJSON(buf []byte) (err error) {
	var t time.Time
	t, err = time.Parse("\"2006-01-02 15:04:05\"", string(buf))

	if err == nil {
		ts.TS = t
	}
	return
}

func (client *Client) GetPower(siteID int, startTime, endTime time.Time) (entries []PowerMeasurement, err error) {
	args := url.Values{}

	args.Set("startTime", startTime.Format("2006-01-02 15:04:05"))
	args.Set("endTime", endTime.Format("2006-01-02 15:04:05"))

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

	err = client.call("/site/"+strconv.Itoa(siteID)+"/power", args, &powerStats)

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

func (client *Client) GetPowerOverview(siteID int) (lifeTime, lastYear, lastMonth, lastDay, current float64, err error) {
	args := url.Values{}

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

	err = client.call("/site/"+strconv.Itoa(siteID)+"/overview", args, &overviewResponse)

	if err == nil {
		lifeTime = overviewResponse.Overview.LifeTimeData.Energy
		lastYear = overviewResponse.Overview.LastYearData.Energy
		lastMonth = overviewResponse.Overview.LastMonthData.Energy
		lastDay = overviewResponse.Overview.LastDayData.Energy
		current = overviewResponse.Overview.CurrentPower.Power
	}
	return
}
