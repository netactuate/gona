package gona

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateVPCIPReservationContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "POST", "/vpcs/1/ip-reservations", `{"description":"vip","address":"10.0.0.5"}`)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{"ipReservationId":2}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.CreateVPCIPReservation(1, &CreateVPCIPReservationRequest{Description: "vip", Address: "10.0.0.5"})
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != 2 {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}
func TestDeleteVPCIPReservationContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "DELETE", "/vpcs/1/ip-reservations/2", ``)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			err := c.DeleteVPCIPReservation(1, 2)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestListVPCIPReservationsByTypeContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "GET", "/vpcs/1/ip-reservations", ``)
				if r.URL.Query().Get("type") != "lbVip" {
					t.Error("missing type query")
				}
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{"reservations":[{"ipReservationId":2,"address":"10.0.0.5","description":"vip","reservedOn":"today","inUse":true}]}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.ListVPCIPReservationsByType(1, "lbVip")
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 || got[0].IPReservationID != 2 || !got[0].InUse || got[0].Address != "10.0.0.5" || got[0].ReservedOn != "today" {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}
func TestGetNKEAPIAccessContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "GET", "/nke/clusters/1/api-access", ``)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{"publicEnabled":true,"privateEndpointEnabled":true,"privateEndpointAddress":"10.0.0.2","privateDnsName":"api.internal","privateEndpointState":"ready","allowListActive":true,"state":"applied","appliedOn":"today","failureReason":"reason"}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.GetNKEAPIAccess(1)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !got.PublicEnabled || !got.PrivateEndpointEnabled || !got.AllowListActive || got.PrivateEndpointAddress != "10.0.0.2" || got.PrivateDNSName != "api.internal" || got.PrivateEndpointState != "ready" || got.State != "applied" || got.AppliedOn != "today" || got.FailureReason != "reason" {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}
func TestUpdateNKEAPIAccessContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "PATCH", "/nke/clusters/1/api-access", `{"publicEnabled":false,"privateEndpointEnabled":true}`)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{"publicEnabled":false,"privateEndpointEnabled":true,"state":"applying"}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.UpdateNKEAPIAccess(1, &UpdateNKEAPIAccessRequest{PublicEnabled: BoolPtr(false), PrivateEndpointEnabled: BoolPtr(true)})
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got.PublicEnabled || !got.PrivateEndpointEnabled || got.State != "applying" {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}
func TestListNKEAPIAllowListContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "GET", "/nke/clusters/1/api-access/allow-list", ``)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":[{"allowListEntryId":"platform-1","network":"192.0.2.0/24","description":"platform","isPlatform":true,"owner":"system","createdOn":"today"}]}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.ListNKEAPIAllowList(1)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 || got[0].AllowListEntryID != "platform-1" || !got[0].IsPlatform || got[0].Owner != "system" || got[0].Network != "192.0.2.0/24" || got[0].CreatedOn != "today" {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}
func TestCreateNKEAPIAllowEntryContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "POST", "/nke/clusters/1/api-access/allow-list", `{"network":"192.0.2.0/24","description":"office"}`)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{"allowListEntryId":"customer-2","network":"192.0.2.0/24","description":"office","isPlatform":false}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.CreateNKEAPIAllowEntry(1, &CreateNKEAPIAllowEntryRequest{Network: "192.0.2.0/24", Description: "office"})
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got.AllowListEntryID != "customer-2" || got.Description != "office" || got.Network != "192.0.2.0/24" || got.IsPlatform {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}
func TestDeleteNKEAPIAllowEntryContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "DELETE", "/nke/clusters/1/api-access/allow-list/customer-2", ``)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			err := c.DeleteNKEAPIAllowEntry(1, "customer-2")
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestListVPCAttachmentsContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "GET", "/vpcs/1/attachments", ``)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":[{"id":2,"type":"customer-router","state":"ready","routerId":3,"vrfId":4}]}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.ListVPCAttachments(1)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 || got[0].ID != 2 || got[0].RouterID != 3 || got[0].VRFID != 4 || got[0].State != "ready" {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}
func TestGetVPCAttachmentContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "GET", "/vpcs/1/attachments/2", ``)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{"id":2,"type":"customer-router","address":"10.0.0.2","description":"router","routerId":3,"cloudRouterAttachmentId":5,"vrfId":4,"routerState":"ready","state":"ready","stateMessage":"ok","readyOn":"today","createdOn":"yesterday","deleteRequestedOn":"tomorrow","accept":["192.0.2.0/24"],"advertise":[]}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.GetVPCAttachment(1, 2)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got.ID != 2 || got.Address != "10.0.0.2" || got.CloudRouterAttachmentID != 5 || got.RouterState != "ready" || got.StateMessage != "ok" || got.ReadyOn != "today" || got.CreatedOn != "yesterday" || got.DeleteRequestedOn != "tomorrow" || len(got.Accept) != 1 {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}
func TestCreateVPCAttachmentContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "POST", "/vpcs/1/attachments", `{"type":"cloud-router","accept":["192.0.2.0/24"],"advertise":["10.0.0.0/24"],"address":"10.0.0.2","routerId":3,"vrfId":4,"description":"router"}`)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{"attachmentId":2}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.CreateVPCAttachment(1, &CreateVPCAttachmentRequest{Type: "cloud-router", Accept: []string{"192.0.2.0/24"}, Advertise: []string{"10.0.0.0/24"}, Address: "10.0.0.2", RouterID: 3, VRFID: 4, Description: "router"})
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != 2 {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}
func TestDeleteVPCAttachmentContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "DELETE", "/vpcs/1/attachments/2", ``)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			err := c.DeleteVPCAttachment(1, 2)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestListVPCAttachmentPrefixesContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "GET", "/vpcs/1/attachments/2/prefixes", ``)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":[{"id":7,"direction":"accept","network":"192.0.2.0/24","ipVersion":4,"createdOn":"today","deleteRequestedOn":"tomorrow"}]}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.ListVPCAttachmentPrefixes(1, 2)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 || got[0].ID != 7 || got[0].Direction != "accept" || got[0].IPVersion != 4 || got[0].CreatedOn != "today" || got[0].DeleteRequestedOn != "tomorrow" {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}
func TestCreateVPCAttachmentPrefixContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "POST", "/vpcs/1/attachments/2/prefixes", `{"network":"192.0.2.0/24","direction":"accept"}`)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{"prefixId":7}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.CreateVPCAttachmentPrefix(1, 2, &CreateVPCAttachmentPrefixRequest{Network: "192.0.2.0/24", Direction: "accept"})
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != 7 {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}
func TestDeleteVPCAttachmentPrefixContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "DELETE", "/vpcs/1/attachments/2/prefixes/7", ``)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			err := c.DeleteVPCAttachmentPrefix(1, 2, 7)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
func TestGetNKEPrivateNetworkContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "GET", "/nke/clusters/1/private-network", ``)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{"clusterId":1,"vpcId":2,"state":"roll_required","nodes":[{"id":3,"name":"node","privateNicState":"roll_required","privateIpAddress":"10.0.0.3"}],"nodesRequiringRoll":[3]}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.GetNKEPrivateNetwork(1)
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got.ClusterID != 1 || got.VPCID != 2 || got.State != "roll_required" || len(got.Nodes) != 1 || got.Nodes[0].PrivateNICState != "roll_required" || got.Nodes[0].PrivateIPAddress != "10.0.0.3" || len(got.NodesRequiringRoll) != 1 {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}
func TestAddNKEPrivateNetworkContract(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(fmt.Sprint(missing), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				checkNKEDrop1Request(t, r, "POST", "/nke/clusters/1/private-network", `{"vpcId":2}`)
				if missing {
					w.WriteHeader(404)
					fmt.Fprint(w, `{"code":404,"data":{}}`)
					return
				}
				fmt.Fprint(w, `{"code":200,"data":{"clusterId":1}}`)
			}))
			defer server.Close()
			c := NewV3Client("test-key", server.URL)
			got, err := c.AddNKEPrivateNetwork(1, &AddNKEPrivateNetworkRequest{VPCID: 2})
			if missing {
				if !IsV3NotFound(err) {
					t.Fatalf("404: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != 1 {
				t.Fatalf("decoded: %#v", got)
			}
		})
	}
}

