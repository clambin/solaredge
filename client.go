package solaredge

import (
	"context"
	"encoding/json"
	"fmt"
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

//go:generate mockery --name API
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
		httpClient = http.DefaultClient
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

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create client request: %w", err)
	}

	var resp *http.Response
	resp, err = client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call server: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to call server: %s", resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, response)

	if err != nil {
		return &InvalidServerResponse{
			Body: string(body),
			Err:  err,
		}
	}

	return
}

var _ error = InvalidServerResponse{}

type InvalidServerResponse struct {
	Body string
	Err  error
}

func (e InvalidServerResponse) Error() string {
	return "invalid server response: " + e.Err.Error()
}

func (e InvalidServerResponse) Is(e2 error) bool {
	_, ok := e2.(*InvalidServerResponse)
	return ok
}

func (e InvalidServerResponse) Unwrap() error {
	return e.Err
}
