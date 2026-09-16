package gona

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
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
	Firewalls          *VPCResponseFirewalls  `json:"firewalls,omitempty"`
	InternalNetwork    *VPCNetwork            `json:"internalNetwork,omitempty"`
	DHCP               *VPCDHCP               `json:"dhcp,omitempty"`
	IPv4NetworkDetails *VPCIPv4NetworkDetails `json:"ipv4NetworkDetails,omitempty"`
	FloatingIPs        VPCResponseFloatingIPs `json:"floatingIps"`
	Gateways           json.RawMessage        `json:"gateways,omitempty"`
	Counters           struct {
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

// UnmarshalJSON accepts both shapes this API returns for a VPC.
//
// GET /vpcs/{id} returns vpcId with the VPC's own fields nested under metadata. GET /vpcs
// returns each row flat, identified by id, with label, description, createdOn and readyOn at
// the top level and no metadata key at all. Describing only the nested form leaves every
// listed VPC with a zero VPCID, so a caller cannot match a VPC it just created.
//
// The flat row is the metadata shape alongside the fields both forms share, so it is decoded
// twice: once through the alias for the shared fields, once into the metadata.
func (v *VPC) UnmarshalJSON(data []byte) error {
	var probe map[string]json.RawMessage
	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}

	type alias VPC
	if err := json.Unmarshal(data, (*alias)(v)); err != nil {
		return err
	}
	if _, ok := probe["metadata"]; ok {
		return nil
	}

	if err := json.Unmarshal(data, &v.Metadata); err != nil {
		return err
	}
	if rawID, ok := probe["id"]; ok {
		if err := json.Unmarshal(rawID, &v.VPCID); err != nil {
			return err
		}
	}
	return nil
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
	Network *VPCLoadBalancerSet `json:"network,omitempty"`
	HTTP    *VPCLoadBalancerSet `json:"http,omitempty"`
}

// VPCLoadBalancerSet holds what the VPC payload says about one class of load balancer.
//
// The field is polymorphic. It arrives either as a list of balancers carrying identifiers, or as
// an object counting how many are applied, requested for deletion and in total. Both shapes are
// valid answers from the same endpoint, so both are decoded here and a caller reads whichever it
// needs: IDs is empty when the platform sent counts, and Total is the list length when it sent a
// list.
type VPCLoadBalancerSet struct {
	Applied         int
	DeleteRequested int
	Total           int
	IDs             []int
	Raw             json.RawMessage
}

func (s *VPCLoadBalancerSet) UnmarshalJSON(data []byte) error {
	s.Raw = append(s.Raw[:0], data...)
	trimmed := bytes.TrimSpace(data)

	if bytes.HasPrefix(trimmed, []byte("[")) {
		var rows []map[string]json.RawMessage
		if err := json.Unmarshal(data, &rows); err != nil {
			return err
		}
		for _, row := range rows {
			for key, raw := range row {
				lower := strings.ToLower(key)
				if !strings.HasSuffix(lower, "lbid") && lower != "id" {
					continue
				}
				var id int
				if err := json.Unmarshal(raw, &id); err == nil && id != 0 {
					s.IDs = append(s.IDs, id)
					break
				}
			}
		}
		s.Total = len(rows)
		s.Applied = len(rows)
		return nil
	}

	var counts struct {
		Applied         int `json:"applied"`
		DeleteRequested int `json:"deleteRequested"`
		Total           int `json:"total"`
	}
	if err := json.Unmarshal(data, &counts); err != nil {
		return err
	}
	s.Applied, s.DeleteRequested, s.Total = counts.Applied, counts.DeleteRequested, counts.Total
	return nil
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

type VPCIPReservations struct {
	Gateways   json.RawMessage `json:"gateways,omitempty"`
	Interfaces json.RawMessage `json:"interfaces,omitempty"`
	VMs        json.RawMessage `json:"vms,omitempty"`
}

// ListVPCs returns all VPCs visible to the account.
func (c *V3Client) ListVPCs() ([]VPC, error) {
	listData, err := c.getList("/vpcs?limit=1000")
	if err != nil {
		return nil, fmt.Errorf("list VPCs: %w", err)
	}
	var vpcs []VPC
	if err := json.Unmarshal(listData.Data, &vpcs); err != nil {
		return nil, fmt.Errorf("list VPCs unmarshal: %w", err)
	}
	return vpcs, nil
}

// ListVPCLocations returns the locations where VPCs can be created.
func (c *V3Client) ListVPCLocations() ([]V3Location, error) {
	resp, err := c.get("/vpcs/locations")
	if err != nil {
		return nil, fmt.Errorf("list VPC locations: %w", err)
	}
	var locations []V3Location
	if err := json.Unmarshal(resp.Data, &locations); err != nil {
		return nil, fmt.Errorf("list VPC locations unmarshal: %w", err)
	}
	return locations, nil
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

// AddVPCStandbyGateway adds a redundant standby gateway to a VPC.
func (c *V3Client) AddVPCStandbyGateway(vpcID int) error {
	path := fmt.Sprintf("/vpcs/%d/gateway/standby", vpcID)
	if _, err := c.post(path, nil); err != nil {
		return fmt.Errorf("add standby gateway for VPC %d: %w", vpcID, err)
	}
	return nil
}

// DeleteVPCStandbyGateway removes the redundant standby gateway from a VPC.
func (c *V3Client) DeleteVPCStandbyGateway(vpcID int) error {
	path := fmt.Sprintf("/vpcs/%d/gateway/standby", vpcID)
	if _, err := c.del(path); err != nil {
		return fmt.Errorf("delete standby gateway for VPC %d: %w", vpcID, err)
	}
	return nil
}

// GetVPCIPReservations returns the gateway, interface and VM IP reservations for a VPC.
func (c *V3Client) GetVPCIPReservations(vpcID int) (*VPCIPReservations, error) {
	path := fmt.Sprintf("/vpcs/%d/ip-reservations", vpcID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get IP reservations for VPC %d: %w", vpcID, err)
	}
	var reservations VPCIPReservations
	if err := json.Unmarshal(resp.Data, &reservations); err != nil {
		return nil, fmt.Errorf("get IP reservations for VPC %d unmarshal: %w", vpcID, err)
	}
	return &reservations, nil
}

// GetVPCSSHSettings returns the bastion SSH settings for a VPC.
func (c *V3Client) GetVPCSSHSettings(vpcID int) (*VPCSSHSettings, error) {
	path := fmt.Sprintf("/vpcs/%d/ssh", vpcID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get SSH settings for VPC %d: %w", vpcID, err)
	}
	var settings VPCSSHSettings
	if err := json.Unmarshal(resp.Data, &settings); err != nil {
		return nil, fmt.Errorf("get SSH settings for VPC %d unmarshal: %w", vpcID, err)
	}
	return &settings, nil
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
