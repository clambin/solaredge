package solaredge_test

import (
	"context"
	"github.com/clambin/solaredge"
	"github.com/clambin/solaredge/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestClient_Site_E2E(t *testing.T) {
	apikey := os.Getenv("SOLAREDGE_APIKEY")
	if apikey == "" {
		t.Skip("SOLAREDGE_APIKEY not set. Skipping")
	}

	c := solaredge.Client{Token: apikey}
	ctx := context.Background()

	sites, err := c.GetSites(ctx)
	require.NoError(t, err)

	for _, site := range sites {
		start, end, err := site.GetDataPeriod(ctx)
		require.NoError(t, err)
		assert.False(t, start.IsZero())
		assert.False(t, end.IsZero())

		energy, err := site.GetEnergy(ctx, "YEAR", start, end)
		require.NoError(t, err)
		assert.NotEmpty(t, energy.Values)

		siteEnergy, err := site.GetTimeFrameEnergy(ctx, end.Add(-365*24*time.Hour), end)
		require.NoError(t, err)
		assert.NotZero(t, siteEnergy.Energy)

		sitePower, err := site.GetPower(ctx, end.Add(-7*24*time.Hour), end)
		require.NoError(t, err)
		assert.Equal(t, "QUARTER_OF_AN_HOUR", sitePower.TimeUnit)
		assert.Equal(t, "W", sitePower.Unit)
		assert.NotEmpty(t, sitePower.Values)

		powerOverview, err := site.GetPowerOverview(ctx)
		require.NoError(t, err)
		assert.NotZero(t, powerOverview.LifeTimeData)
		assert.False(t, time.Time(powerOverview.LastUpdateTime).IsZero())

		powerDetails, err := site.GetPowerDetails(ctx, end.Add(7*24*time.Hour), end)
		require.NoError(t, err)
		assert.NotEmpty(t, powerDetails.Meters)
	}
}

func TestClient_GetSites(t *testing.T) {
	c := solaredge.Client{
		Token:      "1234",
		HTTPClient: &http.Client{Transport: &testutil.Server{Token: "1234"}},
	}

	sites, err := c.GetSites(context.Background())

	require.NoError(t, err)
	require.Len(t, sites, 1)
	assert.Equal(t, 1, sites[0].ID)
	assert.Equal(t, "site 1", sites[0].Name)
}

func TestSites_FindByID(t *testing.T) {
	tests := []struct {
		name  string
		s     solaredge.Sites
		id    int
		want  solaredge.Site
		want1 bool
	}{
		{
			name:  "match",
			s:     solaredge.Sites{{ID: 1, Name: "foo"}, {ID: 2, Name: "bar"}},
			id:    1,
			want:  solaredge.Site{ID: 1, Name: "foo"},
			want1: true,
		},
		{
			name:  "no match",
			s:     solaredge.Sites{{ID: 1, Name: "foo"}, {ID: 2, Name: "bar"}},
			id:    3,
			want:  solaredge.Site{},
			want1: false,
		},
		{
			name:  "empty",
			s:     solaredge.Sites{},
			id:    1,
			want:  solaredge.Site{},
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
		s     solaredge.Sites
		desc  string
		want  solaredge.Site
		want1 bool
	}{
		{
			name:  "match",
			s:     solaredge.Sites{{ID: 1, Name: "foo"}, {ID: 2, Name: "bar"}},
			desc:  "foo",
			want:  solaredge.Site{ID: 1, Name: "foo"},
			want1: true,
		},
		{
			name:  "no match",
			s:     solaredge.Sites{{ID: 1, Name: "foo"}, {ID: 2, Name: "bar"}},
			desc:  "snafu",
			want:  solaredge.Site{},
			want1: false,
		},
		{
			name:  "empty",
			s:     solaredge.Sites{},
			desc:  "foo",
			want:  solaredge.Site{},
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
