package gona

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func newLongtailTestV3Client(t *testing.T, handler http.HandlerFunc) *V3Client {
	t.Helper()
	c := NewV3Client("test-key", "https://example.invalid")
	c.httpClient.Transport = testRoundTripper{handler: handler}
	return c
}

func TestGetGraphUsesDocumentedQueryParameters(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/graphs/graph" {
			t.Fatalf("expected /graphs/graph, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("port"); got != "12" {
			t.Fatalf("expected port 12, got %q", got)
		}
		if got := r.URL.Query().Get("time"); got != "weekly" {
			t.Fatalf("expected time weekly, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":{"series":[1,2,3]}}`))
	})

	graph, err := c.GetGraph(GraphQuery{Port: 12, Time: "weekly"})
	if err != nil {
		t.Fatalf("expected graph request to succeed, got: %v", err)
	}
	if !strings.Contains(string(graph.Raw), `"series"`) {
		t.Fatalf("expected raw graph data, got %s", string(graph.Raw))
	}
}

func TestServicesOptionalFilterIsOmittedWhenNil(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/services/colocation" {
			t.Fatalf("expected /services/colocation, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("service_id"); got != "" {
			t.Fatalf("expected no service_id, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":[{"id":7,"name":"colo"}]}`))
	})

	services, err := c.GetColocationServices(ServiceListOptions{})
	if err != nil {
		t.Fatalf("expected colocation services to succeed, got: %v", err)
	}
	if len(services) != 1 || services[0].ID != 7 {
		t.Fatalf("unexpected colocation services: %#v", services)
	}
}

func TestServicesOptionalFilterIsSentWhenSet(t *testing.T) {
	serviceID := 42
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/services/iptransit/ips" {
			t.Fatalf("expected /services/iptransit/ips, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("service_iptransit_id"); got != "42" {
			t.Fatalf("expected service_iptransit_id 42, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":[{"id":9,"service_iptransit_id":42,"ip":"203.0.113.1"}]}`))
	})

	addresses, err := c.GetIPTransitIPAddresses(&serviceID)
	if err != nil {
		t.Fatalf("expected IP transit addresses to succeed, got: %v", err)
	}
	if len(addresses) != 1 || addresses[0].ServiceIPTransitID != 42 {
		t.Fatalf("unexpected IP transit addresses: %#v", addresses)
	}
}

func TestSupportCreateTicketSendsJSONBody(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/support/tickets" {
			t.Fatalf("expected /support/tickets, got %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("expected application/json, got %q", got)
		}
		var body CreateTicketRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if body.Subject != "Need help" || body.Department != 4 || body.Urgency != "High" {
			t.Fatalf("unexpected request body: %#v", body)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"result":"success","code":201,"data":{"id":"T-1","subject":"Need help"}}`))
	})

	ticket, err := c.CreateTicket(&CreateTicketRequest{
		Subject:    "Need help",
		Message:    "Please assist",
		Department: 4,
		Urgency:    "High",
	})
	if err != nil {
		t.Fatalf("expected create ticket to succeed, got: %v", err)
	}
	if ticket.ID != "T-1" {
		t.Fatalf("unexpected ticket: %#v", ticket)
	}
}

func TestSupportAttachmentMetadataUsesWithoutDataOnlyWhenSet(t *testing.T) {
	withoutData := 1
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/support/tickets/T-1/attachment/reply/R-2/0" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("without_data"); got != "1" {
			t.Fatalf("expected without_data 1, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":{"name":"log.txt","size":12}}`))
	})

	attachment, err := c.GetTicketAttachment("T-1", "reply", "R-2", 0, &withoutData)
	if err != nil {
		t.Fatalf("expected attachment metadata to succeed, got: %v", err)
	}
	if attachment.Name != "log.txt" || attachment.Size != 12 {
		t.Fatalf("unexpected attachment: %#v", attachment)
	}
}

func TestSupportAttachmentDownloadReturnsRawBytes(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/support/tickets/T-1/attachment/ticket/0/1/download" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("file-bytes"))
	})

	body, err := c.DownloadTicketAttachment("T-1", "ticket", "0", 1)
	if err != nil {
		t.Fatalf("expected attachment download to succeed, got: %v", err)
	}
	if string(body) != "file-bytes" {
		t.Fatalf("unexpected body: %q", string(body))
	}
}

func TestSupportAttachmentDownloadReportsHTTPError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("missing"))
	})

	_, err := c.DownloadTicketAttachment("T-1", "ticket", "0", 1)
	if err == nil {
		t.Fatal("expected attachment download error")
	}
	if strings.Contains(err.Error(), "test-key") {
		t.Fatalf("expected redacted URL, got %v", err)
	}
}

func TestPlatformLookingGlassExecuteOmitsUnsetOptions(t *testing.T) {
	target := "example.net"
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/platform/looking-glass/execute" {
			t.Fatalf("expected /platform/looking-glass/execute, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("target"); got != target {
			t.Fatalf("expected target %q, got %q", target, got)
		}
		if got := r.URL.Query().Get("full"); got != "" {
			t.Fatalf("expected no full query value, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":{"output":"ok"}}`))
	})

	result, err := c.ExecutePlatformLookingGlass(PlatformLookingGlassExecuteOptions{Target: &target})
	if err != nil {
		t.Fatalf("expected looking glass execute to succeed, got: %v", err)
	}
	if !strings.Contains(string(result.Raw), `"output"`) {
		t.Fatalf("expected raw looking glass output, got %s", string(result.Raw))
	}
}

func TestV3ListMethodsUsePaginatedData(t *testing.T) {
	var calls int
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path != "/ssl-certificates" {
			t.Fatalf("expected /ssl-certificates, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		if calls == 1 {
			w.Write([]byte(`{"code":200,"data":{"meta":{"limit":1,"offset":0,"total":2,"suggestedFilters":[]},"data":[{"sslCertificateId":1,"name":"one"}]}}`))
			return
		}
		if got := r.URL.Query().Get("offset"); got != "1" {
			t.Fatalf("expected offset 1, got %q", got)
		}
		if r.Body != nil {
			body, _ := io.ReadAll(r.Body)
			if len(body) != 0 {
				t.Fatalf("expected empty GET body, got %q", string(body))
			}
		}
		w.Write([]byte(`{"code":200,"data":{"meta":{"limit":1,"offset":1,"total":2,"suggestedFilters":[]},"data":[{"sslCertificateId":2,"name":"two"}]}}`))
	})

	certs, err := c.GetSSLCertificates()
	if err != nil {
		t.Fatalf("expected SSL certificate list to succeed, got: %v", err)
	}
	if len(certs) != 2 || certs[1].SSLCertificateID != 2 {
		t.Fatalf("unexpected SSL certificates: %#v", certs)
	}
}

func TestV3HTTPLBGroupsList(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/http-loadbalancers/77/groups" {
			t.Fatalf("expected /http-loadbalancers/77/groups, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"meta":{"limit":1000,"offset":0,"total":1,"suggestedFilters":[]},"data":[{"httpGroupId":3,"name":"web","algorithm":"round-robin"}]}}`))
	})

	groups, err := c.GetHTTPLBGroups(77)
	if err != nil {
		t.Fatalf("expected HTTP LB groups list to succeed, got: %v", err)
	}
	if len(groups) != 1 || groups[0].HTTPGroupID != 3 {
		t.Fatalf("unexpected HTTP LB groups: %#v", groups)
	}
}
