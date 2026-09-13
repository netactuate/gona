package gona

import (
	"encoding/json"
	"fmt"
)

type NKEAddonCatalogEntry struct {
	AddonID              int    `json:"addonId"`
	AddonType            string `json:"addonType"`
	Version              string `json:"version"`
	Channel              string `json:"channel"`
	DisplayName          string `json:"displayName"`
	MinKubernetesVersion string `json:"minKubernetesVersion"`
	MaxKubernetesVersion string `json:"maxKubernetesVersion"`
	IsDefault            bool   `json:"isDefault"`
	RequiresVpc          bool   `json:"requiresVpc"`
}

type NKEAddon struct {
	ID               int                    `json:"id"`
	AddonID          int                    `json:"addonId"`
	ClusterID        int                    `json:"clusterId"`
	AddonType        string                 `json:"addonType"`
	Version          string                 `json:"version"`
	Channel          string                 `json:"channel"`
	DisplayName      string                 `json:"displayName"`
	State            string                 `json:"state"`
	UpdateAvailable  bool                   `json:"updateAvailable"`
	Catalog          NKEAddonCatalog        `json:"catalog"`
	Health           NKEAddonHealth         `json:"health"`
	WorkloadHealth   NKEAddonWorkloadHealth `json:"workloadHealth"`
	Config           NKEAddonReadConfig     `json:"config"`
	Timestamps       NKEAddonTimestamps     `json:"timestamps"`
	FailureReason    string                 `json:"failureReason"`
	InstallRetryOn   string                 `json:"installRetryOn"`
	InstallFailureCt int                    `json:"installFailureCt"`
}

type NKEAddonCatalog struct {
	DefaultVersion string `json:"defaultVersion"`
}

type NKEAddonHealth struct {
	LastHeartbeatOn      string              `json:"lastHeartbeatOn"`
	ObservedVersion      string              `json:"observedVersion"`
	Conditions           []NKEAddonCondition `json:"conditions"`
	ConditionsReportedOn string              `json:"conditionsReportedOn"`
	Summary              string              `json:"summary"`
}

type NKEAddonWorkloadHealth struct {
	State      string              `json:"state"`
	Conditions []NKEAddonCondition `json:"conditions"`
}

type NKEAddonCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type NKEAddonTimestamps struct {
	RequestedOn string `json:"requestedOn"`
	InstalledOn string `json:"installedOn"`
	UpdatedOn   string `json:"updatedOn"`
	DeletedOn   string `json:"deletedOn"`
}

// NKEAddonReadConfig is the config an addon reports back, which is NOT the shape it was
// written with. Both known addon types return a LIST keyed by their own concern:
//
//	netactuate-dns  write {zone, mode}                  read {zones: [{dnsZoneId, zone, mode}]}
//	storage         write {poolLabel, capacity, ...}    read {integrations: [{...}]}
//
// Both verified live on 2026-09-10 against clusters built for the purpose.
type NKEAddonReadConfig struct {
	Zones        []NKEDNSAddonZone            `json:"zones"`
	Integrations []NKEStorageAddonIntegration `json:"integrations"`
}

// NKEStorageAddonIntegration is one StorageClass binding as the storage addon reports it.
// Note the field names differ from the write side: makeDefault goes out, isDefaultClass
// comes back, and the read adds ids the write never sees.
type NKEStorageAddonIntegration struct {
	StorageIntegrationID    int    `json:"storageIntegrationId"`
	BlockNamespaceID        int    `json:"blockNamespaceId"`
	StorageClassName        string `json:"storageClassName"`
	VolumeSnapshotClassName string `json:"volumeSnapshotClassName"`
	IsDefaultClass          bool   `json:"isDefaultClass"`
	ReclaimPolicy           string `json:"reclaimPolicy"`
}

type NKEDNSAddonZone struct {
	DNSZoneID int    `json:"dnsZoneId"`
	Zone      string `json:"zone"`
	Mode      string `json:"mode"`
}

type CreateNKEAddonRequest struct {
	AddonType string      `json:"addonType"`
	Version   string      `json:"version,omitempty"`
	Channel   string      `json:"channel,omitempty"`
	Config    interface{} `json:"config,omitempty"`
}

type UpdateNKEAddonRequest struct {
	Version string      `json:"version,omitempty"`
	Channel string      `json:"channel,omitempty"`
	Config  interface{} `json:"config,omitempty"`
}

type NKEDNSAddonWriteConfig struct {
	Zone string `json:"zone"`
	Mode string `json:"mode,omitempty"`
}

