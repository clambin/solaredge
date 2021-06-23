package solaredge_test

import (
	"github.com/clambin/gotools/httpstub"
	"github.com/clambin/solaredge-exporter/pkg/solaredge"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_GetSiteIDs(t *testing.T) {
	server := &Server{}
	client := solaredge.NewClient("", httpstub.NewTestClient(server.serve))

	_, err := client.GetSiteIDs()

	if assert.Error(t, err) {
		assert.Equal(t, "403 Forbidden", err.Error())
	}

	client = solaredge.NewClient("TESTTOKEN", httpstub.NewTestClient(server.serve))

	var siteIDs []int
	siteIDs, err = client.GetSiteIDs()

	assert.NoError(t, err)
	if assert.Len(t, siteIDs, 1) {
		assert.Equal(t, 1, siteIDs[0])
	}
}
