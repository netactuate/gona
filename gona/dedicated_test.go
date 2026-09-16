package gona

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestFilterDedicatedDevices(t *testing.T) {
	perPage := 25
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dedicated/filter-dedicated-devices" {
			t.Fatalf("expected /dedicated/filter-dedicated-devices, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("per_page"); got != "25" {
			t.Fatalf("expected per_page 25, got %q", got)
		}
		if got := r.URL.Query().Get("ram_mb"); got != "32768,65536" {
			t.Fatalf("expected ram_mb range, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
		// The live endpoint nests the rows at data.devices.paginator.data, alongside the column
		// lists the portal uses.
		w.Write([]byte(`{"result":"success","code":200,"data":{"devices":{"paginator":{"current_page":1,"last_page":1,"data":[{"device_id":101,"dc_name":"nyc"}]},"sortable_columns":[],"filterable_columns":[]},"filters":{}}}`))
	})

	devices, err := c.FilterDedicatedDevices(DedicatedDeviceFilterOptions{PerPage: &perPage, RAMMB: "32768,65536"})
	if err != nil {
		t.Fatalf("expected filter to succeed, got: %v", err)
	}
	if len(devices) != 1 || string(devices[0]["device_id"]) != "101" {
		t.Fatalf("unexpected devices: %#v", devices)
	}
}

func TestDedicatedReadMethods(t *testing.T) {
	tests := []struct {
		name string
		call func(*Client) error
		path string
	}{
		{
			name: "locations",
			path: "/dedicated/locations",
			call: func(c *Client) error {
				locations, err := c.ListDedicatedLocations()
				if err != nil {
					return err
				}
				if len(locations) != 1 || locations[0].ShortName != "JNB" {
					t.Fatalf("unexpected locations: %#v", locations)
				}
				return nil
			},
		},
		{
			name: "device os",
			path: "/dedicated/os/device/101",
			call: func(c *Client) error {
				yes := true
				profiles, err := c.ListDedicatedDeviceOSProfiles(101, &yes)
				if err != nil {
					return err
				}
				if len(profiles) != 1 || profiles[0].OSID != 7 {
					t.Fatalf("unexpected profiles: %#v", profiles)
				}
				return nil
			},
		},
		{
			name: "plans",
			path: "/dedicated/plans/3",
			call: func(c *Client) error {
				plans, err := c.ListDedicatedPlans(3)
				if err != nil {
					return err
				}
				if len(plans) != 1 || string(plans[0]["id"]) != "11" {
					t.Fatalf("unexpected plans: %#v", plans)
				}
				return nil
			},
		},
		{
			name: "servers",
			path: "/dedicated/servers",
			call: func(c *Client) error {
				servers, err := c.ListDedicatedServers()
				if err != nil {
					return err
				}
				if len(servers) != 1 || servers[0].MBPKGID != 55 {
					t.Fatalf("unexpected servers: %#v", servers)
				}
				return nil
			},
		},
	}

	responses := map[string]string{
		"/dedicated/locations":     `{"result":"success","code":200,"data":[{"short_name":"JNB","pub_description":"Johannesburg facility","location_id":3}]}`,
		"/dedicated/os/device/101": `{"result":"success","code":200,"data":[{"id":7,"name":"Ubuntu","disklayouts":null,"scripts":null}]}`,
		"/dedicated/plans/3":       `{"result":"success","code":200,"data":[{"id":11,"name":"E3"}]}`,
		"/dedicated/servers":       `{"result":"success","code":200,"data":[{"mbpkgid":55,"hostname":"metal.example.com","ob_id":123,"building":null}]}`,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Fatalf("expected GET, got %s", r.Method)
				}
				if r.URL.Path != tt.path {
					t.Fatalf("expected %s, got %s", tt.path, r.URL.Path)
				}
				if tt.path == "/dedicated/os/device/101" && r.URL.Query().Get("is_buyable") != "1" {
					t.Fatalf("expected is_buyable=1, got %q", r.URL.RawQuery)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(responses[tt.path]))
			})

			if err := tt.call(c); err != nil {
				t.Fatalf("expected read to succeed, got: %v", err)
			}
		})
	}
}

