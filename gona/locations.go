package gona

import (
	"context"

	"github.com/google/go-querystring/query"
)

// Location is an API response message of available deployment locations
type Location struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IATACode  string `json:"iata_code"`
	Continent string `json:"continent"`
	Flag      string `json:"flag"`
	Disabled  int    `json:"disabled"`
}

// GetLocations public method on Client to get a list of locations
func (c *Client) GetLocations(ctx context.Context) ([]Location, error) {
	r := make([]Location, 0)
	if err := c.get(ctx, "cloud/locations", &r); err != nil {
		return nil, err
	}
	return r, nil
}

// GetLocationsForPool specifies locations that are valid for a given CloudPool
func (c *Client) GetLocationForPool(ctx context.Context, pool CloudPool) ([]Location, error) {
	values, err := query.Values(struct {
		CloudPoolID CloudPool `url:"cloud_pool_id"`
	}{
		CloudPoolID: pool,
	})
	if err != nil {
		return nil, err
	}
	var r []Location
	if err := c.get(ctx, "cloud/locations?"+values.Encode(), &r); err != nil {
		return nil, err
	}
	return r, nil
}