type NKEStorageAddonWriteConfig struct {
	BlockNamespaceID  int    `json:"blockNamespaceId,omitempty"`
	PoolLabel         string `json:"poolLabel,omitempty"`
	Capacity          int    `json:"capacity,omitempty"`
	EnableAutoScaling *bool  `json:"enableAutoScaling,omitempty"`
	StorageClassName  string `json:"storageClassName,omitempty"`
	MakeDefault       *bool  `json:"makeDefault,omitempty"`
	ReclaimPolicy     string `json:"reclaimPolicy,omitempty"`
}

type NKEClusterDNSZone struct {
	DNSZoneID     int             `json:"dnsZoneId"`
	ClusterID     int             `json:"clusterId"`
	Zone          string          `json:"zone"`
	Mode          string          `json:"mode"`
	FailureReason string          `json:"failureReason"`
	State         string          `json:"state"`
	Health        json.RawMessage `json:"health"`
	Timestamps    json.RawMessage `json:"timestamps"`
}

func (c *V3Client) ListAddonCatalog() ([]NKEAddonCatalogEntry, error) {
	resp, err := c.get("/nke/addons")
	if err != nil {
		return nil, fmt.Errorf("list NKE addon catalog: %w", err)
	}
	var entries []NKEAddonCatalogEntry
	if err := unmarshalV3MaybeList(resp.Data, &entries); err != nil {
		return nil, fmt.Errorf("list NKE addon catalog data unmarshal: %w", err)
	}
	return entries, nil
}

func (c *V3Client) ListClusterAddons(clusterID int) ([]NKEAddon, error) {
	path := fmt.Sprintf("/nke/clusters/%d/addons", clusterID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("list NKE cluster %d addons: %w", clusterID, err)
	}
	var addons []NKEAddon
	if err := unmarshalV3MaybeList(resp.Data, &addons); err != nil {
		return nil, fmt.Errorf("list NKE cluster %d addons data unmarshal: %w", clusterID, err)
	}
	return addons, nil
}

func (c *V3Client) GetClusterAddon(clusterID int, addonType string) (*NKEAddon, error) {
	path := fmt.Sprintf("/nke/clusters/%d/addons/%s", clusterID, addonType)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get NKE cluster %d addon %s: %w", clusterID, addonType, err)
	}
	var addon NKEAddon
	if err := json.Unmarshal(resp.Data, &addon); err != nil {
		return nil, fmt.Errorf("get NKE cluster %d addon %s unmarshal: %w", clusterID, addonType, err)
	}
	return &addon, nil
}

func (c *V3Client) CreateClusterAddon(clusterID int, req *CreateNKEAddonRequest) (*NKEAddon, error) {
	path := fmt.Sprintf("/nke/clusters/%d/addons", clusterID)
	resp, err := c.post(path, req)
	if err != nil {
		return nil, fmt.Errorf("create NKE cluster %d addon: %w", clusterID, err)
	}
	var addon NKEAddon
	if err := json.Unmarshal(resp.Data, &addon); err != nil {
		return nil, fmt.Errorf("create NKE cluster %d addon unmarshal: %w", clusterID, err)
	}
	return &addon, nil
}

func (c *V3Client) UpdateClusterAddon(clusterID int, addonType string, req *UpdateNKEAddonRequest) (*NKEAddon, error) {
	path := fmt.Sprintf("/nke/clusters/%d/addons/%s", clusterID, addonType)
	resp, err := c.patch(path, req)
	if err != nil {
		return nil, fmt.Errorf("update NKE cluster %d addon %s: %w", clusterID, addonType, err)
	}
	var addon NKEAddon
	if err := json.Unmarshal(resp.Data, &addon); err != nil {
		return nil, fmt.Errorf("update NKE cluster %d addon %s unmarshal: %w", clusterID, addonType, err)
	}
	return &addon, nil
}

func (c *V3Client) DeleteClusterAddon(clusterID int, addonType string) error {
	path := fmt.Sprintf("/nke/clusters/%d/addons/%s", clusterID, addonType)
	if _, err := c.del(path); err != nil {
		return fmt.Errorf("delete NKE cluster %d addon %s: %w", clusterID, addonType, err)
	}
	return nil
}

func (c *V3Client) ListClusterDNSZones(clusterID int) ([]NKEClusterDNSZone, error) {
	path := fmt.Sprintf("/nke/clusters/%d/dns-zones", clusterID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("list NKE cluster %d DNS zones: %w", clusterID, err)
	}
	var zones []NKEClusterDNSZone
	if err := unmarshalV3MaybeList(resp.Data, &zones); err != nil {
		return nil, fmt.Errorf("list NKE cluster %d DNS zones data unmarshal: %w", clusterID, err)
	}
	return zones, nil
}

func unmarshalV3MaybeList(data json.RawMessage, v interface{}) error {
	if err := json.Unmarshal(data, v); err == nil {
		return nil
	}

	var listData V3ListData
	if err := json.Unmarshal(data, &listData); err != nil {
		return err
	}
	return json.Unmarshal(listData.Data, v)
}
