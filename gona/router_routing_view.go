package gona

import (
	"encoding/json"
	"fmt"
)

// RouterRoutingViewRequest selects live routing views to return for a VRF.
type RouterRoutingViewRequest struct {
	Views []RouterRoutingViewSelector `json:"views"`
}

// RouterRoutingViewSelector identifies one routing view and optional filter.
type RouterRoutingViewSelector struct {
	ID        string `json:"id,omitempty"`
	IPVersion int    `json:"ipVersion"`
	Name      string `json:"name"`
	Filter    string `json:"filter,omitempty"`
}

// GetRouterRoutingViews returns the requested live routing views for a VRF.
func (c *V3Client) GetRouterRoutingViews(routerID, vrfID int, req *RouterRoutingViewRequest) (json.RawMessage, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/view/routing/%d", routerID, vrfID)
	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("get routing views for VRF %d on router %d: %w", vrfID, routerID, err)
	}
	return resp.Data, nil
}

// GetRouterRoutingOverview returns the live routing overview for a VRF.
func (c *V3Client) GetRouterRoutingOverview(routerID, vrfID int) (json.RawMessage, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/view/routing/%d/overview", routerID, vrfID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get routing overview for VRF %d on router %d: %w", vrfID, routerID, err)
	}
	return resp.Data, nil
}
