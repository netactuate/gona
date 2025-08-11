package gona

import (
// 	"net/url"
// 	"strconv"
//   "log"
     "fmt"
//	"github.com/google/go-querystring/query"
)

type Datacenter struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    IATA string `json:"iata"`
}

func (c *Client) GetDatacenterByIATA(iata string) (int, error) {
	var resp Datacenter
	if err := c.get(fmt.Sprintf("platform/datacenters-by-iata/%s", iata), &resp); err != nil {
		return 0, err
	}
	return resp.ID, nil
}