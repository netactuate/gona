package gona

import "testing"

func TestCloudPool_Name(t *testing.T) {
	tests := map[string]struct {
		pool CloudPool
		want string
	}{
		"General Compute": {CloudPoolGeneralCompute, "General Compute"},
		"AMD EPYC":        {CloudPoolAMDEPYC, "AMD EPYC"},
		"Default":         {CloudPoolDefault, "Default"},
		"Unknown":         {CloudPool(99), "Unknown"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.pool.Name(); got != tt.want {
				t.Errorf("Name() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCloudPoolFromName(t *testing.T) {
	tests := map[string]struct {
		name string
		want CloudPool
	}{
		"General Compute": {"General Compute", CloudPoolGeneralCompute},
		"AMD EPYC":        {"AMD EPYC", CloudPoolAMDEPYC},
		"Default":         {"Default", CloudPoolDefault},
		"invalid":         {"invalid", CloudPoolDefault},
		"empty":           {"", CloudPoolDefault},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := CloudPoolFromName(tt.name); got != tt.want {
				t.Errorf("CloudPoolFromName() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCloudPool_RoundTrip(t *testing.T) {
	pools := []CloudPool{CloudPoolGeneralCompute, CloudPoolAMDEPYC, CloudPoolDefault}
	for _, pool := range pools {
		if got := CloudPoolFromName(pool.Name()); got != pool {
			t.Errorf("Roundtrip failed: %d -> %q -> %d", pool, pool.Name(), got)
		}
	}
}

func TestInvalidPoolString(t *testing.T) {
	// Test that invalid pool values have sensible Name() output
	invalidPool := CloudPool(42)
	if invalidPool.Name() != "Unknown" {
		t.Errorf("CloudPool(42).Name() = %s, want Unknown", invalidPool.Name())
	}

	// Test CloudPoolFromName with invalid input
	result := CloudPoolFromName("invalid pool name")
	if result != CloudPoolDefault {
		t.Errorf("CloudPoolFromName(\"invalid pool name\") = %d, want %d", result, CloudPoolDefault)
	}
}
