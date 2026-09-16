package gona

import (
	"encoding/json"
	"testing"
	"time"
)

// A stalled cloud router build keeps a null date on every remaining step forever, and never
// sets readyOn. This test pins the shape that detection depends on: the exact body observed
// from a router that stalled must parse into "5 steps done, stuck on Cloud
// Router configured", which is what the wait loop counts.
func TestRouterBuildStallShapeIsDetectable(t *testing.T) {
	const body = `{
	  "name": "stalled-router",
	  "readyOn": null,
	  "build": [
	    {"text": "Began creation.",                        "date": "2026-09-10T15:15:39.247Z"},
	    {"text": "Hardware provisioned.",                  "date": "2026-09-10T15:15:44.973Z"},
	    {"text": "Cloud Router software updated.",         "date": "2026-09-10T15:17:31.425Z"},
	    {"text": "API enabled.",                           "date": "2026-09-10T15:17:41.390Z"},
	    {"text": "Cloud Router connectivity established.", "date": "2026-09-10T15:18:02.164Z"},
	    {"text": "Cloud Router configured.",               "date": null},
	    {"text": "Cloud Router ready.",                    "date": null}
	  ]
	}`

	var r Router
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if r.ReadyOn != nil {
		t.Fatal("readyOn should be nil on a stalled router")
	}

	done, pending := 0, ""
	for _, e := range r.Build {
		if !e.Date.IsZero() {
			done++
		} else if pending == "" {
			pending = e.Text
		}
	}

	// A null date MUST come back as the zero time, or the wait loop counts a stalled step as
	// completed and never reports the stall.
	if done != 5 {
		t.Fatalf("completed steps: got %d, want 5", done)
	}
	if pending != "Cloud Router configured." {
		t.Fatalf("stuck step: got %q, want %q", pending, "Cloud Router configured.")
	}

	// And a healthy build must not look stalled.
	const healthy = `{"readyOn":"2026-09-10T15:20:00.000Z","build":[
	  {"text":"Began creation.","date":"2026-09-10T15:15:39.247Z"}]}`
	var h Router
	if err := json.Unmarshal([]byte(healthy), &h); err != nil {
		t.Fatalf("unmarshal healthy: %v", err)
	}
	if h.ReadyOn == nil {
		t.Fatal("a ready router must report readyOn")
	}
	if h.ReadyOn.Before(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("readyOn parsed wrong: %v", h.ReadyOn)
	}
}
