package solaredge

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Token      string
	HTTPClient *http.Client
}

type API interface {
	GetSiteIDs(ctx context.Context) (ids []int, err error)
	GetPower(ctx context.Context, id int, from time.Time, to time.Time) (measurements []PowerMeasurement, err error)
	GetPowerOverview(ctx context.Context, id int) (lifeTime, lastYear, lastMonth, lastDay, current float64, err error)
}

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

func (client *Client) call(ctx context.Context, endpoint string, args url.Values, response interface{}) (err error) {
	args.Add("api_key", client.Token)

	fullURL := apiURL + endpoint + "?" + args.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	var resp *http.Response

	if resp, err = client.HTTPClient.Do(req); err == nil {
		defer func(body io.ReadCloser) {
			_ = body.Close()
		}(resp.Body)

		if resp.StatusCode == 200 {
			body, _ := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(body, response)
		} else {
			err = errors.New(resp.Status)
		}
	}
	return
}
