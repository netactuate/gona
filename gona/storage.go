package gona

import (
	"encoding/json"
	"fmt"
)

type StorageHardwareClass struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type StorageS3Credentials struct {
	Endpoints []string `json:"endpoints"`
	AccessKey string   `json:"accessKey"`
	SecretKey string   `json:"secretKey"`
	UserKey   string   `json:"userKey"`
}

type StorageBlockCredentials struct {
	Endpoints []string `json:"endpoints"`
	UserKey   string   `json:"userKey,omitempty"`
	SecretKey string   `json:"secretKey,omitempty"`
	Pool      string   `json:"pool"`
	Namespace string   `json:"namespace"`
	ClusterID string   `json:"clusterId"`
	ImageName string   `json:"imageName,omitempty"`
}

type storageCreateResponse struct {
	BucketID         int                         `json:"bucketId,omitempty"`
	ObjectStoreID    int                         `json:"objectStoreId,omitempty"`
	BlockNamespaceID int                         `json:"blockNamespaceId,omitempty"`
	BlockVolumeID    int                         `json:"blockVolumeId,omitempty"`
	Details          storageCreateResponseDetail `json:"details"`
}

type storageCreateResponseDetail struct {
	Ready  bool `json:"ready"`
	Queued bool `json:"queued"`
}

type StorageBucket struct {
	Credentials StorageS3Credentials  `json:"credentials"`
	Metadata    StorageBucketMetadata `json:"metadata"`
}

type StorageBucketMetadata struct {
	BucketID      int                  `json:"bucketId"`
	Label         string               `json:"label"`
	Ready         bool                 `json:"ready"`
	Private       bool                 `json:"private"`
	AssignedOn    string               `json:"assignedOn"`
	Location      V3Location           `json:"location"`
	Capacity      V3Capacity           `json:"capacity"`
	HardwareClass StorageHardwareClass `json:"hardwareClass"`
}

type CreateStorageBucketRequest struct {
	LocationID        int    `json:"locationId"`
	Label             string `json:"label"`
	Capacity          int    `json:"capacity,omitempty"`
	Private           *bool  `json:"private,omitempty"`
	EnableAutoScaling *bool  `json:"enableAutoScaling,omitempty"`
}

type UpdateStorageBucketRequest struct {
	Label             string `json:"label,omitempty"`
	Capacity          int    `json:"capacity,omitempty"`
	EnableAutoScaling *bool  `json:"enableAutoScaling,omitempty"`
	Private           *bool  `json:"private,omitempty"`
}

func (c *V3Client) CreateStorageBucket(req *CreateStorageBucketRequest) (int, error) {
	resp, err := c.post("/storage/buckets", req)
	if err != nil {
		return 0, fmt.Errorf("create storage bucket: %w", err)
	}
	var created storageCreateResponse
	if err := json.Unmarshal(resp.Data, &created); err != nil {
		return 0, fmt.Errorf("create storage bucket unmarshal: %w", err)
	}
	return created.BucketID, nil
}

func (c *V3Client) GetStorageBucket(bucketID int) (*StorageBucket, error) {
	path := fmt.Sprintf("/storage/buckets/%d", bucketID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get storage bucket %d: %w", bucketID, err)
	}
	var bucket StorageBucket
	if err := json.Unmarshal(resp.Data, &bucket); err != nil {
		return nil, fmt.Errorf("get storage bucket %d unmarshal: %w", bucketID, err)
	}
	return &bucket, nil
}

func (c *V3Client) UpdateStorageBucket(bucketID int, req *UpdateStorageBucketRequest) error {
	path := fmt.Sprintf("/storage/buckets/%d", bucketID)
	_, err := c.patch(path, req)
	if err != nil {
		return fmt.Errorf("update storage bucket %d: %w", bucketID, err)
	}
	return nil
}

func (c *V3Client) DeleteStorageBucket(bucketID int) error {
	path := fmt.Sprintf("/storage/buckets/%d", bucketID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete storage bucket %d: %w", bucketID, err)
	}
	return nil
}

func (c *V3Client) WaitForStorageBucketReady(bucketID int) error {
	return c.waitForCondition(func() (bool, error) {
		bucket, err := c.GetStorageBucket(bucketID)
		if err != nil {
			return false, err
		}
		c.debugLog("storage bucket %d ready: %v", bucketID, bucket.Metadata.Ready)
		return bucket.Metadata.Ready, nil
	}, StorageWaitConfig)
}

type StorageObjectStore struct {
	Credentials StorageS3Credentials       `json:"credentials"`
	Metadata    StorageObjectStoreMetadata `json:"metadata"`
}

type StorageObjectStoreMetadata struct {
	ObjectStoreID int                  `json:"objectStoreId"`
	Label         string               `json:"label"`
	Ready         bool                 `json:"ready"`
	AssignedOn    string               `json:"assignedOn"`
	Location      V3Location           `json:"location"`
	Capacity      V3Capacity           `json:"capacity"`
	HardwareClass StorageHardwareClass `json:"hardwareClass"`
}

