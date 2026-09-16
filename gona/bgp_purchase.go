package gona

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// BGPGroup is an account BGP group returned by the vAPI2 BGP endpoints.
type BGPGroup struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupType   string `json:"group_type"`
}

// CreateBGPGroupRequest contains the fields required to create a BGP group.
type CreateBGPGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupType   string `json:"group_type,omitempty"`
}

// BGPPrefix is an account BGP prefix returned by the vAPI2 BGP endpoints.
type BGPPrefix struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Prefix         string `json:"prefix"`
	GroupID        int    `json:"group_id"`
	ASNID          int    `json:"asn_id"`
	AnycastProfile int    `json:"anycast_profile"`
	AgreementID    int    `json:"agreement_id"`
}

// BuyBGPPrefixesRequest contains the fields for purchasing anycast BGP prefixes.
type BuyBGPPrefixesRequest struct {
	Name           string `json:"name"`
	GroupID        int    `json:"group_id,omitempty"`
	ASNID          int    `json:"asn_id,omitempty"`
	AnycastProfile int    `json:"anycast_profile,omitempty"`
	AgreementID    int    `json:"agreement_id"`
}

// BGPASN is an account ASN returned by the vAPI2 BGP endpoints.
type BGPASN struct {
	ID        int    `json:"id"`
	ASN       int    `json:"asn"`
	Name      string `json:"name"`
	GroupType string `json:"group_type"`
}

// AccountAgreement is a legal agreement returned by the account agreements endpoint.
type AccountAgreement struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

// CreateBGPGroup creates an account BGP group.
func (c *Client) CreateBGPGroup(req *CreateBGPGroupRequest) (*BGPGroup, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal BGP group create request: %w", err)
	}
	var group BGPGroup
	if err := c.postJSON(context.Background(), "bgp/bgpgroup", body, &group); err != nil {
		return nil, fmt.Errorf("create BGP group: %w", err)
	}
	return &group, nil
}

// GetBGPGroup returns an account BGP group by ID.
func (c *Client) GetBGPGroup(id int) (*BGPGroup, error) {
	var group BGPGroup
	if err := c.get(context.Background(), "bgp/bgpgroup/"+strconv.Itoa(id), &group); err != nil {
		return nil, fmt.Errorf("get BGP group %d: %w", id, err)
	}
	return &group, nil
}

// ListBGPGroups returns account BGP groups, optionally filtered by group type.
func (c *Client) ListBGPGroups(groupType string) ([]BGPGroup, error) {
	path := "bgp/bgpgroups"
	if groupType != "" {
		values := url.Values{}
		values.Set("group_type", groupType)
		path += "?" + values.Encode()
	}
	var groups []BGPGroup
	if err := c.get(context.Background(), path, &groups); err != nil {
		return nil, fmt.Errorf("list BGP groups: %w", err)
	}
	return groups, nil
}

// BuyBGPPrefixes purchases anycast BGP prefixes for the account.
func (c *Client) BuyBGPPrefixes(req *BuyBGPPrefixesRequest) (*BGPPrefix, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal BGP prefix purchase request: %w", err)
	}
	var prefix BGPPrefix
	if err := c.postJSON(context.Background(), "bgp/bgpbuyprefixes", body, &prefix); err != nil {
		return nil, fmt.Errorf("buy BGP prefixes: %w", err)
	}
	return &prefix, nil
}

// GetBGPPrefix returns an account BGP prefix by ID.
func (c *Client) GetBGPPrefix(id int) (*BGPPrefix, error) {
	var prefix BGPPrefix
	if err := c.get(context.Background(), "bgp/bgpprefix/"+strconv.Itoa(id), &prefix); err != nil {
		return nil, fmt.Errorf("get BGP prefix %d: %w", id, err)
	}
	return &prefix, nil
}

// ListBGPPrefixes returns account BGP prefixes, optionally filtered by group type.
func (c *Client) ListBGPPrefixes(groupType string) ([]BGPPrefix, error) {
	path := "bgp/bgpprefixes"
	if groupType != "" {
		values := url.Values{}
		values.Set("group_type", groupType)
		path += "?" + values.Encode()
	}
	var prefixes []BGPPrefix
	if err := c.get(context.Background(), path, &prefixes); err != nil {
		return nil, fmt.Errorf("list BGP prefixes: %w", err)
	}
	return prefixes, nil
}

// ListBGPASNs returns account ASNs, optionally filtered by group type.
func (c *Client) ListBGPASNs(groupType string) ([]BGPASN, error) {
	path := "bgp/bgpasns"
	if groupType != "" {
		values := url.Values{}
		values.Set("group_type", groupType)
		path += "?" + values.Encode()
	}
	var asns []BGPASN
	if err := c.get(context.Background(), path, &asns); err != nil {
		return nil, fmt.Errorf("list BGP ASNs: %w", err)
	}
	return asns, nil
}

// GetBGPASN returns an account ASN by ID.
func (c *Client) GetBGPASN(id int) (*BGPASN, error) {
	values := url.Values{}
	values.Set("id", strconv.Itoa(id))
	var asn BGPASN
	if err := c.get(context.Background(), "bgp/bgpasn?"+values.Encode(), &asn); err != nil {
		return nil, fmt.Errorf("get BGP ASN %d: %w", id, err)
	}
	return &asn, nil
}

// ListAccountAgreements returns legal agreements available to the account.
func (c *Client) ListAccountAgreements() ([]AccountAgreement, error) {
	var agreements []AccountAgreement
	if err := c.get(context.Background(), "account/agreements", &agreements); err != nil {
		return nil, fmt.Errorf("list account agreements: %w", err)
	}
	return agreements, nil
}
