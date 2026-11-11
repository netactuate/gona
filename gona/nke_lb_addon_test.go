package gona

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"
)

// The nke-addons operator reports health conditions with a boolean status (vapi3 field
// addonConditions.status is type boolean). A string-only Status field made every add-on read
// fail once the operator had sent its first heartbeat.
const cloudLBAddonBody = `{"code":200,"data":{
  "id": 81, "addonId": 12, "clusterId": 512, "addonType": "netactuate-cloud-lb",
  "version": "v0.3.0", "channel": "stable", "displayName": "Cloud Load Balancer",
  "state": "ready", "updateAvailable": false,
  "health": {"lastHeartbeatOn": "2026-11-10 12:00:00", "conditions": [
    {"type": "OperatorReady", "status": true, "reason": "Running", "message": ""},
    {"type": "GatewayReady", "status": false, "reason": "Pending"}
  ], "conditionsReportedOn": "2026-11-10 12:00:00"},
  "workloadHealth": {"conditions": [{"type": "Available", "status": "True"}]},
  "config": null,
  "timestamps": {"requestedOn": "2026-11-10 11:58:00", "installedOn": "2026-11-10 11:59:30"},
  "failureReason": null
}}`

func TestGetClusterAddonDecodesBooleanConditionStatus(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/nke/clusters/512/addons/netactuate-cloud-lb" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cloudLBAddonBody))
	})

	addon, err := c.GetClusterAddon(512, "netactuate-cloud-lb")
	if err != nil {
		t.Fatalf("GetClusterAddon: %v", err)
	}
	if addon.State != "ready" || addon.Version != "v0.3.0" || addon.Timestamps.InstalledOn == "" {
		t.Fatalf("unexpected addon: %+v", addon)
	}
	conds := addon.Health.Conditions
	if len(conds) != 2 || conds[0].Status != "True" || conds[1].Status != "False" {
		t.Fatalf("health conditions: got %+v, want statuses True, False", conds)
	}
	if conds[0].Type != "OperatorReady" || conds[0].Reason != "Running" || conds[1].Reason != "Pending" {
		t.Fatalf("other condition fields lost: %+v", conds)
	}
	if w := addon.WorkloadHealth.Conditions; len(w) != 1 || w[0].Status != "True" {
		t.Fatalf("string status not kept: %+v", w)
	}
}

func TestNKEAddonConditionStatusShapes(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{`{"type":"A","status":true}`, "True"},
		{`{"type":"A","status":false}`, "False"},
		{`{"type":"A","status":"Unknown"}`, "Unknown"},
		{`{"type":"A","status":null}`, ""},
		{`{"type":"A"}`, ""},
	}
	for _, tc := range cases {
		var cond NKEAddonCondition
		if err := json.Unmarshal([]byte(tc.in), &cond); err != nil {
			t.Fatalf("%s: %v", tc.in, err)
		}
		if cond.Status != tc.want || cond.Type != "A" {
			t.Fatalf("%s: got %+v, want status %q", tc.in, cond, tc.want)
		}
	}

	var cond NKEAddonCondition
	if err := json.Unmarshal([]byte(`{"type":"A","status":1}`), &cond); err == nil {
		t.Fatal("a numeric status should be an error")
	}

	// Marshalling is unchanged: Status is still a plain string.
	out, err := json.Marshal(NKEAddonCondition{Type: "A", Status: "True"})
	if err != nil || !strings.Contains(string(out), `"status":"True"`) {
		t.Fatalf("marshal: %s, %v", out, err)
	}
}

