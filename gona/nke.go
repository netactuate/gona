package gona

import (
	"encoding/json"
	"fmt"
)

type NKETag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon,omitempty"`
	Color       string `json:"color,omitempty"`
}

type NKEClusterTagInput struct {
	TagID int `json:"tagId"`
}

type NKEAddons struct {
	KubernetesDashboard bool `json:"kubernetesDashboard,omitempty"`
}

type NKEBilling struct {
	PackageID  int  `json:"packageId"`
	LocationID int  `json:"locationId"`
	ContractID *int `json:"contractId,omitempty"`
}

type NKECluster struct {
	ClusterID           int    `json:"clusterId"`
	Name                string `json:"name"`
	ContractID          int    `json:"contractId"`
	Replicas            int    `json:"replicas"`
	DoAutoscaling       int    `json:"doAutoscaling"` // 0 or 1 from API
	HasHighAvailability bool   `json:"hasHighAvailabity"` // note: API typo preserved
	Status              struct {
		Cluster string `json:"cluster"`
		Scaling string `json:"scaling"`
	} `json:"status"`
	Version struct {
		Active    string  `json:"active"`
		Requested *string `json:"requested"`
	} `json:"version"`
	Location V3Location `json:"location"`
	Package  V3Package  `json:"package"`
	Nodes    struct {
		Total       int  `json:"total"`
		Ready       int  `json:"ready"`
		Building    int  `json:"building"`
		Minimum     int  `json:"minimum"`
		Maximum     *int `json:"maximum"`
		Outdated    int  `json:"outdated"`
		Deleting    int  `json:"deleting"`
		Unevictable int  `json:"unevictable"`
	} `json:"nodes"`
	Networks struct {
		Pod     string `json:"pod"`
		Service string `json:"service"`
	} `json:"networks"`
	URLs struct {
		API                 string `json:"api"`
		Prometheus          string `json:"prometheus"`
		KubernetesDashboard string `json:"kubernetesDashboard"`
	} `json:"urls"`
	Dates struct {
		Created string  `json:"created"`
		Ready   *string `json:"ready"`
	} `json:"dates"`
	KubernetesDashboard struct {
		Requested bool   `json:"requested"`
		URL       string `json:"url,omitempty"`
	} `json:"kubernetesDashboard"`
	Tags []NKETag `json:"tags"`
}

type CreateNKEClusterRequest struct {
	Name             string               `json:"name"`
	Version          string               `json:"version"`
	Replicas         int                  `json:"replicas"`
	MinimumNodes     int                  `json:"minimumNodes"`
	MaximumNodes     int                  `json:"maximumNodes"`
	DoAutoscaling    bool                 `json:"doAutoscaling"`
	DoDualStack      bool                 `json:"doDualStack"`
	HighAvailability bool                 `json:"enable.high.availability"`
	Billing          NKEBilling           `json:"billing"`
	AddonsToInstall  *NKEAddons           `json:"addonsToInstall,omitempty"`
	Tags             []NKEClusterTagInput `json:"tags,omitempty"`
}

type NKEUpdateNodes struct {
	Minimum int  `json:"minimum"`
	Maximum *int `json:"maximum"`
}

type NKEUpdateBilling struct {
	PackageID int `json:"packageId,omitempty"`
}

type UpdateNKEClusterRequest struct {
	Name          string               `json:"name,omitempty"`
	Version       string               `json:"version,omitempty"`
	DoAutoscaling *bool                `json:"doAutoscaling,omitempty"`
	Billing       *NKEUpdateBilling    `json:"billing,omitempty"`
	Nodes         *NKEUpdateNodes      `json:"nodes,omitempty"`
	Tags          []NKEClusterTagInput `json:"tags,omitempty"`
}

type NKEWorkerNode struct {
	WorkerNodeID int    `json:"workerNodeId"`
	ClusterID    int    `json:"clusterId"`
	Name         string `json:"name"`
	MBPkgID      int    `json:"mbpkgid,omitempty"`
	LocationID   int    `json:"locationId,omitempty"`
	Status       struct {
		Ready bool `json:"ready"`
	} `json:"status"`
	Package  V3Package  `json:"package"`
	Location V3Location `json:"location"`
}

type NKELogEntry struct {
	RecordedOn string `json:"recordedOn"`
	Message    string `json:"message"`
}

type nkeCreateResponse struct {
	ClusterID int `json:"clusterId"`
}

func (c *V3Client) ListNKEVersions() ([]string, error) {
	resp, err := c.get("/nke/versions")
	if err != nil {
		return nil, fmt.Errorf("list NKE versions: %w", err)
	}
	var versions []string
	if err := json.Unmarshal(resp.Data, &versions); err != nil {
		return nil, fmt.Errorf("list NKE versions unmarshal: %w", err)
	}
	return versions, nil
}

func (c *V3Client) ListNKEClusters() ([]NKECluster, error) {
	resp, err := c.get("/nke/clusters")
	if err != nil {
		return nil, fmt.Errorf("list NKE clusters: %w", err)
	}
	var listData V3ListData
	if err := json.Unmarshal(resp.Data, &listData); err != nil {
		return nil, fmt.Errorf("list NKE clusters unmarshal: %w", err)
	}
	var clusters []NKECluster
	if err := json.Unmarshal(listData.Data, &clusters); err != nil {
		return nil, fmt.Errorf("list NKE clusters data unmarshal: %w", err)
	}
	return clusters, nil
}

