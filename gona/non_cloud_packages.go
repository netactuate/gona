package gona

import (
	"context"
	"sort"
)

type NonCloudPackageDetails struct {
	DCName       string `json:"dc_name"`
	IATACode     string `json:"iata_code"`
	BWCommit     string `json:"bw_commit"`
	OverageType  string `json:"overage_type"`
	OverageRate  string `json:"overage_rate"`
	AggBWMbPkgID string `json:"agg_bw_mbpkgid"`
}

type ColocationPackage struct {
	MBPkgID       int                    `json:"mbpkgid"`
	PackageStatus string                 `json:"package_status"`
	FQDN          string                 `json:"fqdn"`
	BillingCycle  string                 `json:"billingcycle"`
	NextDueDate   string                 `json:"nextduedate"`
	Amount        string                 `json:"amount"`
	Details       NonCloudPackageDetails `json:"details"`
	Status        string                 `json:"status"`
}

type TransitPackage struct {
	MBPkgID       int                    `json:"mbpkgid"`
	PackageStatus string                 `json:"package_status"`
	FQDN          string                 `json:"fqdn"`
	BillingCycle  string                 `json:"billingcycle"`
	NextDueDate   string                 `json:"nextduedate"`
	Amount        string                 `json:"amount"`
	Details       NonCloudPackageDetails `json:"details"`
	Status        string                 `json:"status"`
}

func (c *Client) GetColocationPackages() ([]ColocationPackage, error) {
	var byID map[string]ColocationPackage
	if err := c.get(context.Background(), "colo/packages", &byID); err != nil {
		return nil, err
	}

	packages := make([]ColocationPackage, 0, len(byID))
	for _, pkg := range byID {
		packages = append(packages, pkg)
	}
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].MBPkgID < packages[j].MBPkgID
	})

	return packages, nil
}

func (c *Client) GetTransitPackages() ([]TransitPackage, error) {
	var byID map[string]TransitPackage
	if err := c.get(context.Background(), "transit/packages", &byID); err != nil {
		return nil, err
	}

	packages := make([]TransitPackage, 0, len(byID))
	for _, pkg := range byID {
		packages = append(packages, pkg)
	}
	sort.Slice(packages, func(i, j int) bool {
		return packages[i].MBPkgID < packages[j].MBPkgID
	})

	return packages, nil
}
