package v2

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
	Target     string
}

const apiURL = "https://monitoringapi.solaredge.com"

func (c *Client) buildRequest(ctx context.Context, endpoint string, args url.Values) (*http.Request, error) {
	if args == nil {
		args = make(url.Values)
	}
	args.Add("api_key", c.SiteKey)
	args.Add("version", "1.0.0")

	fullURL := cmp.Or(c.Target, apiURL) + endpoint + "?" + args.Encode()

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
	resp, err := c.HTTPClient.Do(req)
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
	return strings.Replace(path, "{siteId}", strconv.Itoa(siteId), -1)
}
