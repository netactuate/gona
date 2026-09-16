package gona

import (
	"context"
	"encoding/json"
	"fmt"
)

// AccessControlSubnet models the user access control subnet shape.
// The live account has zero rows, so this has not been measured from a response.
type AccessControlSubnet struct {
	ID     json.Number `json:"id"`
	Label  string      `json:"label"`
	Subnet string      `json:"subnet"`
}

type CreateAccessControlSubnetRequest struct {
	Label  string `json:"label"`
	Subnet string `json:"subnet"`
}

type UpdateAccessControlSubnetRequest struct {
	Label  *string `json:"label,omitempty"`
	Subnet *string `json:"subnet,omitempty"`
}

func (c *Client) GetAccessControlSubnets() ([]AccessControlSubnet, error) {
	var subnets []AccessControlSubnet
	if err := c.get(context.Background(), "account/user-access-control-subnet-list", &subnets); err != nil {
		return nil, err
	}
	return subnets, nil
}

func (c *Client) GetAccessControlSubnet(id string) (*AccessControlSubnet, error) {
	var subnet AccessControlSubnet
	if err := c.get(context.Background(), "account/user-access-control-subnet/"+id, &subnet); err != nil {
		return nil, err
	}
	return &subnet, nil
}

func (c *Client) CreateAccessControlSubnet(req *CreateAccessControlSubnetRequest) (*AccessControlSubnet, error) {
	values, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("create access control subnet marshal: %w", err)
	}

	var subnet AccessControlSubnet
	if err := c.postJSON(context.Background(), "account/user-access-control-subnet", values, &subnet); err != nil {
		return nil, err
	}
	return &subnet, nil
}

func (c *Client) UpdateAccessControlSubnet(id string, req *UpdateAccessControlSubnetRequest) (*AccessControlSubnet, error) {
	values, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("update access control subnet %s marshal: %w", id, err)
	}

	var subnet AccessControlSubnet
	if err := c.patch(context.Background(), "account/user-access-control-subnet/"+id, values, &subnet); err != nil {
		return nil, err
	}
	return &subnet, nil
}

func (c *Client) DeleteAccessControlSubnet(id string) error {
	return c.delete(context.Background(), "account/user-access-control-subnet/"+id, nil, nil)
}
