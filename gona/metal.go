package gona

import (
// 	"net/url"
 	"strconv"
     "log"
//     "fmt"
	"github.com/google/go-querystring/query"
)

type CreateMetalRequest struct {
	Location     int    `url:"location,omitempty"`
	Device       int    `url:"device_id,omitempty"`
	SSHKey       string `url:"ssh_key,omitempty"`
	SSHKeyID     int    `url:"ssh_key_id,omitempty"`
	Password     string `url:"root_password,omitempty"`
	BuildScript  string `url:"build_script,omitempty"`
	DiskLayout   int    `url:"disklayout,omitempty"`
	Profile      int    `url:"profile,omitempty"`
	Hostname     string `url:"fqdn,omitempty"`
}

type BuildMetalRequest struct {
	MBPKGID      int    `url:"mbpkgid"`
	SSHKey       string `url:"ssh_key,omitempty"`
	SSHKeyID     int    `url:"ssh_key_id,omitempty"`
	Password     string `url:"root_password,omitempty"`
	BuildScript  string `url:"build_script,omitempty"`
	DiskLayout   int    `url:"disklayout,omitempty"`
	Profile      int    `url:"profile,omitempty"`
	Hostname     string `url:"fqdn,omitempty"`
}

type MetalBuild struct {
	MBPKGID int    `json:"mbpkgid"`
	Status   string `json:"status"`
	Build    int    `json:"build"`
}
type MetalBuildStatus struct {
    MBPKGID int    `json:"mbpkgid"`
    Response string `json:"response"`
    Status   string `json:"status"`
    Percent   int   `json:"percent"`
    ImageName string `json:"image_name"`
}

type SingleMetal struct {
	ID                 int     `json:"id"`
	DatacenterID       int     `json:"datacenter_id"`
	Canceling          int     `json:"canceling"`
	MBPKGID            int     `json:"mbpkgid"`
	Price              string  `json:"price"`
	Hostname           string  `json:"hostname"`
	Eth0MAC            string  `json:"eth0_mac"`
	Eth1MAC            *string `json:"eth1_mac"`
	IPMIMAC            string  `json:"ipmi_mac"`
	PrimaryIP          string  `json:"primary_ip"`
	PrimaryIPv6        *string `json:"primary_ipv6"`
	NPSInstalled       int     `json:"nps_installed"`
	NPSOS              string  `json:"nps_os"`
	MBModel            *string `json:"mb_model"`
	CPU0Model          string  `json:"cpu0_model"`
	CPU1Model          *string `json:"cpu1_model"`
	TotalRAM           int     `json:"total_ram"`
	IPMIPubIP          string  `json:"ipmi_pubip"`
	IPMICXUser         *string `json:"ipmi_cxuser"`
	IPMICXPass         *string `json:"ipmi_cxpass"`
	IPMIStatus         int     `json:"ipmi_status"`
	Locked             int     `json:"locked"`
	LockedMsg          *string `json:"locked_msg"`
	IPMIStatusTime     string  `json:"ipmi_status_time"`
	OBID               *string `json:"ob_id"`
	Info               string  `json:"info"`
	Title              string  `json:"title"`
	IPMIRefreshEnabled int     `json:"ipmi_refresh_enabled"`
	Location           string  `json:"location"`
	IPSubnetID         int     `json:"ip_subnet_id"`
	IPSubnetName       string  `json:"ip_subnet_name"`
	PackageStatus      string  `json:"package_status"`
	Building           *string `json:"building"`
}


func (c *Client) GetMetalBuildStatus(buildID int) (MetalBuildStatus, error) {
    var buildStatus MetalBuildStatus
    endpoint := "dedicated/server/build_status/" + strconv.Itoa(buildID)
    log.Printf("[DEBUG] GetMetalBuildStatus: Making request to endpoint: %s", endpoint)

    if err := c.get(endpoint, &buildStatus); err != nil {
        log.Printf("[DEBUG] GetMetalBuildStatus: API call failed for endpoint %s: %v", endpoint, err)
        return MetalBuildStatus{}, err
    }

    return buildStatus, nil
}


func (c *Client) CreateMetal(r *CreateMetalRequest) (b MetalBuild, err error) {
	values, err := query.Values(r)
	if err != nil {
		return b, err
	}
	if err := c.post("dedicated/server/buy_build", []byte(values.Encode()), &b); err != nil {
		return b, err
	}

	return b, nil
}

func (c *Client) BuildMetal(id int, r *BuildMetalRequest) (b MetalBuild, err error) {
	values, err := query.Values(r)
	if err != nil {
		return b, err
	}
	if err := c.post("dedicated/server/re_build/" + strconv.Itoa(id), []byte(values.Encode()), &b); err != nil {
		return b, err
	}

	return b, nil
}

func (c *Client) GetMetal(id int) (metal SingleMetal, err error) {
	if err := c.get("dedicated/servers/"+strconv.Itoa(id), &metal); err != nil {
		return metal, err
	}
	return metal, nil
}

