package gona

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// DedicatedDeviceFilterOptions contains optional filters for available
// dedicated devices.
type DedicatedDeviceFilterOptions struct {
	PerPage    *int
	NIC        string
	CPUType    string
	GPUType    string
	DiskType   string
	Cores      string
	RAMMB      string
	DiskMIB    string
	DCName     string
	RegionName string
}

// DedicatedDevice is a dedicated device returned by the filter endpoint.
type DedicatedDevice map[string]json.RawMessage

// DedicatedLocation is a location where dedicated servers can be deployed.
type DedicatedLocation struct {
	ShortName      string `json:"short_name"`
	PubDescription string `json:"pub_description"`
	LocationID     int    `json:"location_id"`
}

// DedicatedPlan is a dedicated server plan available at a location.
type DedicatedPlan map[string]json.RawMessage

// DedicatedServerBuildRequest contains deployment fields for a dedicated server.
type DedicatedServerBuildRequest struct {
	FQDN         string  `json:"fqdn"`
	Profile      int     `json:"profile"`
	DiskLayout   *int    `json:"disklayout,omitempty"`
	RootPassword *string `json:"root_password,omitempty"`
	SSHKey       *string `json:"ssh_key,omitempty"`
	SSHKeyID     *int    `json:"ssh_key_id,omitempty"`
	BuildScript  *string `json:"build_script,omitempty"`
}

// DedicatedServerActionRequest contains optional fields for dedicated server
// power and delete actions.
type DedicatedServerActionRequest struct {
	Force    *bool   `json:"force,omitempty"`
	Password *string `json:"password,omitempty"`
}

// DedicatedIPv4ReverseRequest contains the IPv4 reverse DNS update fields.
type DedicatedIPv4ReverseRequest struct {
	MBPkgID int    `json:"mbpkgid"`
	ID      int    `json:"id"`
	Reverse string `json:"reverse"`
}

// DedicatedPowerStatus is the power status returned for a dedicated server.
type DedicatedPowerStatus map[string]json.RawMessage

