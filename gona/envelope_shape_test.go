package gona

import (
	"net/http"
	"strconv"
	"testing"
)

func TestV2LaravelPaginatorListUnwrapsPageData(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/services/colocation" {
			t.Fatalf("path = %s, want /services/colocation", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","message":null,"meta":[],"data":{"paginator":{"current_page":1,"data":[{"id":7,"name":"colo"}],"first_page_url":"https://example.invalid/services/colocation?page=1","from":1,"last_page":1,"last_page_url":"https://example.invalid/services/colocation?page=1"}}}`))
	})

	services, err := c.GetColocationServices(ServiceListOptions{})
	if err != nil {
		t.Fatalf("GetColocationServices error: %v", err)
	}
	if len(services) != 1 || services[0].ID != 7 {
		t.Fatalf("services = %#v, want service 7", services)
	}
}

func TestV2LaravelPaginatorFollowsPages(t *testing.T) {
	seen := map[int]bool{}
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/services/iptransit" {
			t.Fatalf("path = %s, want /services/iptransit", r.URL.Path)
		}
		page := 1
		if got := r.URL.Query().Get("page"); got != "" {
			parsed, err := strconv.Atoi(got)
			if err != nil {
				t.Fatalf("page query = %q", got)
			}
			page = parsed
		}
		seen[page] = true
		w.WriteHeader(http.StatusOK)
		switch page {
		case 1:
			w.Write([]byte(`{"result":"success","message":null,"meta":[],"data":{"paginator":{"current_page":1,"data":[{"id":1,"name":"first"}],"first_page_url":"https://example.invalid/services/iptransit?page=1","from":1,"last_page":2,"last_page_url":"https://example.invalid/services/iptransit?page=2"}}}`))
		case 2:
			w.Write([]byte(`{"result":"success","message":null,"meta":[],"data":{"paginator":{"current_page":2,"data":[{"id":2,"name":"second"}],"first_page_url":"https://example.invalid/services/iptransit?page=1","from":2,"last_page":2,"last_page_url":"https://example.invalid/services/iptransit?page=2"}}}`))
		default:
			t.Fatalf("unexpected page %d", page)
		}
	})

	services, err := c.GetIPTransitServices(ServiceListOptions{})
	if err != nil {
		t.Fatalf("GetIPTransitServices error: %v", err)
	}
	if len(services) != 2 || services[0].ID != 1 || services[1].ID != 2 {
		t.Fatalf("services = %#v, want both pages", services)
	}
	if !seen[1] || !seen[2] {
		t.Fatalf("seen pages = %#v, want pages 1 and 2", seen)
	}
}

func TestV2LaravelPaginatorUnwrapsDedicatedFilter(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dedicated/filter-dedicated-devices" {
			t.Fatalf("path = %s, want /dedicated/filter-dedicated-devices", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","message":null,"meta":[],"data":{"paginator":{"current_page":1,"data":[{"device_id":101,"dc_name":"nyc"}],"first_page_url":"https://example.invalid/dedicated/filter-dedicated-devices?page=1","from":1,"last_page":1,"last_page_url":"https://example.invalid/dedicated/filter-dedicated-devices?page=1"}}}`))
	})

	devices, err := c.FilterDedicatedDevices(DedicatedDeviceFilterOptions{})
	if err != nil {
		t.Fatalf("FilterDedicatedDevices error: %v", err)
	}
	if len(devices) != 1 || string(devices[0]["device_id"]) != "101" {
		t.Fatalf("devices = %#v, want device 101", devices)
	}
}

func TestV2LaravelPaginatorUnwrapsDirectDedicatedFilter(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dedicated/filter-dedicated-devices" {
			t.Fatalf("path = %s, want /dedicated/filter-dedicated-devices", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","message":null,"meta":[],"data":{"current_page":1,"data":[{"device_id":101,"dc_name":"nyc"}],"first_page_url":"https://example.invalid/dedicated/filter-dedicated-devices?page=1","from":1,"last_page":1,"last_page_url":"https://example.invalid/dedicated/filter-dedicated-devices?page=1"}}`))
	})

	devices, err := c.FilterDedicatedDevices(DedicatedDeviceFilterOptions{})
	if err != nil {
		t.Fatalf("FilterDedicatedDevices error: %v", err)
	}
	if len(devices) != 1 || string(devices[0]["device_id"]) != "101" {
		t.Fatalf("devices = %#v, want device 101", devices)
	}
}

func TestV2LaravelPaginatorUnwrapsDirectDedicatedLocations(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/dedicated/locations" {
			t.Fatalf("path = %s, want /dedicated/locations", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		// The live endpoint answers a plain list of short_name, pub_description and location_id.
		w.Write([]byte(`{"result":"success","message":null,"meta":[],"data":[{"short_name":"JNB","pub_description":"Johannesburg facility","location_id":3}]}`))
	})

	locations, err := c.ListDedicatedLocations()
	if err != nil {
		t.Fatalf("ListDedicatedLocations error: %v", err)
	}
	if len(locations) != 1 || locations[0].LocationID != 3 || locations[0].ShortName != "JNB" {
		t.Fatalf("locations = %#v, want JNB at location 3", locations)
	}
}

