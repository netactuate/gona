package gona

import (
	"net/http"
	"testing"
)

// Minimal API response with synthetic identifiers. contract_id is a JSON
// number; state (power) and status (server) are separate fields.
const getServerResponseFixture = `{
	"result": "success",
	"message": "",
	"code": 200,
	"data": {
		"contract_id": 123,
		"mbpkgid": 456,
		"state": "UP",
		"fqdn": "server.example.test",
		"package": "VR1x1x25",
		"location_id": 789,
		"status": "RUNNING"
	}
}`

func TestGetServerResponseMapping(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(getServerResponseFixture))
	})

	server, err := c.GetServer(456)
	if err != nil {
		t.Fatalf("GetServer: %v", err)
	}
	if server.PackageBillingContractId != 123 {
		t.Fatalf("expected PackageBillingContractId 123 from contract_id, got %d", server.PackageBillingContractId)
	}
	if server.ID != 456 || server.Name != "server.example.test" || server.Package != "VR1x1x25" || server.LocationID != 789 {
		t.Fatalf("unexpected Server fields, possible unrelated regression: %+v", server)
	}
	if server.ServerStatus != "RUNNING" || server.PowerStatus != "UP" {
		t.Fatalf("expected ServerStatus=RUNNING PowerStatus=UP, got ServerStatus=%q PowerStatus=%q", server.ServerStatus, server.PowerStatus)
	}
}
