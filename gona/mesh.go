package gona

import (
	"encoding/json"
	"fmt"
)

type MagicMesh struct {
	MeshID      int     `json:"meshId"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type MeshRouter struct {
	RouterID    int     `json:"routerId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IPv4Address string  `json:"ipv4Address"`
}

type CreateMagicMeshRequest struct {
	Name        string              `json:"name"`
	Description *string             `json:"description,omitempty"`
	Routers     []MeshRouterEntry   `json:"routers,omitempty"`
}

type MeshRouterEntry struct {
	RouterID int `json:"routerId"`
}

type CreateMagicMeshResponse struct {
	MeshID int `json:"meshId"`
}

type UpdateMagicMeshRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type AddMeshRouterRequest struct {
	RouterID int `json:"routerId"`
}

func (c *V3Client) CreateMagicMesh(req *CreateMagicMeshRequest) (*CreateMagicMeshResponse, error) {
	resp, err := c.post("/cloud-routing/meshes", req)
	if err != nil {
		return nil, fmt.Errorf("failed to create magic mesh: %w", err)
	}

	var createResp CreateMagicMeshResponse
	if err := json.Unmarshal(resp.Data, &createResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal create mesh response: %w", err)
	}

	return &createResp, nil
}

func (c *V3Client) GetMagicMesh(meshID int) (*MagicMesh, error) {
	path := fmt.Sprintf("/cloud-routing/meshes/%d", meshID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get magic mesh %d: %w", meshID, err)
	}

	var mesh MagicMesh
	if err := json.Unmarshal(resp.Data, &mesh); err != nil {
		return nil, fmt.Errorf("failed to unmarshal magic mesh: %w", err)
	}

	return &mesh, nil
}

func (c *V3Client) UpdateMagicMesh(meshID int, req *UpdateMagicMeshRequest) error {
	path := fmt.Sprintf("/cloud-routing/meshes/%d", meshID)

	_, err := c.patch(path, req)
	if err != nil {
		return fmt.Errorf("failed to update magic mesh %d: %w", meshID, err)
	}

	return nil
}

func (c *V3Client) DeleteMagicMesh(meshID int) error {
	path := fmt.Sprintf("/cloud-routing/meshes/%d", meshID)

	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("failed to delete magic mesh %d: %w", meshID, err)
	}

	return nil
}

func (c *V3Client) ListMeshRouters(meshID int) ([]MeshRouter, error) {
	path := fmt.Sprintf("/cloud-routing/meshes/%d/routers", meshID)

	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list routers for mesh %d: %w", meshID, err)
	}

	var listData V3ListData
	if err := json.Unmarshal(resp.Data, &listData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mesh routers list envelope: %w", err)
	}

	var routers []MeshRouter
	if err := json.Unmarshal(listData.Data, &routers); err != nil {
		return nil, fmt.Errorf("failed to unmarshal mesh routers: %w", err)
	}

	return routers, nil
}

func (c *V3Client) AddRouterToMesh(meshID int, req *AddMeshRouterRequest) error {
	path := fmt.Sprintf("/cloud-routing/meshes/%d/routers", meshID)

	_, err := c.post(path, req)
	if err != nil {
		return fmt.Errorf("failed to add router %d to mesh %d: %w", req.RouterID, meshID, err)
	}

	return nil
}

func (c *V3Client) RemoveRouterFromMesh(meshID, routerID int) error {
	path := fmt.Sprintf("/cloud-routing/meshes/%d/routers/%d", meshID, routerID)

	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("failed to remove router %d from mesh %d: %w", routerID, meshID, err)
	}

	return nil
}
