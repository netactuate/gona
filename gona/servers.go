// servers.go (in your github.com/netactuate/gona package)

package gona

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
)

// Server struct defines what a VPS looks like
type Server struct {
	Name                     string `json:"fqdn"`
	ID                       int    `json:"mbpkgid"`
	OS                       string `json:"os"`
	OSID                     int    `json:"os_id"`
	PrimaryIPv4              string `json:"ip"`
	PrimaryIPv6              string `json:"ipv6"`
	PlanID                   int    `json:"plan_id"`
	Package                  string `json:"package"`
	PackageBilling           string `json:"package_billing"`
	PackageBillingContractId string `json:"package_billing_contract_id"`
	Location                 string `json:"city"`
	LocationID               int    `json:"location_id"`
	ServerStatus             string `json:"status"`
	PowerStatus              string `json:"state"`
	Installed                int    `json:"installed"`
}

// GetServers external method on Client to list your instances
func (c *Client) GetServers() ([]Server, error) {
	var serverList []Server
	if err := c.get("cloud/servers", &serverList); err != nil {
		return nil, err
	}
	return serverList, nil
}

// GetServer external method on Client to get an instance
func (c *Client) GetServer(id int) (server Server, err error) {
	if err := c.get("cloud/server?mbpkgid="+strconv.Itoa(id), &server); err != nil {
		return server, err
	}
	return server, nil
}

// CreateServerRequest is a set of parameters for a server creation call.
type CreateServerRequest struct {
	Plan                     string `url:"plan,omitempty"`
	Location                 int    `url:"location,omitempty"`
	Image                    int    `url:"image,omitempty"`
	FQDN                     string `url:"fqdn,omitempty"`
	SSHKey                   string `url:"ssh_key,omitempty"`
	SSHKeyID                 int    `url:"ssh_key_id,omitempty"`
	Password                 string `url:"password,omitempty"`
	PackageBilling           string `url:"package_billing,omitempty"`
	PackageBillingContractId string `url:"package_billing_contract_id,omitempty"`
	CloudConfig              string `url:"cloud_config,omitempty"`
	ScriptContent            string `url:"script_content,omitempty"`
	Params                   string `url:"params,omitempty"`
}

// ServerBuild is a server creation response message.
type ServerBuild struct {
	ServerID int    `json:"mbpkgid"`
	Status   string `json:"status"`
	Build    int    `json:"build"`
}

// CreateServer external method on Client to buy and build a new instance.
func (c *Client) CreateServer(r *CreateServerRequest) (b ServerBuild, err error) {
	values, err := query.Values(r)
	if err != nil {
		return b, err
	}
	if values.Has("script_content") {
		values.Add("script_type", "user-data")
	}

	if err := c.post("cloud/server/buy_build", []byte(values.Encode()), &b); err != nil {
		return b, err
	}
	return b, nil
}

// BuildServerRequest is a set of parameters for a server re-building call.
type BuildServerRequest struct {
	Plan                     string `url:"plan,omitempty"`
	Location                 int    `url:"location,omitempty"`
	Image                    int    `url:"image,omitempty"`
	FQDN                     string `url:"fqdn,omitempty"`
	SSHKey                   string `url:"ssh_key,omitempty"`
	SSHKeyID                 int    `url:"ssh_key_id,omitempty"`
	Password                 string `url:"password,omitempty"`
	PackageBilling           string `url:"package_billing,omitempty"`
	PackageBillingContractId string `url:"package_billing_contract_id,omitempty"`
	CloudConfig              string `url:"cloud_config,omitempty"`
	ScriptContent            string `url:"script_content,omitempty"`
	Params                   string `url:"params,omitempty"`
}

