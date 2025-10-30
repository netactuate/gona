package gona

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestGetPlans(t *testing.T) {
	expectedPlans := []Plan{
		{ID: 1, Name: "Small", RAM: "2GB", Disk: "40GB", Transfer: "2TB", Price: "10.00", Available: "1"},
		{ID: 2, Name: "Medium", RAM: "4GB", Disk: "80GB", Transfer: "4TB", Price: "20.00", Available: "1"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cloud/sizes" {
			t.Errorf("Expected path /api/cloud/sizes, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockAPIResponse(expectedPlans))
	}))
	defer server.Close()

	client := NewClientCustom("test-key", server.URL+"/api/")
	ctx := context.Background()

	plans, err := client.GetPlans(ctx)
	if err != nil {
		t.Fatalf("GetPlans() error = %v", err)
	}

	if len(plans) != len(expectedPlans) {
		t.Errorf("GetPlans() returned %d plans, want %d", len(plans), len(expectedPlans))
	}

	if !slices.Equal(expectedPlans, plans) {
		t.Errorf("GetPlans() = %v, want %v", plans, expectedPlans)
	}
}

func TestPlan_JSONRoundTrip(t *testing.T) {
	plan := Plan{
		ID:        1,
		Name:      "Small",
		RAM:       "2GB",
		Disk:      "40GB",
		Transfer:  "2TB",
		Price:     "10.00",
		Available: "1",
	}

	data, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("Failed to marshal Plan: %v", err)
	}

	var decoded Plan
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal Plan: %v", err)
	}

	if decoded.ID != plan.ID || decoded.Name != plan.Name {
		t.Errorf("Roundtrip failed: got %+v, want %+v", decoded, plan)
	}
}
