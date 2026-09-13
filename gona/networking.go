package gona

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type VLAN struct {
	ID                   int                   `json:"id"`
	MBID                 int                   `json:"mbid"`
	Private              int                   `json:"private"`
	AllowSRIOV           int                   `json:"allow_sriov"`
	DisplayName          string                `json:"display_name"`
	Description          string                `json:"description"`
	LastUpdated          string                `json:"last_updated"`
	Created              string                `json:"created"`
	ProvisionedLocations []ProvisionedLocation `json:"provisioned_locations"`
}

type ProvisionedLocation struct {
	Provisioned bool   `json:"provisioned"`
	Name        string `json:"name"`
	LocationID  int    `json:"location_id"`
	Flag        string `json:"flag"`
	IATACode    string `json:"iata_code"`
}

type CloudFloatingIPv4 struct {
	FloatingIPv4ID int                 `json:"floatingIpv4Id"`
	AssignedOn     string              `json:"AssignedOn"`
	Address        string              `json:"address"`
	VLANID         int                 `json:"vlanId"`
	PTRDomain      *string             `json:"ptrDomain"`
	Location       *FloatingIPLocation `json:"location"`
}

type FloatingIPLocation struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Flag      string `json:"flag"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type ServerNIC struct {
	NICID          int `json:"-"`
	MBPkgID        int `json:"-"`
	CustomerVLANID int `json:"-"`
	AttachOrder    int `json:"-"`
}

func (n *ServerNIC) UnmarshalJSON(data []byte) error {
	type Alias ServerNIC
	aux := &struct {
		NICID          interface{} `json:"nic_id"`
		ID             interface{} `json:"id"`
		MBPkgID        interface{} `json:"mbpkgid"`
		CustomerVLANID interface{} `json:"customer_vlan_id"`
		AttachOrder    interface{} `json:"attach_order"`
		*Alias
	}{
		Alias: (*Alias)(n),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	n.NICID = toInt(aux.NICID)
	if n.NICID == 0 {
		n.NICID = toInt(aux.ID)
	}
	n.MBPkgID = toInt(aux.MBPkgID)
	n.CustomerVLANID = toInt(aux.CustomerVLANID)
	n.AttachOrder = toInt(aux.AttachOrder)
	return nil
}

type ServerNICAttachRequest struct {
	CustomerVLANID int `json:"customer_vlan_id"`
}

type ServerNICUpdateRequest struct {
	MBPkgID        int `json:"mbpkgid"`
	CustomerVLANID int `json:"customer_vlan_id"`
	AttachOrder    int `json:"attach_order"`
}

func (c *Client) GetVLANs() ([]VLAN, error) {
	var vlans []VLAN
	if err := c.get(context.Background(), "cloud/networking/vlans", &vlans); err != nil {
		return nil, err
	}
	return vlans, nil
}

func (c *V3Client) ListCloudFloatingIPv4() ([]CloudFloatingIPv4, error) {
	listData, err := c.getList("/cloud/networking/floating-ips/ipv4")
	if err != nil {
		return nil, fmt.Errorf("list cloud floating IPv4 addresses: %w", err)
	}

	var floatingIPs []CloudFloatingIPv4
	if err := json.Unmarshal(listData.Data, &floatingIPs); err != nil {
		return nil, fmt.Errorf("list cloud floating IPv4 addresses unmarshal inner: %w", err)
	}
	return floatingIPs, nil
}

func (c *Client) GetServerNICs(mbpkgID int) ([]ServerNIC, error) {
	var nics []ServerNIC
	path := fmt.Sprintf("cloud/networking/nics/%d", mbpkgID)
	if err := c.get(context.Background(), path, &nics); err != nil {
		return nil, err
	}
	return nics, nil
}

func (c *Client) AttachServerNIC(mbpkgID int, req *ServerNICAttachRequest) (ServerNIC, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return ServerNIC{}, fmt.Errorf("encoding server NIC attach request: %w", err)
	}

	var raw json.RawMessage
	path := fmt.Sprintf("cloud/networking/nics/%d", mbpkgID)
	if err := c.postJSON(context.Background(), path, body, &raw); err != nil {
		return ServerNIC{}, err
	}
	return serverNICFromRaw(raw), nil
}

func (c *Client) UpdateServerNIC(nicID int, req *ServerNICUpdateRequest) (ServerNIC, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return ServerNIC{}, fmt.Errorf("encoding server NIC update request: %w", err)
	}

	var raw json.RawMessage
	path := fmt.Sprintf("cloud/networking/nics/%d", nicID)
	if err := c.putJSON(context.Background(), path, body, &raw); err != nil {
		return ServerNIC{}, err
	}
	return serverNICFromRaw(raw), nil
}

func (c *Client) DetachServerNIC(mbpkgID, nicID int) error {
	values := url.Values{}
	values.Add("mbpkgid", strconv.Itoa(mbpkgID))
	path := fmt.Sprintf("cloud/networking/nics/%d?%s", nicID, values.Encode())
	return c.delete(context.Background(), path, nil, nil)
}

func serverNICFromRaw(raw json.RawMessage) ServerNIC {
	var nic ServerNIC
	if len(raw) == 0 {
		return nic
	}
	if err := json.Unmarshal(raw, &nic); err == nil {
		return nic
	}

	var nics []ServerNIC
	if err := json.Unmarshal(raw, &nics); err == nil && len(nics) == 1 {
		return nics[0]
	}
	return ServerNIC{}
}
