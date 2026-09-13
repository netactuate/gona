package gona

import (
	"context"
	"fmt"
)

type BillingPackage struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	DomuLabel    *string `json:"domu_label"`
	PackageID    int     `json:"packageid"`
	Domain       string  `json:"domain"`
	Amount       string  `json:"amount"`
	BillingCycle string  `json:"billingcycle"`
	DomainStatus string  `json:"domainstatus"`
	NextDueDate  string  `json:"nextduedate"`
	DedicatedIP  string  `json:"dedicatedip"`
}

type CloudPool struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	// A CPU model string such as "EPYC-Milan", not a vcpu count, despite the name. Null on
	// pools that do not constrain it. Modelled as *int on 2026-09-11 from a single sampled row
	// where it happened to be null, which failed the moment a second pool was read.
	RequiredVCPU      *string  `json:"required_vcpu"`
	HardCapabilities  []string `json:"hard_capabilities"`
	SoftCapabilities  []string `json:"soft_capabilities"`
	Private           int      `json:"private"`
	BackupCloudPoolID *int     `json:"backup_cloud_pool_id"`
	DefaultRAMPrice   string   `json:"default_ram_price"`
	DefaultCPUPrice   string   `json:"default_cpu_price"`
	DefaultDiskPrice  string   `json:"default_disk_price"`
	LastUpdated       string   `json:"last_updated"`
	Created           string   `json:"created"`
	ContractID        *int     `json:"contract_id"`
}

type CloudCapacity struct {
	PackageID    int     `json:"pkg_id"`
	PackageName  string  `json:"pkg_name"`
	PackageCPU   int     `json:"pkg_cpu"`
	PackageRAM   int     `json:"pkg_ram"`
	PackageDisk  int     `json:"pkg_disk"`
	PackageNet   int     `json:"pkg_net"`
	PackagePort  int     `json:"pkg_port"`
	MonthlyPrice float64 `json:"monthly_price"`
	Available    int     `json:"available"`
}

type DedicatedCapacity struct {
	DeviceID       int    `json:"device_id"`
	LocationID     int    `json:"location_id"`
	LookingGlass   string `json:"looking_glass"`
	MBPkgID        int    `json:"mbpkgid"`
	Name           string `json:"name"`
	NPSEnabled     int    `json:"nps_enabled"`
	PubDescription string `json:"pub_description"`
}

func (c *Client) GetBillingPackages() ([]BillingPackage, error) {
	var packages []BillingPackage
	if err := c.get(context.Background(), "cloud/billing-packages", &packages); err != nil {
		return nil, err
	}
	return packages, nil
}

func (c *Client) GetCloudPools() ([]CloudPool, error) {
	var pools []CloudPool
	if err := c.get(context.Background(), "cloud/pools", &pools); err != nil {
		return nil, err
	}
	return pools, nil
}

func (c *Client) GetCloudCapacity(cloudPoolID, locationID int) ([]CloudCapacity, error) {
	var capacity []CloudCapacity
	path := fmt.Sprintf("cloud/capacity?cloud_pool_id=%d&location_id=%d", cloudPoolID, locationID)
	if err := c.get(context.Background(), path, &capacity); err != nil {
		return nil, err
	}
	return capacity, nil
}

func (c *Client) GetDedicatedCapacity() ([]DedicatedCapacity, error) {
	var capacity []DedicatedCapacity
	if err := c.get(context.Background(), "dedicated/capacity", &capacity); err != nil {
		return nil, err
	}
	return capacity, nil
}
