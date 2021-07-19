package solaredge

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Token      string
	HTTPClient *http.Client
	APIURL     string
}

type API interface {
	GetSiteIDs(ctx context.Context) (ids []int, err error)
	GetPower(ctx context.Context, id int, from time.Time, to time.Time) (measurements []PowerMeasurement, err error)
	GetPowerOverview(ctx context.Context, id int) (lifeTime, lastYear, lastMonth, lastDay, current float64, err error)
}

// NewClient creates a new API client
//
// Deprecated: this adds little value and will be removed
func NewClient(token string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &Client{
		Token:      token,
		HTTPClient: httpClient,
	}
}

const (
	apiURL = "https://monitoringapi.solaredge.com"
)

func (client *Client) getURL() (response string) {
	response = apiURL
	if client.APIURL != "" {
		response = client.APIURL
	}
	return
}

func (client *Client) call(ctx context.Context, endpoint string, args url.Values, response interface{}) (err error) {
	args.Add("api_key", client.Token)

	fullURL := client.getURL() + endpoint + "?" + args.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)

	var resp *http.Response
	resp, err = client.HTTPClient.Do(req)

	if err == nil {
		if resp.StatusCode == 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(body, response)
		} else {
			err = errors.New(resp.Status)
		}
		_ = resp.Body.Close()
	}

	return
}
