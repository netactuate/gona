package gona

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestBindBGPGroupFirewallSet(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/bgp/bgp-groups/42/firewall-sets" {
			t.Fatalf("expected /bgp/bgp-groups/42/firewall-sets, got %s", r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		for _, want := range []string{`"id":1`, `"firewall_set_id":9`, `"interface_number":2`, `"set_priority":100`} {
			if !strings.Contains(string(body), want) {
				t.Fatalf("expected request body to contain %s, got %s", want, string(body))
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":{"id":1,"bgp2_group_id":42,"firewall_set_id":9,"interface_number":2,"set_priority":100}}`))
	})

	binding, err := c.BindBGPGroupFirewallSet(42, &BindBGPGroupFirewallSetRequest{
		ID:              1,
		FirewallSetID:   9,
		InterfaceNumber: 2,
		SetPriority:     100,
	})
	if err != nil {
		t.Fatalf("expected bind to succeed, got: %v", err)
	}
	if binding.BGPGroupID != 42 || binding.FirewallSetID != 9 {
		t.Fatalf("unexpected binding: %#v", binding)
	}
}

func TestUnbindBGPGroupFirewallSet(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/bgp/bgp-groups/42/firewall-sets/9" {
			t.Fatalf("expected /bgp/bgp-groups/42/firewall-sets/9, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200}`))
	})

	if err := c.UnbindBGPGroupFirewallSet(42, 9); err != nil {
		t.Fatalf("expected unbind to succeed, got: %v", err)
	}
}

func TestBGPGroupAndSessionActions(t *testing.T) {
	tests := []struct {
		name string
		call func(*Client) error
		path string
	}{
		{name: "refresh group", call: func(c *Client) error { return c.RefreshBGPGroupSessions(42) }, path: "/bgp/bgpgroup/42/refresh"},
		{name: "start group", call: func(c *Client) error { return c.StartBGPGroupSessions(42) }, path: "/bgp/bgpgroup/42/start"},
		{name: "stop group", call: func(c *Client) error { return c.StopBGPGroupSessions(42) }, path: "/bgp/bgpgroup/42/stop"},
		{name: "refresh session", call: func(c *Client) error { return c.RefreshBGPSession(99) }, path: "/bgp/bgpsession/99/refresh"},
		{name: "start session", call: func(c *Client) error { return c.StartBGPSession(99) }, path: "/bgp/bgpsession/99/start"},
		{name: "stop session", call: func(c *Client) error { return c.StopBGPSession(99) }, path: "/bgp/bgpsession/99/stop"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Fatalf("expected POST, got %s", r.Method)
				}
				if r.URL.Path != tt.path {
					t.Fatalf("expected %s, got %s", tt.path, r.URL.Path)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"result":"success","code":200}`))
			})

			if err := tt.call(c); err != nil {
				t.Fatalf("expected action to succeed, got: %v", err)
			}
		})
	}
}

func TestGetBGPSummary(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bgp/bgpsummary" {
			t.Fatalf("expected /bgp/bgpsummary, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":{"sessions":2,"prefixes":1}}`))
	})

	summary, err := c.GetBGPSummary()
	if err != nil {
		t.Fatalf("expected summary to succeed, got: %v", err)
	}
	if string(summary["sessions"]) != "2" {
		t.Fatalf("unexpected summary: %#v", summary)
	}
}

func TestGetBGPDashboard(t *testing.T) {
	flapWindow := 15
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bgp/dashboard" {
			t.Fatalf("expected /bgp/dashboard, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("group_type"); got != "anycast" {
			t.Fatalf("expected group_type anycast, got %q", got)
		}
		if got := r.URL.Query().Get("flap_window"); got != "15" {
			t.Fatalf("expected flap_window 15, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":{"flapping_sessions":1}}`))
	})

	dashboard, err := c.GetBGPDashboard(BGPDashboardOptions{GroupType: "anycast", FlapWindow: &flapWindow})
	if err != nil {
		t.Fatalf("expected dashboard to succeed, got: %v", err)
	}
	if string(dashboard["flapping_sessions"]) != "1" {
		t.Fatalf("unexpected dashboard: %#v", dashboard)
	}
}
