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
// to the Accounts List, Meters or Sensors API, feel free to get in touch to get these implemented in this library.
//
// Deprecated: v1 of this package is now replaced by github.com/clambin/solaredge/v2
//
// [API documentation]: https://knowledge-center.solaredge.com/sites/kc/files/se_monitoring_api.pdf
package solaredge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

// Client to interact with the SolarEdge API
//
// Deprecated: v1 of this package is now replaced by github.com/clambin/solaredge/v2
type Client struct {
	// Token contains the SolarEdge authentication token
	Token string
	// HTTPClient specifies the http client to use. Defaults to http.DefaultClient
	HTTPClient *http.Client
}

const (
	apiURL = "https://monitoringapi.solaredge.com"
)

func (c *Client) call(ctx context.Context, endpoint string, args url.Values, response any) error {
	req, err := c.buildRequest(ctx, endpoint, args)
	if err != nil {
		return fmt.Errorf("failed to create client request: %w", err)
	}

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		// hide the API key from the error
		var urlError *url.Error
		if errors.As(err, &urlError) {
			urlError.URL = hideAPIKey(urlError.URL)
		}
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
				Body: body,
			}
		}
	case http.StatusForbidden:
		err = makeAPIError(body)
	default:
		err = &HTTPError{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Body:       body,
		}
	}

	return err
}

func (c *Client) buildRequest(ctx context.Context, endpoint string, args url.Values) (*http.Request, error) {
	if endpoint == "" {
		endpoint = "/"
	}
	args.Add("api_key", c.Token)
	args.Add("version", "1.0.0")
	fullURL := apiURL + endpoint + "?" + args.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err == nil {
		req.Header.Set("Accept", "application/json")
	}
	return req, err
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

func hideAPIKey(input string) string {
	re := regexp.MustCompile(`api_key=(?P<token>\w+)(&|$)`)
	return re.ReplaceAllString(input, "api_key=<REDACTED>$2")
}
