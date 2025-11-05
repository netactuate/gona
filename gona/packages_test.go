package gona

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestGetPackages(t *testing.T) {
	expectedPackages := []Package{
		{ID: 123, Status: "active", Locked: "0", PlanName: "Standard", Installed: 1},
		{ID: 456, Status: "active", Locked: "0", PlanName: "Premium", Installed: 1},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload := slices.Clone(expectedPackages)
		slices.Reverse(payload)
		if r.URL.Path != "/api/cloud/packages" {
			t.Errorf("Expected path /api/cloud/packages, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockAPIResponse(payload))
	}))
	defer server.Close()

	client := NewClientCustom("test-key", server.URL+"/api/")
	ctx := context.Background()

	packages, err := client.GetPackages(ctx)
	if err != nil {
		t.Fatalf("GetPackages() error = %v", err)
	}

	if len(packages) != 2 {
		t.Errorf("GetPackages() returned %d packages, want 2", len(packages))
	}

	slices.SortStableFunc(packages, func(a, b Package) int {
		return cmp.Compare(a.ID, b.ID)
	})
	if !slices.Equal(expectedPackages, packages) {
		t.Errorf("GetPackages() = %v, want %v", packages, expectedPackages)
	}
}

func TestGetPackage(t *testing.T) {
	const PackageID = 123
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != fmt.Sprint("/api/cloud/package/", PackageID) {
			t.Errorf("Expected path /api/cloud/package/%d, got %s", PackageID, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockAPIResponse(Package{
			ID: PackageID, Status: "active", Locked: "0", PlanName: "Standard", Installed: 1,
		}))
	}))
	defer server.Close()

	client := NewClientCustom("test-key", server.URL+"/api/")
	ctx := context.Background()

	pkg, err := client.GetPackage(ctx, PackageID)
	if err != nil {
		t.Fatalf("GetPackage() error = %v", err)
	}

	if pkg.ID != PackageID {
		t.Errorf("Package.ID = %d, want %d", pkg.ID, PackageID)
	}
	if pkg.PlanName != "Standard" {
		t.Errorf("Package.PlanName = %q, want %q", pkg.PlanName, "Standard")
	}
}

func TestPackage_JSON_RoundTrip(t *testing.T) {
	pkg := Package{
		ID:        123,
		Status:    "active",
		Locked:    "0",
		PlanName:  "Standard",
		Installed: 1,
	}

	data, err := json.Marshal(pkg)
	if err != nil {
		t.Fatalf("Failed to marshal Package: %v", err)
	}

	var decoded Package
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal Package: %v", err)
	}

	if decoded.ID != pkg.ID || decoded.PlanName != pkg.PlanName {
		t.Errorf("Roundtrip failed: got %+v, want %+v", decoded, pkg)
	}
}
