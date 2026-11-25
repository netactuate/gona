package gona

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestNKEDrop1RequestZerosAndOmission(t *testing.T) {
	zero := 0
	disabled := false
	labels := map[string]string{}
	taints := []NKENodeTaint{}
	for _, tc := range []struct {
		request interface{}
		want    string
	}{
		{&CreateNKENodePoolRequest{Name: "nah-test", Mode: "managed", MinimumNodes: &zero}, `{"name":"nah-test","mode":"managed","minimumNodes":0}`},
		{&UpdateNKENodePoolRequest{MinimumNodes: &zero, AutoscalingEnabled: &disabled, Labels: &labels, Taints: &taints}, `{"minimumNodes":0,"autoscalingEnabled":false,"labels":{},"taints":[]}`},
		{&NKEAutoscalerConfig{ScaleDownEnabled: &disabled}, `{"scaleDownEnabled":false}`},
		{&NKEClusterNetwork{}, `{}`},
	} {
		body, err := json.Marshal(tc.request)
		if err != nil {
			t.Fatal(err)
		}
		if string(body) != tc.want {
			t.Errorf("body %s, want %s", body, tc.want)
		}
	}
}
func TestNKEDrop1ContractResponseShapes(t *testing.T) {
	var pool NKENodePool
	if err := json.Unmarshal([]byte(`{"nodePoolId":4,"mode":"managed","minimumNodes":0,"maximumNodes":null,"packageId":5,"autoscalingEnabled":0,"desiredNodes":0,"nodes":{"ready":1,"total":2,"building":1,"deleting":0,"outdated":0},"labels":{},"taints":[],"health":{"state":"yellow"}}`), &pool); err != nil {
		t.Fatal(err)
	}
	if pool.Nodes.Total != 2 || pool.Nodes.Ready != 1 || pool.DesiredNodes == nil || *pool.DesiredNodes != 0 || pool.MaximumNodes != nil {
		t.Fatalf("pool: %#v", pool)
	}
	var addon NKEAddon
	if err := json.Unmarshal([]byte(`{"config":{"scaleDownEnabled":false,"balanceSimilarNodeGroups":false,"maxNodesTotal":null,"expander":"least-waste"}}`), &addon); err != nil {
		t.Fatal(err)
	}
	if addon.Config.ScaleDownEnabled == nil || *addon.Config.ScaleDownEnabled || addon.Config.MaxNodesTotal != nil {
		t.Fatalf("config: %#v", addon.Config)
	}
	var rows []NKEStorageAvailability
	if err := json.Unmarshal([]byte(`[{"locationId":3,"storageAvailable":false}]`), &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 1 || rows[0].LocationID != 3 || rows[0].StorageAvailable {
		t.Fatal(rows)
	}
}

func TestNKEDrop1MissingPoolIsNotFound(t *testing.T) {
	c := newSDKCoverageTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/nke/clusters/1/node-pools" {
			t.Fatalf("request %s %s", r.Method, r.URL.Path)
		}
		w.Write([]byte(`{"code":200,"data":{"meta":{"total":0,"limit":1000,"offset":0},"data":[]}}`))
	})
	_, err := c.GetNKENodePool(1, 2)
	if !IsV3NotFound(err) {
		t.Fatalf("missing pool: %v", err)
	}
}
func TestNKEDrop1WaitTimeouts(t *testing.T) {
	c := newSDKCoverageTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/nke/clusters/1" {
			w.Write([]byte(`{"code":200,"data":{"status":{"cluster":"Initializing"}}}`))
			return
		}
		if r.URL.Path != "/nke/clusters/1/worker-nodes" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		w.Write([]byte(`{"code":200,"data":{"meta":{"total":0,"limit":1000,"offset":0},"data":[]}}`))
	})
	for _, wait := range []func() error{func() error { return c.WaitForNKEClusterHealthyTimeout(1, 2*time.Millisecond) }, func() error { return c.WaitForNKEWorkerNodesTimeout(1, 1, 2*time.Millisecond) }} {
		start := time.Now()
		if err := wait(); err == nil {
			t.Fatal("expected timeout")
		}
		if time.Since(start) > time.Second {
			t.Fatal("caller timeout ignored")
		}
	}
}
