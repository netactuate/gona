package gona

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// BindBGPGroupFirewallSetRequest contains the fields for binding a firewall set
// to a BGP group interface.
type BindBGPGroupFirewallSetRequest struct {
	ID              int `json:"id"`
	FirewallSetID   int `json:"firewall_set_id"`
	InterfaceNumber int `json:"interface_number"`
	SetPriority     int `json:"set_priority"`
}

// BGPGroupFirewallSetBinding is a firewall set binding returned for a BGP group.
type BGPGroupFirewallSetBinding struct {
	ID              int `json:"id"`
	BGPGroupID      int `json:"bgp2_group_id"`
	FirewallSetID   int `json:"firewall_set_id"`
	InterfaceNumber int `json:"interface_number"`
	SetPriority     int `json:"set_priority"`
}

// BGPSummary is the summary payload returned by the BGP summary endpoint.
type BGPSummary map[string]json.RawMessage

// BGPDashboardOptions contains optional filters for the BGP dashboard endpoint.
type BGPDashboardOptions struct {
	GroupType  string
	FlapWindow *int
}

// BGPDashboard is the dashboard payload returned by the BGP dashboard endpoint.
type BGPDashboard map[string]json.RawMessage

// BindBGPGroupFirewallSet binds a firewall set to a BGP group.
func (c *Client) BindBGPGroupFirewallSet(groupID int, req *BindBGPGroupFirewallSetRequest) (*BGPGroupFirewallSetBinding, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal BGP group firewall set bind request: %w", err)
	}
	var binding BGPGroupFirewallSetBinding
	path := fmt.Sprintf("bgp/bgp-groups/%d/firewall-sets", groupID)
	if err := c.postJSON(context.Background(), path, body, &binding); err != nil {
		return nil, fmt.Errorf("bind firewall set to BGP group %d: %w", groupID, err)
	}
	return &binding, nil
}

// UnbindBGPGroupFirewallSet removes a firewall set binding from a BGP group.
func (c *Client) UnbindBGPGroupFirewallSet(groupID, firewallSetID int) error {
	path := fmt.Sprintf("bgp/bgp-groups/%d/firewall-sets/%d", groupID, firewallSetID)
	if err := c.delete(context.Background(), path, nil, nil); err != nil {
		return fmt.Errorf("unbind firewall set %d from BGP group %d: %w", firewallSetID, groupID, err)
	}
	return nil
}

// RefreshBGPGroupSessions asks the platform to refresh all sessions in a BGP group.
func (c *Client) RefreshBGPGroupSessions(groupID int) error {
	return c.postBGPGroupAction(groupID, "refresh")
}

// StartBGPGroupSessions asks the platform to start all sessions in a BGP group.
func (c *Client) StartBGPGroupSessions(groupID int) error {
	return c.postBGPGroupAction(groupID, "start")
}

// StopBGPGroupSessions asks the platform to stop all sessions in a BGP group.
func (c *Client) StopBGPGroupSessions(groupID int) error {
	return c.postBGPGroupAction(groupID, "stop")
}

// RefreshBGPSession asks the platform to refresh a BGP session.
func (c *Client) RefreshBGPSession(sessionID int) error {
	return c.postBGPSessionAction(sessionID, "refresh")
}

// StartBGPSession asks the platform to start a BGP session.
func (c *Client) StartBGPSession(sessionID int) error {
	return c.postBGPSessionAction(sessionID, "start")
}

// StopBGPSession asks the platform to stop a BGP session.
func (c *Client) StopBGPSession(sessionID int) error {
	return c.postBGPSessionAction(sessionID, "stop")
}

// GetBGPSummary returns the account BGP summary payload.
func (c *Client) GetBGPSummary() (BGPSummary, error) {
	var summary BGPSummary
	if err := c.get(context.Background(), "bgp/bgpsummary", &summary); err != nil {
		return nil, fmt.Errorf("get BGP summary: %w", err)
	}
	return summary, nil
}

// GetBGPDashboard returns BGP dashboard data, optionally filtered by group type
// and flap window.
func (c *Client) GetBGPDashboard(opts BGPDashboardOptions) (BGPDashboard, error) {
	path := "bgp/dashboard"
	values := url.Values{}
	if opts.GroupType != "" {
		values.Set("group_type", opts.GroupType)
	}
	if opts.FlapWindow != nil {
		values.Set("flap_window", strconv.Itoa(*opts.FlapWindow))
	}
	if encoded := values.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var dashboard BGPDashboard
	if err := c.get(context.Background(), path, &dashboard); err != nil {
		return nil, fmt.Errorf("get BGP dashboard: %w", err)
	}
	return dashboard, nil
}

func (c *Client) postBGPGroupAction(groupID int, action string) error {
	path := fmt.Sprintf("bgp/bgpgroup/%d/%s", groupID, action)
	if err := c.post(context.Background(), path, nil, nil); err != nil {
		return fmt.Errorf("%s BGP group sessions %d: %w", action, groupID, err)
	}
	return nil
}

func (c *Client) postBGPSessionAction(sessionID int, action string) error {
	path := fmt.Sprintf("bgp/bgpsession/%d/%s", sessionID, action)
	if err := c.post(context.Background(), path, nil, nil); err != nil {
		return fmt.Errorf("%s BGP session %d: %w", action, sessionID, err)
	}
	return nil
}
