package gona

import (
	"encoding/json"
	"fmt"
)

func (c *V3Client) GetVPCNameservers(vpcID int) (*VPCNameservers, error) {
	path := fmt.Sprintf("/vpcs/%d/dhcp/nameservers", vpcID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get nameservers for VPC %d: %w", vpcID, err)
	}
	var ns VPCNameservers
	if err := json.Unmarshal(resp.Data, &ns); err != nil {
		return nil, fmt.Errorf("get nameservers unmarshal: %w", err)
	}
	return &ns, nil
}

func (c *V3Client) UpdateVPCNameservers(vpcID int, req *VPCNameservers) (*VPCNameservers, error) {
	path := fmt.Sprintf("/vpcs/%d/dhcp/nameservers", vpcID)
	resp, err := c.patch(path, req)
	if err != nil {
		return nil, fmt.Errorf("update nameservers for VPC %d: %w", vpcID, err)
	}
	var ns VPCNameservers
	if err := json.Unmarshal(resp.Data, &ns); err != nil {
		return nil, fmt.Errorf("update nameservers unmarshal: %w", err)
	}
	return &ns, nil
}
