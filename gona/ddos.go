package gona

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

type DDoSAttack struct {
	AttackID    int    `json:"id"`
	DateStart   string `json:"date_start"`
	DateEnd     string `json:"date_end"`
	Status      int    `json:"status"`
	IP          string `json:"ip"`
	Prefix      string `json:"prefix"`
	Direction   string `json:"direction"`
	PPS         int    `json:"pps"`
	RuleID      int    `json:"rule_id"`
	RuleType    string `json:"rule_type"`
	BanDuration int    `json:"ban_duration"`
	RuleName    string `json:"rule_name"`
}

func (c *Client) GetDDoSAttacks() ([]DDoSAttack, error) {
	var attacks []DDoSAttack
	if err := c.get(context.Background(), "ddos/attacks", &attacks); err != nil {
		return nil, fmt.Errorf("get ddos attacks: %w", err)
	}
	return attacks, nil
}

func (c *Client) GetDDoSActiveAttacks() ([]DDoSAttack, error) {
	var attacks []DDoSAttack
	if err := c.get(context.Background(), "ddos/attacks/active", &attacks); err != nil {
		return nil, fmt.Errorf("get active ddos attacks: %w", err)
	}
	return attacks, nil
}

type DDoSDashboardOptions struct {
	Period       *int
	IncludeEnded *bool
	Limit        *int
}

type DDoSDashboard struct {
	TotalAttacks         int                   `json:"total_attacks"`
	ActiveRules          int                   `json:"active_rules"`
	LongestAttackSeconds int                   `json:"longest_attack_seconds"`
	TopAttacks           []DDoSDashboardAttack `json:"top_attacks"`
	Period               int                   `json:"period"`
}

type DDoSDashboardAttack struct{}

func (c *Client) GetDDoSDashboard(opts DDoSDashboardOptions) (DDoSDashboard, error) {
	path := "ddos/dashboard"
	values := url.Values{}
	if opts.Period != nil {
		values.Set("period", strconv.Itoa(*opts.Period))
	}
	if opts.IncludeEnded != nil {
		values.Set("include_ended", strconv.FormatBool(*opts.IncludeEnded))
	}
	if opts.Limit != nil {
		values.Set("limit", strconv.Itoa(*opts.Limit))
	}
	if encoded := values.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var dashboard DDoSDashboard
	if err := c.get(context.Background(), path, &dashboard); err != nil {
		return DDoSDashboard{}, fmt.Errorf("get ddos dashboard: %w", err)
	}
	return dashboard, nil
}

type DDoSRule struct {
	RuleID      int              `json:"id"`
	RuleName    string           `json:"rule_name"`
	Description string           `json:"description"`
	Prefixes    []DDoSRulePrefix `json:"prefixes"`
	Rules       []DDoSRuleAction `json:"rules"`
}

type DDoSRulePrefix struct {
	PrefixID    int    `json:"id"`
	Prefix      string `json:"prefix"`
	PrefixType  string `json:"prefix_type"`
	Description string `json:"description"`
	AllowedPPS  int    `json:"allowed_pps"`
}

type DDoSRuleAction struct {
	Name              string `json:"name"`
	ActionType        string `json:"action_type"`
	RunOrder          int    `json:"run_order"`
	ActionName        string `json:"action_name"`
	ActionDescription string `json:"action_description"`
}

func (c *Client) GetDDoSRules() ([]DDoSRule, error) {
	var rules []DDoSRule
	if err := c.get(context.Background(), "ddos/rules", &rules); err != nil {
		return nil, fmt.Errorf("get ddos rules: %w", err)
	}
	return rules, nil
}

// GetDDoSRule returns a DDoS rule by id.
func (c *Client) GetDDoSRule(id int) (*DDoSRule, error) {
	var rule DDoSRule
	if err := c.get(context.Background(), "ddos/rule/"+strconv.Itoa(id), &rule); err != nil {
		return nil, fmt.Errorf("get ddos rule %d: %w", id, err)
	}
	return &rule, nil
}