func TestV2LaravelPaginatorUnwrapsDirectFirewallExternalIPSets(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/firewall/external-ipsets" {
			t.Fatalf("path = %s, want /firewall/external-ipsets", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","message":null,"meta":[],"data":{"current_page":1,"data":[{"id":9,"name":"allowlist","description":"allowed ranges"}],"first_page_url":"https://example.invalid/firewall/external-ipsets?page=1","from":1,"last_page":1,"last_page_url":"https://example.invalid/firewall/external-ipsets?page=1"}}`))
	})

	sets, err := c.GetFirewallExternalIPSets()
	if err != nil {
		t.Fatalf("GetFirewallExternalIPSets error: %v", err)
	}
	if len(sets) != 1 || sets[0].ID != 9 {
		t.Fatalf("sets = %#v, want set 9", sets)
	}
}

func TestListStorageTypesUnwrapsObjectValues(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/storage" {
			t.Fatalf("path = %s, want /storage", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"block":{"name":"Block","description":"Block storage","features":[]},"object":{"type":"object","name":"Object","description":"Object storage"}}}`))
	})

	types, err := c.ListStorageTypes()
	if err != nil {
		t.Fatalf("ListStorageTypes error: %v", err)
	}
	if len(types) != 2 || types[0].Type != "block" || types[1].Type != "object" {
		t.Fatalf("types = %#v, want block and object", types)
	}
}

func TestListStorageTypesUnwrapsTwoLevelList(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/storage" {
			t.Fatalf("path = %s, want /storage", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"meta":{"limit":1000,"offset":0,"total":1,"suggestedFilters":[]},"data":[{"type":"block","name":"Block","description":"Block storage","features":[]}]}}`))
	})

	types, err := c.ListStorageTypes()
	if err != nil {
		t.Fatalf("ListStorageTypes error: %v", err)
	}
	if len(types) != 1 || types[0].Type != "block" {
		t.Fatalf("types = %#v, want block", types)
	}
}

func TestV3ListUnwrapsNestedListObject(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/vpcs" {
			t.Fatalf("path = %s, want /vpcs", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"vpcs":{"meta":{"limit":1000,"offset":0,"total":1,"suggestedFilters":[]},"data":[{"vpcId":55,"metadata":{"label":"test","status":"Running"},"location":{"id":3,"name":"New York"},"floatingIps":{}}]}}}`))
	})

	vpcs, err := c.ListVPCs()
	if err != nil {
		t.Fatalf("ListVPCs error: %v", err)
	}
	if len(vpcs) != 1 || vpcs[0].VPCID != 55 {
		t.Fatalf("vpcs = %#v, want VPC 55", vpcs)
	}
}

func TestV3ListUnwrapsTwoLevelEnvelopeForVPCs(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/vpcs" {
			t.Fatalf("path = %s, want /vpcs", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"meta":{"limit":1000,"offset":0,"total":1,"suggestedFilters":[]},"data":[{"vpcId":55,"metadata":{"label":"test","status":"Running"},"location":{"id":3,"name":"New York"},"floatingIps":{},"loadBalancers":{"network":{"applied":0,"deleteRequested":0,"total":0},"http":{"applied":1,"deleteRequested":0,"total":1}}}]}}`))
	})

	vpcs, err := c.ListVPCs()
	if err != nil {
		t.Fatalf("ListVPCs error: %v", err)
	}
	if len(vpcs) != 1 || vpcs[0].VPCID != 55 {
		t.Fatalf("vpcs = %#v, want VPC 55", vpcs)
	}
	if vpcs[0].LoadBalancers == nil || vpcs[0].LoadBalancers.HTTP == nil || vpcs[0].LoadBalancers.HTTP.Total != 1 {
		t.Fatalf("load balancers = %#v, want HTTP total 1", vpcs[0].LoadBalancers)
	}
}

func TestVPCSSHSettingsAcceptsObjectPort(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/vpcs/55/ssh" {
			t.Fatalf("path = %s, want /vpcs/55/ssh", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"enabled":true,"port":{"port":2222},"bastion":{"ipv4":"192.0.2.10","ipv6":"2001:db8::10"}}}`))
	})

	settings, err := c.GetVPCSSHSettings(55)
	if err != nil {
		t.Fatalf("GetVPCSSHSettings error: %v", err)
	}
	if settings.Port == nil || *settings.Port != 2222 || !settings.Enabled {
		t.Fatalf("settings = %#v, want enabled port 2222", settings)
	}
}

func TestVPCSSHSettingsAcceptsObjectKeys(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/vpcs/55/ssh" {
			t.Fatalf("path = %s, want /vpcs/55/ssh", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"enabled":true,"port":2222,"keys":{"123":{"id":123,"sshKeyId":123,"name":"ops","dates":{"created":"2026-01-01","enabled":"2026-01-01"}}},"bastion":{"ipv4":"192.0.2.10","ipv6":"2001:db8::10"}}}`))
	})

	settings, err := c.GetVPCSSHSettings(55)
	if err != nil {
		t.Fatalf("GetVPCSSHSettings error: %v", err)
	}
	key := settings.Keys["123"]
	if len(settings.Keys) != 1 || key.GetID() != 123 {
		t.Fatalf("keys = %#v, want key 123", settings.Keys)
	}
}