func (c *V3Client) GetNKECluster(clusterID int) (*NKECluster, error) {
	path := fmt.Sprintf("/nke/clusters/%d", clusterID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get NKE cluster %d: %w", clusterID, err)
	}
	var cluster NKECluster
	if err := json.Unmarshal(resp.Data, &cluster); err != nil {
		return nil, fmt.Errorf("get NKE cluster %d unmarshal: %w", clusterID, err)
	}
	return &cluster, nil
}

func (c *V3Client) CreateNKECluster(req *CreateNKEClusterRequest) (int, error) {
	resp, err := c.post("/nke/clusters", req)
	if err != nil {
		return 0, fmt.Errorf("create NKE cluster: %w", err)
	}
	var created nkeCreateResponse
	if err := json.Unmarshal(resp.Data, &created); err != nil {
		return 0, fmt.Errorf("create NKE cluster unmarshal: %w", err)
	}
	return created.ClusterID, nil
}

func (c *V3Client) UpdateNKECluster(clusterID int, req *UpdateNKEClusterRequest) error {
	path := fmt.Sprintf("/nke/clusters/%d", clusterID)
	_, err := c.patch(path, req)
	if err != nil {
		return fmt.Errorf("update NKE cluster %d: %w", clusterID, err)
	}
	return nil
}

func (c *V3Client) DeleteNKECluster(clusterID int) error {
	path := fmt.Sprintf("/nke/clusters/%d", clusterID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete NKE cluster %d: %w", clusterID, err)
	}
	return nil
}

// GenerateNKEKubeconfig creates a new access token/kubeconfig for the cluster.
// expirationSeconds controls token lifetime (default 3600, max 3153600000).
func (c *V3Client) GenerateNKEKubeconfig(clusterID int, expirationSeconds int) (string, error) {
	path := fmt.Sprintf("/nke/clusters/%d/kubeconfig", clusterID)
	body := map[string]interface{}{
		"expirationSeconds": expirationSeconds,
	}
	resp, err := c.post(path, body)
	if err != nil {
		return "", fmt.Errorf("generate NKE kubeconfig for cluster %d: %w", clusterID, err)
	}
	var kubeconfig string
	if err := json.Unmarshal(resp.Data, &kubeconfig); err != nil {
		// Fall back to raw data if it is not a JSON string
		return string(resp.Data), nil
	}
	return kubeconfig, nil
}

func (c *V3Client) ListNKEClusterLogs(clusterID int) ([]NKELogEntry, error) {
	path := fmt.Sprintf("/nke/clusters/%d/logs", clusterID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get NKE cluster %d logs: %w", clusterID, err)
	}
	var listData V3ListData
	if err := json.Unmarshal(resp.Data, &listData); err != nil {
		return nil, fmt.Errorf("NKE cluster logs unmarshal: %w", err)
	}
	var entries []NKELogEntry
	if err := json.Unmarshal(listData.Data, &entries); err != nil {
		return nil, fmt.Errorf("NKE cluster logs data unmarshal: %w", err)
	}
	return entries, nil
}

func (c *V3Client) ListNKEWorkerNodes(clusterID int) ([]NKEWorkerNode, error) {
	path := fmt.Sprintf("/nke/clusters/%d/worker-nodes", clusterID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("list NKE worker nodes for cluster %d: %w", clusterID, err)
	}
	var listData V3ListData
	if err := json.Unmarshal(resp.Data, &listData); err != nil {
		return nil, fmt.Errorf("list NKE worker nodes unmarshal: %w", err)
	}
	var nodes []NKEWorkerNode
	if err := json.Unmarshal(listData.Data, &nodes); err != nil {
		return nil, fmt.Errorf("list NKE worker nodes data unmarshal: %w", err)
	}
	return nodes, nil
}

func (c *V3Client) GetNKEWorkerNode(clusterID, workerNodeID int) (*NKEWorkerNode, error) {
	path := fmt.Sprintf("/nke/clusters/%d/worker-nodes/%d", clusterID, workerNodeID)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get NKE worker node %d for cluster %d: %w", workerNodeID, clusterID, err)
	}
	var node NKEWorkerNode
	if err := json.Unmarshal(resp.Data, &node); err != nil {
		return nil, fmt.Errorf("get NKE worker node %d unmarshal: %w", workerNodeID, err)
	}
	return &node, nil
}

func (c *V3Client) DeleteNKEWorkerNode(clusterID, workerNodeID int) error {
	path := fmt.Sprintf("/nke/clusters/%d/worker-nodes/%d", clusterID, workerNodeID)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete NKE worker node %d for cluster %d: %w", workerNodeID, clusterID, err)
	}
	return nil
}

func (c *V3Client) WaitForNKEClusterHealthy(clusterID int) error {
	return c.waitForCondition(func() (bool, error) {
		cluster, err := c.GetNKECluster(clusterID)
		if err != nil {
			return false, err
		}
		c.debugLog("NKE cluster %d status: cluster=%s scaling=%s", clusterID, cluster.Status.Cluster, cluster.Status.Scaling)
		switch cluster.Status.Cluster {
		case "Failed", "Error":
			return false, fmt.Errorf("NKE cluster %d entered failed state: %s", clusterID, cluster.Status.Cluster)
		case "Healthy":
			return true, nil
		}
		return false, nil
	}, NKEWaitConfig)
}
