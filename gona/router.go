package gona

import (
	"encoding/json"
	"fmt"
	"time"
)

type Router struct {
	Name             string       `json:"name"`
	Description      *string      `json:"description"`
	ReadyOn          *time.Time   `json:"readyOn"`
	HasDefaultVRF    bool         `json:"hasDefaultVrf"`
	CanJoinMagicMesh bool         `json:"canJoinMagicMesh"`
	MeshID           *int         `json:"meshId"`
	Build            []BuildEvent `json:"build,omitempty"`
}

type BuildEvent struct {
	Text string    `json:"text"`
	Date time.Time `json:"date"`
}

type RouterConfig struct {
	DefaultVrfID int                        `json:"defaultVrfId"`
	Service      RouterService              `json:"service"`
	PrefixLists  []interface{}              `json:"prefixLists"`
	VRF          map[string]RouterVRFConfig `json:"vrf"`
	IPSec        interface{}                `json:"ipSec"`
	Metadata     RouterConfigMetadata       `json:"metadata"`
}

type FlexibleIPv4 string

type RouterConfigMetadata struct {
	Status           string          `json:"status"`
	Name             string          `json:"name"`
	UpdatedOn        *string         `json:"updatedOn"`
	Version          int             `json:"version"`
	IPv4Address      FlexibleIPv4    `json:"ipv4Address"`
	Location         *RouterLocation `json:"location"`
	HasDefaultVrf    bool            `json:"hasDefaultVrf"`
	MeshID           *int            `json:"meshId"`
	CanJoinMagicMesh bool            `json:"canJoinMagicMesh"`
}

type RouterLocation struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	Flag *string `json:"flag"`
}

type RouterService struct {
	NTP NTPService `json:"ntp"`
}

type NTPService struct {
	Enabled     bool          `json:"enabled"`
	InterfaceID *int          `json:"interfaceId"`
	Upstreams   []NTPUpstream `json:"upstreams"`
}

type NTPUpstream struct {
	Domain string `json:"domain"`
}

type UpdateRouterRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CreateRouterRequest struct {
	PackageID   int     `json:"packageId"`
	LocationID  int     `json:"locationId"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type CreateRouterResponse struct {
	RouterID int `json:"routerId"`
}

func (f *FlexibleIPv4) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexibleIPv4(s)
		return nil
	}
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*f = FlexibleIPv4(fmt.Sprintf("%d", i))
		return nil
	}

	return fmt.Errorf("IPv4Address must be either string or int")
}

func (c *V3Client) GetRouter(routerID int) (*Router, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d", routerID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get router %d: %w", routerID, err)
	}

	var router Router
	if err := json.Unmarshal(resp.Data, &router); err != nil {
		return nil, fmt.Errorf("failed to unmarshal router: %w", err)
	}

	return &router, nil
}

func (c *V3Client) GetRouterConfig(routerID int) (*RouterConfig, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config", routerID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get router config %d: %w", routerID, err)
	}

	var config RouterConfig
	if err := json.Unmarshal(resp.Data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal router config: %w", err)
	}

	return &config, nil
}

func (c *V3Client) CreateRouter(req *CreateRouterRequest) (*CreateRouterResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("create router request cannot be nil")
	}
	if req.PackageID == 0 {
		return nil, fmt.Errorf("packageId is required")
	}
	if req.LocationID == 0 {
		return nil, fmt.Errorf("locationId is required")
	}

	path := "/cloud-routing/routers"

	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create router: %w", err)
	}

	var createResp CreateRouterResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create response: %w", err)
	}

	return &createResp, nil
}

func (c *V3Client) UpdateRouter(routerID int, req *UpdateRouterRequest) (*Router, error) {
	if req == nil {
		return nil, fmt.Errorf("update router request cannot be nil")
	}

	path := fmt.Sprintf("/cloud-routing/routers/%d", routerID)

	resp, err := c.patch(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update router %d: %w", routerID, err)
	}

	var router Router
	if err := json.Unmarshal(resp.Data, &router); err != nil {
		return nil, fmt.Errorf("failed to unmarshal updated router: %w", err)
	}

	return &router, nil
}

func (c *V3Client) DeleteRouter(routerID int) error {
	path := fmt.Sprintf("/cloud-routing/routers/%d", routerID)

	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("failed to delete router %d: %w", routerID, err)
	}

	return nil
}

func (c *V3Client) WaitForRouterReady(routerID int) error {
	return c.waitForCondition(func() (bool, error) {
		router, err := c.GetRouter(routerID)
		if err != nil {
			return false, err
		}
		isReady := router.ReadyOn != nil
		c.debugLog("Router %d ready status: %v (readyOn: %v)", routerID, isReady, router.ReadyOn)
		return isReady, nil
	}, RouterWaitConfig)
}
