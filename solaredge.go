/*
Package solaredge provides a client library for the SolarEdge Cloud-Based Monitoring Platform. The API gives access
to data saved in the monitoring servers for your installed SolarEdge equipment and its performance (i.e. generated power & energy).

The implementation is based on SolarEdge's official [API documentation].

The current version of this library implements the following sections of the API:

  - Site Data API
  - Site Equipment API
  - API Versions

Access to SolarEdge data is determined by the user's API Key & installation. If your situation gives you access
to the Accounts List, Meters or Sensors API, feel free to get in touch to get these implemented in this library.

[API documentation]: https://knowledge-center.solaredge.com/sites/kc/files/se_monitoring_api.pdf
*/
package solaredge

import (
	"cmp"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	SiteKey    string
	HTTPClient *http.Client
	baseURL    string
}

const apiURL = "https://monitoringapi.solaredge.com"

func (c *Client) buildRequest(ctx context.Context, endpoint string, args url.Values) (*http.Request, error) {
	if args == nil {
		args = make(url.Values)
	}
	args.Add("api_key", c.SiteKey)
	args.Add("version", "1.0.0")

	fullURL := cmp.Or(c.baseURL, apiURL) + endpoint + "?" + args.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err == nil {
		req.Header.Set("Accept", "application/json")
	}
	return req, err
}
func call[T any](ctx context.Context, c *Client, path string, args url.Values) (T, error) {
	var response T
	req, err := c.buildRequest(ctx, path, args)
	if err != nil {
		return response, err
	}
	httpClient := cmp.Or(c.HTTPClient, http.DefaultClient)
	resp, err := httpClient.Do(req)
	if err != nil {
		return response, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return response, newResponseError(resp)
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func makePath(path string, siteId int) string {
	return strings.ReplaceAll(path, "{siteId}", strconv.Itoa(siteId))
}
