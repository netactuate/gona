package gona

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"strconv"
)

// BoolPtr returns a pointer to v.
func BoolPtr(v bool) *bool {
	return &v
}

// IntPtr returns a pointer to v.
func IntPtr(v int) *int {
	return &v
}

// StringPtr returns a pointer to v.
func StringPtr(v string) *string {
	return &v
}

// BandwidthStatsRange is the data envelope returned for a bandwidth range.
type BandwidthStatsRange json.RawMessage

// CloudLocation is a cloud deployment location.
type CloudLocation struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Location  string  `json:"location"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	IATACode  string  `json:"iata_code"`
	Flag      string  `json:"flag"`
	Latitude  *string `json:"latitude"`
	Longitude *string `json:"longitude"`
}

// Kernel describes a boot kernel option.
type Kernel struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

// ServerBuildStatus is the status of an asynchronous server build.
type ServerBuildStatus struct {
	ID       int             `json:"id"`
	Status   string          `json:"status"`
	Percent  int             `json:"percent"`
	Response string          `json:"response"`
	Raw      json.RawMessage `json:"-"`
}

// UnmarshalJSON preserves the raw build status payload.
func (s *ServerBuildStatus) UnmarshalJSON(data []byte) error {
	type alias ServerBuildStatus
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0], data...)
	*s = ServerBuildStatus(a)
	return nil
}

// ServerIPAddress is an IPv4 or IPv6 address attached to a server.
type ServerIPAddress struct {
	ID      int             `json:"id"`
	IP      string          `json:"ip"`
	Reverse *string         `json:"reverse"`
	Netmask *string         `json:"netmask"`
	Gateway *string         `json:"gateway"`
	Type    *string         `json:"type"`
	Primary *int            `json:"primary"`
	Raw     json.RawMessage `json:"-"`
}

// UnmarshalJSON preserves address fields that are not modelled yet.
func (a *ServerIPAddress) UnmarshalJSON(data []byte) error {
	type alias ServerIPAddress
	var v alias
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	v.Raw = append(v.Raw[:0], data...)
	*a = ServerIPAddress(v)
	return nil
}

// ServerStatus is a server status payload.
type ServerStatus struct {
	Status string          `json:"status"`
	State  string          `json:"state"`
	Raw    json.RawMessage `json:"-"`
}

// UnmarshalJSON preserves the raw server status payload.
func (s *ServerStatus) UnmarshalJSON(data []byte) error {
	type alias ServerStatus
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0], data...)
	*s = ServerStatus(a)
	return nil
}

// RescueStartRequest contains the credentials for starting rescue mode.
type RescueStartRequest struct {
	RescuePass string  `json:"rescue_pass"`
	Password   *string `json:"password,omitempty"`
}

// ResetRootPasswordRequest contains a new root password.
type ResetRootPasswordRequest struct {
	RootPass string  `json:"rootpass"`
	Password *string `json:"password,omitempty"`
}

// UpdateServerOptionsRequest updates mutable server options.
type UpdateServerOptionsRequest struct {
	FQDN        *string `json:"fqdn,omitempty"`
	Autorescue  *int    `json:"autorescue,omitempty"`
	Description *string `json:"description,omitempty"`
	VCPUs       *int    `json:"vcpus,omitempty"`
	Boot        *string `json:"boot,omitempty"`
	KernelID    *int    `json:"kernel_id,omitempty"`
}

// ReverseDNSRequest updates reverse DNS for an address.
type ReverseDNSRequest struct {
	Reverse string `json:"reverse"`
}

// CreateUsageContractRequest creates a usage contract for an account.
type CreateUsageContractRequest struct {
	MBID int `json:"mb_id"`
}

// AttemptSSHRequest contains credentials used to test SSH access.
type AttemptSSHRequest struct {
	MBPkgID  int    `json:"mbpkgid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// BindFirewallSetRequest binds a firewall set to a server interface.
type BindFirewallSetRequest struct {
	FirewallSetID int `json:"firewall_set_id"`
	InterfaceID   int `json:"interface_id"`
	SetPriority   int `json:"set_priority"`
}

// ReplaceImageRequest selects the image that replaces an existing image.
type ReplaceImageRequest struct {
	ReplaceID int `json:"replace_id"`
}

// ScalingOptionsRequest filters scaling options for a server.
type ScalingOptionsRequest struct {
	IncludeCurrentPlan *bool
	MinRAM             *int
	MaxRAM             *int
	MinCPUs            *int
	MaxCPUs            *int
}

