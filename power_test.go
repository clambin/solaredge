package solaredge_test

import (
	"github.com/clambin/gotools/httpstub"
	"github.com/clambin/solaredge-exporter/pkg/solaredge"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient_GetPower(t *testing.T) {
	server := &Server{}
	client := solaredge.NewClient("TESTTOKEN", httpstub.NewTestClient(server.serve))
	siteIDs, err := client.GetSiteIDs()
	assert.NoError(t, err)
	if assert.Len(t, siteIDs, 1) {
		assert.Equal(t, 1, siteIDs[0])
	}

	entries, err := client.GetPower(siteIDs[0], time.Now().Add(-1*time.Hour), time.Now())

	assert.NoError(t, err)
	if assert.Len(t, entries, 2) {
		assert.Equal(t, 12.0, entries[0].Value)
		assert.Equal(t, 24.0, entries[1].Value)
	}
}

func TestClient_GetPowerOverview(t *testing.T) {
	server := &Server{}
	client := solaredge.NewClient("TESTTOKEN", httpstub.NewTestClient(server.serve))
	siteIDs, err := client.GetSiteIDs()
	assert.NoError(t, err)
	if assert.Len(t, siteIDs, 1) {
		assert.Equal(t, 1, siteIDs[0])
	}

	var lifeTime, lastYear, lastMonth, lastDay, current float64
	lifeTime, lastYear, lastMonth, lastDay, current, err = client.GetPowerOverview(siteIDs[0])

	if assert.NoError(t, err) {
		assert.Equal(t, 10000.0, lifeTime)
		assert.Equal(t, 1000.0, lastYear)
		assert.Equal(t, 100.0, lastMonth)
		assert.Equal(t, 10.0, lastDay)
		assert.Equal(t, 1.0, current)
	}
}
