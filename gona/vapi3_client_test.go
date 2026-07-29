package gona

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsV3NotFound(t *testing.T) {
	notFound := &V3NotFoundError{
		StatusCode: 404,
		Body:       `{"code":404,"message":"VPC does not exist."}`,
	}

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "bare sentinel",
			err:  notFound,
			want: true,
		},
		{
			name: "wrapped once, as the resource helpers wrap",
			err:  fmt.Errorf("list floating IPs for VPC %d: %w", 318, notFound),
			want: true,
		},
		{
			name: "wrapped twice",
			err: fmt.Errorf("read location: %w",
				fmt.Errorf("list floating IPs for VPC %d: %w", 318, notFound)),
			want: true,
		},
		{
			name: "unrelated error",
			err:  errors.New("API error on GET /vpcs/318: HTTP 500"),
			want: false,
		},
		{
			name: "wrapped unrelated error",
			err:  fmt.Errorf("list floating IPs for VPC %d: %w", 318, errors.New("boom")),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsV3NotFound(tt.err); got != tt.want {
				t.Errorf("IsV3NotFound(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

// A deleted parent VPC reaches the caller through ListVPCFloatingIPs, which
// wraps. Callers key state removal off IsV3NotFound, so the wrap must not hide
// the 404: if it does, refresh hard-errors and the resource can only be
// recovered by hand.
func TestListVPCFloatingIPsReportsMissingVPCAsNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"code":404,"message":"VPC does not exist.","bad_fields":{"bad_field":"vpcId"}}`)
	}))
	defer server.Close()

	_, err := NewV3Client("test-key", server.URL).ListVPCFloatingIPs(318)
	if err == nil {
		t.Fatal("ListVPCFloatingIPs(318) = nil error, want a not-found error")
	}
	if !IsV3NotFound(err) {
		t.Errorf("IsV3NotFound(%v) = false, want true", err)
	}
}
