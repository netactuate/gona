package gona

import (
	"encoding/json"
	"testing"
)

// The exact body shape from GET /nke/clusters/{id}/addons/storage. If this does not unmarshal into the fields the
// provider reads, no amount of provider-side fiddling will help.
func TestNKEAddonStorageConfigUnmarshals(t *testing.T) {
	const body = `{
	  "id": 63, "addonId": 63, "clusterId": 397, "addonType": "storage",
	  "version": "v0.3.0", "channel": "stable", "displayName": "Storage",
	  "state": "degraded", "updateAvailable": false,
	  "config": {"integrations": [{"storageIntegrationId": 33, "blockNamespaceId": 50291,
	    "storageClassName": "netactuate-block", "volumeSnapshotClassName": null,
	    "isDefaultClass": true, "reclaimPolicy": "Delete"}]}
	}`

	var a NKEAddon
	if err := json.Unmarshal([]byte(body), &a); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(a.Config.Integrations) != 1 {
		t.Fatalf("expected 1 integration, got %d", len(a.Config.Integrations))
	}
	in := a.Config.Integrations[0]
	if in.BlockNamespaceID != 50291 {
		t.Fatalf("BlockNamespaceID: got %d, want 50291", in.BlockNamespaceID)
	}
	if in.StorageIntegrationID != 33 {
		t.Fatalf("StorageIntegrationID: got %d, want 33", in.StorageIntegrationID)
	}
	if !in.IsDefaultClass {
		t.Fatal("IsDefaultClass: got false, want true (makeDefault goes out, isDefaultClass comes back)")
	}
	if in.StorageClassName != "netactuate-block" {
		t.Fatalf("StorageClassName: got %q", in.StorageClassName)
	}
}
