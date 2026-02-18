package gona

import (
	"encoding/json"
	"fmt"
)

type VPC struct {
	VPCID    int         `json:"vpcId"`
	Metadata VPCMetadata `json:"metadata"`
	Location V3Location  `json:"location"`
	Bastion  struct {
		Enabled   bool `json:"enabled"`
		Addresses struct {
			IPv4 string `json:"ipv4"`
			IPv6 string `json:"ipv6"`
		} `json:"addresses"`
		Port *int `json:"port"`
	} `json:"bastion"`
	Firewalls *VPCResponseFirewalls `json:"firewalls,omitempty"`
	InternalNetwork *VPCNetwork `json:"internalNetwork,omitempty"`
	DHCP *VPCDHCP `json:"dhcp,omitempty"`
	IPv4NetworkDetails *VPCIPv4NetworkDetails `json:"ipv4NetworkDetails,omitempty"`
	FloatingIPs VPCResponseFloatingIPs `json:"floatingIps"`
	Gateways json.RawMessage `json:"gateways,omitempty"`
	Counters struct {
		IPReservations struct {
			Gateways   int `json:"gateways"`
			Interfaces int `json:"interfaces"`
			VMs        int `json:"vms"`
		} `json:"ipReservations"`
		Rules struct {
			DNAT     V3RuleCounter `json:"dnat"`
			SNAT     V3RuleCounter `json:"snat"`
			Firewall V3RuleCounter `json:"firewall"`
		} `json:"rules"`
	} `json:"counters"`
	LoadBalancers      *VPCLoadBalancers `json:"loadBalancers,omitempty"`
	LoadBalancerGroups json.RawMessage   `json:"loadBalancerGroups,omitempty"`
}

type VPCMetadata struct {
	CreatedOn     string `json:"createdOn,omitempty"`
	Label         string `json:"label"`
	Description   string `json:"description"`
	ReadyOn       string `json:"readyOn,omitempty"`
	Status        string `json:"status"`
	UptimeSeconds int    `json:"uptimeSeconds,omitempty"`
}

type VPCResponseFirewalls struct {
	IPv4 *VPCResponseFirewallIP `json:"ipv4,omitempty"`
	IPv6 *VPCResponseFirewallIP `json:"ipv6,omitempty"`
}

type VPCResponseFirewallIP struct {
	Inbound  *VPCResponseFirewallToggle `json:"inbound,omitempty"`
	Outbound *VPCResponseFirewallToggle `json:"outbound,omitempty"`
}

type VPCResponseFirewallToggle struct {
	Enabled bool `json:"enabled"`
}

type VPCDHCP struct {
	Nameservers *VPCDHCPNameservers `json:"nameservers,omitempty"`
}

type VPCDHCPNameservers struct {
	IPv4 []string `json:"ipv4,omitempty"`
	IPv6 []string `json:"ipv6,omitempty"`
}

type VPCIPv4NetworkDetails struct {
	Netmask   string `json:"netmask,omitempty"`
	Broadcast string `json:"broadcast,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
}

type VPCResponseFloatingIPs struct {
	IPv4 []string `json:"ipv4,omitempty"`
	IPv6 []string `json:"ipv6,omitempty"`
}

type V3RuleCounter struct {
	Applied         int `json:"applied"`
	DeleteRequested int `json:"deleteRequested"`
	Total           int `json:"total"`
}

type VPCLoadBalancers struct {
	Network []VPCNetworkLoadBalancer `json:"network,omitempty"`
	HTTP    []VPCHTTPLoadBalancer    `json:"http,omitempty"`
}

type VPCNetworkLoadBalancer struct {
	NetworkLbID int    `json:"networkLbId"`
	CreatedOn   string `json:"createdOn,omitempty"`
}

type VPCHTTPLoadBalancer struct {
	HTTPLbID  int    `json:"httpLbId"`
	CreatedOn string `json:"createdOn,omitempty"`
}

type VPCNetwork struct {
	IPv4 string `json:"ipv4,omitempty"`
	IPv6 string `json:"ipv6,omitempty"`
}

type VPCNameservers struct {
	IPv4 []VPCNameserver `json:"ipv4,omitempty"`
	IPv6 []VPCNameserver `json:"ipv6,omitempty"`
}

type VPCNameserver struct {
	Server string `json:"server"`
}

type VPCFirewalls struct {
	IPv4 *VPCFirewallDirections `json:"ipv4,omitempty"`
	IPv6 *VPCFirewallDirections `json:"ipv6,omitempty"`
}

type VPCFirewallDirections struct {
	Inbound  *bool `json:"inbound,omitempty"`
	Outbound *bool `json:"outbound,omitempty"`
}

type VPCDefaults struct {
	EnableDefaultSnatRule *bool `json:"enableDefaultSnatRule,omitempty"`
}

type VPCPortRange struct {
	Start int `json:"start,omitempty"`
	End   int `json:"end,omitempty"`
}

type CreateVPCRequest struct {
	Label       string          `json:"label"`
	Description string          `json:"description"`
	LocationID  int             `json:"location_id"`
	Network     *VPCNetwork     `json:"network,omitempty"`
	Nameservers *VPCNameservers `json:"nameservers,omitempty"`
	Firewalls   *VPCFirewalls   `json:"firewalls,omitempty"`
	Defaults    *VPCDefaults    `json:"defaults,omitempty"`
}

type UpdateVPCRequest struct {
	Label       string        `json:"label,omitempty"`
	Description string        `json:"description,omitempty"`
	Firewalls   *VPCFirewalls `json:"firewalls,omitempty"`
}

func (c *V3Client) CreateVPC(req *CreateVPCRequest) (*VPC, error) {
	resp, err := c.post("/vpcs", req)
	if err != nil {
		return nil, fmt.Errorf("create VPC: %w", err)
	}
	var vpc VPC
	if err := json.Unmarshal(resp.Data, &vpc); err != nil {
		return nil, fmt.Errorf("create VPC unmarshal: %w", err)
	}
	return &vpc, nil
}

func (c *V3Client) GetVPC(vpcID int) (*VPC, error) {
	path := fmt.Sprintf("/vpcs/%d", vpcID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get VPC %d: %w", vpcID, err)
	}
	var vpc VPC
	if err := json.Unmarshal(resp.Data, &vpc); err != nil {
		return nil, fmt.Errorf("get VPC %d unmarshal: %w", vpcID, err)
	}
	return &vpc, nil
}

func (c *V3Client) UpdateVPC(vpcID int, req *UpdateVPCRequest) (*VPC, error) {
	path := fmt.Sprintf("/vpcs/%d", vpcID)
	resp, err := c.patch(path, req)
	if err != nil {
		return nil, fmt.Errorf("update VPC %d: %w", vpcID, err)
	}
	var vpc VPC
	if err := json.Unmarshal(resp.Data, &vpc); err != nil {
		return nil, fmt.Errorf("update VPC %d unmarshal: %w", vpcID, err)
	}
	return &vpc, nil
}

func (c *V3Client) DeleteVPC(vpcID int) error {
	path := fmt.Sprintf("/vpcs/%d", vpcID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete VPC %d: %w", vpcID, err)
	}
	return nil
}

func (c *V3Client) WaitForVPCReady(vpcID int) error {
	return c.waitForCondition(func() (bool, error) {
		vpc, err := c.GetVPC(vpcID)
		if err != nil {
			return false, err
		}
		c.debugLog("VPC %d status: %s", vpcID, vpc.Metadata.Status)
		return vpc.Metadata.Status == "Running", nil
	}, VPCWaitConfig)
}
