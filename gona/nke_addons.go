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

// UnmarshalJSON accepts status as a string ("True") or as a JSON boolean. The nke-addons
// operator reports conditions with a boolean status, which a plain string field can't
// decode, so any read of an add-on with reported conditions failed. A boolean is stored
// as "True" or "False", the Kubernetes spelling.
func (c *NKEAddonCondition) UnmarshalJSON(data []byte) error {
	type alias NKEAddonCondition
	aux := &struct {
		Status interface{} `json:"status"`
		*alias
	}{
		alias: (*alias)(c),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	switch v := aux.Status.(type) {
	case string:
		c.Status = v
	case bool:
		if v {
			c.Status = "True"
		} else {
			c.Status = "False"
		}
	case nil:
		c.Status = ""
	default:
		return fmt.Errorf("addon condition status: unexpected type %T", v)
	}
	return nil
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
// Both known addon types use this read shape.
type NKEAddonReadConfig struct {
	NKEAutoscalerConfig

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

// NKELBServiceFirewallAllow is one source range the netactuate-cloud-lb add-on opened on the
// VPC gateway for a managed Service or Ingress.
type NKELBServiceFirewallAllow struct {
	FirewallRuleID int    `json:"firewallRuleId"`
	SourceRange    string `json:"sourceRange"`
	Protocol       string `json:"protocol"`
	Port           int    `json:"port"`
}

// NKELBServicePort is one port of a managed Service or Ingress: the public port on the floating
// IP and the node port it forwards to.
type NKELBServicePort struct {
	Port     int `json:"port"`
	NodePort int `json:"nodePort"`
}

// NKELBService is one Kubernetes Service (type LoadBalancer) or Ingress that the
// netactuate-cloud-lb add-on's in-cluster operator reported. Read-only: the operator
// creates these rows from cluster objects.
type NKELBService struct {
	LBServiceID       int                         `json:"lbServiceId"`
	ClusterID         int                         `json:"clusterId"`
	Kind              string                      `json:"kind"`
	Namespace         string                      `json:"namespace"`
	Name              string                      `json:"name"`
	Protocol          string                      `json:"protocol"`
	Ports             []NKELBServicePort          `json:"ports"`
	FloatingIPID      int                         `json:"floatingIpId"`
	L4GroupID         int                         `json:"l4GroupId"`
	L7GroupID         int                         `json:"l7GroupId"`
	SourceRanges      []string                    `json:"sourceRanges"`
	TLSMode           string                      `json:"tlsMode"`
	Hostnames         []string                    `json:"hostnames"`
	VIP               string                      `json:"vip"`
	State             string                      `json:"state"`
	FailureReason     string                      `json:"failureReason"`
	Health            json.RawMessage             `json:"health"`
	Timestamps        json.RawMessage             `json:"timestamps"`
	FirewallAllowList []NKELBServiceFirewallAllow `json:"firewallAllowList"`
}

// ListClusterLBServices lists the Services and Ingresses the netactuate-cloud-lb add-on manages
// on a cluster (GET /nke/clusters/{id}/lb-services).
func (c *V3Client) ListClusterLBServices(clusterID int) ([]NKELBService, error) {
	listData, err := c.getList(fmt.Sprintf("/nke/clusters/%d/lb-services?limit=1000", clusterID))
	if err != nil {
		return nil, fmt.Errorf("list NKE cluster %d LB services: %w", clusterID, err)
	}
	var services []NKELBService
	if err := json.Unmarshal(listData.Data, &services); err != nil {
		return nil, fmt.Errorf("list NKE cluster %d LB services data unmarshal: %w", clusterID, err)
	}
	return services, nil
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

// NKEAutoscalerConfig contains the catalog settings for netactuate-cluster-autoscaler.
type NKEAutoscalerConfig struct {
	ScaleDownEnabled              *bool    `json:"scaleDownEnabled,omitempty"`
	ScaleDownUnneededTime         *string  `json:"scaleDownUnneededTime,omitempty"`
	ScaleDownDelayAfterAdd        *string  `json:"scaleDownDelayAfterAdd,omitempty"`
	ScaleDownUtilizationThreshold *float64 `json:"scaleDownUtilizationThreshold,omitempty"`
	MaxNodeProvisionTime          *string  `json:"maxNodeProvisionTime,omitempty"`
	Expander                      *string  `json:"expander,omitempty"`
	SkipNodesWithLocalStorage     *bool    `json:"skipNodesWithLocalStorage,omitempty"`
	BalanceSimilarNodeGroups      *bool    `json:"balanceSimilarNodeGroups,omitempty"`
	MaxNodesTotal                 *int     `json:"maxNodesTotal,omitempty"`
}
