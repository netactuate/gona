package gona

import (
	"encoding/json"
	"fmt"
)

type VPCServer struct {
	ServerID   int    `json:"serverId"`
	MBPkgID    int    `json:"mbpkgid"`
	FQDN       string `json:"fqdn"`
	Membership string `json:"membership"`
	SubnetID   int    `json:"subnetId"`
	NICID      int    `json:"nicId"`
	State      string `json:"state"`
	Address    string `json:"address"`
}
type AttachVPCServerRequest struct {
	MBPkgID  int `json:"mbpkgid"`
	SubnetID int `json:"subnetId,omitempty"`
}

func (c *V3Client) ListVPCServers(vpcID int) ([]VPCServer, error) {
	r, err := c.getListUnder(fmt.Sprintf("/vpcs/%d/servers", vpcID), "")
	if err != nil {
		return nil, err
	}
	var out []VPCServer
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *V3Client) AttachVPCServer(vpcID int, req *AttachVPCServerRequest) (*VPCServer, error) {
	r, err := c.post(fmt.Sprintf("/vpcs/%d/servers", vpcID), req)
	if err != nil {
		return nil, err
	}
	var out VPCServer
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
func (c *V3Client) DetachVPCServer(vpcID, mbpkgid int) error {
	_, err := c.del(fmt.Sprintf("/vpcs/%d/servers/%d", vpcID, mbpkgid))
	return err
}
