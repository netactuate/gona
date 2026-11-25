package gona

import (
	"encoding/json"
	"fmt"
)

type NKENodeTaint struct {
	Key    string `json:"key"`
	Value  string `json:"value,omitempty"`
	Effect string `json:"effect"`
}

type NKENodePool struct {
	NodePoolID                int               `json:"nodePoolId"`
	Name                      string            `json:"name"`
	Mode                      string            `json:"mode"`
	PackageID                 *int              `json:"packageId"`
	PackageName               *string           `json:"packageName"`
	Version                   *string           `json:"version"`
	AcceleratorTypeID         *int              `json:"acceleratorTypeId"`
	AcceleratorAttachmentType *string           `json:"acceleratorAttachmentType"`
	AcceleratorFirmwareID     *int              `json:"acceleratorFirmwareId"`
	AutoscalingEnabled        int               `json:"autoscalingEnabled"`
	DesiredNodes              *int              `json:"desiredNodes"`
	LastScaleOn               *string           `json:"lastScaleOn"`
	Generation                int               `json:"generation"`
	MinimumNodes              int               `json:"minimumNodes"`
	MaximumNodes              *int              `json:"maximumNodes"`
	Labels                    map[string]string `json:"labels"`
	Taints                    []NKENodeTaint    `json:"taints"`
	CreatedOn                 string            `json:"createdOn"`
	Nodes                     struct {
		Ready    int `json:"ready"`
		Total    int `json:"total"`
		Building int `json:"building"`
		Deleting int `json:"deleting"`
		Outdated int `json:"outdated"`
	} `json:"nodes"`
	BuildFailureReason *string         `json:"buildFailureReason"`
	LastBuildFailureOn *string         `json:"lastBuildFailureOn"`
	RetryBuildOn       *string         `json:"retryBuildOn"`
	BuildFailureCt     int             `json:"buildFailureCt"`
	Health             json.RawMessage `json:"health"`
}

type CreateNKENodePoolRequest struct {
	Name                      string            `json:"name"`
	Mode                      string            `json:"mode"`
	PackageID                 *int              `json:"packageId,omitempty"`
	MinimumNodes              *int              `json:"minimumNodes,omitempty"`
	MaximumNodes              *int              `json:"maximumNodes,omitempty"`
	Version                   string            `json:"version,omitempty"`
	Labels                    map[string]string `json:"labels,omitempty"`
	Taints                    []NKENodeTaint    `json:"taints,omitempty"`
	AcceleratorTypeID         *int              `json:"acceleratorTypeId,omitempty"`
	AcceleratorAttachmentType *string           `json:"acceleratorAttachmentType,omitempty"`
	AcceleratorFirmwareID     *int              `json:"acceleratorFirmwareId,omitempty"`
}

type UpdateNKENodePoolRequest struct {
	Name                      *string            `json:"name,omitempty"`
	PackageID                 *int               `json:"packageId,omitempty"`
	MinimumNodes              *int               `json:"minimumNodes,omitempty"`
	MaximumNodes              *int               `json:"maximumNodes,omitempty"`
	AutoscalingEnabled        *bool              `json:"autoscalingEnabled,omitempty"`
	Labels                    *map[string]string `json:"labels,omitempty"`
	Taints                    *[]NKENodeTaint    `json:"taints,omitempty"`
	AcceleratorTypeID         *int               `json:"acceleratorTypeId,omitempty"`
	AcceleratorAttachmentType *string            `json:"acceleratorAttachmentType,omitempty"`
	AcceleratorFirmwareID     *int               `json:"acceleratorFirmwareId,omitempty"`
}

func (c *V3Client) ListNKENodePools(clusterID int) ([]NKENodePool, error) {
	r, err := c.getList(fmt.Sprintf("/nke/clusters/%d/node-pools", clusterID))
	if err != nil {
		return nil, fmt.Errorf("list NKE node pools for cluster %d: %w", clusterID, err)
	}
	var pools []NKENodePool
	if err = json.Unmarshal(r.Data, &pools); err != nil {
		return nil, fmt.Errorf("decode NKE node pools: %w", err)
	}
	return pools, nil
}

// GetNKENodePool uses the list endpoint because the API has no individual GET route.
func (c *V3Client) GetNKENodePool(clusterID, poolID int) (*NKENodePool, error) {
	pools, err := c.ListNKENodePools(clusterID)
	if err != nil {
		return nil, fmt.Errorf("get NKE node pool %d: %w", poolID, err)
	}
	for i := range pools {
		if pools[i].NodePoolID == poolID {
			return &pools[i], nil
		}
	}
	return nil, &V3NotFoundError{StatusCode: 404, Body: fmt.Sprintf("NKE node pool %d not found", poolID)}
}
func (c *V3Client) CreateNKENodePool(clusterID int, req *CreateNKENodePoolRequest) (int, error) {
	r, err := c.post(fmt.Sprintf("/nke/clusters/%d/node-pools", clusterID), req)
	if err != nil {
		return 0, fmt.Errorf("create NKE node pool: %w", err)
	}
	var result struct {
		NodePoolID int `json:"nodePoolId"`
	}
	if err = json.Unmarshal(r.Data, &result); err != nil {
		return 0, fmt.Errorf("decode created NKE node pool: %w", err)
	}
	return result.NodePoolID, nil
}
func (c *V3Client) UpdateNKENodePool(clusterID, poolID int, req *UpdateNKENodePoolRequest) error {
	_, err := c.patch(fmt.Sprintf("/nke/clusters/%d/node-pools/%d", clusterID, poolID), req)
	if err != nil {
		return fmt.Errorf("update NKE node pool %d: %w", poolID, err)
	}
	return nil
}
func (c *V3Client) SetNKENodePoolSize(clusterID, poolID, size int) error {
	_, err := c.put(fmt.Sprintf("/nke/clusters/%d/node-pools/%d/size", clusterID, poolID), struct {
		DesiredNodes int `json:"desiredNodes"`
	}{size})
	if err != nil {
		return fmt.Errorf("size NKE node pool %d: %w", poolID, err)
	}
	return nil
}
func (c *V3Client) DeleteNKENodePool(clusterID, poolID int) error {
	_, err := c.del(fmt.Sprintf("/nke/clusters/%d/node-pools/%d", clusterID, poolID))
	if err != nil {
		return fmt.Errorf("delete NKE node pool %d: %w", poolID, err)
	}
	return nil
}
func (c *V3Client) RemoveNKEWorkerNode(clusterID, workerNodeID int) error {
	node, err := c.GetNKEWorkerNode(clusterID, workerNodeID)
	if err != nil {
		return fmt.Errorf("resolve pool for NKE worker node %d: %w", workerNodeID, err)
	}
	_, err = c.post(fmt.Sprintf("/nke/clusters/%d/worker-nodes/%d/remove", clusterID, workerNodeID), struct {
		NodePoolID int `json:"nodePoolId"`
	}{node.PoolID})
	if err != nil {
		return fmt.Errorf("remove NKE worker node %d: %w", workerNodeID, err)
	}
	return nil
}
