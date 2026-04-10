package gona

import (
	"encoding/json"
	"fmt"
	"time"
)

type VPCFloatingIP struct {
	FloatingIPID int    `json:"floatingIpId"`
	Address      string `json:"address"`
	IPVersion    int    `json:"ipVersion"`
	PTR          string `json:"ptr,omitempty"`
	IsPrimary    bool   `json:"isPrimary"`
}

type CreateVPCFloatingIPRequest struct {
	IPVersion int    `json:"ipVersion"`
	PTR       string `json:"ptr,omitempty"`
}

type UpdateVPCFloatingIPRequest struct {
	PTR string `json:"ptr"`
}

type CreateVPCFloatingIPResponse struct {
	FloatingIPID int `json:"floatingIpId"`
}

func (c *V3Client) CreateVPCFloatingIP(vpcID int, req *CreateVPCFloatingIPRequest) (*CreateVPCFloatingIPResponse, error) {
	path := fmt.Sprintf("/vpcs/%d/floating-ips", vpcID)
	const maxAttempts = 7
	const retryDelay = 10 * time.Second

	var resp *V3APIResponse
	var err error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			c.debugLog("CreateVPCFloatingIP attempt %d/%d after transient/not-ready error: %v", attempt+1, maxAttempts, err)
			time.Sleep(retryDelay)
		}
		resp, err = c.post(path, req)
		if err == nil || (!isTransientServerError(err) && !isVPCNotReadyError(err)) {
			break
		}
	}
	if err != nil {
		return nil, fmt.Errorf("create floating IP for VPC %d: %w", vpcID, err)
	}
	var result CreateVPCFloatingIPResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("create floating IP unmarshal: %w", err)
	}
	return &result, nil
}

func (c *V3Client) GetVPCFloatingIP(vpcID, floatingIPID int) (*VPCFloatingIP, error) {
	fips, err := c.ListVPCFloatingIPs(vpcID)
	if err != nil {
		return nil, err
	}
	for _, f := range fips {
		if f.FloatingIPID == floatingIPID {
			return &f, nil
		}
	}
	return nil, &V3NotFoundError{StatusCode: 404, Body: fmt.Sprintf("floating IP %d not found in VPC %d", floatingIPID, vpcID)}
}

func (c *V3Client) UpdateVPCFloatingIP(vpcID, floatingIPID int, req *UpdateVPCFloatingIPRequest) error {
	path := fmt.Sprintf("/vpcs/%d/floating-ips/%d", vpcID, floatingIPID)
	_, err := c.patch(path, req)
	if err != nil {
		return fmt.Errorf("update floating IP %d for VPC %d: %w", floatingIPID, vpcID, err)
	}
	return nil
}

func (c *V3Client) DeleteVPCFloatingIP(vpcID, floatingIPID int) error {
	path := fmt.Sprintf("/vpcs/%d/floating-ips/%d", vpcID, floatingIPID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete floating IP %d for VPC %d: %w", floatingIPID, vpcID, err)
	}
	return nil
}

func (c *V3Client) ListVPCFloatingIPs(vpcID int) ([]VPCFloatingIP, error) {
	path := fmt.Sprintf("/vpcs/%d/floating-ips", vpcID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("list floating IPs for VPC %d: %w", vpcID, err)
	}
	var listData V3ListData
	if err := json.Unmarshal(resp.Data, &listData); err != nil {
		return nil, fmt.Errorf("list floating IPs unmarshal outer: %w", err)
	}
	var fips []VPCFloatingIP
	if err := json.Unmarshal(listData.Data, &fips); err != nil {
		return nil, fmt.Errorf("list floating IPs unmarshal inner: %w", err)
	}
	return fips, nil
}
