package solaredge

import (
	"net/url"
)

func (client *Client) GetSiteIDs() (sites []int, err error) {
	var sitesResponse struct {
		Sites struct {
			Count int
			Site  []struct {
				ID int
			}
		}
	}

	args := url.Values{}
	err = client.call("/sites/list", args, &sitesResponse)

	if err == nil {
		for _, site := range sitesResponse.Sites.Site {
			sites = append(sites, site.ID)
		}
	}

	return
}
