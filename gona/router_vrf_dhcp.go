package gona

import (
	"encoding/json"
	"fmt"
)

type RouterDHCPRange struct {
	FirstAddress string `json:"firstAddress"`
	LastAddress  string `json:"lastAddress"`
}

type RouterDHCPServer struct {
	Address string `json:"address"`
}

type RouterDHCPStaticRoute struct {
	Network string `json:"network"`
	NextHop string `json:"nextHop"`
}

type RouterVRFDHCPConfig struct {
	Enabled              bool                    `json:"enabled"`
	InterfaceID          int                     `json:"interfaceId"`
	Subnet               string                  `json:"subnet"`
	DefaultRouterAddress string                  `json:"defaultRouterAddress"`
	ClientDomainName     string                  `json:"clientDomainName"`
	LeaseTimeout         int                     `json:"leaseTimeout"`
	DoPingCheck          bool                    `json:"doPingCheck"`
	Range                *RouterDHCPRange        `json:"range"`
	DomainNameServers    []RouterDHCPServer      `json:"domainNameServers"`
	NTPServers           []RouterDHCPServer      `json:"ntpServers"`
	StaticRoutes         []RouterDHCPStaticRoute `json:"staticRoutes"`
}

type UpdateRouterVRFDHCPRequest struct {
	Enabled              bool                    `json:"enabled"`
	InterfaceID          int                     `json:"interfaceId"`
	Subnet               string                  `json:"subnet"`
	DefaultRouterAddress string                  `json:"defaultRouterAddress"`
	ClientDomainName     string                  `json:"clientDomainName"`
	LeaseTimeout         int                     `json:"leaseTimeout"`
	DoPingCheck          bool                    `json:"doPingCheck"`
	Range                *RouterDHCPRange        `json:"range"`
	DomainNameServers    []RouterDHCPServer      `json:"domainNameServers"`
	NTPServers           []RouterDHCPServer      `json:"ntpServers"`
	StaticRoutes         []RouterDHCPStaticRoute `json:"staticRoutes"`
}

type UpdateRouterVRFDHCPResponse struct {
	RouterID int `json:"routerId"`
}

func (c *V3Client) GetRouterVRFDHCP(routerID, vrfID int) (*RouterVRFDHCPConfig, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/services/dhcp", routerID, vrfID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get DHCP config for router %d VRF %d: %w", routerID, vrfID, err)
	}

	var dhcpConfig RouterVRFDHCPConfig
	if err := json.Unmarshal(resp.Data, &dhcpConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal DHCP config: %w", err)
	}

	return &dhcpConfig, nil
}

func (c *V3Client) UpdateRouterVRFDHCP(routerID, vrfID int, req *UpdateRouterVRFDHCPRequest) (*UpdateRouterVRFDHCPResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/services/dhcp", routerID, vrfID)

	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update DHCP config for router %d VRF %d: %w", routerID, vrfID, err)
	}

	var updateResp UpdateRouterVRFDHCPResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update DHCP response: %w", err)
	}

	return &updateResp, nil
}
