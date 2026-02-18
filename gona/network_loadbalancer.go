package gona

import (
	"encoding/json"
	"fmt"
)

// --- Nested types ---

type NLBGroupMatch struct {
	Address string `json:"address"`
}

type NLBGroupHealthCheck struct {
	Enabled  bool   `json:"enabled"`
	Method   string `json:"method"`
	Interval int    `json:"interval"`
	Retries  int    `json:"retries"`
	Delay    int    `json:"delay"`
	Timeout  int    `json:"timeout"`
}

type NLBGroupRulePorts struct {
	Match    int `json:"match"`
	Internal int `json:"internal"`
}

type NLBGroupRule struct {
	Protocol      string            `json:"protocol"`
	NetworkRuleID int               `json:"networkRuleId,omitempty"`
	Ports         NLBGroupRulePorts `json:"ports"`
}

type NLBGroupBackend struct {
	Name             string `json:"name"`
	InternalAddress  string `json:"internalAddress"`
	IsOnline         bool   `json:"isOnline,omitempty"`
	NetworkBackendID int    `json:"networkBackendId,omitempty"`
}

// --- Main types ---

type NLBGroup struct {
	NetworkGroupID int                 `json:"networkGroupId"`
	Name           string              `json:"name"`
	Description    string              `json:"description"`
	IPVersion      int                 `json:"ipVersion"`
	Algorithm      string              `json:"algorithm"`
	IsOnline       bool                `json:"isOnline"`
	Match          NLBGroupMatch       `json:"match"`
	HealthCheck    NLBGroupHealthCheck `json:"healthCheck"`
	Rules          []NLBGroupRule      `json:"rules"`
	Backends       []NLBGroupBackend   `json:"backends"`
}

type CreateNLBGroupRequest struct {
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	IPVersion   int                 `json:"ipVersion"`
	Algorithm   string              `json:"algorithm"`
	Match       NLBGroupMatch       `json:"match"`
	HealthCheck NLBGroupHealthCheck `json:"healthCheck"`
	Rules       []NLBGroupRule      `json:"rules"`
	Backends    []NLBGroupBackend   `json:"backends"`
}

type ReplaceNLBGroupRequest struct {
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	IPVersion   int                 `json:"ipVersion"`
	Algorithm   string              `json:"algorithm"`
	Match       NLBGroupMatch       `json:"match"`
	HealthCheck NLBGroupHealthCheck `json:"healthCheck"`
	Rules       []NLBGroupRule      `json:"rules"`
	Backends    []NLBGroupBackend   `json:"backends"`
}

// --- CRUD ---

func (c *V3Client) CreateNLBGroup(nlbID int, req *CreateNLBGroupRequest) (*NLBGroup, error) {
	path := fmt.Sprintf("/network-loadbalancers/%d/groups", nlbID)
	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("create NLB group for NLB %d: %w", nlbID, err)
	}
	var group NLBGroup
	if err := json.Unmarshal(resp.Data, &group); err != nil {
		return nil, fmt.Errorf("create NLB group unmarshal: %w", err)
	}
	return &group, nil
}

func (c *V3Client) GetNLBGroup(nlbID, groupID int) (*NLBGroup, error) {
	path := fmt.Sprintf("/network-loadbalancers/%d/groups/%d", nlbID, groupID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get NLB group %d for NLB %d: %w", groupID, nlbID, err)
	}
	var group NLBGroup
	if err := json.Unmarshal(resp.Data, &group); err != nil {
		return nil, fmt.Errorf("get NLB group unmarshal: %w", err)
	}
	return &group, nil
}

func (c *V3Client) ListNLBGroups(nlbID int) ([]NLBGroup, error) {
	path := fmt.Sprintf("/network-loadbalancers/%d/groups", nlbID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("list NLB groups for NLB %d: %w", nlbID, err)
	}
	var listData V3ListData
	if err := json.Unmarshal(resp.Data, &listData); err != nil {
		return nil, fmt.Errorf("list NLB groups unmarshal outer: %w", err)
	}
	var groups []NLBGroup
	if err := json.Unmarshal(listData.Data, &groups); err != nil {
		return nil, fmt.Errorf("list NLB groups unmarshal inner: %w", err)
	}
	return groups, nil
}

func (c *V3Client) ReplaceNLBGroup(nlbID, groupID int, req *ReplaceNLBGroupRequest) (*NLBGroup, error) {
	path := fmt.Sprintf("/network-loadbalancers/%d/groups/%d", nlbID, groupID)
	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("replace NLB group %d for NLB %d: %w", groupID, nlbID, err)
	}
	var group NLBGroup
	if err := json.Unmarshal(resp.Data, &group); err != nil {
		return nil, fmt.Errorf("replace NLB group unmarshal: %w", err)
	}
	return &group, nil
}

func (c *V3Client) DeleteNLBGroup(nlbID, groupID int) error {
	path := fmt.Sprintf("/network-loadbalancers/%d/groups/%d", nlbID, groupID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete NLB group %d for NLB %d: %w", groupID, nlbID, err)
	}
	return nil
}
