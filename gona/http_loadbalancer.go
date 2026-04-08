package gona

import (
	"encoding/json"
	"fmt"
)

// --- Nested types ---

type HTTPLBGroupMatch struct {
	Address string `json:"address"`
	Ports   string `json:"ports"`
}

type HTTPLBGroupHealthCheck struct {
	Active  HTTPLBGroupHealthCheckActive  `json:"active"`
	Passive HTTPLBGroupHealthCheckPassive `json:"passive"`
}

type HTTPLBGroupHealthCheckActive struct {
	Enabled  bool   `json:"enabled"`
	Interval int    `json:"interval,omitempty"`
	Retries  int    `json:"retries,omitempty"`
	Delay    int    `json:"delay,omitempty"`
	Timeout  *int   `json:"timeout"`
	Path     string `json:"path,omitempty"`
}

type HTTPLBGroupHealthCheckPassive struct {
	Enabled bool `json:"enabled"`
}

type HTTPLBGroupRuleMatch struct {
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type HTTPLBGroupRuleSSL struct {
	Enabled          bool `json:"enabled"`
	SSLCertificateID *int `json:"sslCertificateId"`
}

type HTTPLBGroupRule struct {
	HTTPRuleID           int                  `json:"httpRuleId,omitempty"`
	HTTPSRedirectEnabled bool                 `json:"httpsRedirectEnabled"`
	Match                HTTPLBGroupRuleMatch `json:"match"`
	SSL                  HTTPLBGroupRuleSSL   `json:"ssl"`
}

type HTTPLBGroupBackend struct {
	Name            string `json:"name"`
	InternalAddress string `json:"internalAddress"`
	IsOnline        bool   `json:"isOnline,omitempty"`
	HTTPBackendID   int    `json:"httpBackendId,omitempty"`
}

// --- Main types ---

type HTTPLBGroup struct {
	HTTPGroupID           int                    `json:"httpGroupId"`
	Name                  string                 `json:"name"`
	Description           string                 `json:"description"`
	Algorithm             string                 `json:"algorithm"`
	StickySessionsEnabled bool                   `json:"stickySessionsEnabled"`
	SSLToBackendEnabled   bool                   `json:"sslToBackendEnabled"`
	InternalPort          int                    `json:"internalPort"`
	IsOnline              bool                   `json:"isOnline"`
	Match                 HTTPLBGroupMatch       `json:"match"`
	HealthCheck           HTTPLBGroupHealthCheck `json:"healthCheck"`
	Rules                 []HTTPLBGroupRule      `json:"rules"`
	Backends              []HTTPLBGroupBackend   `json:"backends"`
}

type CreateHTTPLBGroupRequest struct {
	Name                  string                 `json:"name"`
	Description           string                 `json:"description,omitempty"`
	Algorithm             string                 `json:"algorithm"`
	StickySessionsEnabled bool                   `json:"stickySessionsEnabled"`
	SSLToBackendEnabled   bool                   `json:"sslToBackendEnabled"`
	InternalPort          int                    `json:"internalPort"`
	Match                 HTTPLBGroupMatch       `json:"match"`
	HealthCheck           HTTPLBGroupHealthCheck `json:"healthCheck"`
	Rules                 []HTTPLBGroupRule      `json:"rules"`
	Backends              []HTTPLBGroupBackend   `json:"backends"`
}

type ReplaceHTTPLBGroupRequest struct {
	Name                  string                 `json:"name"`
	Description           string                 `json:"description,omitempty"`
	Algorithm             string                 `json:"algorithm"`
	StickySessionsEnabled bool                   `json:"stickySessionsEnabled"`
	SSLToBackendEnabled   bool                   `json:"sslToBackendEnabled"`
	InternalPort          int                    `json:"internalPort"`
	Match                 HTTPLBGroupMatch       `json:"match"`
	HealthCheck           HTTPLBGroupHealthCheck `json:"healthCheck"`
	Rules                 []HTTPLBGroupRule      `json:"rules"`
	Backends              []HTTPLBGroupBackend   `json:"backends"`
}

// --- CRUD ---

func (c *V3Client) CreateHTTPLBGroup(httpLbID int, req *CreateHTTPLBGroupRequest) (*HTTPLBGroup, error) {
	path := fmt.Sprintf("/http-loadbalancers/%d/groups", httpLbID)
	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("create HTTP LB group for LB %d: %w", httpLbID, err)
	}
	var group HTTPLBGroup
	if err := json.Unmarshal(resp.Data, &group); err != nil {
		return nil, fmt.Errorf("create HTTP LB group unmarshal: %w", err)
	}
	return &group, nil
}

func (c *V3Client) GetHTTPLBGroup(httpLbID, groupID int) (*HTTPLBGroup, error) {
	path := fmt.Sprintf("/http-loadbalancers/%d/groups/%d", httpLbID, groupID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get HTTP LB group %d for LB %d: %w", groupID, httpLbID, err)
	}
	var group HTTPLBGroup
	if err := json.Unmarshal(resp.Data, &group); err != nil {
		return nil, fmt.Errorf("get HTTP LB group unmarshal: %w", err)
	}
	return &group, nil
}

func (c *V3Client) ReplaceHTTPLBGroup(httpLbID, groupID int, req *ReplaceHTTPLBGroupRequest) (*HTTPLBGroup, error) {
	path := fmt.Sprintf("/http-loadbalancers/%d/groups/%d", httpLbID, groupID)
	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("replace HTTP LB group %d for LB %d: %w", groupID, httpLbID, err)
	}
	var group HTTPLBGroup
	if err := json.Unmarshal(resp.Data, &group); err != nil {
		return nil, fmt.Errorf("replace HTTP LB group unmarshal: %w", err)
	}
	return &group, nil
}

func (c *V3Client) DeleteHTTPLBGroup(httpLbID, groupID int) error {
	path := fmt.Sprintf("/http-loadbalancers/%d/groups/%d", httpLbID, groupID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete HTTP LB group %d for LB %d: %w", groupID, httpLbID, err)
	}
	return nil
}
