package gona

import (
	"encoding/json"
	"fmt"
)

type WireguardPeerAllowedIP struct {
	Network string `json:"network"`
}

type CreateRouterVRFInterfaceWireguardPeerRequest struct {
	AllowedIPs   []WireguardPeerAllowedIP `json:"allowedIps"`
	PublicKey    *string                  `json:"publicKey"`
	PreSharedKey *string                  `json:"preSharedKey"`
	Name         *string                  `json:"name"`
	Description  *string                  `json:"description"`
	Remote       *string                  `json:"remote"`
}

type CreateRouterVRFInterfaceWireguardPeerResponse struct {
	WireguardPeerID int `json:"wireguardPeerId"`
}

type RouterVRFInterfaceWireguardPeer struct {
	WireguardPeerID int                      `json:"wireguardPeerId"`
	AllowedIPs      []WireguardPeerAllowedIP `json:"allowedIps"`
	PublicKey       string                   `json:"publicKey"`
	PrivateKey      string                   `json:"privateKey"`
	PreSharedKey    *string                  `json:"preSharedKey"`
	Remote          *string                  `json:"remote"`
	Name            *string                  `json:"name"`
	Description     *string                  `json:"description"`
}

func (c *V3Client) CreateRouterVRFInterfaceWireguardPeer(routerID int, vrfID int, interfaceID int, req CreateRouterVRFInterfaceWireguardPeerRequest) (*CreateRouterVRFInterfaceWireguardPeerResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/interfaces/%d/wireguard-peers", routerID, vrfID, interfaceID)

	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create wireguard peer on interface %d in VRF %d on router %d: %w", interfaceID, vrfID, routerID, err)
	}

	var createResp CreateRouterVRFInterfaceWireguardPeerResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create wireguard peer response: %w", err)
	}

	return &createResp, nil
}

func (c *V3Client) GetRouterVRFInterfaceWireguardPeer(routerID int, vrfID int, interfaceID int, wireguardPeerID int) (*RouterVRFInterfaceWireguardPeer, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/interfaces/%d/wireguard-peers/%d", routerID, vrfID, interfaceID, wireguardPeerID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get wireguard peer %d on interface %d in VRF %d on router %d: %w", wireguardPeerID, interfaceID, vrfID, routerID, err)
	}

	var peer RouterVRFInterfaceWireguardPeer
	if err := json.Unmarshal(resp.Data, &peer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wireguard peer response: %w", err)
	}

	return &peer, nil
}

func (c *V3Client) DeleteRouterVRFInterfaceWireguardPeer(routerID int, vrfID int, interfaceID int, wireguardPeerID int) error {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/interfaces/%d/wireguard-peers/%d", routerID, vrfID, interfaceID, wireguardPeerID)

	if _, err := c.del(path); err != nil {
		return fmt.Errorf("failed to delete wireguard peer %d on interface %d in VRF %d on router %d: %w", wireguardPeerID, interfaceID, vrfID, routerID, err)
	}

	return nil
}
