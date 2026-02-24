package gona

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

func toBool(v interface{}) bool {
	switch v := v.(type) {
	case bool:
		return v
	case float64:
		return v == 1
	case string:
		return v == "1" || v == "true"
	}
	return false
}

func toInt(v interface{}) int {
	switch v := v.(type) {
	case float64:
		return int(v)
	case string:
		n, _ := strconv.Atoi(v)
		return n
	}
	return 0
}

// FirewallSet
type FirewallSet struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Enabled            bool   `json:"-"`
	IsDraft            bool   `json:"-"`
	DraftFirewallSetID *int   `json:"draft_firewall_set_id"`
	Created            string `json:"created"`
	LastUpdated        string `json:"last_updated"`
}

func (s *FirewallSet) UnmarshalJSON(data []byte) error {
	type Alias FirewallSet
	aux := &struct {
		Enabled interface{} `json:"enabled"`
		IsDraft interface{} `json:"is_draft"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	s.Enabled = toBool(aux.Enabled)
	s.IsDraft = toBool(aux.IsDraft)
	return nil
}

// FirewallMatchOptions contains optional match criteria options
type FirewallMatchOptions struct {
	ICMPType string `json:"icmp_type,omitempty"`
}

// FirewallMatchCriteria
type FirewallMatchCriteria struct {
	Protocol             string                `json:"protocol,omitempty"`
	SourceNet            []string              `json:"source_net,omitempty"`
	DestinationNet       []string              `json:"destination_net,omitempty"`
	SourcePortStart      *int                  `json:"source_port_start"`
	SourcePortEnd        *int                  `json:"source_port_end"`
	DestinationPortStart *int                  `json:"destination_port_start"`
	DestinationPortEnd   *int                  `json:"destination_port_end"`
	IPVersionNumber      *int                  `json:"ip_version_number,omitempty"`
	Options              *FirewallMatchOptions `json:"options,omitempty"`
}

// FirewallRule
type FirewallRule struct {
	ID            int                    `json:"id"`
	FirewallSetID int                    `json:"-"`
	IPVersion     string                 `json:"ip_version"`
	Direction     string                 `json:"direction"`
	Action        string                 `json:"action"`
	Enabled       bool                   `json:"-"`
	MatchCriteria *FirewallMatchCriteria `json:"match_criteria"`
	AdminComment  string                 `json:"admin_comment"`
	RulePriority  int                    `json:"rule_priority"`
	Created       string                 `json:"created"`
	LastUpdated   string                 `json:"last_updated"`
}

func (r *FirewallRule) UnmarshalJSON(data []byte) error {
	type Alias FirewallRule
	aux := &struct {
		FirewallSetID interface{} `json:"firewall_set_id"`
		Enabled       interface{} `json:"enabled"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	r.FirewallSetID = toInt(aux.FirewallSetID)
	r.Enabled = toBool(aux.Enabled)
	return nil
}

// CreateFirewallRuleRequest
type CreateFirewallRuleRequest struct {
	IPVersion     string                 `json:"ip_version"`
	Direction     string                 `json:"direction,omitempty"`
	Action        string                 `json:"action"`
	Enabled       bool                   `json:"enabled"`
	RulePriority  *int                   `json:"rule_priority,omitempty"`
	AdminComment  string                 `json:"admin_comment,omitempty"`
	MatchCriteria *FirewallMatchCriteria `json:"match_criteria"`
}

func (c *Client) GetFirewallSets() ([]FirewallSet, error) {
	var sets []FirewallSet
	if err := c.get("firewall/sets", &sets); err != nil {
		return nil, err
	}
	return sets, nil
}

func (c *Client) GetFirewallSet(id int) (FirewallSet, error) {
	var set FirewallSet
	if err := c.get("firewall/sets/"+strconv.Itoa(id), &set); err != nil {
		return FirewallSet{}, err
	}
	return set, nil
}

func (c *Client) CreateFirewallSet(name, description string, enabled bool) (FirewallSet, error) {
	values := url.Values{}
	values.Add("name", name)
	values.Add("description", description)
	if enabled {
		values.Add("enabled", "1")
	} else {
		values.Add("enabled", "0")
	}

	var set FirewallSet
	if err := c.post("firewall/sets", []byte(values.Encode()), &set); err != nil {
		return FirewallSet{}, err
	}
	return set, nil
}

func (c *Client) UpdateFirewallSet(id int, name, description string, enabled bool) (FirewallSet, error) {
	values := url.Values{}
	values.Add("name", name)
	values.Add("description", description)
	if enabled {
		values.Add("enabled", "1")
	} else {
		values.Add("enabled", "0")
	}

	var set FirewallSet
	if err := c.put("firewall/sets/"+strconv.Itoa(id), []byte(values.Encode()), &set); err != nil {
		return FirewallSet{}, err
	}
	return set, nil
}

func (c *Client) DeleteFirewallSet(id int) error {
	return c.delete("firewall/sets/"+strconv.Itoa(id), nil, nil)
}

func (c *Client) EnableFirewallSet(id int) error {
	return c.put("firewall/sets/"+strconv.Itoa(id)+"/enable", []byte{}, nil)
}

func (c *Client) DisableFirewallSet(id int) error {
	return c.put("firewall/sets/"+strconv.Itoa(id)+"/disable", []byte{}, nil)
}

func (c *Client) CreateDraftFirewallSet(id int) (FirewallSet, error) {
	var set FirewallSet
	if err := c.post("firewall/sets/"+strconv.Itoa(id)+"/create-draft", []byte{}, &set); err != nil {
		return FirewallSet{}, err
	}
	return set, nil
}

func (c *Client) PublishDraftFirewallSet(draftID int) (FirewallSet, error) {
	var set FirewallSet
	if err := c.post("firewall/sets/publish-draft/"+strconv.Itoa(draftID), []byte{}, &set); err != nil {
		return FirewallSet{}, err
	}
	return set, nil
}

func (c *Client) DeleteDraftFirewallSet(draftID int) error {
	return c.delete("firewall/sets/delete-draft/"+strconv.Itoa(draftID), nil, nil)
}

func (c *Client) SyncFirewallSetRules(setID int) error {
	return c.post("firewall/sets/"+strconv.Itoa(setID)+"/vm/sync-all", []byte{}, nil)
}

func (c *Client) GetFirewallRules(setID int) ([]FirewallRule, error) {
	var rules []FirewallRule
	path := fmt.Sprintf("firewall/sets/%d/rules", setID)
	if err := c.get(path, &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

func (c *Client) GetFirewallRule(setID, ruleID int) (FirewallRule, error) {
	var rule FirewallRule
	path := fmt.Sprintf("firewall/sets/%d/rules/%d", setID, ruleID)
	if err := c.get(path, &rule); err != nil {
		return FirewallRule{}, err
	}
	return rule, nil
}

func (c *Client) CreateFirewallRule(setID int, req *CreateFirewallRuleRequest) (FirewallRule, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return FirewallRule{}, fmt.Errorf("encoding rule: %w", err)
	}

	var rule FirewallRule
	path := fmt.Sprintf("firewall/sets/%d/rules", setID)
	if err := c.postJSON(path, body, &rule); err != nil {
		return FirewallRule{}, err
	}
	return rule, nil
}

func (c *Client) UpdateFirewallRule(setID, ruleID int, req *CreateFirewallRuleRequest) (FirewallRule, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return FirewallRule{}, fmt.Errorf("encoding rule: %w", err)
	}

	var rule FirewallRule
	path := fmt.Sprintf("firewall/%d/%d", setID, ruleID)
	if err := c.putJSON(path, body, &rule); err != nil {
		return FirewallRule{}, err
	}
	return rule, nil
}

func (c *Client) DeleteFirewallRule(setID, ruleID int) error {
	path := fmt.Sprintf("firewall/%d/rules/%d", setID, ruleID)
	return c.delete(path, nil, nil)
}

// FirewallSetVM represents a VM attached to a firewall set
type FirewallSetVM struct {
	ID            int    `json:"-"`
	Mbpkgid       int    `json:"-"`
	InterfaceID   int    `json:"-"`
	FirewallSetID int    `json:"-"`
	SetPriority   int    `json:"-"`
	Created       string `json:"created"`
	LastUpdated   string `json:"last_updated"`
	IATACode      string `json:"iata_code"`
	Location      string `json:"location"`
	Hostname      string `json:"hostname"`
}

func (v *FirewallSetVM) UnmarshalJSON(data []byte) error {
	type Alias FirewallSetVM
	aux := &struct {
		ID            interface{} `json:"id"`
		Mbpkgid       interface{} `json:"mbpkgid"`
		InterfaceID   interface{} `json:"interface_id"`
		FirewallSetID interface{} `json:"firewall_set_id"`
		SetPriority   interface{} `json:"set_priority"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	v.ID = toInt(aux.ID)
	v.Mbpkgid = toInt(aux.Mbpkgid)
	v.InterfaceID = toInt(aux.InterfaceID)
	v.FirewallSetID = toInt(aux.FirewallSetID)
	v.SetPriority = toInt(aux.SetPriority)
	return nil
}

type attachVMRequest struct {
	VMList []attachVMEntry `json:"vm_list"`
}

type attachVMEntry struct {
	Mbpkgid     int `json:"mbpkgid"`
	InterfaceID int `json:"interface_id"`
	SetPriority int `json:"set_priority"`
}

func (c *Client) GetFirewallSetVMs(setID int) ([]FirewallSetVM, error) {
	var vms []FirewallSetVM
	path := fmt.Sprintf("firewall/sets/%d/vm-list", setID)
	if err := c.get(path, &vms); err != nil {
		return nil, err
	}
	return vms, nil
}

func (c *Client) AttachFirewallSetVM(setID, mbpkgid, interfaceID, setPriority int) ([]FirewallSetVM, error) {
	req := attachVMRequest{
		VMList: []attachVMEntry{
			{
				Mbpkgid:     mbpkgid,
				InterfaceID: interfaceID,
				SetPriority: setPriority,
			},
		},
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("encoding attach request: %w", err)
	}

	var vms []FirewallSetVM
	path := fmt.Sprintf("firewall/sets/%d/vm/attach", setID)
	if err := c.postJSON(path, body, &vms); err != nil {
		return nil, err
	}
	return vms, nil
}

func (c *Client) DetachFirewallSetVM(setID, mbpkgid int) error {
	path := fmt.Sprintf("firewall/sets/%d/vm/detach/%d", setID, mbpkgid)
	return c.post(path, []byte{}, nil)
}
