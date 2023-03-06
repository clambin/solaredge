package solaredge

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestBuildURL(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/sites/list":
			_, _ = w.Write([]byte(`{ "sites": { "Count": 1, "site": [ { "id": 1 }] } }`))
		default:
			http.Error(w, "", http.StatusNotFound)
		}
	}))
	defer s.Close()

	tests := []struct {
		name     string
		endpoint string
		args     url.Values
		want     string
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name:     "no siteID",
			endpoint: "/version/current",
			args:     url.Values{},
			want:     s.URL + "/version/current?api_key=123&version=1.0.0",
			wantErr:  assert.NoError,
		},
		{
			name:     "with siteID",
			endpoint: "/site/%d/power",
			args:     make(url.Values),
			want:     s.URL + "/site/1/power?api_key=123&version=1.0.0",
			wantErr:  assert.NoError,
		},
		{
			name:     "with args",
			endpoint: "/site/%d/power",
			args:     url.Values{"foo": []string{"bar"}},
			want:     s.URL + "/site/1/power?api_key=123&foo=bar&version=1.0.0",
			wantErr:  assert.NoError,
		},
		{
			name:     "no endpoint",
			endpoint: "",
			args:     url.Values{},
			want:     s.URL + "/?api_key=123&version=1.0.0",
			wantErr:  assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{Token: "123", apiURL: s.URL}
			result, err := c.buildURL(context.TODO(), tt.endpoint, tt.args)

			assert.Equal(t, tt.want, result)
			tt.wantErr(t, err)
		})
	}
}

func TestBuildURL_Init_Error(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	s.Close()

	c := Client{Token: "123", apiURL: s.URL}
	_, err := c.buildURL(context.TODO(), "/site/%d/power", url.Values{})
	assert.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "init: "))
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

func TestClient_Authentication(t *testing.T) {
	c, s, _ := makeTestServer(nil)
	defer s.Close()

	goodToken := c.Token
	c.Token = "BADTOKEN"
	c.apiURL = s.URL

	_, err := c.GetSites(context.Background())

	require.Error(t, err)
	require.ErrorIs(t, err, &APIError{})

	c.Token = goodToken
	_, err = c.GetSites(context.Background())

	assert.NoError(t, err)
}

func TestClient_Timeout(t *testing.T) {
	c, s, h := makeTestServer(nil)
	defer s.Close()

	c.HTTPClient = &http.Client{Timeout: 100 * time.Millisecond}
	h.slow = true

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// this should finish after 100 ms (http.Client timeout)
	start := time.Now()
	_, err := c.GetSites(ctx)
	require.Error(t, err)
	assert.Less(t, time.Since(start), 400*time.Millisecond)
}

func TestClient_Fail(t *testing.T) {
	c, s, h := makeTestServer(nil)
	defer s.Close()

	h.fail = true
	ctx := context.Background()

	_, err := c.GetSites(ctx)
	require.Error(t, err)
	assert.ErrorIs(t, err, &HTTPError{})
	assert.Equal(t, "500 Internal Server Error", err.Error())
}

func TestClient_Errors(t *testing.T) {
	c, s, h := makeTestServer(nil)

	goodToken := c.Token
	c.Token = "BADTOKEN"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := c.GetSites(ctx)
	require.Error(t, err)
	assert.ErrorIs(t, err, &APIError{})

	c.Token = goodToken
	h.garbage = true

	_, err = c.GetSites(ctx)
	require.Error(t, err)
	assert.Equal(t, `json parse error: invalid character '=' after object key`, err.Error())
	assert.ErrorIs(t, err, &ParseError{})

	err = errors.Unwrap(err)
	require.Error(t, err)
	var err3 *json.SyntaxError
	require.ErrorAs(t, err, &err3)
	assert.Equal(t, int64(8), err3.Offset)

	s.Close()
	_, err = c.GetSites(ctx)
	require.Error(t, err)
	//t.Log(err)

	c.apiURL = "invalid url"
	_, err = c.GetSites(ctx)
	require.Error(t, err)
	assert.Equal(t, `Get "/sites/list?api_key=<REDACTED>&version=1.0.0": unsupported protocol scheme ""`, err.Error())
	var err4 *url.Error
	require.ErrorAs(t, err, &err4)
	assert.Equal(t, "/sites/list?api_key=<REDACTED>&version=1.0.0", err4.URL)
}

func TestClient_SetActiveSiteID(t *testing.T) {
	var c Client
	assert.Zero(t, c.getActiveSiteID())
	c.SetActiveSiteID(1)
	assert.Equal(t, 1, c.getActiveSiteID())
}

func makeTestServer(response any) (*Client, *httptest.Server, *Server) {
	h := Server{token: "TESTTOKEN", response: response}
	s := httptest.NewServer(&h)
	c := &Client{
		Token:      "TESTTOKEN",
		HTTPClient: http.DefaultClient,
		apiURL:     s.URL,
	}
	return c, s, &h
}

type Server struct {
	fail     bool
	token    string
	slow     bool
	garbage  bool
	response any
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.fail {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if s.slow && wait(r.Context(), 5*time.Second) == false {
		http.Error(w, "", http.StatusRequestTimeout)
		return
	}

	if s.authenticate(r) == false {
		http.Error(w, "", http.StatusForbidden)
		return
	}

	if s.garbage {
		_, _ = w.Write([]byte(`[{"foo"="bar"}]`))
		return
	}

	var response any

	switch r.URL.Path {
	case "/sites/list":
		response = struct {
			Sites struct {
				Count int    `json:"count"`
				Site  []Site `json:"site"`
			} `json:"sites"`
		}{
			Sites: struct {
				Count int    `json:"count"`
				Site  []Site `json:"site"`
			}{

				Count: 1,
				Site:  []Site{{ID: 1, Name: "home"}},
			},
		}
	default:
		response = s.response
	}
	_ = json.NewEncoder(w).Encode(response)
}

func wait(ctx context.Context, duration time.Duration) (passed bool) {
	timer := time.NewTimer(duration)
	for {
		select {
		case <-timer.C:
			return true
		case <-ctx.Done():
			return false
		}
	}
}

func (s *Server) authenticate(req *http.Request) bool {
	values := req.URL.Query()
	value, ok := values["api_key"]

	return ok && len(value) > 0 && value[0] == s.token
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
