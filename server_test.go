package solaredge_test

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Server struct {
	token string
	slow  bool
}

func (server *Server) apiHandler(w http.ResponseWriter, req *http.Request) {
	log.Debug("apiHandler: " + req.URL.Path)

	if server.slow && wait(req.Context(), 5*time.Second) == false {
		http.Error(w, "context exceeded", http.StatusRequestTimeout)
		return
	}

	if server.authenticate(req) == false {
		http.Error(w, "403 Forbidden", http.StatusForbidden)
		return
	}

	response, ok := responses[req.URL.Path]

	if ok == false {
		http.Error(w, "endpoint not implemented: "+req.URL.Path, http.StatusNotFound)
		return
	}

	_, _ = w.Write([]byte(response))
}

func wait(ctx context.Context, duration time.Duration) (passed bool) {
	timer := time.NewTimer(duration)
loop:
	for {
		select {
		case <-timer.C:
			break loop
		case <-ctx.Done():
			return false
		}
	}
	return true
}

func (server *Server) authenticate(req *http.Request) bool {
	values := req.URL.Query()
	value, ok := values["api_key"]

	return ok && len(value) > 0 && value[0] == server.token
}

var responses = map[string]string{
	"/sites/list": `{ "sites": { "count": 1, "site": [ { "id": 1, "name": "foo" } ] } }`,
	"/site/1/overview": `{ "overview": { "lastUpdateTime": "2021-05-19 17:08:23", 
			"lifeTimeData": { "energy": 10000.0 },
    		"lastYearData": { "energy": 1000.0 },
		    "lastMonthData": { "energy": 100.0 },
		    "lastDayData": { "energy": 10.0 },
			"currentPower": { "power": 1.0 },
    		"measuredBy": "INVERTER" } }`,
	"/site/1/power": `{ "power": { "timeUnit": "QUARTER_OF_AN_HOUR", "unit": "W", "measuredBy": "INVERTER", "values": [
			{ "date": "2021-05-18 00:00:00", "value": 12.0 },
      		{ "date": "2021-05-18 00:15:00", "value": 24.0 },
      		{ "date": "2021-05-18 00:15:00", "value": null } ] } }`,
}
