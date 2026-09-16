package gona

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func newSDKCoverageTestV3Client(t *testing.T, handler http.HandlerFunc) *V3Client {
	t.Helper()
	c := NewV3Client("test-key", "https://example.invalid")
	c.httpClient.Transport = testRoundTripper{handler: handler}
	return c
}

func TestListMagicMeshesUsesRegisteredPath(t *testing.T) {
	c := newSDKCoverageTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cloud-routing/meshes" {
			t.Fatalf("path = %s, want /cloud-routing/meshes", r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "1000" {
			t.Fatalf("limit = %q, want 1000", r.URL.Query().Get("limit"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"meta":{"limit":1000,"offset":0,"total":1,"suggestedFilters":[]},"data":[{"meshId":17,"name":"mesh","description":"test"}]}}`))
	})

	meshes, err := c.ListMagicMeshes()
	if err != nil {
		t.Fatalf("ListMagicMeshes error: %v", err)
	}
	if len(meshes) != 1 || meshes[0].MeshID != 17 {
		t.Fatalf("meshes = %#v, want mesh 17", meshes)
	}
}

func TestAddOIDCClientBareMetalServersUsesServersBody(t *testing.T) {
	c := newSDKCoverageTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/oidc/clients/44/allow-list/bare-metal" {
			t.Fatalf("path = %s, want bare-metal allow-list path", r.URL.Path)
		}
		var body addOIDCClientBareMetalServersRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if len(body.Servers) != 2 || body.Servers[0].MBPkgID != 100 || body.Servers[1].MBPkgID != 101 {
			t.Fatalf("body = %#v, want two servers", body)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{}}`))
	})

	if err := c.AddOIDCClientBareMetalServers(44, []int{100, 101}); err != nil {
		t.Fatalf("AddOIDCClientBareMetalServers error: %v", err)
	}
}

func TestGetOIDCClientBareMetalServersReadsObjectAndIntegerRows(t *testing.T) {
	c := newSDKCoverageTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oidc/clients/44/allow-list/bare-metal" {
			t.Fatalf("path = %s, want bare-metal allow-list path", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"servers":[{"mbpkgid":100},101]}}`))
	})

	servers, err := c.GetOIDCClientBareMetalServers(44)
	if err != nil {
		t.Fatalf("GetOIDCClientBareMetalServers error: %v", err)
	}
	if len(servers) != 2 || servers[0].MBPkgID != 100 || servers[1].MBPkgID != 101 {
		t.Fatalf("servers = %#v, want mbpkgids 100 and 101", servers)
	}
}

func TestReplaceVPCNameserversUsesPutEnvelope(t *testing.T) {
	c := newSDKCoverageTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/vpcs/55/dhcp/nameservers" {
			t.Fatalf("path = %s, want nameserver path", r.URL.Path)
		}
		var body ReplaceVPCNameserversRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if len(body.Nameservers) != 1 || body.Nameservers[0].Server != "192.0.2.53" {
			t.Fatalf("body = %#v, want one nameserver", body)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"nameservers":[{"server":"192.0.2.53"}]}}`))
	})

	resp, err := c.ReplaceVPCNameservers(55, &ReplaceVPCNameserversRequest{
		Nameservers: []VPCNameserver{{Server: "192.0.2.53"}},
	})
	if err != nil {
		t.Fatalf("ReplaceVPCNameservers error: %v", err)
	}
	if len(resp.Nameservers) != 1 || resp.Nameservers[0].Server != "192.0.2.53" {
		t.Fatalf("response = %#v, want one nameserver", resp)
	}
}

func TestFirewallAvailableVMOptionsOmitUnsetValues(t *testing.T) {
	disable := true
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/firewall/sets/12/available-vm-list" {
			t.Fatalf("path = %s, want available-vm-list path", r.URL.Path)
		}
		if r.URL.Query().Get("disable_interface_id_filter") != "1" {
			t.Fatalf("disable_interface_id_filter = %q, want 1", r.URL.Query().Get("disable_interface_id_filter"))
		}
		if r.URL.Query().Has("bw") || r.URL.Query().Has("vpc_id") {
			t.Fatalf("unexpected query values: %s", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":[{"id":7,"mbpkgid":123,"hostname":"vm.example"}]}`))
	})

	vms, err := c.GetFirewallSetAvailableVMs(12, &FirewallAvailableVMOptions{DisableInterfaceIDFilter: &disable})
	if err != nil {
		t.Fatalf("GetFirewallSetAvailableVMs error: %v", err)
	}
	if len(vms) != 1 || vms[0].Mbpkgid != 123 {
		t.Fatalf("vms = %#v, want VM 123", vms)
	}
}

func TestReorderFirewallRulesUsesJSONBody(t *testing.T) {
	after := 3
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/firewall/sets/12/rules/re-order" {
			t.Fatalf("path = %s, want rule reorder path", r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if string(body) != `{"move_id":9,"after_id":3}` {
			t.Fatalf("body = %s, want move and after ids", string(body))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":{}}`))
	})

	if err := c.ReorderFirewallRules(12, &ReorderFirewallRulesRequest{MoveID: 9, AfterID: &after}); err != nil {
		t.Fatalf("ReorderFirewallRules error: %v", err)
	}
}
