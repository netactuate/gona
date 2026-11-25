package gona

import (
	"encoding/json"
	"fmt"
)

type VPCAttachment struct {
	BGP *VPCAttachmentBGP `json:"bgp,omitempty"`

	ID                      int      `json:"id"`
	Type                    string   `json:"type"`
	Address                 string   `json:"address"`
	Description             string   `json:"description"`
	RouterID                int      `json:"routerId"`
	CloudRouterAttachmentID int      `json:"cloudRouterAttachmentId"`
	VRFID                   int      `json:"vrfId"`
	RouterState             string   `json:"routerState"`
	State                   string   `json:"state"`
	StateMessage            string   `json:"stateMessage"`
	ReadyOn                 string   `json:"readyOn"`
	CreatedOn               string   `json:"createdOn"`
	DeleteRequestedOn       string   `json:"deleteRequestedOn"`
	Accept                  []string `json:"accept"`
	Advertise               []string `json:"advertise"`
}
type CreateVPCAttachmentRequest struct {
	BGP *VPCAttachmentBGPRequest `json:"bgp,omitempty"`

	Type        string   `json:"type"`
	Accept      []string `json:"accept"`
	Advertise   []string `json:"advertise,omitempty"`
	Address     string   `json:"address,omitempty"`
	RouterID    int      `json:"routerId,omitempty"`
	VRFID       int      `json:"vrfId,omitempty"`
	Description string   `json:"description,omitempty"`
}
type VPCAttachmentPrefix struct {
	ID                int    `json:"id"`
	Direction         string `json:"direction"`
	IPVersion         int    `json:"ipVersion"`
	Network           string `json:"network"`
	CreatedOn         string `json:"createdOn"`
	DeleteRequestedOn string `json:"deleteRequestedOn"`
}
type CreateVPCAttachmentPrefixRequest struct {
	Network   string `json:"network"`
	Direction string `json:"direction"`
}

func (c *V3Client) ListVPCAttachments(vpcID int) ([]VPCAttachment, error) {
	r, err := c.getListUnder(fmt.Sprintf("/vpcs/%d/attachments", vpcID), "")
	if err != nil {
		return nil, err
	}
	var out []VPCAttachment
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *V3Client) GetVPCAttachment(vpcID, attachmentID int) (*VPCAttachment, error) {
	r, err := c.get(fmt.Sprintf("/vpcs/%d/attachments/%d", vpcID, attachmentID))
	if err != nil {
		return nil, err
	}
	var out VPCAttachment
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *V3Client) CreateVPCAttachment(vpcID int, req *CreateVPCAttachmentRequest) (int, error) {
	r, err := c.drop3Request("POST", fmt.Sprintf("/vpcs/%d/attachments", vpcID), req)
	if err != nil {
		return 0, err
	}
	var out struct {
		ID int `json:"attachmentId"`
	}
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}

func (c *V3Client) DeleteVPCAttachment(vpcID, attachmentID int) error {
	_, err := c.del(fmt.Sprintf("/vpcs/%d/attachments/%d", vpcID, attachmentID))
	if err != nil {
		return err
	}
	return nil
}

func (c *V3Client) ListVPCAttachmentPrefixes(vpcID, attachmentID int) ([]VPCAttachmentPrefix, error) {
	r, err := c.getListUnder(fmt.Sprintf("/vpcs/%d/attachments/%d/prefixes", vpcID, attachmentID), "")
	if err != nil {
		return nil, err
	}
	var out []VPCAttachmentPrefix
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *V3Client) CreateVPCAttachmentPrefix(vpcID, attachmentID int, req *CreateVPCAttachmentPrefixRequest) (int, error) {
	r, err := c.post(fmt.Sprintf("/vpcs/%d/attachments/%d/prefixes", vpcID, attachmentID), req)
	if err != nil {
		return 0, err
	}
	var out struct {
		ID int `json:"prefixId"`
	}
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}

func (c *V3Client) DeleteVPCAttachmentPrefix(vpcID, attachmentID, prefixID int) error {
	_, err := c.del(fmt.Sprintf("/vpcs/%d/attachments/%d/prefixes/%d", vpcID, attachmentID, prefixID))
	if err != nil {
		return err
	}
	return nil
}

type VPCAttachmentBGPRequest struct {
	PeerASN   int    `json:"peerAsn"`
	MaxPrefix int    `json:"maxPrefix,omitempty"`
	Password  string `json:"password,omitempty"`
}
type VPCAttachmentBGP struct {
	PeerASN        int    `json:"peerAsn"`
	MaxPrefix      int    `json:"maxPrefix"`
	HasPassword    bool   `json:"hasPassword"`
	GatewayASN     int    `json:"gatewayAsn"`
	GatewayAddress string `json:"gatewayAddress"`
	PeerAddress    string `json:"peerAddress"`
	State          string `json:"state"`
	StateMessage   string `json:"stateMessage"`
	AppliedOn      string `json:"appliedOn"`
}

func (c *V3Client) PutVPCAttachmentBGP(vpcID, attachmentID int, req *VPCAttachmentBGPRequest) error {
	_, err := c.drop3Request("PUT", fmt.Sprintf("/vpcs/%d/attachments/%d/bgp", vpcID, attachmentID), req)
	return err
}
func (c *V3Client) DeleteVPCAttachmentBGP(vpcID, attachmentID int) error {
	_, err := c.del(fmt.Sprintf("/vpcs/%d/attachments/%d/bgp", vpcID, attachmentID))
	return err
}
