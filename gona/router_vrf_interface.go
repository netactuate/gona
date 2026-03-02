package gona

import (
	"encoding/json"
	"fmt"
)

type CreateRouterVRFInterfaceRequest struct {
	Type               string  `json:"type"`
	Name               string  `json:"name"`
	Description        *string `json:"description,omitempty"`
	IPv4CIDR           *string `json:"ipv4Cidr,omitempty"`
	IPv6CIDR           *string `json:"ipv6Cidr,omitempty"`
	EthernetHardwareID *string `json:"ethernetHardwareId,omitempty"`
	WireguardPort      *int    `json:"wireguardPort,omitempty"`
}

type CreateRouterVRFInterfaceResponse struct {
	InterfaceID int `json:"interfaceId"`
}

type UpdateRouterVRFInterfaceRequest struct {
	Type               string  `json:"type"`
	Name               string  `json:"name"`
	Description        *string `json:"description,omitempty"`
	IPv4CIDR           *string `json:"ipv4Cidr,omitempty"`
	IPv6CIDR           *string `json:"ipv6Cidr,omitempty"`
	EthernetHardwareID *string `json:"ethernetHardwareId,omitempty"`
	WireguardPort      *int    `json:"wireguardPort,omitempty"`
}

type UpdateRouterVRFInterfaceResponse struct {
	InterfaceID int `json:"interfaceId"`
}

type RouterVRFInterface struct {
	InterfaceID        int                                `json:"interfaceId"`
	VrfID              int                                `json:"vrfId"`
	Type               string                             `json:"type"`
	Name               string                             `json:"name"`
	Description        *string                            `json:"description"`
	IPv4CIDR           *string                            `json:"ipv4Cidr"`
	IPv6CIDR           *string                            `json:"ipv6Cidr"`
	EthernetHardwareID *string                            `json:"ethernetHardwareId"`
	WireguardPort      *int                               `json:"wireguardPort"`
	StaticRoutes       []interface{}                      `json:"staticRoutes"`
	Peers              []RouterVRFInterfaceWireguardPeer  `json:"peers,omitempty"`
}

type GetRouterVRFInterfacesResponse map[string]RouterVRFInterface

func (c *V3Client) CreateRouterVRFInterface(routerID int, vrfID int, req CreateRouterVRFInterfaceRequest) (*CreateRouterVRFInterfaceResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/interfaces", routerID, vrfID)

	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create interface on VRF %d, router %d: %w", vrfID, routerID, err)
	}

	var createResp CreateRouterVRFInterfaceResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create interface response: %w", err)
	}

	return &createResp, nil
}

func (c *V3Client) GetRouterVRFInterface(routerID int, vrfID int, interfaceID int) (*RouterVRFInterface, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/interfaces/%d", routerID, vrfID, interfaceID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get interface %d on VRF %d, router %d: %w", interfaceID, vrfID, routerID, err)
	}

	var vrfInterface RouterVRFInterface
	if err := json.Unmarshal(resp.Data, &vrfInterface); err != nil {
		return nil, fmt.Errorf("failed to unmarshal interface response: %w", err)
	}

	return &vrfInterface, nil
}

func (c *V3Client) UpdateRouterVRFInterface(routerID int, vrfID int, interfaceID int, req UpdateRouterVRFInterfaceRequest) (*UpdateRouterVRFInterfaceResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/interfaces/%d", routerID, vrfID, interfaceID)

	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update interface %d on VRF %d, router %d: %w", interfaceID, vrfID, routerID, err)
	}

	var updateResp UpdateRouterVRFInterfaceResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update interface response: %w", err)
	}

	return &updateResp, nil
}

func (c *V3Client) DeleteRouterVRFInterface(routerID int, vrfID int, interfaceID int) error {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/interfaces/%d", routerID, vrfID, interfaceID)

	if _, err := c.del(path); err != nil {
		return fmt.Errorf("failed to delete interface %d on VRF %d, router %d: %w", interfaceID, vrfID, routerID, err)
	}

	return nil
}