func TestDrop2ListsPagination(t *testing.T) {
	for _, kind := range []string{"reservations", "allow-list", "attachments", "prefixes"} {
		t.Run(kind, func(t *testing.T) {
			calls := 0
			s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.Method != "GET" {
					t.Error(r.Method)
				}
				offset := r.URL.Query().Get("offset")
				row := `{"id":2,"ipReservationId":2,"allowListEntryId":"customer-2"}`
				n := 0
				if offset != "" {
					if offset != "1" {
						t.Error(offset)
					}
					n = 1
					row = `{"id":3,"ipReservationId":3,"allowListEntryId":"customer-3"}`
				}
				if kind == "reservations" && r.URL.Query().Get("type") != "customerReserved" {
					t.Error("type lost across pages")
				}
				body := fmt.Sprintf(`{"meta":{"total":2,"limit":1,"offset":%d},"data":[%s]}`, n, row)
				if kind == "reservations" {
					body = `{"reservations":` + body + `}`
				}
				fmt.Fprintf(w, `{"data":%s}`, body)
			}))
			defer s.Close()
			c := NewV3Client("test-key", s.URL)
			switch kind {
			case "reservations":
				rows, err := c.ListVPCIPReservationsByType(1, "customerReserved")
				if err != nil || len(rows) != 2 || rows[1].IPReservationID != 3 {
					t.Fatalf("%v %v", rows, err)
				}
			case "allow-list":
				rows, err := c.ListNKEAPIAllowList(1)
				if err != nil || len(rows) != 2 || rows[1].AllowListEntryID != "customer-3" {
					t.Fatalf("%v %v", rows, err)
				}
			case "attachments":
				rows, err := c.ListVPCAttachments(1)
				if err != nil || len(rows) != 2 || rows[1].ID != 3 {
					t.Fatalf("%v %v", rows, err)
				}
			case "prefixes":
				rows, err := c.ListVPCAttachmentPrefixes(1, 2)
				if err != nil || len(rows) != 2 || rows[1].ID != 3 {
					t.Fatalf("%v %v", rows, err)
				}
			}
			if calls != 2 {
				t.Fatalf("calls %d", calls)
			}
		})
	}
}
func TestVPCQ31CountsJSON(t *testing.T) {
	var v VPC
	if err := json.Unmarshal([]byte(`{"vpcId":1,"metadata":{},"maxFloatingIpv4s":20,"usedFloatingIpv4s":3}`), &v); err != nil {
		t.Fatal(err)
	}
	if v.MaxFloatingIPv4s != 20 || v.UsedFloatingIPv4s != 3 {
		t.Fatal(v)
	}
	raw, err := json.Marshal(VPC{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), "FloatingIpv4s") {
		t.Fatalf("zero counts must be omitted: %s", raw)
	}
}
