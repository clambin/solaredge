// Package solaredge is a client library for the SolarEdge Cloud-Based Monitoring Platform. The API gives access
// to data saved in the monitoring servers for your installed SolarEdge equipment and its performance (i.e. generated power & energy).
//
// The implementation is based on SolarEdge's official [API documentation].
//
// The current version of this library implements the following sections of the API:
//
//   - Site Data API
//   - Site Equipment API
//   - API Versions
//
// Access to SolarEdge data is determined by the user's API Key & installation. If your situation gives you access
// to the Accounts List, Meters or Sensors API, feel free to get in touch with me to get these implemented in this library.
//
// [API documentation]: https://knowledge-center.solaredge.com/sites/kc/files/se_monitoring_api.pdf
package solaredge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Client to interact with the SolarEdge API
type Client struct {
	// Token contains the SolarEdge authentication token
	Token string
	// HTTPClient specifies the http client to use. Defaults to http.DefaultClient
	HTTPClient *http.Client

	lock         sync.Mutex
	activeSiteID int
	apiURL       string
}

const (
	apiURL = "https://monitoringapi.solaredge.com"
)

func (c *Client) call(ctx context.Context, endpoint string, args url.Values, response any) error {
	fullURL, err := c.buildURL(ctx, endpoint, args)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create client request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	httpClient := c.HTTPClient
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

	switch resp.StatusCode {
	case http.StatusOK:
		if err = json.Unmarshal(body, response); err != nil {
			err = &ParseError{
				Err:  err,
				Body: string(body),
			}
		}
	case http.StatusForbidden:
		err = makeAPIError(body)
	default:
		err = &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       string(body),
		}
	}

	return err
}

func (c *Client) initialize(ctx context.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.activeSiteID > 0 {
		return nil
	}
	sites, err := c.GetSites(ctx)
	if err != nil {
		return err
	}
	if len(sites) == 0 {
		return fmt.Errorf("no sites found")
	}
	c.activeSiteID = sites[0].ID
	return nil
}

// SetActiveSiteID sets the active site
func (c *Client) SetActiveSiteID(id int) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.activeSiteID = id
}

func (c *Client) getActiveSiteID() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.activeSiteID
}

func (c *Client) buildURL(ctx context.Context, endpoint string, args url.Values) (string, error) {
	target := apiURL
	if c.apiURL != "" {
		target = c.apiURL
	}

	fullURL, err := url.Parse(target)
	if err != nil {
		return "", err
	}

	if endpoint == "" {
		endpoint = "/"
	} else if strings.Contains(endpoint, "%d") {
		if err = c.initialize(ctx); err != nil {
			return "", fmt.Errorf("init: %w", err)
		}
		endpoint = strings.Replace(endpoint, "%d", strconv.Itoa(c.getActiveSiteID()), 1)
	}

	fullURL.Path = endpoint
	args.Add("api_key", c.Token)
	args.Add("version", "1.0.0")
	fullURL.RawQuery = args.Encode()
	return fullURL.String(), nil
}

func buildArgsFromTimeRange(start, end time.Time, label, layout string) (url.Values, error) {
	if start.IsZero() {
		return nil, fmt.Errorf("start cannot be zero")
	}
	if end.IsZero() {
		return nil, fmt.Errorf("end cannot be zero")
	}
	args := make(url.Values)
	args.Set("start"+label, start.Format(layout))
	args.Set("end"+label, end.Format(layout))
	return args, nil
}
