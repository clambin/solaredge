package solaredge

import (
	"encoding/json"
	"time"
)

type Value struct {
	Date  Time    `json:"date"`
	Value float64 `json:"value,omitempty"`
}

type Date time.Time

var _ json.Marshaler = Date{}
var _ json.Unmarshaler = &Date{}

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

func (d Date) String() string {
	return time.Time(d).String()
}

var _ json.Unmarshaler = &Time{}
var _ json.Marshaler = &Time{}

type Time time.Time

func (t *Time) UnmarshalJSON(bytes []byte) error {
	date, err := time.Parse(`"2006-01-02 15:04:05"`, string(bytes))
	if err == nil {
		*t = Time(date)
	}
	return err
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(`"2006-01-02 15:04:05"`)), nil
}

func (t Time) String() string {
	return time.Time(t).String()
}
