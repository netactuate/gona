package gona

import (
	"net/http"
	"testing"
)

// ListAPIACLs decodes entries with and without resources.
func TestListAPIACLs(t *testing.T) {
	var gotQuery string
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/account/api-acls" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		gotQuery = r.URL.Query().Get("sub_user_id")
		w.Write([]byte(`{"result":"success","message":null,"meta":[],"code":200,"data":[
			{"id":11,"type":4,"start_ip":"192.0.2.10","end_ip":"192.0.2.10","bits":32,"mb_id":5,"sub_id":0,
			 "created_at":"2026-11-17 10:00:00","description":"VPC 901 egress","editable":true,
			 "resources":[{"resource_type":"vpc","resource_id":901,"resource_label":"web"},
			              {"resource_type":"nke_cluster","resource_id":587,"resource_label":null}]},
			{"id":12,"type":4,"start_ip":"203.0.113.0","end_ip":"203.0.113.255","bits":24,"mb_id":5,"sub_id":0,
			 "created_at":"2026-11-01 09:00:00","description":null,"editable":false,"resources":[]}
		]}`))
	})

	acls, err := c.ListAPIACLs(nil)
	if err != nil {
		t.Fatalf("ListAPIACLs: %v", err)
	}
	if gotQuery != "" {
		t.Fatalf("sub_user_id sent without being set: %q", gotQuery)
	}
	if len(acls) != 2 {
		t.Fatalf("want 2 entries, got %d", len(acls))
	}
	a := acls[0]
	if a.ID != 11 || a.Bits != 32 || a.StartIP != "192.0.2.10" || !a.Editable || len(a.Resources) != 2 {
		t.Fatalf("first entry: %+v", a)
	}
	if a.Resources[0].ResourceType != "vpc" || a.Resources[0].ResourceID != 901 || a.Resources[0].ResourceLabel != "web" {
		t.Fatalf("vpc resource: %+v", a.Resources[0])
	}
	if a.Resources[1].ResourceID != 587 || a.Resources[1].ResourceLabel != "" {
		t.Fatalf("nke resource: %+v", a.Resources[1])
	}
	if acls[1].ID != 12 || acls[1].Description != "" || len(acls[1].Resources) != 0 {
		t.Fatalf("customer entry: %+v", acls[1])
	}

	sub := 7
	if _, err := c.ListAPIACLs(&sub); err != nil {
		t.Fatalf("ListAPIACLs(sub): %v", err)
	}
	if gotQuery != "7" {
		t.Fatalf("sub_user_id: got %q", gotQuery)
	}
}
