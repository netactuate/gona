package gona

import (
	"encoding/json"
	"testing"
)

func TestOIDCClientListMeasuredShape(t *testing.T) {
	raw := []byte(`{
		"tenant": 288,
		"clients": {
			"meta": {"total": 1, "offset": 0, "limit": 1000, "suggestedFilters": []},
			"data": [{
				"clientId": 123,
				"createdOn": "2026-09-11T00:00:00Z",
				"lastUsedOn": null,
				"label": "nah-test",
				"description": "test",
				"jwksHttpsUrl": "https://keys.example.com/.well-known/jwks.json",
				"accountDefault": false,
				"defaultAudience": "https://api.example.com/oidc-audience",
				"ttl": 300,
				"enforceAllowList": 0
			}]
		}
	}`)

	var wrapped oidcClientsResponse
	if err := json.Unmarshal(raw, &wrapped); err != nil {
		t.Fatalf("unmarshal OIDC clients response: %v", err)
	}
	if wrapped.Tenant.String() != "288" {
		t.Fatalf("tenant = %q, want 288", wrapped.Tenant.String())
	}

	var rows []oidcClientListRow
	if err := json.Unmarshal(wrapped.Clients.Data, &rows); err != nil {
		t.Fatalf("unmarshal OIDC client rows: %v", err)
	}
	client := rows[0].toClient(wrapped.Tenant)
	if client.JWKSURI == nil || *client.JWKSURI != "https://keys.example.com/.well-known/jwks.json" {
		t.Fatalf("jwks_uri = %#v, want list jwksHttpsUrl value", client.JWKSURI)
	}
	if client.LastUsedOn != nil {
		t.Fatalf("last used = %#v, want nil", client.LastUsedOn)
	}
	if client.EnforceAllowList.Bool() {
		t.Fatalf("enforce allow list = true, want false from numeric 0")
	}
}

func TestOIDCBoolAcceptsBoolAndNumericValues(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want bool
	}{
		{`true`, true},
		{`false`, false},
		{`1`, true},
		{`0`, false},
		{`"1"`, true},
		{`"0"`, false},
	} {
		var got OIDCBool
		if err := json.Unmarshal([]byte(tc.raw), &got); err != nil {
			t.Fatalf("unmarshal %s: %v", tc.raw, err)
		}
		if got.Bool() != tc.want {
			t.Fatalf("unmarshal %s = %t, want %t", tc.raw, got.Bool(), tc.want)
		}
	}
}
