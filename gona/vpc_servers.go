package gona

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

// UnmarshalJSON accepts subnetId as a number, a numeric string or null. vapi3 lists attached
// rows with a string ("9", "0") and native rows with a number, and POST can answer null (LD-216, LD-242).
func (s *VPCServer) UnmarshalJSON(data []byte) error {
	type plain VPCServer
	var raw struct {
		*plain
		SubnetID json.RawMessage `json:"subnetId"`
	}
	raw.plain = (*plain)(s)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	id, err := flexibleInt(raw.SubnetID)
	if err != nil {
		return fmt.Errorf("subnetId: %w", err)
	}
	s.SubnetID = id
	return nil
}

func flexibleInt(raw json.RawMessage) (int, error) {
	text := strings.TrimSpace(string(raw))
	if text == "" || text == "null" {
		return 0, nil
	}
	if strings.HasPrefix(text, "\"") {
		var str string
		if err := json.Unmarshal(raw, &str); err != nil {
			return 0, err
		}
		if str = strings.TrimSpace(str); str == "" {
			return 0, nil
		}
		return strconv.Atoi(str)
	}
	var n int
	err := json.Unmarshal(raw, &n)
	return n, err
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
