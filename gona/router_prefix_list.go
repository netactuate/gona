package gona

import (
	"encoding/json"
	"fmt"
)

type PrefixListRule struct {
	Action string `json:"action"`
	Prefix string `json:"prefix"`
}

type RouterPrefixList struct {
	PrefixListID int              `json:"prefixListId"`
	Name         string           `json:"name"`
	IPVersion    int              `json:"ipVersion"`
	Description  string           `json:"description,omitempty"`
	Rules        []PrefixListRule `json:"rules"`
}

type CreateRouterPrefixListRequest struct {
	Name        string           `json:"name"`
	IPVersion   int              `json:"ipVersion"`
	Description string           `json:"description,omitempty"`
	Rules       []PrefixListRule `json:"rules"`
}

type CreateRouterPrefixListResponse struct {
	PrefixListID int `json:"prefixListId"`
}

type UpdateRouterPrefixListRequest struct {
	Name        string           `json:"name"`
	IPVersion   int              `json:"ipVersion"`
	Description string           `json:"description,omitempty"`
	Rules       []PrefixListRule `json:"rules"`
}

type UpdateRouterPrefixListResponse struct {
	PrefixListID int `json:"prefixListId"`
}

func (c *V3Client) CreateRouterPrefixList(routerID int, req *CreateRouterPrefixListRequest) (*CreateRouterPrefixListResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/prefix-lists", routerID)

	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create prefix list for router %d: %w", routerID, err)
	}

	var createResp CreateRouterPrefixListResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create prefix list response: %w", err)
	}

	return &createResp, nil
}

func (c *V3Client) ListRouterPrefixLists(routerID int) ([]RouterPrefixList, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/prefix-lists", routerID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list prefix lists for router %d: %w", routerID, err)
	}

	var lists []RouterPrefixList
	if err := json.Unmarshal(resp.Data, &lists); err != nil {
		return nil, fmt.Errorf("failed to unmarshal prefix lists: %w", err)
	}

	return lists, nil
}

func (c *V3Client) GetRouterPrefixList(routerID, prefixListID int) (*RouterPrefixList, error) {
	lists, err := c.ListRouterPrefixLists(routerID)
	if err != nil {
		return nil, err
	}

	for _, list := range lists {
		if list.PrefixListID == prefixListID {
			return &list, nil
		}
	}

	return nil, fmt.Errorf("prefix list %d not found on router %d", prefixListID, routerID)
}

func (c *V3Client) UpdateRouterPrefixList(routerID, prefixListID int, req *UpdateRouterPrefixListRequest) (*UpdateRouterPrefixListResponse, error) {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/prefix-lists/%d", routerID, prefixListID)

	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update prefix list %d on router %d: %w", prefixListID, routerID, err)
	}

	var updateResp UpdateRouterPrefixListResponse
	if err := json.Unmarshal(resp.Data, &updateResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal update prefix list response: %w", err)
	}

	return &updateResp, nil
}

func (c *V3Client) DeleteRouterPrefixList(routerID, prefixListID int) error {
	path := fmt.Sprintf("/cloud-routing/routers/%d/config/prefix-lists/%d", routerID, prefixListID)

	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("failed to delete prefix list %d on router %d: %w", prefixListID, routerID, err)
	}

	return nil
}
