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

// CloudNetworkingLocation maps a cloud location to a datacenter.
type CloudNetworkingLocation struct {
	LocationID   int `json:"locationId"`
	DatacenterID int `json:"datacenterId"`
}

// CreateCloudFloatingIPv4Request creates a floating IPv4 address.
type CreateCloudFloatingIPv4Request struct {
	PTRDomain *string `json:"ptrDomain,omitempty"`
	VLANID    *int    `json:"vlanId,omitempty"`
}

// CloudFloatingIPv4VM is a virtual machine allowed to access a floating IPv4 address.
type CloudFloatingIPv4VM struct {
	MBPkgID int    `json:"mbpkgid"`
	FQDN    string `json:"fqdn,omitempty"`
	IP      string `json:"ip,omitempty"`
}

// CloudFloatingIPv4VMRef identifies a virtual machine for floating IPv4 grants.
type CloudFloatingIPv4VMRef struct {
	MBPkgID int `json:"mbpkgid"`
}

// GrantCloudFloatingIPv4VMsRequest grants VMs access to a floating IPv4 address.
type GrantCloudFloatingIPv4VMsRequest struct {
	RevokeExisting *bool                    `json:"revokeExisting,omitempty"`
	VMs            []CloudFloatingIPv4VMRef `json:"vms,omitempty"`
}

// RevokeCloudFloatingIPv4VMsRequest revokes VM access to a floating IPv4 address.
type RevokeCloudFloatingIPv4VMsRequest struct {
	VMs []CloudFloatingIPv4VMRef `json:"vms,omitempty"`
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

// GetCustomerVLAN returns a customer VLAN by id.
func (c *Client) GetCustomerVLAN(customerVLANID int) (VLAN, error) {
	var vlan VLAN
	if err := c.get(context.Background(), fmt.Sprintf("cloud/networking/vlans/%d", customerVLANID), &vlan); err != nil {
		return VLAN{}, err
	}
	return vlan, nil
}

// ListCustomerVLANsAtLocation returns customer VLANs at a location.
func (c *Client) ListCustomerVLANsAtLocation(locationID int) ([]VLAN, error) {
	var vlans []VLAN
	if err := c.get(context.Background(), fmt.Sprintf("cloud/networking/locations/%d/vlans", locationID), &vlans); err != nil {
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

// CreateCloudFloatingIPv4 adds a floating IPv4 address to the account.
func (c *V3Client) CreateCloudFloatingIPv4(req *CreateCloudFloatingIPv4Request) (CloudFloatingIPv4, error) {
	resp, err := c.post("/cloud/networking/floating-ips/ipv4", req)
	if err != nil {
		return CloudFloatingIPv4{}, fmt.Errorf("create cloud floating IPv4 address: %w", err)
	}

	var floatingIP CloudFloatingIPv4
	if err := json.Unmarshal(resp.Data, &floatingIP); err != nil {
		return CloudFloatingIPv4{}, fmt.Errorf("create cloud floating IPv4 address unmarshal: %w", err)
	}
	return floatingIP, nil
}

// DeleteCloudFloatingIPv4 deletes a floating IPv4 address.
func (c *V3Client) DeleteCloudFloatingIPv4(floatingIPv4ID int) error {
	if _, err := c.del(fmt.Sprintf("/cloud/networking/floating-ips/ipv4/%d", floatingIPv4ID)); err != nil {
		return fmt.Errorf("delete cloud floating IPv4 address: %w", err)
	}
	return nil
}

// ListCloudFloatingIPv4VMs returns VMs allowed to access a floating IPv4 address.
func (c *V3Client) ListCloudFloatingIPv4VMs(floatingIPv4ID int) ([]CloudFloatingIPv4VM, error) {
	listData, err := c.getList(fmt.Sprintf("/cloud/networking/floating-ips/ipv4/%d/vms", floatingIPv4ID))
	if err != nil {
		return nil, fmt.Errorf("list cloud floating IPv4 VMs: %w", err)
	}

	var vms []CloudFloatingIPv4VM
	if err := json.Unmarshal(listData.Data, &vms); err != nil {
		return nil, fmt.Errorf("list cloud floating IPv4 VMs unmarshal inner: %w", err)
	}
	return vms, nil
}

// GrantCloudFloatingIPv4VMs grants VMs access to a floating IPv4 address.
func (c *V3Client) GrantCloudFloatingIPv4VMs(floatingIPv4ID int, req *GrantCloudFloatingIPv4VMsRequest) error {
	if _, err := c.post(fmt.Sprintf("/cloud/networking/floating-ips/ipv4/%d/vms/mass-grant", floatingIPv4ID), req); err != nil {
		return fmt.Errorf("grant cloud floating IPv4 VMs: %w", err)
	}
	return nil
}

// RevokeCloudFloatingIPv4VMs revokes VM access to a floating IPv4 address.
func (c *V3Client) RevokeCloudFloatingIPv4VMs(floatingIPv4ID int, req *RevokeCloudFloatingIPv4VMsRequest) error {
	if _, err := c.post(fmt.Sprintf("/cloud/networking/floating-ips/ipv4/%d/vms/mass-revoke", floatingIPv4ID), req); err != nil {
		return fmt.Errorf("revoke cloud floating IPv4 VMs: %w", err)
	}
	return nil
}

// ListCloudNetworkingLocations returns cloud location to datacenter mappings.
func (c *V3Client) ListCloudNetworkingLocations() ([]CloudNetworkingLocation, error) {
	resp, err := c.get("/cloud/networking/locations")
	if err != nil {
		return nil, fmt.Errorf("list cloud networking locations: %w", err)
	}

	var locations []CloudNetworkingLocation
	if err := json.Unmarshal(resp.Data, &locations); err != nil {
		return nil, fmt.Errorf("list cloud networking locations unmarshal: %w", err)
	}
	return locations, nil
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