// DeploySizesRequest filters deploy sizes for a location.
type DeploySizesRequest struct {
	MinCPU *int
	MinRAM *int
}

// GetBandwidthStats returns bandwidth statistics for mbpkgID. When date is not
// empty it is sent as the API date query parameter.
func (c *Client) GetBandwidthStats(mbpkgID int, date string) (json.RawMessage, error) {
	path := fmt.Sprintf("cloud/bw_stats/%d", mbpkgID)
	if date != "" {
		q := url.Values{}
		q.Set("date", date)
		path += "?" + q.Encode()
	}
	return c.getAPIRaw(path)
}

// GetBandwidthStatsRange returns bandwidth statistics over the API's default range.
func (c *Client) GetBandwidthStatsRange(mbpkgID int) (json.RawMessage, error) {
	return c.getAPIRaw(fmt.Sprintf("cloud/bw_stats_range/%d", mbpkgID))
}

// GetImagesProvisioningJobsCount returns the count of image provisioning jobs.
func (c *Client) GetImagesProvisioningJobsCount() (json.RawMessage, error) {
	return c.getAPIRaw("cloud/images-provisioning-jobs-count")
}

// GetBaseImages returns base images available for server deployment.
func (c *Client) GetBaseImages() ([]Image, error) {
	var images []Image
	if err := c.get(context.Background(), "cloud/images/base", &images); err != nil {
		return nil, err
	}
	return images, nil
}

// GetPrivateImages returns private images available to the account.
func (c *Client) GetPrivateImages() ([]Image, error) {
	var images []Image
	if err := c.get(context.Background(), "cloud/images/private", &images); err != nil {
		return nil, err
	}
	return images, nil
}

// ReplaceImage replaces image id with the image identified by req.ReplaceID.
func (c *Client) ReplaceImage(id int, req *ReplaceImageRequest) (json.RawMessage, error) {
	return c.postJSONRaw(fmt.Sprintf("cloud/images/%d/replace_image", id), req)
}

// GetIPLimits returns IP limits for a server.
func (c *Client) GetIPLimits(mbpkgID int) (json.RawMessage, error) {
	return c.getAPIRaw(fmt.Sprintf("cloud/iplimits/%d", mbpkgID))
}

// GetCloudExtras returns optional extras for a server package.
func (c *Client) GetCloudExtras(mbpkgID int) (json.RawMessage, error) {
	return c.getAPIRaw(fmt.Sprintf("cloud/extras/%d", mbpkgID))
}

// UpdateCloudIPv4ReverseDNS updates reverse DNS for an IPv4 address.
func (c *Client) UpdateCloudIPv4ReverseDNS(id int, reverse string) error {
	_, err := c.putJSONRaw(fmt.Sprintf("cloud/ipv4/%d", id), &ReverseDNSRequest{Reverse: reverse})
	return err
}

// UpdateCloudIPv6ReverseDNS updates reverse DNS for an IPv6 address.
func (c *Client) UpdateCloudIPv6ReverseDNS(id int, reverse string) error {
	_, err := c.putJSONRaw(fmt.Sprintf("cloud/ipv6/%d", id), &ReverseDNSRequest{Reverse: reverse})
	return err
}

// GetKernels returns boot kernels available for cloud servers.
func (c *Client) GetKernels() ([]Kernel, error) {
	var kernels []Kernel
	if err := c.get(context.Background(), "cloud/kernels", &kernels); err != nil {
		return nil, err
	}
	return kernels, nil
}

// GetCloudLocation returns a single cloud location.
func (c *Client) GetCloudLocation(id int) (CloudLocation, error) {
	var location CloudLocation
	if err := c.get(context.Background(), fmt.Sprintf("cloud/locations/%d", id), &location); err != nil {
		return CloudLocation{}, err
	}
	return location, nil
}

// GetCloudPool returns a single cloud pool.
func (c *Client) GetCloudPool(cloudPoolID int) (CloudPool, error) {
	var pool CloudPool
	if err := c.get(context.Background(), fmt.Sprintf("cloud/pools/%d", cloudPoolID), &pool); err != nil {
		return CloudPool{}, err
	}
	return pool, nil
}

