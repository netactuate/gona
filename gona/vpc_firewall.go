package gona

import (
	"encoding/json"
	"fmt"
)

type VPCFirewallRule struct {
	FirewallRuleID int           `json:"firewallRuleId"`
	IPVersion      int           `json:"ipVersion"`
	Direction      string        `json:"direction"`
	Protocol       string        `json:"protocol,omitempty"`
	Description    string        `json:"description,omitempty"`
	Network        string        `json:"network,omitempty"`
	Address        string        `json:"address,omitempty"`
	PrefixLength   int           `json:"prefixLength,omitempty"`
	Port           *VPCPortRange `json:"port,omitempty"`
}

type CreateVPCFirewallRuleResponse struct {
	FirewallRuleID int `json:"firewallRuleId"`
}

type CreateVPCFirewallRuleRequest struct {
	IPVersion   int           `json:"ipVersion"`
	Direction   string        `json:"direction"`
	Protocol    string        `json:"protocol,omitempty"`
	Description string        `json:"description,omitempty"`
	Network     string        `json:"network,omitempty"`
	Port        *VPCPortRange `json:"port,omitempty"`
}

type UpdateVPCFirewallRuleRequest struct {
	Direction   string        `json:"direction,omitempty"`
	IPVersion   int           `json:"ipVersion,omitempty"`
	Protocol    string        `json:"protocol,omitempty"`
	Description string        `json:"description,omitempty"`
	Network     string        `json:"network,omitempty"`
	Port        *VPCPortRange `json:"port,omitempty"`
}

func (c *V3Client) CreateVPCFirewallRule(vpcID int, req *CreateVPCFirewallRuleRequest) (*CreateVPCFirewallRuleResponse, error) {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/firewall", vpcID)
	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("create firewall rule for VPC %d: %w", vpcID, err)
	}
	var result CreateVPCFirewallRuleResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("create firewall rule unmarshal: %w", err)
	}
	return &result, nil
}

func (c *V3Client) ListVPCFirewallRules(vpcID, ipVersion int) ([]VPCFirewallRule, error) {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/firewall/ipv%d", vpcID, ipVersion)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("list firewall rules for VPC %d: %w", vpcID, err)
	}
	var listData V3ListData
	if err := json.Unmarshal(resp.Data, &listData); err != nil {
		return nil, fmt.Errorf("list firewall rules unmarshal outer: %w", err)
	}
	var rules []VPCFirewallRule
	if err := json.Unmarshal(listData.Data, &rules); err != nil {
		return nil, fmt.Errorf("list firewall rules unmarshal inner: %w", err)
	}
	return rules, nil
}

func (c *V3Client) GetVPCFirewallRule(vpcID, ruleID, ipVersion int) (*VPCFirewallRule, error) {
	rules, err := c.ListVPCFirewallRules(vpcID, ipVersion)
	if err != nil {
		return nil, err
	}
	for _, r := range rules {
		if r.FirewallRuleID == ruleID {
			return &r, nil
		}
	}
	return nil, &V3NotFoundError{StatusCode: 404, Body: fmt.Sprintf("firewall rule %d not found in VPC %d", ruleID, vpcID)}
}

func (c *V3Client) UpdateVPCFirewallRule(vpcID, ruleID int, req *UpdateVPCFirewallRuleRequest) (*VPCFirewallRule, error) {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/firewall/%d", vpcID, ruleID)
	resp, err := c.patch(path, req)
	if err != nil {
		return nil, fmt.Errorf("update firewall rule %d for VPC %d: %w", ruleID, vpcID, err)
	}
	var rule VPCFirewallRule
	if err := json.Unmarshal(resp.Data, &rule); err != nil {
		return nil, fmt.Errorf("update firewall rule unmarshal: %w", err)
	}
	return &rule, nil
}

func (c *V3Client) DeleteVPCFirewallRule(vpcID, ruleID int) error {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/firewall/%d", vpcID, ruleID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete firewall rule %d for VPC %d: %w", ruleID, vpcID, err)
	}
	return nil
}

func (c *V3Client) ApplyVPCFirewallChanges(vpcID int) error {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/firewall/apply-changes", vpcID)
	_, err := c.post(path, nil)
	if err != nil {
		return fmt.Errorf("apply firewall changes for VPC %d: %w", vpcID, err)
	}
	return nil
}
