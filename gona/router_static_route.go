package gona

import (
	"encoding/json"
	"fmt"
)

type StaticRouteVia struct {
	NextHop     string `json:"nextHop,omitempty"`
	InterfaceID *int   `json:"interfaceId,omitempty"`
	TunnelID    *int   `json:"tunnelId,omitempty"`
	IPSecPeerID *int   `json:"ipSecPeerId,omitempty"`
}

type RouterStaticRoute struct {
	RouteID     int            `json:"staticRouteId"`
	Network     string         `json:"network"`
	Via         StaticRouteVia `json:"via"`
	Description string         `json:"description,omitempty"`
	Distance    *int           `json:"distance,omitempty"`
}

type CreateRouterStaticRouteRequest struct {
	Network     string         `json:"network"`
	Via         StaticRouteVia `json:"via"`
	Description string         `json:"description,omitempty"`
	Distance    *int           `json:"distance,omitempty"`
}

type CreateRouterStaticRouteResponse struct {
	RouteID int `json:"staticRouteId"`
}

type UpdateRouterStaticRouteRequest struct {
	Network     string         `json:"network"`
	Via         StaticRouteVia `json:"via"`
	Description string         `json:"description,omitempty"`
	Distance    *int           `json:"distance,omitempty"`
}

type UpdateRouterStaticRouteResponse struct {
	RouteID int `json:"staticRouteId"`
}

func (c *V3Client) ListRouterStaticRoutes(routerID int, vrfID int) ([]RouterStaticRoute, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/static-routes", routerID, vrfID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list static routes for VRF %d on router %d: %w", vrfID, routerID, err)
	}

	var routes []RouterStaticRoute
	if err := json.Unmarshal(resp.Data, &routes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal static routes response: %w", err)
	}

	return routes, nil
}

func (c *V3Client) GetRouterStaticRoute(routerID int, vrfID int, routeID int) (*RouterStaticRoute, error) {
	routes, err := c.ListRouterStaticRoutes(routerID, vrfID)
	if err != nil {
		return nil, err
	}

	for _, route := range routes {
		if route.RouteID == routeID {
			return &route, nil
		}
	}

	return nil, fmt.Errorf("static route %d not found for VRF %d on router %d", routeID, vrfID, routerID)
}

func (c *V3Client) CreateRouterStaticRoute(routerID int, vrfID int, req CreateRouterStaticRouteRequest) (*CreateRouterStaticRouteResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/static-routes", routerID, vrfID)

	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create static route for VRF %d on router %d: %w", vrfID, routerID, err)
	}

	var createResp CreateRouterStaticRouteResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create static route response: %w", err)
	}

	return &createResp, nil
}

func (c *V3Client) UpdateRouterStaticRoute(routerID int, vrfID int, routeID int, req UpdateRouterStaticRouteRequest) (*UpdateRouterStaticRouteResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/static-routes/%d", routerID, vrfID, routeID)

	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update static route %d for VRF %d on router %d: %w", routeID, vrfID, routerID, err)
	}

	var updateResp UpdateRouterStaticRouteResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update static route response: %w", err)
	}

	return &updateResp, nil
}

func (c *V3Client) DeleteRouterStaticRoute(routerID int, vrfID int, routeID int) error {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/vrfs/%d/static-routes/%d", routerID, vrfID, routeID)

	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("failed to delete static route %d for VRF %d on router %d: %w", routeID, vrfID, routerID, err)
	}

	return nil
}