// Body shape from vapi3 GET /nke/clusters/{id}/lb-services (getter.parent.js): paginated,
// json columns parsed, and floatingIpId/vip/failureReason null until assigned.
func TestListClusterLBServices(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/nke/clusters/512/lb-services" || r.URL.Query().Get("limit") != "1000" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"data":[
		  {"lbServiceId": 7, "clusterId": 512, "kind": "service", "namespace": "default", "name": "web-lb",
		   "protocol": "tcp", "ports": [{"port": 80, "nodePort": 31080}], "floatingIpId": 901,
		   "l4GroupId": 33, "l7GroupId": null, "sourceRanges": ["203.0.113.0/24"], "tlsMode": null,
		   "hostnames": [], "vip": "198.51.100.20", "state": "ready", "failureReason": null,
		   "health": {"lastHeartbeatOn": "2026-11-10 12:00:00"},
		   "timestamps": {"requestedOn": "2026-11-10 12:00:01", "installedOn": null, "readyOn": "2026-11-10 12:00:09", "createdOn": "2026-11-10 12:00:01"},
		   "firewallAllowList": [{"firewallRuleId": 4401, "sourceRange": "203.0.113.0/24", "protocol": "TCP", "port": 80}]},
		  {"lbServiceId": 8, "clusterId": 512, "kind": "ingress", "namespace": "default", "name": "web",
		   "protocol": "tcp", "ports": [{"port": 80}], "floatingIpId": null, "l4GroupId": null, "l7GroupId": null,
		   "sourceRanges": null, "tlsMode": "none", "hostnames": ["app.example.com"], "vip": null,
		   "state": "pending", "failureReason": null, "health": {"lastHeartbeatOn": null},
		   "timestamps": {"requestedOn": null, "installedOn": null, "readyOn": null, "createdOn": "2026-11-10 12:00:02"},
		   "firewallAllowList": []}
		],"meta":{"total":2}}}`))
	})

	svcs, err := c.ListClusterLBServices(512)
	if err != nil {
		t.Fatalf("ListClusterLBServices: %v", err)
	}
	if len(svcs) != 2 {
		t.Fatalf("got %d services, want 2", len(svcs))
	}
	s := svcs[0]
	if s.Kind != "service" || s.Name != "web-lb" || s.VIP != "198.51.100.20" || s.FloatingIPID != 901 || s.State != "ready" {
		t.Fatalf("service row: %+v", s)
	}
	if len(s.Ports) != 1 || s.Ports[0].Port != 80 || s.Ports[0].NodePort != 31080 {
		t.Fatalf("ports: got %+v", s.Ports)
	}
	if len(s.SourceRanges) != 1 || s.SourceRanges[0] != "203.0.113.0/24" {
		t.Fatalf("source ranges: got %v", s.SourceRanges)
	}
	if len(s.FirewallAllowList) != 1 || s.FirewallAllowList[0].SourceRange != "203.0.113.0/24" {
		t.Fatalf("firewall allow list: %+v", s.FirewallAllowList)
	}
	in := svcs[1]
	if len(in.Hostnames) != 1 || in.Hostnames[0] != "app.example.com" || in.SourceRanges != nil {
		t.Fatalf("ingress hostnames/source ranges: %v %v", in.Hostnames, in.SourceRanges)
	}
	if in.Kind != "ingress" || in.VIP != "" || in.FloatingIPID != 0 || in.TLSMode != "none" || in.State != "pending" {
		t.Fatalf("ingress row (nulls must decode to zero values): %+v", in)
	}
}

func TestListClusterLBServicesEmpty(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"data":[],"meta":{"total":0}}}`))
	})
	svcs, err := c.ListClusterLBServices(512)
	if err != nil || len(svcs) != 0 {
		t.Fatalf("got %v, %v; want empty list, no error", svcs, err)
	}
}

// The list follows pages: a cluster with more Services than one page must come back whole.
func TestListClusterLBServicesFollowsPages(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		switch r.URL.Query().Get("offset") {
		case "", "0":
			w.Write([]byte(`{"code":200,"data":{"data":[{"lbServiceId":1,"name":"a"}],"meta":{"limit":1,"offset":0,"total":2}}}`))
		case "1":
			w.Write([]byte(`{"code":200,"data":{"data":[{"lbServiceId":2,"name":"b"}],"meta":{"limit":1,"offset":1,"total":2}}}`))
		default:
			t.Fatalf("unexpected offset %q", r.URL.Query().Get("offset"))
		}
	})
	svcs, err := c.ListClusterLBServices(512)
	if err != nil || len(svcs) != 2 || svcs[1].Name != "b" {
		t.Fatalf("got %+v, %v; want both pages", svcs, err)
	}
}

// WaitForVPCReadyTimeout must follow the timeout it is given, not VPCWaitConfig's fixed 10 minutes.
func TestWaitForVPCReadyTimeout(t *testing.T) {
	old := VPCWaitConfig
	VPCWaitConfig = V3WaitConfig{Interval: 5 * time.Millisecond, Timeout: time.Hour}
	defer func() { VPCWaitConfig = old }()

	calls, readyAfter := 0, 0 // readyAfter 0 = never ready
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		status := "Unavailable"
		if readyAfter > 0 && calls >= readyAfter {
			status = "Running"
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"vpcId":769,"metadata":{"status":"` + status + `"}}}`))
	})

	t.Run("times out at the given timeout", func(t *testing.T) {
		start := time.Now()
		err := c.WaitForVPCReadyTimeout(769, 30*time.Millisecond)
		if err == nil || !strings.Contains(err.Error(), "30ms") || time.Since(start) > 5*time.Second {
			t.Fatalf("want a timeout after 30ms, got %v after %s", err, time.Since(start))
		}
	})
	t.Run("ready", func(t *testing.T) {
		calls, readyAfter = 0, 3
		if err := c.WaitForVPCReadyTimeout(769, time.Minute); err != nil {
			t.Fatalf("want ready once the status is Running, got %v", err)
		}
	})
}
