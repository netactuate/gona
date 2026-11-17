package gona

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func checkCreateIdempotency(t *testing.T, key, path, wantBody string, create func(*Client, string) error) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != path {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Idempotency-Key"); got != key {
			t.Errorf("Idempotency-Key = %q, want %q", got, key)
		}
		if key == "" {
			if _, exists := r.Header["Idempotency-Key"]; exists {
				t.Error("empty key must omit the header")
			}
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
		}
		if string(body) != wantBody {
			t.Errorf("body = %q, want %q", body, wantBody)
		}
		if strings.Contains(string(body), "Idempotency") || (key != "" && strings.Contains(string(body), key)) {
			t.Error("idempotency key leaked into request body")
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"result":"success","code":200,"data":{"mbpkgid":7,"status":"success","build":8}}`)
	}))
	defer server.Close()
	if err := create(NewClientCustom("test-key", server.URL+"/"), key); err != nil {
		t.Fatalf("create: %v", err)
	}
}

func createIdempotencyServer(c *Client, key string) error {
	_, err := c.CreateServer(&CreateServerRequest{Plan: "basic", Location: 2, Image: 3, FQDN: "server.example", IdempotencyKey: key})
	return err
}

func TestIdempotencyKeySentOnCreate(t *testing.T) {
	checkCreateIdempotency(t, "create-retry", "/cloud/server/buy_build", "fqdn=server.example&image=3&location=2&plan=basic", createIdempotencyServer)
}

func TestEmptyIdempotencyKeyPreservesCreateBody(t *testing.T) {
	checkCreateIdempotency(t, "", "/cloud/server/buy_build", "fqdn=server.example&image=3&location=2&plan=basic", createIdempotencyServer)
}

func createIdempotencyMetal(c *Client, key string) error {
	_, err := c.CreateMetal(&CreateMetalRequest{Location: 2, Device: 3, Hostname: "metal.example", IdempotencyKey: key})
	return err
}

func TestIdempotencyKeySentOnCreateMetal(t *testing.T) {
	checkCreateIdempotency(t, "metal-retry", "/dedicated/server/buy_build", "device_id=3&fqdn=metal.example&location=2", createIdempotencyMetal)
}

func TestEmptyIdempotencyKeyPreservesCreateMetalBody(t *testing.T) {
	checkCreateIdempotency(t, "", "/dedicated/server/buy_build", "device_id=3&fqdn=metal.example&location=2", createIdempotencyMetal)
}

func createIdempotencyDedicated(c *Client, key string) error {
	_, err := c.BuyAndDeployDedicatedServer(3, &DedicatedServerBuildRequest{FQDN: "dedicated.example", Profile: 2, IdempotencyKey: key})
	return err
}

func TestIdempotencyKeySentOnDedicatedBuyBuild(t *testing.T) {
	checkCreateIdempotency(t, "dedicated-retry", "/dedicated/server/buy_build/3", `{"fqdn":"dedicated.example","profile":2}`, createIdempotencyDedicated)
}

func TestEmptyIdempotencyKeyPreservesDedicatedBuyBuildBody(t *testing.T) {
	checkCreateIdempotency(t, "", "/dedicated/server/buy_build/3", `{"fqdn":"dedicated.example","profile":2}`, createIdempotencyDedicated)
}
