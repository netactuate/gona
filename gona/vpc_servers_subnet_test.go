package gona

import (
	"encoding/json"
	"testing"
)

// LD-242: vapi3 sends subnetId as a string for attached rows, a number for native rows and null on
// some POST answers (LD-216). Every shape must decode.
func TestVPCServerSubnetIDAcceptsStringNumberAndNull(t *testing.T) {
	body := `[
		{"serverId":1,"mbpkgid":11,"fqdn":"a","membership":"attached","subnetId":"9","nicId":2,"state":"ready","address":"10.0.0.4"},
		{"serverId":2,"mbpkgid":12,"membership":"native","subnetId":0},
		{"serverId":3,"mbpkgid":13,"membership":"attached","subnetId":null},
		{"serverId":4,"mbpkgid":14,"membership":"attached","subnetId":""},
		{"serverId":5,"mbpkgid":15,"membership":"attached"}
	]`
	var out []VPCServer
	if err := json.Unmarshal([]byte(body), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	want := []int{9, 0, 0, 0, 0}
	for i, w := range want {
		if out[i].SubnetID != w {
			t.Fatalf("row %d subnetId %d, want %d", i, out[i].SubnetID, w)
		}
	}
	if out[0].MBPkgID != 11 || out[0].FQDN != "a" || out[0].NICID != 2 || out[0].Address != "10.0.0.4" || out[0].State != "ready" {
		t.Fatalf("other fields lost: %+v", out[0])
	}
	var bad VPCServer
	if err := json.Unmarshal([]byte(`{"subnetId":"nine"}`), &bad); err == nil {
		t.Fatal("non-numeric subnetId decoded without an error")
	}
}
