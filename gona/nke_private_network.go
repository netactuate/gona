package gona

import (
	"encoding/json"
	"fmt"
)

type NKEPrivateNetworkNode struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	PrivateNICState  string `json:"privateNicState"`
	PrivateIPAddress string `json:"privateIpAddress"`
}
type NKEPrivateNetwork struct {
	ClusterID          int                     `json:"clusterId"`
	VPCID              int                     `json:"vpcId"`
	State              string                  `json:"state"`
	Nodes              []NKEPrivateNetworkNode `json:"nodes"`
	NodesRequiringRoll []int                   `json:"nodesRequiringRoll"`
}
type AddNKEPrivateNetworkRequest struct {
	VPCID int `json:"vpcId"`
}

func (c *V3Client) GetNKEPrivateNetwork(clusterID int) (*NKEPrivateNetwork, error) {
	r, err := c.get(fmt.Sprintf("/nke/clusters/%d/private-network", clusterID))
	if err != nil {
		return nil, err
	}
	var out NKEPrivateNetwork
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *V3Client) AddNKEPrivateNetwork(clusterID int, req *AddNKEPrivateNetworkRequest) (int, error) {
	r, err := c.post(fmt.Sprintf("/nke/clusters/%d/private-network", clusterID), req)
	if err != nil {
		return 0, err
	}
	var out struct {
		ID int `json:"clusterId"`
	}
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}
