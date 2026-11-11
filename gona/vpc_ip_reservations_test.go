package gona

import (
	"net/http"
	"testing"
)

// GET /vpcs/{id}/ip-reservations wraps the paginated rows in a "reservations" key.
const vpcIPReservationsBody = `{"code":200,"data":{"reservations":{
  "meta": {"total": 3, "offset": 0, "limit": 1000},
  "data": [
    {"network": "10.11.73.1/32", "reservedOn": "2026-11-11T11:40:00.000Z", "mac": null,
     "type": "VPC Gateway", "description": "VPC Gateway",
     "vm": {"name": null, "hostname": null, "mbpkgid": null}},
    {"network": "10.11.73.4/32", "reservedOn": "2026-11-11T11:46:00.000Z", "mac": "52:54:00:aa:bb:01",
     "type": "Primary Interface", "description": "Primary Interface for nke-tcp-2037",
     "vm": {"name": "nke-tcp-2037", "hostname": "nke-tcp-2037", "mbpkgid": 990001}},
    {"network": "10.11.73.5/32", "reservedOn": "2026-11-11T11:47:00.000Z", "mac": "52:54:00:aa:bb:02",
     "type": "Primary Interface", "description": "Primary Interface for nke-tcp-2038",
     "vm": {"name": "nke-tcp-2038", "hostname": "nke-tcp-2038", "mbpkgid": 990002}}
  ]}}}`

func TestListVPCIPReservations(t *testing.T) {
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/vpcs/778/ip-reservations" {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(vpcIPReservationsBody))
	})

	reservations, err := c.ListVPCIPReservations(778)
	if err != nil {
		t.Fatalf("ListVPCIPReservations: %v", err)
	}
	if len(reservations) != 3 {
		t.Fatalf("got %d reservations, want 3", len(reservations))
	}
	if gw := reservations[0]; gw.Type != "VPC Gateway" || gw.VM.MBPkgID != 0 || gw.Network != "10.11.73.1/32" {
		t.Fatalf("gateway row: %+v", gw)
	}
	vm := reservations[2]
	if vm.Network != "10.11.73.5/32" || vm.Type != "Primary Interface" || vm.MAC != "52:54:00:aa:bb:02" {
		t.Fatalf("vm row: %+v", vm)
	}
	if vm.VM.MBPkgID != 990002 || vm.VM.Name != "nke-tcp-2038" || vm.VM.Hostname != "nke-tcp-2038" {
		t.Fatalf("vm row vm fields: %+v", vm.VM)
	}
}
