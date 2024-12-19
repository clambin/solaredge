package solaredge

import (
	"time"
)

type Date time.Time

func (d *Date) UnmarshalJSON(bytes []byte) error {
	date, err := time.Parse(`"2006-01-02"`, string(bytes))
	if err == nil {
		*d = Date(date)
	}
	return err
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(d).Format(`"2006-01-02"`)), nil
}

type Time time.Time

func (t *Time) UnmarshalJSON(bytes []byte) error {
	date, err := time.Parse(`"`+timeFormat+`"`, string(bytes))
	if err == nil {
		*t = Time(date)
	}
	return err
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(`"2006-01-02 15:04:05"`)), nil
}

// Value is a common data type in the SolarEdge API. It represents a measurement at a moment in time.
type Value struct {
	Date  Time    `json:"date"`
	Value float64 `json:"value"`
}

// TimeUnit defines the granularity of the data to be returned.
//
// Note: the chosen TimeUnit may impose limits the start & end times/dates. See the relevant API for details.
type TimeUnit string

const (
	TimeUnitQuarter TimeUnit = "QUARTER_OF_AN_HOUR"
	TimeUnitHour    TimeUnit = "HOUR"
	TimeUnitDay     TimeUnit = "DAY"
	TimeUnitWeek    TimeUnit = "WEEK"
	TimeUnitMonth   TimeUnit = "MONTH"
	TimeUnitYear    TimeUnit = "YEAR"
)
