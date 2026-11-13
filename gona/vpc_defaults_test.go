package gona

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// Create sends defaults.addEgressIpToApiAcl, and GET /vpcs/{id} returns it under defaults.
func TestVPCAddEgressIPToAPIACL(t *testing.T) {
	var created map[string]interface{}
	c := newLongtailTestV3Client(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/vpcs":
			body, _ := io.ReadAll(r.Body)
			json.Unmarshal(body, &created)
			w.Write([]byte(`{"code":200,"data":{"vpcId":901}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/vpcs/901":
			w.Write([]byte(`{"code":200,"data":{"vpcId":901,"metadata":{"label":"x"},"defaults":{"addEgressIpToApiAcl":true}}}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	if _, err := c.CreateVPC(&CreateVPCRequest{Label: "x", LocationID: 1, Defaults: &VPCDefaults{AddEgressIPToAPIACL: BoolPtr(true)}}); err != nil {
		t.Fatalf("CreateVPC: %v", err)
	}
	defaults, _ := created["defaults"].(map[string]interface{})
	if defaults["addEgressIpToApiAcl"] != true {
		t.Fatalf("create body defaults: %v", created["defaults"])
	}

	vpc, err := c.GetVPC(901)
	if err != nil {
		t.Fatalf("GetVPC: %v", err)
	}
	if vpc.Defaults == nil || vpc.Defaults.AddEgressIPToAPIACL == nil || !*vpc.Defaults.AddEgressIPToAPIACL {
		t.Fatalf("defaults read back: %+v", vpc.Defaults)
	}
}
