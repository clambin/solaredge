package solaredge

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		args     url.Values
		want     string
	}{
		{
			name:     "no siteID",
			endpoint: "/version/current",
			args:     url.Values{},
			want:     "https://monitoringapi.solaredge.com/version/current?api_key=123&version=1.0.0",
		},
		{
			name:     "with siteID",
			endpoint: "/site/1/power",
			args:     make(url.Values),
			want:     "https://monitoringapi.solaredge.com/site/1/power?api_key=123&version=1.0.0",
		},
		{
			name:     "with args",
			endpoint: "/site/1/power",
			args:     url.Values{"foo": []string{"bar"}},
			want:     "https://monitoringapi.solaredge.com/site/1/power?api_key=123&foo=bar&version=1.0.0",
		},
		{
			name:     "no endpoint",
			endpoint: "",
			args:     url.Values{},
			want:     "https://monitoringapi.solaredge.com/?api_key=123&version=1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{Token: "123"}
			req, _ := c.buildRequest(context.TODO(), tt.endpoint, tt.args)

			assert.Equal(t, tt.want, req.URL.String())
		})
	}
}

func Test_buildArgsFromTimeRange(t *testing.T) {
	type args struct {
		start  time.Time
		end    time.Time
		label  string
		layout string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "valid Time",
			args: args{
				start:  time.Date(2023, time.February, 1, 12, 0, 0, 0, time.UTC),
				end:    time.Date(2023, time.February, 26, 23, 0, 0, 0, time.UTC),
				label:  "Time",
				layout: "2006-01-02 15:04:05",
			},
			want:    "endTime=2023-02-26+23%3A00%3A00&startTime=2023-02-01+12%3A00%3A00",
			wantErr: assert.NoError,
		},
		{
			name: "valid Date",
			args: args{
				start:  time.Date(2023, time.February, 1, 12, 0, 0, 0, time.UTC),
				end:    time.Date(2023, time.February, 26, 23, 0, 0, 0, time.UTC),
				label:  "Date",
				layout: "2006-01-02",
			},
			want:    "endDate=2023-02-26&startDate=2023-02-01",
			wantErr: assert.NoError,
		},
		{
			name: "no end date",
			args: args{
				start:  time.Date(2023, time.February, 1, 12, 0, 0, 0, time.UTC),
				label:  "Date",
				layout: "2006-01-02",
			},
			wantErr: assert.Error,
		},
		{
			name: "no start date",
			args: args{
				end:    time.Date(2023, time.February, 26, 12, 0, 0, 0, time.UTC),
				label:  "Date",
				layout: "2006-01-02",
			},
			wantErr: assert.Error,
		},
		{
			name: "no dates",
			args: args{
				label:  "Date",
				layout: "2006-01-02",
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildArgsFromTimeRange(tt.args.start, tt.args.end, tt.args.label, tt.args.layout)
			assert.Equal(t, tt.want, got.Encode())
			tt.wantErr(t, err)
		})
	}
}

func Test_hideAPIKey(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "valid",
			args: args{input: "https://example.com/list/sites?api_key=SECRET&api_version=1.0.0"},
			want: "https://example.com/list/sites?api_key=<REDACTED>&api_version=1.0.0",
		},
		{
			name: "last",
			args: args{input: "https://example.com/list/sites?api_key=SECRET"},
			want: "https://example.com/list/sites?api_key=<REDACTED>",
		},
		{
			name: "no token",
			args: args{input: "https://example.com/list/sites?api_version=1.0.0"},
			want: "https://example.com/list/sites?api_version=1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, hideAPIKey(tt.args.input))
		})
	}
}
