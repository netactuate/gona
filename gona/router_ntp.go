package gona

import (
	"encoding/json"
	"fmt"
)

type RouterNTPUpstream struct {
	Domain string `json:"domain"`
}

type RouterNTPConfig struct {
	Enabled     bool                `json:"enabled"`
	InterfaceID *int                `json:"interfaceId"`
	Upstreams   []RouterNTPUpstream `json:"upstreams"`
}

type UpdateRouterNTPConfigRequest struct {
	Enabled     *bool               `json:"enabled"`
	InterfaceID *int                `json:"interfaceId"`
	Upstreams   []RouterNTPUpstream `json:"upstreams"`
}

type UpdateRouterNTPConfigResponse struct {
	RouterID int `json:"routerId"`
}

func (c *V3Client) GetRouterNTPConfig(routerID int) (*RouterNTPConfig, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/services/ntp", routerID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get NTP config for router %d: %w", routerID, err)
	}

	var ntpConfig RouterNTPConfig
	if err := json.Unmarshal(resp.Data, &ntpConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal NTP config response: %w", err)
	}

	return &ntpConfig, nil
}

func (c *V3Client) UpdateRouterNTPConfig(routerID int, req UpdateRouterNTPConfigRequest) (*UpdateRouterNTPConfigResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/services/ntp", routerID)

	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update NTP config for router %d: %w", routerID, err)
	}

	var updateResp UpdateRouterNTPConfigResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update NTP config response: %w", err)
	}

	return &updateResp, nil
}
