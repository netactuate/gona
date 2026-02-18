package gona

import (
	"encoding/json"
	"fmt"
)

type VPCBackendTemplate struct {
	BackendTemplateID int          `json:"backendTemplateId"`
	Name              string       `json:"name,omitempty"`
	Description       string       `json:"description,omitempty"`
	BackendHosts      []VPCBackend `json:"backendHosts,omitempty"`
}

type VPCBackend struct {
	BackendHostID   int    `json:"backendHostId,omitempty"`
	Name            string `json:"name,omitempty"`
	Address         string `json:"address,omitempty"`
	InternalAddress string `json:"internalAddress,omitempty"`
}

type CreateVPCBackendTemplateRequest struct {
	Name         string       `json:"name,omitempty"`
	Description  string       `json:"description,omitempty"`
	BackendHosts []VPCBackend `json:"backendHosts,omitempty"`
}

type UpdateVPCBackendTemplateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ReplaceVPCBackendTemplateRequest struct {
	Name         string       `json:"name,omitempty"`
	Description  string       `json:"description,omitempty"`
	BackendHosts []VPCBackend `json:"backendHosts"`
}

type CreateVPCBackendRequest struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address"`
}

type UpdateVPCBackendRequest struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"address,omitempty"`
}

func (c *V3Client) CreateVPCBackendTemplate(vpcID int, req *CreateVPCBackendTemplateRequest) (*VPCBackendTemplate, error) {
	path := fmt.Sprintf("/vpcs/%d/backend-templates", vpcID)
	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("create backend template for VPC %d: %w", vpcID, err)
	}
	var tmpl VPCBackendTemplate
	if err := json.Unmarshal(resp.Data, &tmpl); err != nil {
		return nil, fmt.Errorf("create backend template unmarshal: %w", err)
	}
	return &tmpl, nil
}

func (c *V3Client) GetVPCBackendTemplate(vpcID, templateID int) (*VPCBackendTemplate, error) {
	path := fmt.Sprintf("/vpcs/%d/backend-templates/%d", vpcID, templateID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get backend template %d for VPC %d: %w", templateID, vpcID, err)
	}
	var tmpl VPCBackendTemplate
	if err := json.Unmarshal(resp.Data, &tmpl); err != nil {
		return nil, fmt.Errorf("get backend template unmarshal: %w", err)
	}
	return &tmpl, nil
}

func (c *V3Client) ListVPCBackendTemplates(vpcID int) ([]VPCBackendTemplate, error) {
	path := fmt.Sprintf("/vpcs/%d/backend-templates", vpcID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("list backend templates for VPC %d: %w", vpcID, err)
	}
	var listData V3ListData
	if err := json.Unmarshal(resp.Data, &listData); err != nil {
		return nil, fmt.Errorf("list backend templates unmarshal outer: %w", err)
	}
	var tmpls []VPCBackendTemplate
	if err := json.Unmarshal(listData.Data, &tmpls); err != nil {
		return nil, fmt.Errorf("list backend templates unmarshal inner: %w", err)
	}
	return tmpls, nil
}

func (c *V3Client) UpdateVPCBackendTemplate(vpcID, templateID int, req *UpdateVPCBackendTemplateRequest) (*VPCBackendTemplate, error) {
	path := fmt.Sprintf("/vpcs/%d/backend-templates/%d", vpcID, templateID)
	resp, err := c.patch(path, req)
	if err != nil {
		return nil, fmt.Errorf("update backend template %d for VPC %d: %w", templateID, vpcID, err)
	}
	var tmpl VPCBackendTemplate
	if err := json.Unmarshal(resp.Data, &tmpl); err != nil {
		return nil, fmt.Errorf("update backend template unmarshal: %w", err)
	}
	return &tmpl, nil
}

func (c *V3Client) ReplaceVPCBackendTemplate(vpcID, templateID int, req *ReplaceVPCBackendTemplateRequest) (*VPCBackendTemplate, error) {
	path := fmt.Sprintf("/vpcs/%d/backend-templates/%d", vpcID, templateID)
	resp, err := c.put(path, req)
	if err != nil {
		return nil, fmt.Errorf("replace backend template %d for VPC %d: %w", templateID, vpcID, err)
	}
	var tmpl VPCBackendTemplate
	if err := json.Unmarshal(resp.Data, &tmpl); err != nil {
		return nil, fmt.Errorf("replace backend template unmarshal: %w", err)
	}
	return &tmpl, nil
}

func (c *V3Client) DeleteVPCBackendTemplate(vpcID, templateID int) error {
	path := fmt.Sprintf("/vpcs/%d/backend-templates/%d", vpcID, templateID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete backend template %d for VPC %d: %w", templateID, vpcID, err)
	}
	return nil
}

func (c *V3Client) CreateVPCBackend(vpcID, templateID int, req *CreateVPCBackendRequest) (*VPCBackend, error) {
	path := fmt.Sprintf("/vpcs/%d/backend-templates/%d/backends", vpcID, templateID)
	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("create backend for template %d VPC %d: %w", templateID, vpcID, err)
	}
	var backend VPCBackend
	if err := json.Unmarshal(resp.Data, &backend); err != nil {
		return nil, fmt.Errorf("create backend unmarshal: %w", err)
	}
	return &backend, nil
}

func (c *V3Client) ListVPCBackends(vpcID, templateID int) ([]VPCBackend, error) {
	path := fmt.Sprintf("/vpcs/%d/backend-templates/%d/backends", vpcID, templateID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("list backends for template %d VPC %d: %w", templateID, vpcID, err)
	}
	var listData V3ListData
	if err := json.Unmarshal(resp.Data, &listData); err != nil {
		return nil, fmt.Errorf("list backends unmarshal outer: %w", err)
	}
	var backends []VPCBackend
	if err := json.Unmarshal(listData.Data, &backends); err != nil {
		return nil, fmt.Errorf("list backends unmarshal inner: %w", err)
	}
	return backends, nil
}

func (c *V3Client) GetVPCBackend(vpcID, templateID, backendID int) (*VPCBackend, error) {
	backends, err := c.ListVPCBackends(vpcID, templateID)
	if err != nil {
		return nil, err
	}
	for _, b := range backends {
		if b.BackendHostID == backendID {
			return &b, nil
		}
	}
	return nil, &V3NotFoundError{
		StatusCode: 404,
		Body:       fmt.Sprintf("backend %d not found in template %d VPC %d", backendID, templateID, vpcID),
	}
}

func (c *V3Client) UpdateVPCBackend(vpcID, templateID, backendID int, req *UpdateVPCBackendRequest) (*VPCBackend, error) {
	path := fmt.Sprintf("/vpcs/%d/backend-templates/%d/backends/%d", vpcID, templateID, backendID)
	resp, err := c.patch(path, req)
	if err != nil {
		return nil, fmt.Errorf("update backend %d for template %d VPC %d: %w", backendID, templateID, vpcID, err)
	}
	var backend VPCBackend
	if err := json.Unmarshal(resp.Data, &backend); err != nil {
		return nil, fmt.Errorf("update backend unmarshal: %w", err)
	}
	return &backend, nil
}

func (c *V3Client) DeleteVPCBackend(vpcID, templateID, backendID int) error {
	path := fmt.Sprintf("/vpcs/%d/backend-templates/%d/backends/%d", vpcID, templateID, backendID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete backend %d for template %d VPC %d: %w", backendID, templateID, vpcID, err)
	}
	return nil
}