// GetScalingOptions returns scaling options for a server.
func (c *Client) GetScalingOptions(mbpkgID int, req *ScalingOptionsRequest) (json.RawMessage, error) {
	q := url.Values{}
	if req != nil {
		if req.IncludeCurrentPlan != nil {
			q.Set("include_current_plan", strconv.FormatBool(*req.IncludeCurrentPlan))
		}
		if req.MinRAM != nil {
			q.Set("min_ram", strconv.Itoa(*req.MinRAM))
		}
		if req.MaxRAM != nil {
			q.Set("max_ram", strconv.Itoa(*req.MaxRAM))
		}
		if req.MinCPUs != nil {
			q.Set("min_cpus", strconv.Itoa(*req.MinCPUs))
		}
		if req.MaxCPUs != nil {
			q.Set("max_cpus", strconv.Itoa(*req.MaxCPUs))
		}
	}
	path := fmt.Sprintf("cloud/scaling/%d", mbpkgID)
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}
	return c.getAPIRaw(path)
}

// GetServerBuildStatus returns the status of a server build.
func (c *Client) GetServerBuildStatus(buildID int) (ServerBuildStatus, error) {
	var status ServerBuildStatus
	if err := c.get(context.Background(), fmt.Sprintf("cloud/server/build_status/%d", buildID), &status); err != nil {
		return ServerBuildStatus{}, err
	}
	return status, nil
}

// GetServerDeploymentInfo returns deployment metadata. When contractType is
// not empty it is sent as the contract_type query parameter.
func (c *Client) GetServerDeploymentInfo(contractType string) (json.RawMessage, error) {
	path := "cloud/server/deploy/info"
	if contractType != "" {
		q := url.Values{}
		q.Set("contract_type", contractType)
		path += "?" + q.Encode()
	}
	return c.getAPIRaw(path)
}

// GetServerVNCStatus returns VNC status for a server.
func (c *Client) GetServerVNCStatus(mbpkgID int) (json.RawMessage, error) {
	return c.getAPIRaw(fmt.Sprintf("cloud/server/vnc-status/%d", mbpkgID))
}

// UpdateServerOptions updates mutable options for a server.
func (c *Client) UpdateServerOptions(mbpkgID int, req *UpdateServerOptionsRequest) (json.RawMessage, error) {
	return c.putJSONRaw(fmt.Sprintf("cloud/options/%d", mbpkgID), req)
}

// DeleteServerWithOptions deletes a server using the registered delete path.
func (c *Client) DeleteServerWithOptions(mbpkgID int, req *DeleteServerRequest) (DeleteServerResponse, error) {
	var resp DeleteServerResponse
	body, err := json.Marshal(req)
	if err != nil {
		return DeleteServerResponse{}, fmt.Errorf("encode delete server request: %w", err)
	}
	if err := c.postJSON(context.Background(), fmt.Sprintf("cloud/server/%d/delete", mbpkgID), body, &resp); err != nil {
		return DeleteServerResponse{}, err
	}
	return resp, nil
}

// RunServerFSCK starts a filesystem check for a server.
func (c *Client) RunServerFSCK(mbpkgID int) error {
	return c.post(context.Background(), fmt.Sprintf("cloud/server/%d/fsck", mbpkgID), nil, nil)
}

// GetServerIPv4 returns IPv4 addresses attached to a server.
func (c *Client) GetServerIPv4(mbpkgID int) ([]ServerIPAddress, error) {
	var addresses []ServerIPAddress
	if err := c.get(context.Background(), fmt.Sprintf("cloud/server/%d/ipv4", mbpkgID), &addresses); err != nil {
		return nil, err
	}
	return addresses, nil
}

// GetServerIPv6 returns IPv6 addresses attached to a server.
func (c *Client) GetServerIPv6(mbpkgID int) ([]ServerIPAddress, error) {
	var addresses []ServerIPAddress
	if err := c.get(context.Background(), fmt.Sprintf("cloud/server/%d/ipv6", mbpkgID), &addresses); err != nil {
		return nil, err
	}
	return addresses, nil
}

