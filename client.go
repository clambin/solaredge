package solaredge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client structure to interact with the SolarEdge API
type Client struct {
	// Token contains the SolarEdge authentication token
	Token string
	// HTTPClient specifies the http Client to use. Defaults to http.DefaultClient
	HTTPClient *http.Client
	// APIURL specifies the target URL. Only used to override the target during unit testing
	APIURL string
}

// API interface exposes the supported API calls
//
//go:generate mockery --name API
type API interface {
	GetSiteIDs(ctx context.Context) (ids []int, err error)
	GetPower(ctx context.Context, id int, from time.Time, to time.Time) (measurements []PowerMeasurement, err error)
	GetPowerOverview(ctx context.Context, id int) (lifeTime, lastYear, lastMonth, lastDay, current float64, err error)
}

const (
	apiURL = "https://monitoringapi.solaredge.com"
)

func (client *Client) call(ctx context.Context, endpoint string, args url.Values, response any) error {
	args.Add("api_key", client.Token)

	fullURL, err := client.buildURL(endpoint, args)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create client request: %w", err)
	}

	httpClient := client.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
		}
	}

	if err = json.Unmarshal(body, response); err != nil {
		err = &ParseError{
			Body: string(body),
			Err:  err,
		}
	}

	return err
}

func (client *Client) buildURL(endpoint string, args url.Values) (string, error) {
	target := client.APIURL
	if target == "" {
		target = apiURL
	}
	fullURL, err := url.Parse(target)
	if err != nil {
		return "", err
	}
	fullURL.Path = endpoint
	q := fullURL.Query()
	for key, vals := range args {
		for _, val := range vals {
			q.Add(key, val)
		}
	}
	fullURL.RawQuery = q.Encode()
	return fullURL.String(), nil
}