type CreateStorageObjectStoreRequest struct {
	LocationID        int    `json:"locationId"`
	Label             string `json:"label"`
	Capacity          int    `json:"capacity,omitempty"`
	EnableAutoScaling *bool  `json:"enableAutoScaling,omitempty"`
}

type UpdateStorageObjectStoreRequest struct {
	Label             string `json:"label,omitempty"`
	Capacity          int    `json:"capacity,omitempty"`
	EnableAutoScaling *bool  `json:"enableAutoScaling,omitempty"`
}

func (c *V3Client) CreateStorageObjectStore(req *CreateStorageObjectStoreRequest) (int, error) {
	resp, err := c.post("/storage/object-stores", req)
	if err != nil {
		return 0, fmt.Errorf("create storage object store: %w", err)
	}
	var created storageCreateResponse
	if err := json.Unmarshal(resp.Data, &created); err != nil {
		return 0, fmt.Errorf("create storage object store unmarshal: %w", err)
	}
	return created.ObjectStoreID, nil
}

func (c *V3Client) GetStorageObjectStore(objectStoreID int) (*StorageObjectStore, error) {
	path := fmt.Sprintf("/storage/object-stores/%d", objectStoreID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get storage object store %d: %w", objectStoreID, err)
	}
	var store StorageObjectStore
	if err := json.Unmarshal(resp.Data, &store); err != nil {
		return nil, fmt.Errorf("get storage object store %d unmarshal: %w", objectStoreID, err)
	}
	return &store, nil
}

func (c *V3Client) UpdateStorageObjectStore(objectStoreID int, req *UpdateStorageObjectStoreRequest) error {
	path := fmt.Sprintf("/storage/object-stores/%d", objectStoreID)
	_, err := c.patch(path, req)
	if err != nil {
		return fmt.Errorf("update storage object store %d: %w", objectStoreID, err)
	}
	return nil
}

func (c *V3Client) DeleteStorageObjectStore(objectStoreID int) error {
	path := fmt.Sprintf("/storage/object-stores/%d", objectStoreID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete storage object store %d: %w", objectStoreID, err)
	}
	return nil
}

func (c *V3Client) WaitForStorageObjectStoreReady(objectStoreID int) error {
	return c.waitForCondition(func() (bool, error) {
		store, err := c.GetStorageObjectStore(objectStoreID)
		if err != nil {
			return false, err
		}
		c.debugLog("storage object store %d ready: %v", objectStoreID, store.Metadata.Ready)
		return store.Metadata.Ready, nil
	}, StorageWaitConfig)
}

type StorageBlockNamespace struct {
	Credentials StorageBlockCredentials       `json:"credentials"`
	Metadata    StorageBlockNamespaceMetadata `json:"metadata"`
}

type StorageBlockNamespaceMetadata struct {
	BlockNamespaceID int                  `json:"blockNamespaceId"`
	Label            string               `json:"label"`
	Ready            bool                 `json:"ready"`
	AssignedOn       string               `json:"assignedOn"`
	Location         V3Location           `json:"location"`
	Capacity         V3Capacity           `json:"capacity"`
	HardwareClass    StorageHardwareClass `json:"hardwareClass"`
}

type CreateStorageBlockNamespaceRequest struct {
	LocationID        int    `json:"locationId"`
	Label             string `json:"label"`
	Capacity          int    `json:"capacity,omitempty"`
	EnableAutoScaling *bool  `json:"enableAutoScaling,omitempty"`
}

type UpdateStorageBlockNamespaceRequest struct {
	Label             string `json:"label,omitempty"`
	Capacity          int    `json:"capacity,omitempty"`
	EnableAutoScaling *bool  `json:"enableAutoScaling,omitempty"`
}

func (c *V3Client) CreateStorageBlockNamespace(req *CreateStorageBlockNamespaceRequest) (int, error) {
	resp, err := c.post("/storage/block-namespaces", req)
	if err != nil {
		return 0, fmt.Errorf("create storage block namespace: %w", err)
	}
	var created storageCreateResponse
	if err := json.Unmarshal(resp.Data, &created); err != nil {
		return 0, fmt.Errorf("create storage block namespace unmarshal: %w", err)
	}
	return created.BlockNamespaceID, nil
}

func (c *V3Client) GetStorageBlockNamespace(blockNamespaceID int) (*StorageBlockNamespace, error) {
	path := fmt.Sprintf("/storage/block-namespaces/%d", blockNamespaceID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get storage block namespace %d: %w", blockNamespaceID, err)
	}
	var ns StorageBlockNamespace
	if err := json.Unmarshal(resp.Data, &ns); err != nil {
		return nil, fmt.Errorf("get storage block namespace %d unmarshal: %w", blockNamespaceID, err)
	}
	return &ns, nil
}

