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

func TestSites_FindByID(t *testing.T) {
	tests := []struct {
		name  string
		s     v2.Sites
		id    int
		want  v2.SiteDetails
		want1 bool
	}{
		{
			name:  "match",
			s:     v2.Sites{{Id: 1, Name: "foo"}, {Id: 2, Name: "bar"}},
			id:    1,
			want:  v2.SiteDetails{Id: 1, Name: "foo"},
			want1: true,
		},
		{
			name:  "no match",
			s:     v2.Sites{{Id: 1, Name: "foo"}, {Id: 2, Name: "bar"}},
			id:    3,
			want:  v2.SiteDetails{},
			want1: false,
		},
		{
			name:  "empty",
			s:     v2.Sites{},
			id:    1,
			want:  v2.SiteDetails{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.FindByID(tt.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestSites_FindByName(t *testing.T) {
	tests := []struct {
		name  string
		s     v2.Sites
		desc  string
		want  v2.SiteDetails
		want1 bool
	}{
		{
			name:  "match",
			s:     v2.Sites{{Id: 1, Name: "foo"}, {Id: 2, Name: "bar"}},
			desc:  "foo",
			want:  v2.SiteDetails{Id: 1, Name: "foo"},
			want1: true,
		},
		{
			name:  "no match",
			s:     v2.Sites{{Id: 1, Name: "foo"}, {Id: 2, Name: "bar"}},
			desc:  "snafu",
			want:  v2.SiteDetails{},
			want1: false,
		},
		{
			name:  "empty",
			s:     v2.Sites{},
			desc:  "foo",
			want:  v2.SiteDetails{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.s.FindByName(tt.desc)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestClient_GetSites(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSites(ctx)
	require.NoError(t, err)
	assert.Equal(t, responses["/sites/list"], resp)
}

func TestClient_GetSiteDetails(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSiteDetails(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, responses["/site/1/details"], resp)
}

func TestClient_GetSiteDataPeriod(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSiteDataPeriod(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, responses["/site/1/dataPeriod"], resp)
}

func TestClient_GetSiteEnergy(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSiteEnergy(ctx, 1, v2.TimeUnitDay, time.Time{}, time.Time{})
	require.NoError(t, err)
	assert.Equal(t, responses["/site/1/energy"], resp)
}

func TestClient_GetSiteEnergyForTimeFrame(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSiteEnergyForTimeFrame(ctx, 1, time.Time{}, time.Time{})
	require.NoError(t, err)
	assert.Equal(t, responses["/site/1/timeFrameEnergy"], resp)
}

func TestClient_GetSitePower(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSitePower(ctx, 1, time.Time{}, time.Time{})
	require.NoError(t, err)
	assert.Equal(t, responses["/site/1/power"], resp)
}

func TestClient_GetSitePowerOverview(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSitePowerOverview(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, responses["/site/1/overview"], resp)
}

func TestClient_GetSitePowerDetails(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSitePowerDetails(ctx, 1, time.Time{}, time.Time{})
	require.NoError(t, err)
	assert.Equal(t, responses["/site/1/powerDetails"], resp)
}

func TestClient_GetSiteEnergyDetails(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSiteEnergyDetails(ctx, 1, v2.TimeUnitQuarter, time.Time{}, time.Time{})
	require.NoError(t, err)
	assert.Equal(t, responses["/site/1/energyDetails"], resp)
}

func TestClient_GetSitePowerFlow(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSitePowerFlow(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, responses["/site/1/currentPowerFlow"], resp)
}

func TestClient_GetSiteStorageData(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSiteStorageData(ctx, 1, time.Time{}, time.Time{})
	require.NoError(t, err)
	assert.Equal(t, responses["/site/1/storageData"], resp)
}

func TestClient_GetSiteEnvBenefits(t *testing.T) {
	c := v2.Client{Target: testServer.URL, HTTPClient: http.DefaultClient}
	ctx := context.Background()
	resp, err := c.GetSiteEnvBenefits(ctx, 1)
	require.NoError(t, err)
	assert.Equal(t, responses["/site/1/envBenefits"], resp)
}
