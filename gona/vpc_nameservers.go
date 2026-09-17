package gona

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ReplaceVPCNameserversRequest struct {
	Nameservers []VPCNameserver `json:"nameservers"`
}

type ReplaceVPCNameserversResponse struct {
	Nameservers []VPCNameserver `json:"nameservers"`
}

func (c *V3Client) GetVPCNameservers(vpcID int) (*VPCNameservers, error) {
	// The /vpcs/{id}/dhcp/nameservers path is PUT only; GET on it returns HTTP 405.
	// The current nameservers are read from the VPC object instead, where they arrive
	// under dhcp.nameservers as plain string lists.
	vpc, err := c.GetVPC(vpcID)
	if err != nil {
		return nil, fmt.Errorf("get nameservers for VPC %d: %w", vpcID, err)
	}
	ns := &VPCNameservers{}
	if vpc.DHCP != nil && vpc.DHCP.Nameservers != nil {
		for _, s := range vpc.DHCP.Nameservers.IPv4 {
			ns.IPv4 = append(ns.IPv4, VPCNameserver{Server: s})
		}
		for _, s := range vpc.DHCP.Nameservers.IPv6 {
			ns.IPv6 = append(ns.IPv6, VPCNameserver{Server: s})
		}
	}
	return ns, nil
}

// ReplaceVPCNameservers replaces the DHCP nameservers announced by a VPC.
func (c *V3Client) ReplaceVPCNameservers(vpcID int, req *ReplaceVPCNameserversRequest) (*ReplaceVPCNameserversResponse, error) {
	path := fmt.Sprintf("/vpcs/%d/dhcp/nameservers", vpcID)
	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("replace nameservers for VPC %d: %w", vpcID, err)
	}
	var ns ReplaceVPCNameserversResponse
	// A successful replace returns an empty body; skip decoding when there is no
	// payload rather than failing with "unexpected end of JSON input".
	if trimmed := bytes.TrimSpace(resp.Data); len(trimmed) > 0 {
		if err := json.Unmarshal(trimmed, &ns); err != nil {
			return nil, fmt.Errorf("replace nameservers unmarshal: %w", err)
		}
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
	// PATCH may echo an empty body on success; decode only when a payload is present.
	if trimmed := bytes.TrimSpace(resp.Data); len(trimmed) > 0 {
		if err := json.Unmarshal(trimmed, &ns); err != nil {
			return nil, fmt.Errorf("update nameservers unmarshal: %w", err)
		}
	}
	return &ns, nil
}
