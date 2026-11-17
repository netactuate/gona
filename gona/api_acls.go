package gona

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// APIACL is one API key allowlist entry on the account.
type APIACL struct {
	ID          int              `json:"id"`
	Type        int              `json:"type"`
	StartIP     string           `json:"start_ip"`
	EndIP       string           `json:"end_ip"`
	Bits        int              `json:"bits"`
	SubID       int              `json:"sub_id"`
	CreatedAt   string           `json:"created_at"`
	Description string           `json:"description"`
	Editable    bool             `json:"editable"`
	Resources   []APIACLResource `json:"resources"`
}

// APIACLResource is a platform resource holding an APIACL entry.
type APIACLResource struct {
	ResourceType  string `json:"resource_type"`
	ResourceID    int    `json:"resource_id"`
	ResourceLabel string `json:"resource_label"`
}

// ListAPIACLs returns the account's API key allowlist entries, or one team member's with subUserID.
func (c *Client) ListAPIACLs(subUserID *int) ([]APIACL, error) {
	path := "account/api-acls"
	if subUserID != nil {
		values := url.Values{}
		values.Set("sub_user_id", strconv.Itoa(*subUserID))
		path += "?" + values.Encode()
	}

	var acls []APIACL
	if err := c.get(context.Background(), path, &acls); err != nil {
		return nil, fmt.Errorf("list API ACLs: %w", err)
	}
	return acls, nil
}
