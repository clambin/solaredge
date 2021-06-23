package solaredge_test

import (
	"bytes"
	"io"
	"net/http"
)

type Server struct {
}

func (server *Server) serve(req *http.Request) *http.Response {
	values := req.URL.Query()

	var value []string
	var ok bool

	value, ok = values["api_key"]

	if ok == false || len(value) == 0 || value[0] != "TESTTOKEN" {
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Status:     "403 Forbidden",
		}
	}

	var body string

	switch req.URL.Path {
	case "/sites/list":
		body = `{ "sites": { "count": 1, "site": [ { "id": 1, "name": "foo" } ] } }`
	case "/site/1/overview":
		body = `{ "overview": { "lastUpdateTime": "2021-05-19 17:08:23", 
			"lifeTimeData": { "energy": 10000.0 },
    		"lastYearData": { "energy": 1000.0 },
		    "lastMonthData": { "energy": 100.0 },
		    "lastDayData": { "energy": 10.0 },
			"currentPower": { "power": 1.0 },
    		"measuredBy": "INVERTER" } }`
	case "/site/1/power":
		body = `{ "power": { "timeUnit": "QUARTER_OF_AN_HOUR", "unit": "W", "measuredBy": "INVERTER", "values": [
			{ "date": "2021-05-18 00:00:00", "value": 12.0 },
      		{ "date": "2021-05-18 00:15:00", "value": 24.0 },
      		{ "date": "2021-05-18 00:15:00", "value": null } ] } }`
	default:
		return &http.Response{
			StatusCode: http.StatusNotFound,
			Status:     "API " + req.URL.Path + " not implemented",
		}
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}
