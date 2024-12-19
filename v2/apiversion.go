package v2

import (
	"context"
)

// This file implements the "API Versions" section of the SolarEdge API specifications.
// https://knowledge-center.solaredge.com/sites/kc/files/se_monitoring_api.pdf

// GetCurrentAPIVersion returns the current API version used by the SolarEdge API server.
func (c *Client) GetCurrentAPIVersion(ctx context.Context) (GetCurrentAPIVersionResponse, error) {
	return call[GetCurrentAPIVersionResponse](ctx, c, "/version/current", nil)
}

type GetCurrentAPIVersionResponse struct {
	Version APIRelease `json:"version"`
}

type APIRelease struct {
	Release string `json:"release"`
}

// GetSupportedAPIVersions returns all API versions supported by the SolarEdge API server.
func (c *Client) GetSupportedAPIVersions(ctx context.Context) (GetSupportedAPIVersionsResponse, error) {
	return call[GetSupportedAPIVersionsResponse](ctx, c, "/version/supported", nil)
}

type GetSupportedAPIVersionsResponse struct {
	Supported []APIRelease `json:"supported"`
}
