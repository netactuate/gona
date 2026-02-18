package gona

import (
	"encoding/json"
	"fmt"
)

type VPCDNATRule struct {
	DNATRuleID  int    `json:"dnatRuleId"`
	IPVersion   int    `json:"ipVersion"`
	Protocol    string `json:"protocol,omitempty"`
	Description string `json:"description,omitempty"`
	Match       *struct {
		Address string        `json:"address,omitempty"`
		Port    *VPCPortRange `json:"port,omitempty"`
	} `json:"match,omitempty"`
	Translation *struct {
		Address string        `json:"address"`
		Port    *VPCPortRange `json:"port,omitempty"`
	} `json:"translation,omitempty"`
}

type CreateVPCDNATRuleRequest struct {
	IPVersion   int    `json:"ipVersion"`
	Protocol    string `json:"protocol,omitempty"`
	Description string `json:"description,omitempty"`
	Match       *struct {
		Address string        `json:"address,omitempty"`
		Port    *VPCPortRange `json:"port,omitempty"`
	} `json:"match,omitempty"`
	Translation *struct {
		Address string        `json:"address"`
		Port    *VPCPortRange `json:"port,omitempty"`
	} `json:"translation"`
	Priority *struct {
		Location        string `json:"location,omitempty"`
		AfterDnatRuleId int    `json:"afterDnatRuleId,omitempty"`
	} `json:"priority,omitempty"`
}

type UpdateVPCDNATRuleRequest struct {
	Protocol    string `json:"protocol,omitempty"`
	Description string `json:"description,omitempty"`
	Match       *struct {
		Address string        `json:"address,omitempty"`
		Port    *VPCPortRange `json:"port,omitempty"`
	} `json:"match,omitempty"`
	Translation *struct {
		Address string        `json:"address,omitempty"`
		Port    *VPCPortRange `json:"port,omitempty"`
	} `json:"translation,omitempty"`
	Priority *struct {
		Location        string `json:"location,omitempty"`
		AfterDnatRuleId int    `json:"afterDnatRuleId,omitempty"`
	} `json:"priority,omitempty"`
}

func (c *V3Client) CreateVPCDNATRule(vpcID int, req *CreateVPCDNATRuleRequest) (*VPCDNATRule, error) {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/dnat", vpcID)
	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("create DNAT rule for VPC %d: %w", vpcID, err)
	}
	var rule VPCDNATRule
	if err := json.Unmarshal(resp.Data, &rule); err != nil {
		return nil, fmt.Errorf("create DNAT rule unmarshal: %w", err)
	}
	return &rule, nil
}

func (c *V3Client) ListVPCDNATRules(vpcID, ipVersion int) ([]VPCDNATRule, error) {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/dnat/ipv%d", vpcID, ipVersion)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("list DNAT rules for VPC %d: %w", vpcID, err)
	}
	var listData V3ListData
	if err := json.Unmarshal(resp.Data, &listData); err != nil {
		return nil, fmt.Errorf("list DNAT rules unmarshal outer: %w", err)
	}
	var rules []VPCDNATRule
	if err := json.Unmarshal(listData.Data, &rules); err != nil {
		return nil, fmt.Errorf("list DNAT rules unmarshal inner: %w", err)
	}
	return rules, nil
}

func (c *V3Client) GetVPCDNATRule(vpcID, ruleID, ipVersion int) (*VPCDNATRule, error) {
	rules, err := c.ListVPCDNATRules(vpcID, ipVersion)
	if err != nil {
		return nil, err
	}
	for _, r := range rules {
		if r.DNATRuleID == ruleID {
			return &r, nil
		}
	}
	return nil, &V3NotFoundError{StatusCode: 404, Body: fmt.Sprintf("DNAT rule %d not found in VPC %d", ruleID, vpcID)}
}

func (c *V3Client) UpdateVPCDNATRule(vpcID, ruleID int, req *UpdateVPCDNATRuleRequest) (*VPCDNATRule, error) {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/dnat/%d", vpcID, ruleID)
	resp, err := c.patch(path, req)
	if err != nil {
		return nil, fmt.Errorf("update DNAT rule %d for VPC %d: %w", ruleID, vpcID, err)
	}
	var rule VPCDNATRule
	if err := json.Unmarshal(resp.Data, &rule); err != nil {
		return nil, fmt.Errorf("update DNAT rule unmarshal: %w", err)
	}
	return &rule, nil
}

func (c *V3Client) DeleteVPCDNATRule(vpcID, ruleID int) error {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/dnat/%d", vpcID, ruleID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete DNAT rule %d for VPC %d: %w", ruleID, vpcID, err)
	}
	return nil
}

func (c *V3Client) ApplyVPCDNATChanges(vpcID int) error {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/dnat/apply-changes", vpcID)
	_, err := c.post(path, nil)
	if err != nil {
		return fmt.Errorf("apply DNAT changes for VPC %d: %w", vpcID, err)
	}
	return nil
}
