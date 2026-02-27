package gona

import (
	"encoding/json"
	"fmt"
)

type RouterVRFTunnel struct {
	TunnelID        int                     `json:"tunnelId"`
	Name            string                  `json:"name"`
	Description     *string                 `json:"description"`
	IPKey           int                     `json:"ipKey"`
	MTU             string                  `json:"mtu"`
	IPv4CIDR        *string                 `json:"ipv4Cidr"`
	IPv6CIDR        *string                 `json:"ipv6Cidr"`
	IPVersion       int                     `json:"ipVersion"`
	EndpointAddress RouterVRFTunnelEndpoint `json:"endpointAddress"`
}

type RouterVRFTunnelEndpoint struct {
	Source string `json:"source"`
	Remote string `json:"remote"`
}

type CreateRouterVRFTunnelRequest struct {
	IPKey           int                           `json:"ipKey"`
	Name            string                        `json:"name"`
	Description     *string                       `json:"description"`
	MTU             int                           `json:"mtu"`
	IPv4CIDR        *string                       `json:"ipv4Cidr"`
	IPv6CIDR        *string                       `json:"ipv6Cidr"`
	EndpointAddress CreateRouterVRFTunnelEndpoint `json:"endpointAddress"`
}

type CreateRouterVRFTunnelEndpoint struct {
	Remote string `json:"remote"`
}

type CreateRouterVRFTunnelResponse struct {
	TunnelID int `json:"tunnelId"`
}

type UpdateRouterVRFTunnelRequest struct {
	IPKey           int                           `json:"ipKey"`
	Name            string                        `json:"name"`
	Description     *string                       `json:"description"`
	MTU             int                           `json:"mtu"`
	IPv4CIDR        *string                       `json:"ipv4Cidr"`
	IPv6CIDR        *string                       `json:"ipv6Cidr"`
	EndpointAddress UpdateRouterVRFTunnelEndpoint `json:"endpointAddress"`
}

type UpdateRouterVRFTunnelEndpoint struct {
	Remote string `json:"remote"`
}

type UpdateRouterVRFTunnelResponse struct {
	TunnelID int `json:"tunnelId"`
}

func (c *V3Client) GetRouterVRFTunnel(routerID int, vrfID int, tunnelID int) (*RouterVRFTunnel, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/tunnels/%d", routerID, vrfID, tunnelID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get tunnel %d for VRF %d on router %d: %w", tunnelID, vrfID, routerID, err)
	}

	var tunnel RouterVRFTunnel
	if err := json.Unmarshal(resp.Data, &tunnel); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tunnel response: %w", err)
	}

	return &tunnel, nil
}

func (c *V3Client) CreateRouterVRFTunnel(routerID int, vrfID int, req CreateRouterVRFTunnelRequest) (*CreateRouterVRFTunnelResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/tunnels", routerID, vrfID)

	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create tunnel for VRF %d on router %d: %w", vrfID, routerID, err)
	}

	var createResp CreateRouterVRFTunnelResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create tunnel response: %w", err)
	}

	return &createResp, nil
}

func (c *V3Client) UpdateRouterVRFTunnel(routerID int, vrfID int, tunnelID int, req UpdateRouterVRFTunnelRequest) (*UpdateRouterVRFTunnelResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/tunnels/%d", routerID, vrfID, tunnelID)

	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update tunnel %d for VRF %d on router %d: %w", tunnelID, vrfID, routerID, err)
	}

	var updateResp UpdateRouterVRFTunnelResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update tunnel response: %w", err)
	}

	return &updateResp, nil
}

func (c *V3Client) DeleteRouterVRFTunnel(routerID int, vrfID int, tunnelID int) error {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/tunnels/%d", routerID, vrfID, tunnelID)

	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("failed to delete tunnel %d for VRF %d on router %d: %w", tunnelID, vrfID, routerID, err)
	}

	return nil
}
