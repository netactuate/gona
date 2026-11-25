package gona

import (
	"encoding/json"
	"fmt"
)

type VPCSubnet struct {
	SubnetID          int    `json:"subnetId"`
	Network           string `json:"network"`
	PrefixLength      int    `json:"prefixLength"`
	Label             string `json:"label"`
	IsDefault         bool   `json:"isDefault"`
	State             string `json:"state"`
	StateMessage      string `json:"stateMessage"`
	CreatedOn         string `json:"createdOn"`
	ReadyOn           string `json:"readyOn"`
	DeleteRequestedOn string `json:"deleteRequestedOn"`
	DeletedOn         string `json:"deletedOn"`
}
type CreateVPCSubnetRequest struct {
	Network string `json:"network"`
	Label   string `json:"label"`
}

func (c *V3Client) ListVPCSubnets(vpcID int) ([]VPCSubnet, error) {
	r, err := c.getListUnder(fmt.Sprintf("/vpcs/%d/subnets", vpcID), "")
	if err != nil {
		return nil, err
	}
	var out []VPCSubnet
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *V3Client) GetVPCSubnet(vpcID, subnetID int) (*VPCSubnet, error) {
	r, err := c.get(fmt.Sprintf("/vpcs/%d/subnets/%d", vpcID, subnetID))
	if err != nil {
		return nil, err
	}
	var out VPCSubnet
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
func (c *V3Client) CreateVPCSubnet(vpcID int, req *CreateVPCSubnetRequest) (int, error) {
	r, err := c.post(fmt.Sprintf("/vpcs/%d/subnets", vpcID), req)
	if err != nil {
		return 0, err
	}
	var out struct {
		ID int `json:"subnetId"`
	}
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}
func (c *V3Client) DeleteVPCSubnet(vpcID, subnetID int) error {
	_, err := c.del(fmt.Sprintf("/vpcs/%d/subnets/%d", vpcID, subnetID))
	return err
}