// FilterDedicatedDevices returns available dedicated devices matching the
// optional filters.
func (c *Client) FilterDedicatedDevices(opts DedicatedDeviceFilterOptions) ([]DedicatedDevice, error) {
	path := "dedicated/filter-dedicated-devices"
	values := url.Values{}
	if opts.PerPage != nil {
		values.Set("per_page", strconv.Itoa(*opts.PerPage))
	}
	setQuery(values, "nic", opts.NIC)
	setQuery(values, "cpu_type", opts.CPUType)
	setQuery(values, "gpu_type", opts.GPUType)
	setQuery(values, "disk_type", opts.DiskType)
	setQuery(values, "cores", opts.Cores)
	setQuery(values, "ram_mb", opts.RAMMB)
	setQuery(values, "disk_mib", opts.DiskMIB)
	setQuery(values, "dc_name", opts.DCName)
	setQuery(values, "region_name", opts.RegionName)
	if encoded := values.Encode(); encoded != "" {
		path += "?" + encoded
	}

	// The rows sit at data.devices.paginator.data: a paginated list inside a named sub object,
	// alongside the column lists the portal uses. When the paginator is at the top of data
	// instead, the client's own unwrap has already reduced it to a bare list, so both arrive
	// here and both are accepted.
	var raw json.RawMessage
	if err := c.get(context.Background(), path, &raw); err != nil {
		return nil, fmt.Errorf("filter dedicated devices: %w", err)
	}

	var devices []DedicatedDevice
	if err := json.Unmarshal(raw, &devices); err == nil {
		return devices, nil
	}

	var envelope struct {
		Devices struct {
			Paginator struct {
				Data []DedicatedDevice `json:"data"`
			} `json:"paginator"`
		} `json:"devices"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		return nil, fmt.Errorf("filter dedicated devices unmarshal: %w", err)
	}
	return envelope.Devices.Paginator.Data, nil
}

// ListDedicatedLocations returns locations that support dedicated servers.
func (c *Client) ListDedicatedLocations() ([]DedicatedLocation, error) {
	var locations []DedicatedLocation
	if err := c.get(context.Background(), "dedicated/locations", &locations); err != nil {
		return nil, fmt.Errorf("list dedicated locations: %w", err)
	}
	return locations, nil
}

// ListDedicatedDeviceOSProfiles returns OS profiles compatible with a device.
func (c *Client) ListDedicatedDeviceOSProfiles(deviceID int, isBuyable *bool) ([]DedicatedOSProfile, error) {
	path := fmt.Sprintf("dedicated/os/device/%d", deviceID)
	values := url.Values{}
	if isBuyable != nil {
		if *isBuyable {
			values.Set("is_buyable", "1")
		} else {
			values.Set("is_buyable", "0")
		}
	}
	if encoded := values.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var profiles []DedicatedOSProfile
	if err := c.get(context.Background(), path, &profiles); err != nil {
		return nil, fmt.Errorf("list dedicated device %d OS profiles: %w", deviceID, err)
	}
	return profiles, nil
}

// ListDedicatedPlans returns dedicated server plans for a location.
func (c *Client) ListDedicatedPlans(locationID int) ([]DedicatedPlan, error) {
	var plans []DedicatedPlan
	path := fmt.Sprintf("dedicated/plans/%d", locationID)
	if err := c.get(context.Background(), path, &plans); err != nil {
		return nil, fmt.Errorf("list dedicated plans for location %d: %w", locationID, err)
	}
	return plans, nil
}

// DeployDedicatedServer deploys an already purchased dedicated server package.
func (c *Client) DeployDedicatedServer(mbPkgID int, req *DedicatedServerBuildRequest) (MetalBuild, error) {
	var build MetalBuild
	if err := c.postDedicatedJSON(fmt.Sprintf("dedicated/server/build/%d", mbPkgID), req, &build); err != nil {
		return MetalBuild{}, fmt.Errorf("deploy dedicated server %d: %w", mbPkgID, err)
	}
	return build, nil
}

// BuyDedicatedServer purchases a dedicated device without deploying it.
func (c *Client) BuyDedicatedServer(deviceID int) (MetalBuild, error) {
	var build MetalBuild
	path := fmt.Sprintf("dedicated/server/buy/%d", deviceID)
	if err := c.post(context.Background(), path, nil, &build); err != nil {
		return MetalBuild{}, fmt.Errorf("buy dedicated server device %d: %w", deviceID, err)
	}
	return build, nil
}

// BuyAndDeployDedicatedServer purchases a dedicated device and deploys it.
func (c *Client) BuyAndDeployDedicatedServer(deviceID int, req *DedicatedServerBuildRequest) (MetalBuild, error) {
	var build MetalBuild
	if err := c.postDedicatedJSON(fmt.Sprintf("dedicated/server/buy_build/%d", deviceID), req, &build); err != nil {
		return MetalBuild{}, fmt.Errorf("buy and deploy dedicated server device %d: %w", deviceID, err)
	}
	return build, nil
}

// UpdateDedicatedServerIPv4Reverse updates reverse DNS for a dedicated server IPv4 address.
func (c *Client) UpdateDedicatedServerIPv4Reverse(req *DedicatedIPv4ReverseRequest) error {
	if err := c.putDedicatedJSON("dedicated/server/ipv4_reverse", req, nil); err != nil {
		return fmt.Errorf("update dedicated server IPv4 reverse DNS: %w", err)
	}
	return nil
}

// SoftResetDedicatedServer asks the platform to soft reset a dedicated server.
func (c *Client) SoftResetDedicatedServer(mbPkgID int) error {
	return c.postDedicatedAction(mbPkgID, "soft-reset", nil)
}

// DeleteDedicatedServer deletes a dedicated server package.
func (c *Client) DeleteDedicatedServer(mbPkgID int, req *DedicatedServerActionRequest) error {
	return c.postDedicatedAction(mbPkgID, "delete", req)
}

// RebootDedicatedServer asks the platform to reboot a dedicated server.
func (c *Client) RebootDedicatedServer(mbPkgID int, req *DedicatedServerActionRequest) error {
	return c.postDedicatedAction(mbPkgID, "reboot", req)
}

// ShutdownDedicatedServer asks the platform to shut down a dedicated server.
func (c *Client) ShutdownDedicatedServer(mbPkgID int, req *DedicatedServerActionRequest) error {
	return c.postDedicatedAction(mbPkgID, "shutdown", req)
}

// StartDedicatedServer asks the platform to start a dedicated server.
func (c *Client) StartDedicatedServer(mbPkgID int, req *DedicatedServerActionRequest) error {
	return c.postDedicatedAction(mbPkgID, "start", req)
}

// GetDedicatedServerPowerStatus returns the power status for a dedicated server.
func (c *Client) GetDedicatedServerPowerStatus(mbPkgID int, req *DedicatedServerActionRequest) (DedicatedPowerStatus, error) {
	var status DedicatedPowerStatus
	if err := c.postDedicatedActionWithResponse(mbPkgID, "status", req, &status); err != nil {
		return nil, err
	}
	return status, nil
}

// ListDedicatedServers returns dedicated servers for the account.
func (c *Client) ListDedicatedServers() ([]SingleMetal, error) {
	var servers []SingleMetal
	if err := c.get(context.Background(), "dedicated/servers", &servers); err != nil {
		return nil, fmt.Errorf("list dedicated servers: %w", err)
	}
	return servers, nil
}

func setQuery(values url.Values, key, value string) {
	if value != "" {
		values.Set(key, value)
	}
}

func (c *Client) postDedicatedAction(mbPkgID int, action string, req *DedicatedServerActionRequest) error {
	return c.postDedicatedActionWithResponse(mbPkgID, action, req, nil)
}

func (c *Client) postDedicatedActionWithResponse(mbPkgID int, action string, req *DedicatedServerActionRequest, data interface{}) error {
	path := fmt.Sprintf("dedicated/server/%d/%s", mbPkgID, action)
	if action == "soft-reset" {
		path = fmt.Sprintf("dedicated/server/soft-reset/%d", mbPkgID)
	}
	var body []byte
	if req != nil {
		var err error
		body, err = json.Marshal(req)
		if err != nil {
			return fmt.Errorf("marshal dedicated server %s request: %w", action, err)
		}
	}
	if err := c.postJSON(context.Background(), path, body, data); err != nil {
		return fmt.Errorf("%s dedicated server %d: %w", action, mbPkgID, err)
	}
	return nil
}

func (c *Client) postDedicatedJSON(path string, req interface{}, data interface{}) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal dedicated server request: %w", err)
	}
	return c.postJSON(context.Background(), path, body, data)
}

func (c *Client) putDedicatedJSON(path string, req interface{}, data interface{}) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal dedicated server request: %w", err)
	}
	return c.putJSON(context.Background(), path, body, data)
}
