package gona

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testRoundTripper struct {
	handler http.HandlerFunc
}

func (rt testRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	rt.handler(rec, req)
	return rec.Result(), nil
}

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	c := NewClientCustom("test-key", "https://example.invalid/")
	c.client.Transport = testRoundTripper{handler: handler}
	return c
}

// DeleteServer, and now DeleteBGPSession, must be safe to retry: a session
// already gone (a retry, or one cleared out of band) 422s with this message
// instead of a real error.
func TestDeleteBGPSessionTreatsAlreadyGoneAsSuccess(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"result":"error","message":"validation failed","code":422,"fields":{"id":["The bgp id could not be found"]}}`))
	})

	if err := c.DeleteBGPSession(999); err != nil {
		t.Fatalf("expected an already-gone BGP session to be treated as success, got: %v", err)
	}
}

// A genuine validation error on the same endpoint must still surface --
// the special case only matches the specific "already gone" message.
func TestDeleteBGPSessionReturnsRealErrors(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"result":"error","message":"validation failed","code":422,"fields":{"id":["some other real problem"]}}`))
	})

	if err := c.DeleteBGPSession(999); err == nil {
		t.Fatal("expected a genuine validation error to propagate")
	}
}

func TestDeleteBGPSessionSucceedsOnOK(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200}`))
	})

	if err := c.DeleteBGPSession(999); err != nil {
		t.Fatalf("expected a normal delete to succeed, got: %v", err)
	}
}
