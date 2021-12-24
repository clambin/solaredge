package solaredge

import (
	"context"
	"net/url"
)

// GetSiteIDs returns all site IDs registered under the account associated with the supplied client token
func (client *Client) GetSiteIDs(ctx context.Context) (sites []int, err error) {
	var sitesResponse struct {
		Sites struct {
			Count int
			Site  []struct {
				ID int
			}
		}
	}

	args := url.Values{}
	err = client.call(ctx, "/sites/list", args, &sitesResponse)

	if err == nil {
		for _, site := range sitesResponse.Sites.Site {
			sites = append(sites, site.ID)
		}
	}

	return
}
