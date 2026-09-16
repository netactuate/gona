package gona

import (
	"net/http"
	"testing"
)

func TestCreateBGPGroup(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/bgp/bgpgroup" {
			t.Fatalf("expected /bgp/bgpgroup, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":{"id":42,"name":"edge","description":"edge group","group_type":"anycast"}}`))
	})

	group, err := c.CreateBGPGroup(&CreateBGPGroupRequest{Name: "edge", Description: "edge group", GroupType: "anycast"})
	if err != nil {
		t.Fatalf("expected BGP group create to succeed, got: %v", err)
	}
	if group.ID != 42 || group.Name != "edge" || group.GroupType != "anycast" {
		t.Fatalf("unexpected group: %#v", group)
	}
}

func TestBuyBGPPrefixesContractError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/bgp/bgpbuyprefixes" {
			t.Fatalf("expected /bgp/bgpbuyprefixes, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusPreconditionFailed)
		w.Write([]byte(`{"result":"error","message":"agreement is not enabled for this account","data":[],"code":412}`))
	})

	_, err := c.BuyBGPPrefixes(&BuyBGPPrefixesRequest{Name: "edge", GroupID: 42, AgreementID: 7})
	if err == nil {
		t.Fatal("expected contract error")
	}
	if !IsContractError(err) {
		t.Fatalf("expected IsContractError to match, got: %v", err)
	}
}

func TestListBGPASNs(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bgp/bgpasns" {
			t.Fatalf("expected /bgp/bgpasns, got %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("group_type"); got != "anycast" {
			t.Fatalf("expected group_type anycast, got %q", got)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":[{"id":9,"asn":64500,"name":"primary","group_type":"anycast"}]}`))
	})

	asns, err := c.ListBGPASNs("anycast")
	if err != nil {
		t.Fatalf("expected BGP ASN list to succeed, got: %v", err)
	}
	if len(asns) != 1 || asns[0].ID != 9 || asns[0].ASN != 64500 {
		t.Fatalf("unexpected ASNs: %#v", asns)
	}
}

func TestListAccountAgreements(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/account/agreements" {
			t.Fatalf("expected /account/agreements, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success","code":200,"data":[{"id":7,"name":"anycast","title":"Anycast Terms","version":"1"}]}`))
	})

	agreements, err := c.ListAccountAgreements()
	if err != nil {
		t.Fatalf("expected account agreements list to succeed, got: %v", err)
	}
	if len(agreements) != 1 || agreements[0].ID != 7 || agreements[0].Title != "Anycast Terms" {
		t.Fatalf("unexpected agreements: %#v", agreements)
	}
}
