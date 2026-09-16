package gona

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type VPCSSHSettings struct {
	Port    *int                 `json:"port,omitempty"`
	Enabled bool                 `json:"enabled"`
	Keys    map[string]VPCSSHKey `json:"keys,omitempty"`
	Bastion *VPCSSHBastion       `json:"bastion,omitempty"`
}

func (s *VPCSSHSettings) UnmarshalJSON(data []byte) error {
	type alias VPCSSHSettings
	var aux struct {
		*alias
		Port json.RawMessage `json:"port"`
	}
	aux.alias = (*alias)(s)
	if err := json.Unmarshal(data, &aux); err != nil {
		// The keys field arrives either as a map of key id to key, or wrapped in the list
		// envelope this API uses elsewhere: {meta, data:[...]}. Decoding it as a map alone fails
		// against the second shape.
		var wrapped struct {
			*alias
			Port json.RawMessage `json:"port"`
			Keys struct {
				Data []VPCSSHKey `json:"data"`
			} `json:"keys"`
		}
		wrapped.alias = (*alias)(s)
		if werr := json.Unmarshal(data, &wrapped); werr != nil {
			return err
		}
		s.Keys = make(map[string]VPCSSHKey, len(wrapped.Keys.Data))
		for _, k := range wrapped.Keys.Data {
			s.Keys[strconv.Itoa(k.ID)] = k
		}
		aux.Port = wrapped.Port
	}
	if len(aux.Port) == 0 || string(aux.Port) == "null" {
		s.Port = nil
		return nil
	}
	port, err := decodeVPCSSHPort(aux.Port)
	if err != nil {
		return err
	}
	s.Port = port
	return nil
}

func decodeVPCSSHPort(raw json.RawMessage) (*int, error) {
	var port int
	if err := json.Unmarshal(raw, &port); err == nil {
		return &port, nil
	}

	var portString string
	if err := json.Unmarshal(raw, &portString); err == nil {
		parsed, err := strconv.Atoi(portString)
		if err != nil {
			return nil, err
		}
		return &parsed, nil
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil, err
	}
	for _, key := range []string{"port", "value", "number"} {
		if value, ok := fields[key]; ok {
			return decodeVPCSSHPort(value)
		}
	}
	return nil, fmt.Errorf("VPC SSH port object has no port value")
}

type VPCSSHBastion struct {
	IPv4 string `json:"ipv4,omitempty"`
	IPv6 string `json:"ipv6,omitempty"`
}

type VPCSSHKey struct {
	ID          int    `json:"id"`
	SSHKeyID    int    `json:"sshKeyId"`
	Name        string `json:"name,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	PublicKey   string `json:"publicKey,omitempty"`
	Dates       *struct {
		Created string  `json:"created,omitempty"`
		Enabled *string `json:"enabled,omitempty"`
	} `json:"dates,omitempty"`
}

func (k *VPCSSHKey) GetID() int {
	if k.ID != 0 {
		return k.ID
	}
	return k.SSHKeyID
}

// IsEnabled returns true if the key has been enabled for the VPC
func (k *VPCSSHKey) IsEnabled() bool {
	return k.Dates != nil && k.Dates.Enabled != nil
}

type UpdateVPCSSHSettingsRequest struct {
	Port *int `json:"port"`
}

type EnableVPCSSHKeyRequest struct {
	Enabled bool `json:"enabled"`
}

func (c *V3Client) UpdateVPCSSHSettings(vpcID int, req *UpdateVPCSSHSettingsRequest) (*VPCSSHSettings, error) {
	path := fmt.Sprintf("/vpcs/%d/ssh", vpcID)
	resp, err := c.patch(path, req)
	if err != nil {
		return nil, fmt.Errorf("update SSH settings for VPC %d: %w", vpcID, err)
	}
	var settings VPCSSHSettings
	if err := json.Unmarshal(resp.Data, &settings); err != nil {
		return nil, fmt.Errorf("update SSH settings unmarshal: %w", err)
	}
	return &settings, nil
}

func (c *V3Client) ListVPCSSHKeys(vpcID int) ([]VPCSSHKey, error) {
	path := fmt.Sprintf("/vpcs/%d/ssh/keys", vpcID)
	listData, err := c.getList(path)
	if err != nil {
		return nil, fmt.Errorf("list SSH keys for VPC %d: %w", vpcID, err)
	}
	var keys []VPCSSHKey
	if err := json.Unmarshal(listData.Data, &keys); err != nil {
		return nil, fmt.Errorf("list SSH keys unmarshal inner: %w", err)
	}
	return keys, nil
}

func (c *V3Client) GetVPCSSHKey(vpcID, sshKeyID int) (*VPCSSHKey, error) {
	keys, err := c.ListVPCSSHKeys(vpcID)
	if err != nil {
		return nil, err
	}
	for _, k := range keys {
		if k.GetID() == sshKeyID {
			return &k, nil
		}
	}
	return nil, &V3NotFoundError{StatusCode: 404, Body: fmt.Sprintf("SSH key %d not found in VPC %d", sshKeyID, vpcID)}
}

func (c *V3Client) EnableVPCSSHKey(vpcID, sshKeyID int, enabled bool) error {
	path := fmt.Sprintf("/vpcs/%d/ssh/keys/%d", vpcID, sshKeyID)
	req := &EnableVPCSSHKeyRequest{Enabled: enabled}
	_, err := c.patch(path, req)
	if err != nil {
		return fmt.Errorf("enable/disable SSH key %d for VPC %d: %w", sshKeyID, vpcID, err)
	}
	return nil
}

func (c *V3Client) DeleteVPCSSHKey(vpcID, sshKeyID int) error {
	return c.EnableVPCSSHKey(vpcID, sshKeyID, false)
}
