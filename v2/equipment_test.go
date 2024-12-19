package v2_test

import (
	"context"
	v2 "github.com/clambin/solaredge/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestClient_GetComponents(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetComponents(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, responses["/equipment/1/list"], resp)
}

func TestClient_GetInventory(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetInventory(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, responses["/site/1/inventory"], resp)
}

func TestClient_GetInverterTechnicalData(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetInverterTechnicalData(ctx, 1, "SN1", time.Time{}, time.Time{})
	require.NoError(t, err)
	assert.Equal(t, responses["/equipment/1/SN1/data"], resp)
}

func TestClient_GetEquipmentChangeLog(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetEquipmentChangeLog(ctx, 1, "SN1")
	require.NoError(t, err)
	assert.Equal(t, responses["/equipment/1/SN1/changeLog"], resp)
}
