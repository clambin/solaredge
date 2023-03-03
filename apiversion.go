package solaredge

import (
	"context"
	"net/url"
)

// This file contains all supported APIs from the API Versions section of the SolarEdge API specifications.
// https://knowledge-center.solaredge.com/sites/kc/files/se_monitoring_api.pdf

// GetCurrentAPIVersion returns the current API version used by the SolarEdge API server.
func (c *Client) GetCurrentAPIVersion(ctx context.Context) (string, error) {
	var current struct {
		Version struct {
			Release string `json:"release"`
		} `json:"version"`
	}
	err := c.call(ctx, "/version/current", url.Values{}, &current)
	return current.Version.Release, err
}

// GetSupportedAPIVersions returns all API versions supported by the SolarEdge API server.
func (c *Client) GetSupportedAPIVersions(ctx context.Context) ([]string, error) {
	var supported struct {
		Supported []struct {
			Release string `json:"release"`
		} `json:"supported"`
	}

	var apiVersions []string
	err := c.call(ctx, "/version/supported", url.Values{}, &supported)
	if err == nil {
		for _, version := range supported.Supported {
			apiVersions = append(apiVersions, version.Release)
		}
	}

	return apiVersions, err
}
