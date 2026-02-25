package gona

import (
	"encoding/json"
	"fmt"
)

type RouterVRFSNATRule struct {
	SNATRuleID  int    `json:"snatRuleId"`
	IPVersion   int    `json:"ipVersion"`
	Protocol    string `json:"protocol"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Match       *struct {
		InterfaceID int           `json:"interfaceId"`
		Network     string        `json:"network"`
		Port        *VPCPortRange `json:"port"`
	} `json:"match"`
	Translation *struct {
		Network string        `json:"network"`
		Port    *VPCPortRange `json:"port"`
	} `json:"translation"`
	Priority *struct {
		Location        string `json:"location"`
		AfterSnatRuleId *int   `json:"afterSnatRuleId"`
	} `json:"priority"`
}

type CreateRouterVRFSNATRuleRequest struct {
	IPVersion   int    `json:"ipVersion"`
	Protocol    string `json:"protocol"`
	Description string `json:"description,omitempty"`
	Match       *struct {
		InterfaceID int           `json:"interfaceId"`
		Network     string        `json:"network"`
		Port        *VPCPortRange `json:"port,omitempty"`
	} `json:"match"`
	Translation *struct {
		Network string        `json:"network"`
		Port    *VPCPortRange `json:"port,omitempty"`
	} `json:"translation"`
	Priority *struct {
		Location        string `json:"location,omitempty"`
		AfterSnatRuleId *int   `json:"afterSnatRuleId,omitempty"`
	} `json:"priority,omitempty"`
}

type UpdateRouterVRFSNATRuleRequest struct {
	IPVersion   int    `json:"ipVersion"`
	Protocol    string `json:"protocol"`
	Description string `json:"description"`
	Match       *struct {
		InterfaceID int           `json:"interfaceId"`
		Network     string        `json:"network"`
		Port        *VPCPortRange `json:"port,omitempty"`
	} `json:"match,omitempty"`
	Translation *struct {
		Network string        `json:"network"`
		Port    *VPCPortRange `json:"port,omitempty"`
	} `json:"translation,omitempty"`
	Priority *struct {
		Location        string `json:"location,omitempty"`
		AfterSnatRuleId *int   `json:"afterSnatRuleId,omitempty"`
	} `json:"priority,omitempty"`
}

type CreateRouterVRFSNATRuleResponse struct {
	SNATRuleID int `json:"snatRuleId"`
}

type UpdateRouterVRFSNATRuleResponse struct {
	SNATRuleID int `json:"snatRuleId"`
}

func (c *V3Client) CreateRouterVRFSNATRule(routerID, vrfID int, req *CreateRouterVRFSNATRuleRequest) (*CreateRouterVRFSNATRuleResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/snat-rules", routerID, vrfID)

	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create SNAT rule for router %d VRF %d: %w", routerID, vrfID, err)
	}

	var createResp CreateRouterVRFSNATRuleResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create SNAT rule response: %w", err)
	}

	return &createResp, nil
}

func (c *V3Client) ListRouterVRFSNATRules(routerID, vrfID int) ([]RouterVRFSNATRule, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/snat-rules", routerID, vrfID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list SNAT rules for router %d VRF %d: %w", routerID, vrfID, err)
	}

	var rules []RouterVRFSNATRule
	if err := json.Unmarshal(resp.Data, &rules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SNAT rules: %w", err)
	}

	return rules, nil
}

func (c *V3Client) GetRouterVRFSNATRule(routerID, vrfID, snatRuleID int) (*RouterVRFSNATRule, error) {
	rules, err := c.ListRouterVRFSNATRules(routerID, vrfID)
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {
		if rule.SNATRuleID == snatRuleID {
			return &rule, nil
		}
	}

	return nil, fmt.Errorf("SNAT rule %d not found in router %d VRF %d", snatRuleID, routerID, vrfID)
}

func (c *V3Client) UpdateRouterVRFSNATRule(routerID, vrfID, snatRuleID int, req *UpdateRouterVRFSNATRuleRequest) (*UpdateRouterVRFSNATRuleResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/snat-rules/%d", routerID, vrfID, snatRuleID)

	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update SNAT rule %d for router %d VRF %d: %w", snatRuleID, routerID, vrfID, err)
	}

	var updateResp UpdateRouterVRFSNATRuleResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update SNAT rule response: %w", err)
	}

	return &updateResp, nil
}

func (c *V3Client) DeleteRouterVRFSNATRule(routerID, vrfID, snatRuleID int) error {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/snat-rules/%d", routerID, vrfID, snatRuleID)

	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("failed to delete SNAT rule %d for router %d VRF %d: %w", snatRuleID, routerID, vrfID, err)
	}

	return nil
}