func (c *V3Client) UpdateStorageBlockNamespace(blockNamespaceID int, req *UpdateStorageBlockNamespaceRequest) error {
	path := fmt.Sprintf("/storage/block-namespaces/%d", blockNamespaceID)
	_, err := c.patch(path, req)
	if err != nil {
		return fmt.Errorf("update storage block namespace %d: %w", blockNamespaceID, err)
	}
	return nil
}

func (c *V3Client) DeleteStorageBlockNamespace(blockNamespaceID int) error {
	path := fmt.Sprintf("/storage/block-namespaces/%d", blockNamespaceID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete storage block namespace %d: %w", blockNamespaceID, err)
	}
	return nil
}

func (c *V3Client) WaitForStorageBlockNamespaceReady(blockNamespaceID int) error {
	return c.waitForCondition(func() (bool, error) {
		ns, err := c.GetStorageBlockNamespace(blockNamespaceID)
		if err != nil {
			return false, err
		}
		c.debugLog("storage block namespace %d ready: %v", blockNamespaceID, ns.Metadata.Ready)
		return ns.Metadata.Ready, nil
	}, StorageWaitConfig)
}

type StorageBlockVolume struct {
	Credentials StorageBlockCredentials    `json:"credentials"`
	Metadata    StorageBlockVolumeMetadata `json:"metadata"`
}

type StorageBlockVolumeMetadata struct {
	BlockVolumeID int                  `json:"blockVolumeId"`
	Label         string               `json:"label"`
	Ready         bool                 `json:"ready"`
	AssignedOn    string               `json:"assignedOn"`
	Location      V3Location           `json:"location"`
	Capacity      V3Capacity           `json:"capacity"`
	HardwareClass StorageHardwareClass `json:"hardwareClass"`
}

type CreateStorageBlockVolumeRequest struct {
	LocationID int    `json:"locationId"`
	Label      string `json:"label"`
	Capacity   int    `json:"capacity,omitempty"`
}

type UpdateStorageBlockVolumeRequest struct {
	Label    string `json:"label,omitempty"`
	Capacity int    `json:"capacity,omitempty"`
}

func (c *V3Client) CreateStorageBlockVolume(req *CreateStorageBlockVolumeRequest) (int, error) {
	resp, err := c.post("/storage/block-volumes", req)
	if err != nil {
		return 0, fmt.Errorf("create storage block volume: %w", err)
	}
	var created storageCreateResponse
	if err := json.Unmarshal(resp.Data, &created); err != nil {
		return 0, fmt.Errorf("create storage block volume unmarshal: %w", err)
	}
	return created.BlockVolumeID, nil
}

func (c *V3Client) GetStorageBlockVolume(blockVolumeID int) (*StorageBlockVolume, error) {
	path := fmt.Sprintf("/storage/block-volumes/%d", blockVolumeID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get storage block volume %d: %w", blockVolumeID, err)
	}
	var vol StorageBlockVolume
	if err := json.Unmarshal(resp.Data, &vol); err != nil {
		return nil, fmt.Errorf("get storage block volume %d unmarshal: %w", blockVolumeID, err)
	}
	return &vol, nil
}

func (c *V3Client) UpdateStorageBlockVolume(blockVolumeID int, req *UpdateStorageBlockVolumeRequest) error {
	path := fmt.Sprintf("/storage/block-volumes/%d", blockVolumeID)
	_, err := c.patch(path, req)
	if err != nil {
		return fmt.Errorf("update storage block volume %d: %w", blockVolumeID, err)
	}
	return nil
}

func (c *V3Client) DeleteStorageBlockVolume(blockVolumeID int) error {
	path := fmt.Sprintf("/storage/block-volumes/%d", blockVolumeID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete storage block volume %d: %w", blockVolumeID, err)
	}
	return nil
}

func (c *V3Client) WaitForStorageBlockVolumeReady(blockVolumeID int) error {
	return c.waitForCondition(func() (bool, error) {
		vol, err := c.GetStorageBlockVolume(blockVolumeID)
		if err != nil {
			return false, err
		}
		c.debugLog("storage block volume %d ready: %v", blockVolumeID, vol.Metadata.Ready)
		return vol.Metadata.Ready, nil
	}, StorageWaitConfig)
}

type StorageLocation struct {
	Location V3Location           `json:"location"`
	Hardware StorageHardwareClass `json:"hardware"`
}

func (c *V3Client) ListStorageLocations() ([]StorageLocation, error) {
	resp, err := c.get("/storage/locations")
	if err != nil {
		return nil, fmt.Errorf("list storage locations: %w", err)
	}
	var locations []StorageLocation
	if err := json.Unmarshal(resp.Data, &locations); err != nil {
		return nil, fmt.Errorf("list storage locations: %w", err)
	}
	return locations, nil
}
