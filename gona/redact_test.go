package gona

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRedactBodyMasksCredentials(t *testing.T) {
	keys := []string{"password", "root_password", "kubeconfig", "token", "access_key", "secret_key", "secret", "private_key", "api_key", "key"}
	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			for _, body := range []string{
				`{"items":[{"` + strings.ToUpper(key) + `":"planted-credential","label":"visible"}]}`,
				strings.ToUpper(key) + "=planted-credential&label=visible",
			} {
				got := redactBody("ordinary", []byte(body))
				if strings.Contains(got, "planted-credential") || !strings.Contains(got, "visible") {
					t.Fatalf("unsafe or overbroad masking: %s", got)
				}
				if !strings.Contains(got, "[redacted]") && !strings.Contains(got, "%5Bredacted%5D") {
					t.Fatalf("missing redaction marker: %s", got)
				}
			}
		})
	}
	for _, body := range []string{`{"items":[{"VaLuE":"planted-value"}]}`, "VALUE=planted-value"} {
		if got := redactBody("/secrets/list", []byte(body)); strings.Contains(got, "planted-value") {
			t.Fatalf("secret value leaked: %s", got)
		}
		if got := redactBody("/ordinary", []byte(body)); !strings.Contains(got, "planted-value") {
			t.Fatalf("ordinary value masked: %s", got)
		}
	}
	for _, path := range []string{"/kubeconfig/download", "/secrets/item", "/storage/item"} {
		if got := redactBody(path, []byte("unstructured confidential body")); got != "[body redacted]" {
			t.Errorf("%s: unexpected fallback %q", path, got)
		}
	}
}

func TestS7_18_DebugLogMasksKubeconfigAndSecretValues(t *testing.T) {
	t.Setenv("NA_API_DEBUG", "1")
	var logs bytes.Buffer
	oldOutput := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(oldOutput) })
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/kubeconfig":
			io.WriteString(w, `{"code":200,"data":{"kubeconfig":"planted-kubeconfig"}}`)
		case "/secrets/list":
			body, err := io.ReadAll(r.Body)
			if err != nil || !bytes.Contains(body, []byte("planted-secret-value")) {
				t.Errorf("request body was changed: %s, %v", body, err)
			}
			io.WriteString(w, `{"code":200,"data":{}}`)
		default:
			if err := r.ParseForm(); err != nil || r.Form.Get("password") != "planted-password" {
				t.Errorf("form body was changed: %v, %v", r.Form, err)
			}
			io.WriteString(w, `{"result":"success","code":200,"data":{}}`)
		}
	}))
	defer server.Close()
	v3 := NewV3Client("test-key", server.URL+"/")
	response, err := v3.doRequest(http.MethodGet, "/kubeconfig", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(response.Data, []byte("planted-kubeconfig")) {
		t.Fatal("response body was changed")
	}
	if _, err := v3.doRequest(http.MethodPost, "/secrets/list", map[string]interface{}{"items": []interface{}{map[string]string{"value": "planted-secret-value"}}}); err != nil {
		t.Fatal(err)
	}
	v2 := NewClientCustom("test-key", server.URL+"/")
	var responseData json.RawMessage
	if err := v2.post(context.Background(), "cloud/server/buy_build", []byte("password=planted-password"), &responseData); err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"planted-kubeconfig", "planted-secret-value", "planted-password"} {
		if strings.Contains(logs.String(), secret) {
			t.Errorf("debug log leaked %s", secret)
		}
	}
	if !strings.Contains(logs.String(), "[redacted]") {
		t.Error("debug log has no redaction marker")
	}
}

// The V2 and V3 APIs both take the API key as a query parameter, so any URL that reaches an
// error message or a log carries a live credential. Terraform prints those errors to the
// operator's terminal, their CI logs, and whatever support ticket follows. This test exists
// so that a future change cannot quietly put the key back.
func TestRedactURLRemovesAPIKey(t *testing.T) {
	const secret = "EXAMPLE-NOT-A-REAL-KEY-0000000000000000000000000000"

	u, err := url.Parse("https://vapi2.netactuate.com/api/dns/zone?key=" + secret)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	got := redactURL(u)
	if strings.Contains(got, secret) {
		t.Fatalf("redactURL leaked the API key: %s", got)
	}
	if !strings.Contains(got, "REDACTED") {
		t.Fatalf("redactURL did not mark the key as redacted: %s", got)
	}
	if !strings.Contains(got, "vapi2.netactuate.com") {
		t.Fatalf("redactURL destroyed the useful part of the URL: %s", got)
	}

	gotV3 := redactV3URL(u)
	if strings.Contains(gotV3, secret) {
		t.Fatalf("redactV3URL leaked the API key: %s", gotV3)
	}

	// A URL with no key must survive untouched, so ordinary diagnostics stay readable.
	plain, _ := url.Parse("https://vapi2.netactuate.com/api/dns/zones?type=NATIVE")
	if redactURL(plain) != plain.String() {
		t.Fatalf("redactURL altered a URL that had no key: %s", redactURL(plain))
	}
}

// A kubeconfig served as YAML is not JSON, and its base64 padding makes it look like form data.
// It must be withheld, not logged url-encoded.
func TestS7_18_DebugLogWithholdsNonJSONKubeconfig(t *testing.T) {
	yaml := "apiVersion: v1\nclusters:\n- cluster:\n    certificate-authority-data: PLANTEDCA==\nusers:\n- user:\n    token: PLANTEDTOKEN\n"
	got := redactBody("nke/clusters/1/kubeconfig", []byte(yaml))
	if strings.Contains(got, "PLANTED") {
		t.Fatalf("kubeconfig reached the log: %s", got)
	}
}
