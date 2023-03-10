package solaredge

import (
	"context"
	"net/url"
)

// This file contains APIs from the Site Data API section of the SolarEdge API specifications.
// https://knowledge-center.solaredge.com/sites/kc/files/se_monitoring_api.pdf

// GetSites returns all sites registered under the account associated with the supplied API Token.
func (c *Client) GetSites(ctx context.Context) (Sites, error) {
	var response struct {
		Sites struct {
			Count int    `json:"count"`
			Site  []Site `json:"site"`
		} `json:"sites"`
	}

	err := c.call(ctx, "/sites/list", url.Values{}, &response)
	if err == nil {
		for index := range response.Sites.Site {
			response.Sites.Site[index].client = c
		}
	}
	return response.Sites.Site, err
}

// Sites is a list of Site items.
type Sites []Site

// FindByID returns the Site with the specified ID. Returns false if the Site could not be found.
func (s Sites) FindByID(id int) (Site, bool) {
	for _, site := range s {
		if site.ID == id {
			return site, true
		}
	}
	return Site{}, false
}

// FindByName returns the Site with the specified name. Returns false if the Site could not be found.
func (s Sites) FindByName(name string) (Site, bool) {
	for _, site := range s {
		if site.Name == name {
			return site, true
		}
	}
	return Site{}, false
}
