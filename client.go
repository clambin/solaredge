package solaredge

import (
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
	GetSiteIDs() ([]int, error)
	GetPower(int, time.Time, time.Time) ([]PowerMeasurement, error)
	GetPowerOverview(int) (float64, float64, float64, float64, float64, error)
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

func (client *Client) call(endpoint string, args url.Values, response interface{}) (err error) {
	args.Add("api_key", client.Token)

	fullURL := apiURL + endpoint + "?" + args.Encode()

	req, _ := http.NewRequest(http.MethodGet, fullURL, nil)
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
