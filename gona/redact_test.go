package gona

import (
	"net/url"
	"strings"
	"testing"
)

// The V2 and V3 APIs both take the API key as a query parameter, so any URL that reaches an
// error message or a log carries a live credential. Terraform prints those errors to the
// operator's terminal, their CI logs, and whatever support ticket follows. This test exists
// so that a future change cannot quietly put the key back.
func TestRedactURLRemovesAPIKey(t *testing.T) {
	const secret = "EXAMPLE-NOT-A-REAL-KEY-0000000000000000000000000000"

	u, err := url.Parse("https://vapi2.netactuate.com/api/dns/zone?key=" + secret)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	got := redactURL(u)
	if strings.Contains(got, secret) {
		t.Fatalf("redactURL leaked the API key: %s", got)
	}
	if !strings.Contains(got, "REDACTED") {
		t.Fatalf("redactURL did not mark the key as redacted: %s", got)
	}
	if !strings.Contains(got, "vapi2.netactuate.com") {
		t.Fatalf("redactURL destroyed the useful part of the URL: %s", got)
	}

	gotV3 := redactV3URL(u)
	if strings.Contains(gotV3, secret) {
		t.Fatalf("redactV3URL leaked the API key: %s", gotV3)
	}

	// A URL with no key must survive untouched, so ordinary diagnostics stay readable.
	plain, _ := url.Parse("https://vapi2.netactuate.com/api/dns/zones?type=NATIVE")
	if redactURL(plain) != plain.String() {
		t.Fatalf("redactURL altered a URL that had no key: %s", redactURL(plain))
	}
}
