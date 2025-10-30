package gona

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
)

func TestGetOSs(t *testing.T) {
	expected := []OS{
		{ID: 1, Os: "Ubuntu 22.04", Type: "linux", Subtype: "ubuntu", Size: "10GB", Bits: "64", Tech: "kvm"},
		{ID: 2, Os: "Debian 12", Type: "linux", Subtype: "debian", Size: "8GB", Bits: "64", Tech: "kvm"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/cloud/images" {
			t.Errorf("Expected path /api/cloud/images, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockAPIResponse(expected))
	}))
	defer server.Close()

	client := NewClientCustom("test-key", server.URL+"/api/")
	ctx := context.Background()

	osList, err := client.GetOSs(ctx)
	if err != nil {
		t.Fatalf("GetOSs() error = %v", err)
	}

	if len(osList) != 2 {
		t.Errorf("GetOSs() returned %d OSs, want 2", len(osList))
	}

	if !slices.Equal(expected, osList) {
		t.Errorf("GetOSs() = %v, want %v", osList, expected)
	}
}

func TestOS_JSON_RoundTrip(t *testing.T) {
	os := OS{
		ID:      1,
		Os:      "Ubuntu 22.04",
		Type:    "linux",
		Subtype: "ubuntu",
		Size:    "10GB",
		Bits:    "64",
		Tech:    "kvm",
	}

	data, err := json.Marshal(os)
	if err != nil {
		t.Fatalf("Failed to marshal OS: %v", err)
	}

	var decoded OS
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal OS: %v", err)
	}

	if decoded.ID != os.ID || decoded.Os != os.Os {
		t.Errorf("Roundtrip failed: got %+v, want %+v", decoded, os)
	}
}
