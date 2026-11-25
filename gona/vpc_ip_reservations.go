package gona

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type CreateVPCIPReservationRequest struct {
	Description string `json:"description"`
	Address     string `json:"address,omitempty"`
}
type VPCTypedIPReservation struct {
	IPReservationID int    `json:"ipReservationId"`
	Address         string `json:"address"`
	Description     string `json:"description"`
	ReservedOn      string `json:"reservedOn"`
	InUse           bool   `json:"inUse"`
}

func (c *V3Client) CreateVPCIPReservation(vpcID int, req *CreateVPCIPReservationRequest) (int, error) {
	r, err := c.post(fmt.Sprintf("/vpcs/%d/ip-reservations", vpcID), req)
	if err != nil {
		return 0, err
	}
	var out struct {
		ID int `json:"ipReservationId"`
	}
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}

func (c *V3Client) DeleteVPCIPReservation(vpcID, reservationID int) error {
	_, err := c.del(fmt.Sprintf("/vpcs/%d/ip-reservations/%d", vpcID, reservationID))
	if err != nil {
		return err
	}
	return nil
}

func (c *V3Client) ListVPCIPReservationsByType(vpcID int, reservationType string) ([]VPCTypedIPReservation, error) {
	r, err := c.getListUnder(fmt.Sprintf("/vpcs/%d/ip-reservations?type=%s", vpcID, url.QueryEscape(reservationType)), "reservations")
	if err != nil {
		return nil, err
	}
	var out []VPCTypedIPReservation
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}
