package gona

import (
// 	"net/url"
// 	"strconv"
//     "log"
//     "fmt"
	"github.com/google/go-querystring/query"
)

type CreateMetalRequest struct {
	Location     int    `url:"location,omitempty"`
	Device       int    `url:"device_id,omitempty"`
	SSHKey       string `url:"ssh_key,omitempty"`
	SSHKeyID     int    `url:"ssh_key_id,omitempty"`
	Password     string `url:"password,omitempty"`
	BuildScript  string `url:"build_script,omitempty"`
	DiskLayout   int    `url:"disklayout,omitempty"`
	Profile      int    `url:"profile,omitempty"`
	Hostname     string `url:"fqdn,omitempty"`
}
type MetalBuild struct {
	ServerID int    `json:"mbpkgid"`
	Status   string `json:"status"`
	Build    int    `json:"build"`
}


func (c *Client) CreateMetal(r *CreateMetalRequest) (b MetalBuild, err error) {
	values, err := query.Values(r)
	if err != nil {
		return b, err
	}
	if values.Has("script_content") {
		values.Add("script_type", "user-data")
	}

	if err := c.post("dediacated/server/build", []byte(values.Encode()), &b); err != nil {
		return b, err
	}

	return b, nil
}