// ListServerJobs returns queued jobs for a server.
func (c *Client) ListServerJobs(mbpkgID int) ([]JobStatus, error) {
	var jobs []JobStatus
	if err := c.get(context.Background(), fmt.Sprintf("cloud/server/%d/jobs", mbpkgID), &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetServerJob returns a queued job for a server.
func (c *Client) GetServerJob(mbpkgID, jobID int) (JobStatus, error) {
	var job JobStatus
	if err := c.get(context.Background(), fmt.Sprintf("cloud/server/%d/jobs/%d", mbpkgID, jobID), &job); err != nil {
		return JobStatus{}, err
	}
	return job, nil
}

// ReconfigureServerNetwork starts network reconfiguration for a server.
func (c *Client) ReconfigureServerNetwork(mbpkgID int) error {
	return c.post(context.Background(), fmt.Sprintf("cloud/server/%d/netconfig", mbpkgID), nil, nil)
}

// GetServerNetworkIPs returns network IPs attached to a server.
func (c *Client) GetServerNetworkIPs(mbpkgID int) (json.RawMessage, error) {
	return c.getAPIRaw(fmt.Sprintf("cloud/server/%d/networkips", mbpkgID))
}

// ResetServerRootPassword resets the root password for a server.
func (c *Client) ResetServerRootPassword(mbpkgID int, req *ResetRootPasswordRequest) (json.RawMessage, error) {
	return c.postJSONRaw(fmt.Sprintf("cloud/server/%d/password", mbpkgID), req)
}

// RebootServer reboots a server.
func (c *Client) RebootServer(mbpkgID int) error {
	return c.RebootServerWithOptions(mbpkgID, nil)
}

// RebootServerWithOptions reboots a server with optional force behaviour.
func (c *Client) RebootServerWithOptions(mbpkgID int, req *ServerActionRequest) error {
	_, err := c.postJSONRaw(fmt.Sprintf("cloud/server/%d/reboot", mbpkgID), req)
	return err
}

// StartServerRescue starts rescue mode for a server.
func (c *Client) StartServerRescue(mbpkgID int, req *RescueStartRequest) (json.RawMessage, error) {
	return c.postJSONRaw(fmt.Sprintf("cloud/server/%d/rescue_start", mbpkgID), req)
}

// StopServerRescue stops rescue mode for a server.
func (c *Client) StopServerRescue(mbpkgID int) error {
	return c.post(context.Background(), fmt.Sprintf("cloud/server/%d/rescue_stop", mbpkgID), nil, nil)
}

// GetServerBGPSessions returns BGP sessions for a server. When groupType is not
// empty it is sent as the group_type query parameter.
func (c *Client) GetServerBGPSessions(mbpkgID int, groupType string) (json.RawMessage, error) {
	path := fmt.Sprintf("cloud/server/%d/sessions", mbpkgID)
	if groupType != "" {
		q := url.Values{}
		q.Set("group_type", groupType)
		path += "?" + q.Encode()
	}
	return c.getAPIRaw(path)
}

// ShutdownServerWithOptions shuts down a server with optional force behaviour.
func (c *Client) ShutdownServerWithOptions(mbpkgID int, req *ServerActionRequest) error {
	_, err := c.postJSONRaw(fmt.Sprintf("cloud/server/%d/shutdown", mbpkgID), req)
	return err
}

// StartServerWithOptions starts a server with optional force behaviour.
func (c *Client) StartServerWithOptions(mbpkgID int, req *ServerActionRequest) error {
	_, err := c.postJSONRaw(fmt.Sprintf("cloud/server/%d/start", mbpkgID), req)
	return err
}

// GetServerStatus returns status for a server.
func (c *Client) GetServerStatus(mbpkgID int) (ServerStatus, error) {
	var status ServerStatus
	if err := c.get(context.Background(), fmt.Sprintf("cloud/server/%d/status", mbpkgID), &status); err != nil {
		return ServerStatus{}, err
	}
	return status, nil
}

// StartServerVNC starts a VNC session for a server.
func (c *Client) StartServerVNC(mbpkgID int) (json.RawMessage, error) {
	return c.postJSONRaw(fmt.Sprintf("cloud/server/%d/vnc", mbpkgID), nil)
}

// GetServerMonthlyBandwidth returns monthly bandwidth data for a server.
func (c *Client) GetServerMonthlyBandwidth(mbpkgID int) (json.RawMessage, error) {
	return c.getAPIRaw(fmt.Sprintf("cloud/servermonthlybw/%d", mbpkgID))
}

// AttemptSSHConnection attempts an SSH login to a server.
func (c *Client) AttemptSSHConnection(req *AttemptSSHRequest) (json.RawMessage, error) {
	return c.postJSONRaw("cloud/servers/attempt-ssh", req)
}

// GetCurrentServer returns the server associated with the caller IP.
func (c *Client) GetCurrentServer() (Server, error) {
	var server Server
	if err := c.get(context.Background(), "cloud/servers/current", &server); err != nil {
		return Server{}, err
	}
	return server, nil
}

// GetUnprovisionedPackages returns unprovisioned packages available to build.
func (c *Client) GetUnprovisionedPackages() (json.RawMessage, error) {
	return c.getAPIRaw("cloud/servers/unprovisioned")
}

// GetCloudServersUsageInfo returns usage stats for cloud servers.
func (c *Client) GetCloudServersUsageInfo() (json.RawMessage, error) {
	return c.getAPIRaw("cloud/servers/usage/info")
}

// GetVirtualServerContract returns contract data for a virtual server.
func (c *Client) GetVirtualServerContract(mbpkgID int) (ContractUsage, error) {
	var contract ContractUsage
	if err := c.get(context.Background(), fmt.Sprintf("cloud/servers/%d/contract", mbpkgID), &contract); err != nil {
		return ContractUsage{}, err
	}
	return contract, nil
}

// GetServerSummary returns summary data for a server.
func (c *Client) GetServerSummary(mbpkgID int) (json.RawMessage, error) {
	return c.getAPIRaw(fmt.Sprintf("cloud/serversummary/%d", mbpkgID))
}

// GetPlanID returns a plan id by plan name.
func (c *Client) GetPlanID(planName string) (json.RawMessage, error) {
	return c.getAPIRaw("cloud/sizes/plan-id/" + url.PathEscape(planName))
}

// GetDeploySizes returns deploy sizes for a location.
func (c *Client) GetDeploySizes(location string, req *DeploySizesRequest) ([]Size, error) {
	q := url.Values{}
	if req != nil {
		if req.MinCPU != nil {
			q.Set("min_cpu", strconv.Itoa(*req.MinCPU))
		}
		if req.MinRAM != nil {
			q.Set("min_ram", strconv.Itoa(*req.MinRAM))
		}
	}
	path := "cloud/sizes/" + url.PathEscape(location)
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var sizes []Size
	if err := c.get(context.Background(), path, &sizes); err != nil {
		return nil, err
	}
	return sizes, nil
}

// GetStorageLocations returns storage locations. When cloudPoolID is non-nil it
// is sent as the cloud_pool_id query parameter.
func (c *Client) GetStorageLocations(cloudPoolID *int) (json.RawMessage, error) {
	path := "cloud/storage-locations"
	if cloudPoolID != nil {
		q := url.Values{}
		q.Set("cloud_pool_id", strconv.Itoa(*cloudPoolID))
		path += "?" + q.Encode()
	}
	return c.getAPIRaw(path)
}

// BindCloudFirewallSet binds a firewall set to a cloud virtual server.
func (c *Client) BindCloudFirewallSet(mbpkgID int, req *BindFirewallSetRequest) (json.RawMessage, error) {
	return c.postJSONRaw(fmt.Sprintf("cloud/%d/firewall-sets", mbpkgID), req)
}

// UnbindCloudFirewallSet unbinds a firewall set from a cloud virtual server.
func (c *Client) UnbindCloudFirewallSet(mbpkgID int, firewallSet string) error {
	return c.delete(context.Background(), fmt.Sprintf("cloud/%d/firewall-sets/%s", mbpkgID, url.PathEscape(firewallSet)), nil, nil)
}

// CreateUsageContract creates a usage contract.
func (c *Client) CreateUsageContract(req *CreateUsageContractRequest) (ContractUsage, error) {
	var contract ContractUsage
	body, err := json.Marshal(req)
	if err != nil {
		return ContractUsage{}, fmt.Errorf("encode usage contract request: %w", err)
	}
	if err := c.postJSON(context.Background(), "cloud/contract/usage", body, &contract); err != nil {
		return ContractUsage{}, err
	}
	return contract, nil
}

// ParseCloudInit uploads and parses a cloud-init script.
func (c *Client) ParseCloudInit(filename string, r io.Reader) (json.RawMessage, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("create cloud-init form file: %w", err)
	}
	if _, err := io.Copy(part, r); err != nil {
		return nil, fmt.Errorf("write cloud-init form file: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close cloud-init form: %w", err)
	}

	req, err := c.newRequest(context.Background(), "POST", "cloud/parse-cloud-init", &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var raw json.RawMessage
	if err := c.do(req, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *Client) getAPIRaw(path string) (json.RawMessage, error) {
	var raw json.RawMessage
	if err := c.get(context.Background(), path, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *Client) postJSONRaw(path string, req interface{}) (json.RawMessage, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}
	var raw json.RawMessage
	if err := c.postJSON(context.Background(), path, body, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *Client) putJSONRaw(path string, req interface{}) (json.RawMessage, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}
	var raw json.RawMessage
	if err := c.putJSON(context.Background(), path, body, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}