// BuildServer external method on Client to re-build an instance
func (c *Client) BuildServer(id int, r *BuildServerRequest) (b ServerBuild, err error) {
	values, err := query.Values(r)
	if err != nil {
		return b, err
	}
	if values.Has("script_content") {
		values.Add("script_type", "user-data")
	}

	if err := c.post("cloud/server/build/"+strconv.Itoa(id), []byte(values.Encode()), &b); err != nil {
		return b, err
	}
	return b, nil
}

// DeleteServer sends a delete request, surfaces any API‐reported errors,
// and returns the asynchronous job ID.
func (c *Client) DeleteServer(id int, cancelBilling bool) (int, error) {
	values := url.Values{}
	values.Set("mbpkgid", fmt.Sprint(id))
	if cancelBilling {
		values.Set("cancel_billing", "1")
	}
	body := []byte(values.Encode())

	// Build request manually to capture raw API errors
	req, err := c.newRequest("POST", "cloud/server/delete", bytes.NewBuffer(body))
	if err != nil {
		return 0, fmt.Errorf("DeleteServer(%d) build request: %w", id, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("DeleteServer(%d) HTTP error: %w", id, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("DeleteServer(%d) read body: %w", id, err)
	}
	c.debugLog("DeleteServer response: %s", string(raw))

	// Unmarshal API envelope
	var apiResp struct {
		Result  string                   `json:"result"`
		Message string                   `json:"message"`
		Fields  map[string][]interface{} `json:"fields"`
		Data    struct {
			ID int `json:"id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		return 0, fmt.Errorf("DeleteServer(%d) unmarshal: %w\nraw: %s", id, err, string(raw))
	}

	if apiResp.Result != "success" {
		return 0, fmt.Errorf(
			"API error deleting server %d: %s; details: %+v",
			id, apiResp.Message, apiResp.Fields,
		)
	}
	if apiResp.Data.ID == 0 {
		return 0, fmt.Errorf("delete request returned zero job ID for server %d", id)
	}

	return apiResp.Data.ID, nil
}

// UnlinkServer external method on Client to unlink a billing package from a location
func (c *Client) UnlinkServer(id int) error {
	return c.post("cloud/server/unlink/"+strconv.Itoa(id), nil, nil)
}

// StartServer external method on Client to boot up an instance
func (c *Client) StartServer(id int) error {
	return c.post("cloud/server/start/"+strconv.Itoa(id), nil, nil)
}

// StopServer external method on Client to shut down an instance
func (c *Client) StopServer(id int) error {
	return c.post("cloud/server/shutdown/"+strconv.Itoa(id), nil, nil)
}

// Job represents the API's asynchronous job object.
type Job struct {
	ID       int    `json:"id"`
	TSInsert string `json:"ts_insert"`
	Command  string `json:"command"`
	Status   int    `json:"status"`
}

// GetJob fetches a specific job for a server, returning its details.
func (c *Client) GetJob(serverID, jobID int) (Job, error) {
	var job Job
	path := fmt.Sprintf("cloud/server/%d/jobs/%d/", serverID, jobID)
	if err := c.get(path, &job); err != nil {
		return Job{}, fmt.Errorf("failed to fetch job %d for server %d: %w", jobID, serverID, err)
	}
	return job, nil
}

// WaitForJob polls the given job until it matches the expected command and reaches status 5.
func (c *Client) WaitForJob(serverID, jobID int, expectedCommand string) error {
	const (
		maxAttempts  = 200
		intervalSecs = 5
	)

	for i := 0; i < maxAttempts; i++ {
		job, err := c.GetJob(serverID, jobID)
		if err != nil {
			return fmt.Errorf("error polling job %d: %w", jobID, err)
		}
		if job.Command != expectedCommand {
			return fmt.Errorf(
				"job %d command mismatch: got %q, want %q",
				jobID, job.Command, expectedCommand,
			)
		}
		if job.Status == 5 {
			return nil
		}
		time.Sleep(intervalSecs * time.Second)
	}
	return fmt.Errorf(
		"timed out waiting for job %d (command %q) to reach status 5",
		jobID, expectedCommand,
	)
}
