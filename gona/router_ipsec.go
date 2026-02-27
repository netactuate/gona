package gona

import (
	"encoding/json"
	"fmt"
)

type RouterIPSecConfig struct {
	IKEGroup RouterIPSecIKEGroup `json:"ikeGroup"`
	ESPGroup RouterIPSecESPGroup `json:"espGroup"`
}

type RouterIPSecIKEGroup struct {
	DoAutoRenegotiation bool   `json:"doAutoRenegotiation"`
	KeyExchangeVersion  int    `json:"keyExchangeVersion"`
	LifetimeSeconds     int    `json:"lifetimeSeconds"`
	DHGroupNumber       int    `json:"dhGroupNumber"`
	Encryption          string `json:"encryption"`
	Hash                string `json:"hash"`
	PRF                 string `json:"prf"`
}

type RouterIPSecESPGroup struct {
	LifetimeSeconds int    `json:"lifetimeSeconds"`
	Encryption      string `json:"encryption"`
	Hash            string `json:"hash"`
}

type UpdateRouterIPSecConfigRequest struct {
	IKEGroup RouterIPSecIKEGroup `json:"ikeGroup"`
	ESPGroup RouterIPSecESPGroup `json:"espGroup"`
}

type RouterVRFIPSecPeer struct {
	IPSecPeerID          int                          `json:"ipSecPeerId"`
	Name                 string                       `json:"name"`
	Description          *string                      `json:"description"`
	RemoteID             string                       `json:"remoteId"`
	PSKSecret            string                       `json:"pskSecret"`
	DoInitiateConnection bool                         `json:"doInitiateConnection"`
	PeerAddress          string                       `json:"peerAddress"`
	LocalID              string                       `json:"localId"`
	OverlayNetwork       RouterVRFIPSecOverlayNetwork `json:"overlayNetwork"`
}

type RouterVRFIPSecOverlayNetwork struct {
	IPv4 *string `json:"ipv4"`
	IPv6 *string `json:"ipv6"`
}

type CreateRouterVRFIPSecPeerRequest struct {
	Name                 string                              `json:"name"`
	Description          *string                             `json:"description,omitempty"`
	RemoteID             string                              `json:"remoteId"`
	PSKSecret            string                              `json:"pskSecret"`
	DoInitiateConnection bool                                `json:"doInitiateConnection"`
	PeerAddress          string                              `json:"peerAddress"`
	OverlayNetwork       CreateRouterVRFIPSecOverlayNetwork  `json:"overlayNetwork"`
}

type CreateRouterVRFIPSecOverlayNetwork struct {
	IPv4 *string `json:"ipv4,omitempty"`
	IPv6 *string `json:"ipv6,omitempty"`
}

type CreateRouterVRFIPSecPeerResponse struct {
	IPSecPeerID int `json:"ipSecPeerId"`
}

type UpdateRouterVRFIPSecPeerRequest = CreateRouterVRFIPSecPeerRequest
type UpdateRouterVRFIPSecPeerResponse = CreateRouterVRFIPSecPeerResponse

func (c *V3Client) GetRouterIPSecConfig(routerID int) (*RouterIPSecConfig, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/ipSec", routerID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get IPSec config for router %d: %w", routerID, err)
	}

	var config RouterIPSecConfig
	if err := json.Unmarshal(resp.Data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal IPSec config response: %w", err)
	}

	return &config, nil
}

func (c *V3Client) UpdateRouterIPSecConfig(routerID int, req UpdateRouterIPSecConfigRequest) error {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/ipSec", routerID)

	_, err := c.put(path, req)
	if err != nil {
		return fmt.Errorf("failed to update IPSec config for router %d: %w", routerID, err)
	}

	return nil
}

func (c *V3Client) ListRouterVRFIPSecPeers(routerID, vrfID int) ([]RouterVRFIPSecPeer, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/ipSec/peers", routerID, vrfID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list IPSec peers for router %d VRF %d: %w", routerID, vrfID, err)
	}

	var peers []RouterVRFIPSecPeer
	if err := json.Unmarshal(resp.Data, &peers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal IPSec peers: %w", err)
	}

	return peers, nil
}

func (c *V3Client) GetRouterVRFIPSecPeer(routerID, vrfID, peerID int) (*RouterVRFIPSecPeer, error) {
	peers, err := c.ListRouterVRFIPSecPeers(routerID, vrfID)
	if err != nil {
		return nil, err
	}

	for _, peer := range peers {
		if peer.IPSecPeerID == peerID {
			return &peer, nil
		}
	}

	return nil, fmt.Errorf("IPSec peer %d not found in router %d VRF %d", peerID, routerID, vrfID)
}

func (c *V3Client) CreateRouterVRFIPSecPeer(routerID, vrfID int, req CreateRouterVRFIPSecPeerRequest) (*CreateRouterVRFIPSecPeerResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/ipSec/peers", routerID, vrfID)

	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create IPSec peer for router %d VRF %d: %w", routerID, vrfID, err)
	}

	var createResp CreateRouterVRFIPSecPeerResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create IPSec peer response: %w", err)
	}

	return &createResp, nil
}

func (c *V3Client) UpdateRouterVRFIPSecPeer(routerID, vrfID, peerID int, req UpdateRouterVRFIPSecPeerRequest) (*UpdateRouterVRFIPSecPeerResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/ipSec/peers/%d", routerID, vrfID, peerID)

	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update IPSec peer %d for router %d VRF %d: %w", peerID, routerID, vrfID, err)
	}

	var updateResp UpdateRouterVRFIPSecPeerResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update IPSec peer response: %w", err)
	}

	return &updateResp, nil
}

func (c *V3Client) DeleteRouterVRFIPSecPeer(routerID, vrfID, peerID int) error {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/ipSec/peers/%d", routerID, vrfID, peerID)

	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("failed to delete IPSec peer %d for router %d VRF %d: %w", peerID, routerID, vrfID, err)
	}

	return nil
}
