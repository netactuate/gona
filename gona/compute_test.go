package gona

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

const serverStatusResponseFixture = `{
	"result": "success",
	"message": "",
	"code": 200,
	"data": {
		"status": "RUNNING",
		"state": "UP"
	}
}`

const cloudFloatingIPv4ListFixture = `{
	"code": 200,
	"data": {
		"meta": {
			"limit": 1000,
			"offset": 0,
			"total": 1,
			"suggestedFilters": []
		},
		"data": [
			{
				"floatingIpv4Id": 77,
				"AssignedOn": "2026-09-15T00:00:00Z",
				"address": "203.0.113.10",
				"vlanId": 22,
				"ptrDomain": "host.example.test"
			}
		]
	}
}`

func TestGetServerStatusUsesRegisteredPath(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/cloud/server/456/status" {
			t.Fatalf("expected registered status path, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(serverStatusResponseFixture))
	})

	status, err := c.GetServerStatus(456)
	if err != nil {
		t.Fatalf("GetServerStatus: %v", err)
	}
	if status.Status != "RUNNING" || status.State != "UP" || len(status.Raw) == 0 {
		t.Fatalf("unexpected status: %#v", status)
	}
}

func TestDeleteServerWithOptionsOmitsUnsetPointerFields(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/cloud/server/456/delete" {
			t.Fatalf("expected registered delete path, got %s", r.URL.Path)
		}
		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if _, ok := body["cancel_billing"]; !ok {
			t.Fatal("expected cancel_billing to be present")
		}
		if _, ok := body["password"]; ok {
			t.Fatal("expected nil password pointer to be omitted")
		}
		if _, ok := body["force_password"]; ok {
			t.Fatal("expected nil force_password pointer to be omitted")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":{"id":456}}`))
	})

	resp, err := c.DeleteServerWithOptions(456, &DeleteServerRequest{CancelBilling: BoolPtr(true)})
	if err != nil {
		t.Fatalf("DeleteServerWithOptions: %v", err)
	}
	if resp.ID != 456 {
		t.Fatalf("unexpected delete response: %#v", resp)
	}
}

func TestGetScalingOptionsPreservesPointerZero(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cloud/scaling/456" {
			t.Fatalf("expected /cloud/scaling/456, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("min_ram"); got != "0" {
			t.Fatalf("expected min_ram=0 from non-nil pointer, got %q", got)
		}
		if r.URL.Query().Has("max_ram") {
			t.Fatal("expected nil max_ram pointer to be omitted")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":[{"plan":"VR1x1x25"}]}`))
	})

	zero := 0
	raw, err := c.GetScalingOptions(456, &ScalingOptionsRequest{MinRAM: &zero})
	if err != nil {
		t.Fatalf("GetScalingOptions: %v", err)
	}
	if !strings.Contains(string(raw), "VR1x1x25") {
		t.Fatalf("unexpected scaling payload: %s", string(raw))
	}
}

func TestParseCloudInitUploadsMultipartFile(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/cloud/parse-cloud-init" {
			t.Fatalf("expected parse-cloud-init path, got %s", r.URL.Path)
		}
		if err := r.ParseMultipartForm(1024); err != nil {
			t.Fatalf("parse multipart form: %v", err)
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			t.Fatalf("get file part: %v", err)
		}
		defer file.Close()
		content, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("read file part: %v", err)
		}
		if string(content) != "#cloud-config\n" {
			t.Fatalf("unexpected file content: %q", string(content))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":{"valid":true}}`))
	})

	raw, err := c.ParseCloudInit("user-data.yaml", strings.NewReader("#cloud-config\n"))
	if err != nil {
		t.Fatalf("ParseCloudInit: %v", err)
	}
	if !strings.Contains(string(raw), "valid") {
		t.Fatalf("unexpected parse response: %s", string(raw))
	}
}

func TestListCloudFloatingIPv4RecordedResponse(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/cloud/networking/floating-ips/ipv4" {
			t.Fatalf("expected floating IPv4 list path, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cloudFloatingIPv4ListFixture))
	})

	floatingIPs, err := c.ListCloudFloatingIPv4()
	if err != nil {
		t.Fatalf("ListCloudFloatingIPv4: %v", err)
	}
	if len(floatingIPs) != 1 || floatingIPs[0].FloatingIPv4ID != 77 || floatingIPs[0].Address != "203.0.113.10" {
		t.Fatalf("unexpected floating IPs: %#v", floatingIPs)
	}
}

func TestCreateCloudFloatingIPv4OmitsUnsetPointerFields(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/cloud/networking/floating-ips/ipv4" {
			t.Fatalf("expected floating IPv4 create path, got %s", r.URL.Path)
		}
		var body map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if _, ok := body["ptrDomain"]; ok {
			t.Fatal("expected nil ptrDomain pointer to be omitted")
		}
		if _, ok := body["vlanId"]; !ok {
			t.Fatal("expected vlanId pointer to be present")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":200,"data":{"floatingIpv4Id":77,"address":"203.0.113.10","vlanId":22}}`))
	})

	vlanID := 22
	floatingIP, err := c.CreateCloudFloatingIPv4(&CreateCloudFloatingIPv4Request{VLANID: &vlanID})
	if err != nil {
		t.Fatalf("CreateCloudFloatingIPv4: %v", err)
	}
	if floatingIP.FloatingIPv4ID != 77 || floatingIP.VLANID != 22 {
		t.Fatalf("unexpected floating IP: %#v", floatingIP)
	}
}

func TestServerAddressUnmarshalErrorsPropagate(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cloud/server/456/ipv4" {
			t.Fatalf("expected IPv4 path, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":[{"id":"not-an-int","ip":"203.0.113.20"}]}`))
	})

	if _, err := c.GetServerIPv4(456); err == nil {
		t.Fatal("expected unmarshal error to propagate")
	}
}
