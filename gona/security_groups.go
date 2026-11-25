package gona

import (
	"encoding/json"
	"fmt"
)

type SecurityGroup struct {
	SecurityGroupID int                      `json:"securityGroupId"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description"`
	IsDefault       bool                     `json:"isDefault"`
	DefaultForVPCID int                      `json:"defaultForVpcId"`
	MemberCount     int                      `json:"memberCount"`
	CreatedOn       string                   `json:"createdOn"`
	ModifiedOn      string                   `json:"modifiedOn"`
	Rules           []SecurityGroupRule      `json:"rules"`
	References      []SecurityGroupReference `json:"references"`
}
type SecurityGroupReference struct {
	SecurityGroupID int `json:"securityGroupId"`
	RuleID          int `json:"ruleId"`
}
type CreateSecurityGroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
type UpdateSecurityGroupRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}
type SecurityGroupPort struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
type SecurityGroupRule struct {
	RuleID          int                `json:"ruleId"`
	IPVersion       int                `json:"ipVersion"`
	Direction       string             `json:"direction"`
	Protocol        string             `json:"protocol"`
	Port            *SecurityGroupPort `json:"port,omitempty"`
	PortStart       int                `json:"portStart"`
	PortEnd         int                `json:"portEnd"`
	Network         string             `json:"network,omitempty"`
	SecurityGroupID int                `json:"securityGroupId,omitempty"`
	Description     string             `json:"description,omitempty"`
	CreatedOn       string             `json:"createdOn"`
}
type CreateSecurityGroupRuleRequest struct {
	IPVersion       int                `json:"ipVersion"`
	Direction       string             `json:"direction"`
	Protocol        string             `json:"protocol,omitempty"`
	Port            *SecurityGroupPort `json:"port,omitempty"`
	Network         string             `json:"network,omitempty"`
	SecurityGroupID int                `json:"securityGroupId,omitempty"`
	Description     string             `json:"description,omitempty"`
}
type SecurityGroupMember struct {
	MemberID         int    `json:"memberId"`
	VPCID            int    `json:"vpcId"`
	MBPkgID          int    `json:"mbpkgid"`
	NodePoolID       int    `json:"nodePoolId"`
	CreatedOn        string `json:"createdOn"`
	EnforcementState string `json:"enforcementState"`
}
type CreateSecurityGroupMemberRequest struct {
	VPCID      int `json:"vpcId"`
	MBPkgID    int `json:"mbpkgid,omitempty"`
	NodePoolID int `json:"nodePoolId,omitempty"`
}
type VPCSecurityGroupEnforcement struct {
	VPCID     int    `json:"vpcId"`
	Enabled   bool   `json:"enabled"`
	EnabledOn string `json:"enabledOn"`
}
type SetVPCSecurityGroupEnforcementRequest struct {
	Enabled bool `json:"enabled"`
}

func (c *V3Client) ListSecurityGroups() ([]SecurityGroup, error) {
	r, err := c.getListUnder("/cloud/networking/security-groups", "")
	if err != nil {
		return nil, err
	}
	var out []SecurityGroup
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *V3Client) GetSecurityGroup(groupID int) (*SecurityGroup, error) {
	r, err := c.get(fmt.Sprintf("/cloud/networking/security-groups/%d", groupID))
	if err != nil {
		return nil, err
	}
	var out SecurityGroup
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
func (c *V3Client) CreateSecurityGroup(req *CreateSecurityGroupRequest) (int, error) {
	r, err := c.post("/cloud/networking/security-groups", req)
	if err != nil {
		return 0, err
	}
	var out struct {
		ID int `json:"securityGroupId"`
	}
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}
func (c *V3Client) UpdateSecurityGroup(groupID int, req *UpdateSecurityGroupRequest) error {
	_, err := c.patch(fmt.Sprintf("/cloud/networking/security-groups/%d", groupID), req)
	return err
}
func (c *V3Client) DeleteSecurityGroup(groupID int) error {
	_, err := c.del(fmt.Sprintf("/cloud/networking/security-groups/%d", groupID))
	return err
}
func (c *V3Client) ListSecurityGroupRules(groupID int) ([]SecurityGroupRule, error) {
	r, err := c.getListUnder(fmt.Sprintf("/cloud/networking/security-groups/%d/rules", groupID), "")
	if err != nil {
		return nil, err
	}
	var out []SecurityGroupRule
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *V3Client) CreateSecurityGroupRule(groupID int, req *CreateSecurityGroupRuleRequest) (int, error) {
	r, err := c.post(fmt.Sprintf("/cloud/networking/security-groups/%d/rules", groupID), req)
	if err != nil {
		return 0, err
	}
	var out struct {
		ID int `json:"ruleId"`
	}
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}
func (c *V3Client) DeleteSecurityGroupRule(groupID, entryID int) error {
	_, err := c.del(fmt.Sprintf("/cloud/networking/security-groups/%d/rules/%d", groupID, entryID))
	return err
}
func (c *V3Client) ListSecurityGroupMembers(groupID int) ([]SecurityGroupMember, error) {
	r, err := c.getListUnder(fmt.Sprintf("/cloud/networking/security-groups/%d/members", groupID), "")
	if err != nil {
		return nil, err
	}
	var out []SecurityGroupMember
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *V3Client) CreateSecurityGroupMember(groupID int, req *CreateSecurityGroupMemberRequest) (int, error) {
	r, err := c.post(fmt.Sprintf("/cloud/networking/security-groups/%d/members", groupID), req)
	if err != nil {
		return 0, err
	}
	var out struct {
		ID int `json:"memberId"`
	}
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}
func (c *V3Client) DeleteSecurityGroupMember(groupID, entryID int) error {
	_, err := c.del(fmt.Sprintf("/cloud/networking/security-groups/%d/members/%d", groupID, entryID))
	return err
}
func (c *V3Client) ListVPCSecurityGroups(vpcID int) ([]SecurityGroup, error) {
	r, err := c.getListUnder(fmt.Sprintf("/vpcs/%d/security-groups", vpcID), "")
	if err != nil {
		return nil, err
	}
	var out []SecurityGroup
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return out, nil
}
func (c *V3Client) GetVPCSecurityGroupEnforcement(vpcID int) (*VPCSecurityGroupEnforcement, error) {
	r, err := c.get(fmt.Sprintf("/vpcs/%d/security-groups/enforcement", vpcID))
	if err != nil {
		return nil, err
	}
	var out VPCSecurityGroupEnforcement
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
func (c *V3Client) SetVPCSecurityGroupEnforcement(vpcID int, req *SetVPCSecurityGroupEnforcementRequest) (*VPCSecurityGroupEnforcement, error) {
	r, err := c.put(fmt.Sprintf("/vpcs/%d/security-groups/enforcement", vpcID), req)
	if err != nil {
		return nil, err
	}
	var out VPCSecurityGroupEnforcement
	if err = json.Unmarshal(r.Data, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
