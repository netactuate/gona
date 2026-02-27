package gona

import (
	"encoding/json"
	"fmt"
)

type RouterVRFDNATRule struct {
	DNATRuleID  int    `json:"dnatRuleId"`
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
		AfterDnatRuleId *int   `json:"afterDnatRuleId"`
	} `json:"priority"`
}

type CreateRouterVRFDNATRuleRequest struct {
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
		AfterDnatRuleId *int   `json:"afterDnatRuleId,omitempty"`
	} `json:"priority,omitempty"`
}

type UpdateRouterVRFDNATRuleRequest struct {
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
		AfterDnatRuleId *int   `json:"afterDnatRuleId,omitempty"`
	} `json:"priority,omitempty"`
}

type CreateRouterVRFDNATRuleResponse struct {
	DNATRuleID int `json:"dnatRuleId"`
}

type UpdateRouterVRFDNATRuleResponse struct {
	DNATRuleID int `json:"dnatRuleId"`
}

func (c *V3Client) CreateRouterVRFDNATRule(routerID, vrfID int, req *CreateRouterVRFDNATRuleRequest) (*CreateRouterVRFDNATRuleResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/dnat-rules", routerID, vrfID)

	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create DNAT rule for router %d VRF %d: %w", routerID, vrfID, err)
	}

	var createResp CreateRouterVRFDNATRuleResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create DNAT rule response: %w", err)
	}

	return &createResp, nil
}

func (c *V3Client) ListRouterVRFDNATRules(routerID, vrfID int) ([]RouterVRFDNATRule, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/dnat-rules", routerID, vrfID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list DNAT rules for router %d VRF %d: %w", routerID, vrfID, err)
	}

	var rules []RouterVRFDNATRule
	if err := json.Unmarshal(resp.Data, &rules); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DNAT rules: %w", err)
	}

	return rules, nil
}

func (c *V3Client) GetRouterVRFDNATRule(routerID, vrfID, dnatRuleID int) (*RouterVRFDNATRule, error) {
	rules, err := c.ListRouterVRFDNATRules(routerID, vrfID)
	if err != nil {
		return nil, err
	}

	for _, rule := range rules {
		if rule.DNATRuleID == dnatRuleID {
			return &rule, nil
		}
	}

	return nil, fmt.Errorf("DNAT rule %d not found in router %d VRF %d", dnatRuleID, routerID, vrfID)
}

func (c *V3Client) UpdateRouterVRFDNATRule(routerID, vrfID, dnatRuleID int, req *UpdateRouterVRFDNATRuleRequest) (*UpdateRouterVRFDNATRuleResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/dnat-rules/%d", routerID, vrfID, dnatRuleID)

	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update DNAT rule %d for router %d VRF %d: %w", dnatRuleID, routerID, vrfID, err)
	}

	var updateResp UpdateRouterVRFDNATRuleResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update DNAT rule response: %w", err)
	}

	return &updateResp, nil
}

func (c *V3Client) DeleteRouterVRFDNATRule(routerID, vrfID, dnatRuleID int) error {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/dnat-rules/%d", routerID, vrfID, dnatRuleID)

	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("failed to delete DNAT rule %d for router %d VRF %d: %w", dnatRuleID, routerID, vrfID, err)
	}

	return nil
}
