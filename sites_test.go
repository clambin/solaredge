package solaredge_test

import (
	"context"
	"github.com/clambin/gotools/httpstub"
	"github.com/clambin/solaredge"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestClient_GetSiteIDs(t *testing.T) {
	server := &Server{}
	client := solaredge.NewClient("", httpstub.NewTestClient(server.serve))

	_, err := client.GetSiteIDs(context.Background())

	if assert.Error(t, err) {
		assert.Equal(t, "403 Forbidden", err.Error())
	}

	client = solaredge.NewClient("TESTTOKEN", httpstub.NewTestClient(server.serve))

	var siteIDs []int
	siteIDs, err = client.GetSiteIDs(context.Background())

	assert.NoError(t, err)
	if assert.Len(t, siteIDs, 1) {
		assert.Equal(t, 1, siteIDs[0])
	}
}

func TestClient_Timeout(t *testing.T) {
	server := &Server{}
	client := solaredge.NewClient("", httpstub.NewTestClient(server.slowserve))

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := client.GetSiteIDs(ctx)

	assert.Error(t, err)
}
