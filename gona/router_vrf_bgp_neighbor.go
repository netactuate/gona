package gona

import (
	"encoding/json"
	"fmt"
)

type BGPNeighborSource struct {
	Address string `json:"address,omitempty"`
}

type BGPNeighborEnabledIPVersion struct {
	IPv4 bool `json:"ipv4"`
	IPv6 bool `json:"ipv6"`
}

type BGPNeighborASN struct {
	Remote int `json:"remote"`
}

type BGPNeighborRouteMapRule struct {
	PrefixListID       int    `json:"prefixListId"`
	Action             string `json:"action"` // "permit", "deny", or "next"
	SetLocalPreference *int   `json:"setLocalPreference,omitempty"`
	PrependLastAsn     *int   `json:"prependLastAsn,omitempty"`
}

type BGPNeighborRouteMap struct {
	DoDefaultDrop bool                      `json:"doDefaultDrop"`
	Rules         []BGPNeighborRouteMapRule `json:"rules,omitempty"`
}

type RouterVRFBGPNeighbor struct {
	NeighborID       int                         `json:"neighborId,omitempty"`
	Address          string                      `json:"address"`
	IsShutdown       bool                        `json:"isShutdown"`
	DoAsOverride     bool                        `json:"doAsOverride"`
	DoNextHelpSelf   bool                        `json:"doNextHelpSelf"`
	Source           *BGPNeighborSource          `json:"source,omitempty"`
	EnabledIPVersion BGPNeighborEnabledIPVersion `json:"enabledIpVersion"`
	EbgpMultihop     *int                        `json:"ebgpMultihop,omitempty"`
	ASN              BGPNeighborASN              `json:"asn"`
	MD5Secret        string                      `json:"md5Secret,omitempty"`
	Import           *BGPNeighborRouteMap        `json:"import,omitempty"`
	Export           *BGPNeighborRouteMap        `json:"export,omitempty"`
	Name             string                      `json:"name,omitempty"`
	Description      string                      `json:"description,omitempty"`
}

type CreateRouterVRFBGPNeighborRequest struct {
	Address          string                      `json:"address"`
	IsShutdown       bool                        `json:"isShutdown"`
	DoAsOverride     bool                        `json:"doAsOverride"`
	DoNextHelpSelf   bool                        `json:"doNextHelpSelf"`
	Source           *BGPNeighborSource          `json:"source,omitempty"`
	EnabledIPVersion BGPNeighborEnabledIPVersion `json:"enabledIpVersion"`
	EbgpMultihop     *int                        `json:"ebgpMultihop,omitempty"`
	ASN              BGPNeighborASN              `json:"asn"`
	MD5Secret        string                      `json:"md5Secret,omitempty"`
	Import           *BGPNeighborRouteMap        `json:"import,omitempty"`
	Export           *BGPNeighborRouteMap        `json:"export,omitempty"`
	Name             string                      `json:"name,omitempty"`
	Description      string                      `json:"description,omitempty"`
}

type CreateRouterVRFBGPNeighborResponse struct {
	NeighborID int `json:"neighborId"`
}

type UpdateRouterVRFBGPNeighborRequest struct {
	Address          string                      `json:"address"`
	IsShutdown       bool                        `json:"isShutdown"`
	DoAsOverride     bool                        `json:"doAsOverride"`
	DoNextHelpSelf   bool                        `json:"doNextHelpSelf"`
	Source           *BGPNeighborSource          `json:"source,omitempty"`
	EnabledIPVersion BGPNeighborEnabledIPVersion `json:"enabledIpVersion"`
	EbgpMultihop     *int                        `json:"ebgpMultihop,omitempty"`
	ASN              BGPNeighborASN              `json:"asn"`
	MD5Secret        string                      `json:"md5Secret,omitempty"`
	Import           *BGPNeighborRouteMap        `json:"import,omitempty"`
	Export           *BGPNeighborRouteMap        `json:"export,omitempty"`
	Name             string                      `json:"name,omitempty"`
	Description      string                      `json:"description,omitempty"`
}

type UpdateRouterVRFBGPNeighborResponse struct {
	NeighborID int `json:"neighborId"`
}

func (c *V3Client) CreateRouterVRFBGPNeighbor(routerID int, vrfID int, req CreateRouterVRFBGPNeighborRequest) (*CreateRouterVRFBGPNeighborResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/bgp/neighbors", routerID, vrfID)

	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create BGP neighbor for VRF %d on router %d: %w", vrfID, routerID, err)
	}

	var createResp CreateRouterVRFBGPNeighborResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create BGP neighbor response: %w", err)
	}

	return &createResp, nil
}

func (c *V3Client) GetRouterVRFBGPNeighbor(routerID int, vrfID int, neighborID int) (*RouterVRFBGPNeighbor, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/bgp/neighbors/%d", routerID, vrfID, neighborID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get BGP neighbor %d for VRF %d on router %d: %w", neighborID, vrfID, routerID, err)
	}

	var neighbor RouterVRFBGPNeighbor
	if err := json.Unmarshal(resp.Data, &neighbor); err != nil {
		return nil, fmt.Errorf("failed to unmarshal BGP neighbor response: %w", err)
	}

	return &neighbor, nil
}

func (c *V3Client) UpdateRouterVRFBGPNeighbor(routerID int, vrfID int, neighborID int, req UpdateRouterVRFBGPNeighborRequest) (*UpdateRouterVRFBGPNeighborResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/bgp/neighbors/%d", routerID, vrfID, neighborID)

	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update BGP neighbor %d for VRF %d on router %d: %w", neighborID, vrfID, routerID, err)
	}

	var updateResp UpdateRouterVRFBGPNeighborResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update BGP neighbor response: %w", err)
	}

	return &updateResp, nil
}

func (c *V3Client) DeleteRouterVRFBGPNeighbor(routerID int, vrfID int, neighborID int) error {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/bgp/neighbors/%d", routerID, vrfID, neighborID)

	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("failed to delete BGP neighbor %d for VRF %d on router %d: %w", neighborID, vrfID, routerID, err)
	}

	return nil
}
