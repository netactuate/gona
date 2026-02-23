package gona

import (
	"encoding/json"
	"fmt"
)

type RouterVRFBGPConfig struct {
	LocalAsn  *int                   `json:"localAsn"`
	RouterID  string                 `json:"routerId"`
	Networks  []RouterVRFBGPNetwork  `json:"networks"`
	Neighbors []RouterVRFBGPNeighbor `json:"neighbors"`
}

type RouterVRFBGPNetwork struct {
	Subnet string `json:"subnet"`
}

type UpdateRouterVRFBGPRequest struct {
	Networks []RouterVRFBGPNetwork `json:"networks,omitempty"`
	ASN      *RouterVRFBGPASN      `json:"asn,omitempty"`
}

type RouterVRFBGPASN struct {
	Local *int `json:"local,omitempty"`
}

type UpdateRouterVRFBGPResponse struct {
	LocalAsn  *int                   `json:"localAsn"`
	RouterID  string                 `json:"routerId"`
	Networks  []RouterVRFBGPNetwork  `json:"networks"`
	Neighbors []RouterVRFBGPNeighbor `json:"neighbors"`
}

func (c *V3Client) GetRouterVRFBGP(routerID int, vrfID int) (*RouterVRFBGPConfig, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/bgp", routerID, vrfID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get BGP config for VRF %d on router %d: %w", vrfID, routerID, err)
	}

	var bgpConfig RouterVRFBGPConfig
	if err := json.Unmarshal(resp.Data, &bgpConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal BGP config response: %w", err)
	}

	return &bgpConfig, nil
}

func (c *V3Client) UpdateRouterVRFBGP(routerID int, vrfID int, req UpdateRouterVRFBGPRequest) (*UpdateRouterVRFBGPResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/bgp", routerID, vrfID)

	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update BGP config for VRF %d on router %d: %w", vrfID, routerID, err)
	}

	var updateResp UpdateRouterVRFBGPResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update BGP config response: %w", err)
	}

	return &updateResp, nil
}
