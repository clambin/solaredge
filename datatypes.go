package solaredge

import (
	"encoding/json"
	"time"
)

// Value is a common data type in the SolarEdge API. It represents a measurement at a moment in time.
type Value struct {
	Date  Time    `json:"date"`
	Value float64 `json:"value,omitempty"`
}

// Date represents a date (YYYY-MM-DD) in the SolarEdge API.
type Date time.Time

var _ json.Marshaler = Date{}
var _ json.Unmarshaler = &Date{}

// UnmarshalJSON unmarshals a date in a JSON object to a Date
func (d *Date) UnmarshalJSON(bytes []byte) error {
	date, err := time.Parse(`"2006-01-02"`, string(bytes))
	if err == nil {
		*d = Date(date)
	}
	return err
}

// MarshalJSON marshals a Date into a JSON object
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(d).Format(`"2006-01-02"`)), nil
}

// String returns a string representation of a Date. This is equivalent to time.Time{}.String()
func (d Date) String() string {
	return time.Time(d).String()
}

var _ json.Unmarshaler = &Time{}
var _ json.Marshaler = &Time{}

// Time represents a timestamp (YYYY-MM-DD HH:MI:SS) in the SolarEdge API.
type Time time.Time

// UnmarshalJSON unmarshals a timestamp in a JSON object to a Time
func (t *Time) UnmarshalJSON(bytes []byte) error {
	date, err := time.Parse(`"2006-01-02 15:04:05"`, string(bytes))
	if err == nil {
		*t = Time(date)
	}
	return err
}

// MarshalJSON marshals a Time into a JSON object
func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(`"2006-01-02 15:04:05"`)), nil
}

// String returns a string representation of a Time. This is equivalent to time.Time{}.String()
func (t Time) String() string {
	return time.Time(t).String()
}
