package gona

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type NKEAPIAccess struct {
	PublicEnabled          bool   `json:"publicEnabled"`
	PrivateEndpointEnabled bool   `json:"privateEndpointEnabled"`
	PrivateEndpointAddress string `json:"privateEndpointAddress"`
	PrivateDNSName         string `json:"privateDnsName"`
	PrivateEndpointState   string `json:"privateEndpointState"`
	AllowListActive        bool   `json:"allowListActive"`
	State                  string `json:"state"`
	AppliedOn              string `json:"appliedOn"`
	FailureReason          string `json:"failureReason"`
}
type UpdateNKEAPIAccessRequest struct {
	PublicEnabled          *bool `json:"publicEnabled,omitempty"`
	PrivateEndpointEnabled *bool `json:"privateEndpointEnabled,omitempty"`
}
type NKEAPIAllowEntry struct {
	AllowListEntryID string `json:"allowListEntryId"`
	Network          string `json:"network"`
	Description      string `json:"description"`
	IsPlatform       bool   `json:"isPlatform"`
	Owner            string `json:"owner"`
	CreatedOn        string `json:"createdOn"`
}
type CreateNKEAPIAllowEntryRequest struct {
	Network     string `json:"network"`
	Description string `json:"description,omitempty"`
}

func (c *V3Client) GetNKEAPIAccess(clusterID int) (*NKEAPIAccess, error) {
	r, err := c.get(fmt.Sprintf("/nke/clusters/%d/api-access", clusterID))
	if err != nil {
		return nil, err
	}
	var out NKEAPIAccess
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *V3Client) UpdateNKEAPIAccess(clusterID int, req *UpdateNKEAPIAccessRequest) (*NKEAPIAccess, error) {
	r, err := c.patch(fmt.Sprintf("/nke/clusters/%d/api-access", clusterID), req)
	if err != nil {
		return nil, err
	}
	var out NKEAPIAccess
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *V3Client) ListNKEAPIAllowList(clusterID int) ([]NKEAPIAllowEntry, error) {
	r, err := c.getListUnder(fmt.Sprintf("/nke/clusters/%d/api-access/allow-list", clusterID), "")
	if err != nil {
		return nil, err
	}
	var out []NKEAPIAllowEntry
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *V3Client) CreateNKEAPIAllowEntry(clusterID int, req *CreateNKEAPIAllowEntryRequest) (*NKEAPIAllowEntry, error) {
	r, err := c.post(fmt.Sprintf("/nke/clusters/%d/api-access/allow-list", clusterID), req)
	if err != nil {
		return nil, err
	}
	var out NKEAPIAllowEntry
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *V3Client) DeleteNKEAPIAllowEntry(clusterID int, entryID string) error {
	_, err := c.del(fmt.Sprintf("/nke/clusters/%d/api-access/allow-list/%s", clusterID, url.PathEscape(entryID)))
	if err != nil {
		return err
	}
	return nil
}
