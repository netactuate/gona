package gona

import (
	"encoding/json"
	"fmt"
)

type CreateRouterVRFRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CreateRouterVRFResponse struct {
	VrfID int `json:"vrfId"`
}

type UpdateRouterVRFRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateRouterVRFResponse struct {
	VrfID int `json:"vrfId"`
}

type GetRouterVRFResponse map[string]RouterVRFConfig

type RouterVRFConfig struct {
	DNATRules   []interface{}      `json:"dnatRules"`
	SNATRules   []interface{}      `json:"snatRules"`
	Services    RouterVRFServices  `json:"services"`
	Tunnels     []interface{}      `json:"tunnels"`
	BGP         RouterVRFBGPConfig `json:"bgp"`
	Routes      RoutesConfig       `json:"routes"`
	Interfaces  []interface{}      `json:"interfaces"`
	IPSec       IPSecConfig        `json:"ipSec"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	VrfID       int                `json:"vrfId"`
}

type IPSecConfig struct {
	Peers []interface{} `json:"peers"`
}

type RoutesConfig struct {
	Static []interface{} `json:"static"`
}

type RouterVRFServices struct {
	DHCP DHCPService `json:"dhcp"`
}

type DHCPService struct {
	Enabled              bool          `json:"enabled"`
	InterfaceID          *int          `json:"interfaceId"`
	Subnet               *string       `json:"subnet"`
	DefaultRouterAddress *string       `json:"defaultRouterAddress"`
	ClientDomainName     *string       `json:"clientDomainName"`
	LeaseTimeout         *int          `json:"leaseTimeout"`
	DoPingCheck          *bool         `json:"doPingCheck"`
	Range                *string       `json:"range"`
	DomainNameServers    []string      `json:"domainNameServers"`
	NTPServers           []string      `json:"ntpServers"`
	StaticRoutes         []interface{} `json:"staticRoutes"`
}

func (c *V3Client) CreateRouterVRF(routerID int, req CreateRouterVRFRequest) (*CreateRouterVRFResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs", routerID)

	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create VRF on router %d: %w", routerID, err)
	}

	var createResp CreateRouterVRFResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create VRF response: %w", err)
	}

	return &createResp, nil
}

func (c *V3Client) GetRouterVRF(routerID int, vrfID int) (*RouterVRFConfig, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d", routerID, vrfID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get VRF %d on router %d: %w", vrfID, routerID, err)
	}

	var vrf RouterVRFConfig
	if err := json.Unmarshal(resp.Data, &vrf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal VRF response: %w", err)
	}

	return &vrf, nil
}

func (c *V3Client) UpdateRouterVRF(routerID int, vrfID int, req UpdateRouterVRFRequest) (*UpdateRouterVRFResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d", routerID, vrfID)

	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update VRF %d on router %d: %w", vrfID, routerID, err)
	}

	var updateResp UpdateRouterVRFResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update VRF response: %w", err)
	}

	return &updateResp, nil
}

func (c *V3Client) DeleteRouterVRF(routerID int, vrfID int) error {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d", routerID, vrfID)

	if _, err := c.del(path); err != nil {
		return fmt.Errorf("failed to delete VRF %d on router %d: %w", vrfID, routerID, err)
	}

	return nil
}
