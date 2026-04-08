package gona

import (
	"encoding/json"
	"fmt"
	"time"
)

type VPCSNATRule struct {
	SNATRuleID  int    `json:"snatRuleId"`
	IPVersion   int    `json:"ipVersion"`
	Protocol    string `json:"protocol,omitempty"`
	Description string `json:"description,omitempty"`
	Match       *struct {
		InternalCidr string `json:"internalCidr,omitempty"`
	} `json:"match,omitempty"`
	Translation *struct {
		Address *struct {
			Start string `json:"start,omitempty"`
			End   string `json:"end,omitempty"`
		} `json:"address,omitempty"`
		Port *VPCPortRange `json:"port,omitempty"`
	} `json:"translation,omitempty"`
}

type CreateVPCSNATRuleRequest struct {
	IPVersion   int    `json:"ipVersion"`
	Protocol    string `json:"protocol,omitempty"`
	Description string `json:"description,omitempty"`
	Match       *struct {
		InternalCidr string `json:"internalCidr,omitempty"`
	} `json:"match,omitempty"`
	Translation *struct {
		Address *struct {
			Start string `json:"start,omitempty"`
			End   string `json:"end,omitempty"`
		} `json:"address,omitempty"`
		Port *VPCPortRange `json:"port,omitempty"`
	} `json:"translation,omitempty"`
	Priority *struct {
		Location        string `json:"location,omitempty"`
		AfterSnatRuleId int    `json:"afterSnatRuleId,omitempty"`
	} `json:"priority,omitempty"`
}

type UpdateVPCSNATRuleRequest struct {
	Protocol    string `json:"protocol,omitempty"`
	Description string `json:"description,omitempty"`
	Match       *struct {
		InternalCidr string `json:"internalCidr,omitempty"`
	} `json:"match,omitempty"`
	Translation *struct {
		Address *struct {
			Start string `json:"start,omitempty"`
			End   string `json:"end,omitempty"`
		} `json:"address,omitempty"`
		Port *VPCPortRange `json:"port,omitempty"`
	} `json:"translation,omitempty"`
	Priority *struct {
		Location        string `json:"location,omitempty"`
		AfterSnatRuleId int    `json:"afterSnatRuleId,omitempty"`
	} `json:"priority,omitempty"`
}

func (c *V3Client) CreateVPCSNATRule(vpcID int, req *CreateVPCSNATRuleRequest) (*VPCSNATRule, error) {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/snat", vpcID)
	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("create SNAT rule for VPC %d: %w", vpcID, err)
	}
	var rule VPCSNATRule
	if err := json.Unmarshal(resp.Data, &rule); err != nil {
		return nil, fmt.Errorf("create SNAT rule unmarshal: %w", err)
	}
	return &rule, nil
}

func (c *V3Client) ListVPCSNATRules(vpcID, ipVersion int) ([]VPCSNATRule, error) {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/snat/ipv%d", vpcID, ipVersion)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("list SNAT rules for VPC %d: %w", vpcID, err)
	}
	var listData V3ListData
	if err := json.Unmarshal(resp.Data, &listData); err != nil {
		return nil, fmt.Errorf("list SNAT rules unmarshal outer: %w", err)
	}
	var rules []VPCSNATRule
	if err := json.Unmarshal(listData.Data, &rules); err != nil {
		return nil, fmt.Errorf("list SNAT rules unmarshal inner: %w", err)
	}
	return rules, nil
}

func (c *V3Client) GetVPCSNATRule(vpcID, ruleID, ipVersion int) (*VPCSNATRule, error) {
	rules, err := c.ListVPCSNATRules(vpcID, ipVersion)
	if err != nil {
		return nil, err
	}
	for _, r := range rules {
		if r.SNATRuleID == ruleID {
			return &r, nil
		}
	}
	return nil, &V3NotFoundError{StatusCode: 404, Body: fmt.Sprintf("SNAT rule %d not found in VPC %d", ruleID, vpcID)}
}

func (c *V3Client) UpdateVPCSNATRule(vpcID, ruleID int, req *UpdateVPCSNATRuleRequest) (*VPCSNATRule, error) {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/snat/%d", vpcID, ruleID)
	resp, err := c.patch(path, req)
	if err != nil {
		return nil, fmt.Errorf("update SNAT rule %d for VPC %d: %w", ruleID, vpcID, err)
	}
	var rule VPCSNATRule
	if err := json.Unmarshal(resp.Data, &rule); err != nil {
		return nil, fmt.Errorf("update SNAT rule unmarshal: %w", err)
	}
	return &rule, nil
}

func (c *V3Client) DeleteVPCSNATRule(vpcID, ruleID int) error {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/snat/%d", vpcID, ruleID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete SNAT rule %d for VPC %d: %w", ruleID, vpcID, err)
	}
	return nil
}

func (c *V3Client) ApplyVPCSNATChanges(vpcID int) error {
	path := fmt.Sprintf("/vpcs/%d/gateway/rules/snat/apply-changes", vpcID)
	_, err := c.postWithRetry(path, 6, 10*time.Second)
	if err != nil {
		return fmt.Errorf("apply SNAT changes for VPC %d: %w", vpcID, err)
	}
	return nil
}
