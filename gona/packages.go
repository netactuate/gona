package gona

import (
	"context"
	"strconv"
)

// Package struct stores the purchaced package values
type Package struct {
	ID        int    `json:"mbpkgid,string"`
	Status    string `json:"package_status"`
	Locked    string `json:"locked"`
	PlanName  string `json:"name"`
	Installed int    `json:"installed,string"`
}

// GetPackages external method on Client that returns a
// list of Package object from the API
func (c *Client) GetPackages(ctx context.Context) ([]Package, error) {

	var packageList []Package

	if err := c.get(ctx, "cloud/packages", &packageList); err != nil {
		return nil, err
	}

	return packageList, nil
}

// GetPackage external method on Client that takes an id (int) as it's sole
// argument and returns a single Package object
func (c *Client) GetPackage(ctx context.Context, id int) (pkg Package, err error) {
	if err := c.get(ctx, "cloud/package/"+strconv.Itoa(id), &pkg); err != nil {
		return Package{Installed: 0}, err
	}
	return pkg, nil
}