func TestDedicatedBuildAndBuyMethods(t *testing.T) {
	tests := []struct {
		name string
		call func(*Client) (MetalBuild, error)
		path string
	}{
		{name: "deploy", call: func(c *Client) (MetalBuild, error) {
			password := "secret"
			return c.DeployDedicatedServer(55, &DedicatedServerBuildRequest{FQDN: "metal.example.com", Profile: 7, RootPassword: &password})
		}, path: "/dedicated/server/build/55"},
		{name: "buy", call: func(c *Client) (MetalBuild, error) {
			return c.BuyDedicatedServer(101)
		}, path: "/dedicated/server/buy/101"},
		{name: "buy build", call: func(c *Client) (MetalBuild, error) {
			sshKeyID := 22
			return c.BuyAndDeployDedicatedServer(101, &DedicatedServerBuildRequest{FQDN: "metal.example.com", Profile: 7, SSHKeyID: &sshKeyID})
		}, path: "/dedicated/server/buy_build/101"},
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
				w.Write([]byte(`{"result":"success","code":200,"data":{"mbpkgid":55,"status":"queued","build":99}}`))
			})

			build, err := tt.call(c)
			if err != nil {
				t.Fatalf("expected build method to succeed, got: %v", err)
			}
			if build.MBPKGID != 55 || build.Build != 99 {
				t.Fatalf("unexpected build response: %#v", build)
			}
		})
	}
}

func TestDedicatedBuildRequestOmitsUnsetOptionalFields(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		for _, notWant := range []string{"disklayout", "root_password", "ssh_key", "ssh_key_id", "build_script"} {
			if strings.Contains(string(body), notWant) {
				t.Fatalf("expected %s to be omitted from %s", notWant, string(body))
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":{"mbpkgid":55,"status":"queued","build":99}}`))
	})

	if _, err := c.DeployDedicatedServer(55, &DedicatedServerBuildRequest{FQDN: "metal.example.com", Profile: 7}); err != nil {
		t.Fatalf("expected deploy to succeed, got: %v", err)
	}
}

func TestDedicatedPowerAndDeleteActions(t *testing.T) {
	tests := []struct {
		name string
		call func(*Client) error
		path string
	}{
		{name: "soft reset", call: func(c *Client) error { return c.SoftResetDedicatedServer(55) }, path: "/dedicated/server/soft-reset/55"},
		{name: "delete", call: func(c *Client) error {
			password := "confirm"
			return c.DeleteDedicatedServer(55, &DedicatedServerActionRequest{Password: &password})
		}, path: "/dedicated/server/55/delete"},
		{name: "reboot", call: func(c *Client) error {
			force := false
			return c.RebootDedicatedServer(55, &DedicatedServerActionRequest{Force: &force})
		}, path: "/dedicated/server/55/reboot"},
		{name: "shutdown", call: func(c *Client) error { return c.ShutdownDedicatedServer(55, nil) }, path: "/dedicated/server/55/shutdown"},
		{name: "start", call: func(c *Client) error { return c.StartDedicatedServer(55, nil) }, path: "/dedicated/server/55/start"},
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

func TestDedicatedPowerStatus(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/dedicated/server/55/status" {
			t.Fatalf("expected /dedicated/server/55/status, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":{"state":"on"}}`))
	})

	status, err := c.GetDedicatedServerPowerStatus(55, nil)
	if err != nil {
		t.Fatalf("expected status to succeed, got: %v", err)
	}
	if string(status["state"]) != `"on"` {
		t.Fatalf("unexpected status: %#v", status)
	}
}

func TestUpdateDedicatedServerIPv4Reverse(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/dedicated/server/ipv4_reverse" {
			t.Fatalf("expected /dedicated/server/ipv4_reverse, got %s", r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !strings.Contains(string(body), `"reverse":"host.example.com"`) {
			t.Fatalf("unexpected body: %s", string(body))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200}`))
	})

	err := c.UpdateDedicatedServerIPv4Reverse(&DedicatedIPv4ReverseRequest{
		MBPkgID: 55,
		ID:      1,
		Reverse: "host.example.com",
	})
	if err != nil {
		t.Fatalf("expected reverse update to succeed, got: %v", err)
	}
}
