package gona

import (
	"encoding/json"
	"fmt"
	"net"
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
	NTP RouterNTPConfig `json:"ntp"`
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
	var i uint32
	if err := json.Unmarshal(data, &i); err == nil {
		ip := net.IP{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
		*f = FlexibleIPv4(ip.String())
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

// WaitForRouterReady waits with the package default. Prefer WaitForRouterReadyTimeout so the
// caller can supply the resource's own timeout.
//
// A cloud router provisions in about five minutes. A wait that runs much beyond that is
// sitting on a STALLED build, not a slow one, and should fail rather than keep waiting.
func (c *V3Client) WaitForRouterReady(routerID int) error {
	return c.WaitForRouterReadyTimeout(routerID, RouterWaitConfig.Timeout)
}

// WaitForRouterReadyTimeout waits up to the supplied timeout.
func (c *V3Client) WaitForRouterReadyTimeout(routerID int, timeout time.Duration) error {
	cfg := RouterWaitConfig
	if timeout > 0 {
		cfg.Timeout = timeout
	}
	// Watch PROGRESS, not just elapsed time.
	//
	// A cloud router build reports seven timestamped steps, and a healthy one completes in
	// about five minutes with each step a minute or two apart. When the build stalls, the
	// remaining steps simply keep a null date forever: there is no failure state and
	// readyOn never gets set, so a stalled build and a slow one are indistinguishable to
	// any client that only watches the clock.
	//
	// That cost real time on 2026-09-10, when three routers stalled at "Cloud Router
	// configured" and a test harness sat on one for over three hours before anyone looked.
	// Waiting longer never produces information.
	//
	// So track how many steps have completed. If none completes within RouterStallAfter,
	// give up and say WHICH step is stuck, which is the one fact worth reporting.
	var lastCount int
	lastProgress := time.Now()

	err := c.waitForCondition(func() (bool, error) {
		router, err := c.GetRouter(routerID)
		if err != nil {
			return false, err
		}
		if router.ReadyOn != nil {
			return true, nil
		}

		done, pending := 0, ""
		for _, e := range router.Build {
			if !e.Date.IsZero() {
				done++
			} else if pending == "" {
				pending = e.Text
			}
		}
		if done > lastCount {
			lastCount = done
			lastProgress = time.Now()
		}
		c.debugLog("Router %d build: %d steps done, waiting on %q", routerID, done, pending)

		if stalled := time.Since(lastProgress); stalled > RouterStallAfter {
			return false, fmt.Errorf(
				"router %d build has made no progress for %s: %d of %d steps complete, stuck on %q. "+
					"A healthy build finishes in about five minutes, so this is a stalled build rather than a slow one",
				routerID, stalled.Round(time.Second), done, len(router.Build), pending)
		}
		return false, nil
	}, cfg)

	return err
}
