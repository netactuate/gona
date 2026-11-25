package gona

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func drop3Contract(t *testing.T, method, path, body, response string, call func(*V3Client)) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method || r.URL.Path != path {
			t.Errorf("request: %s %s", r.Method, r.URL.Path)
		}
		raw, e := io.ReadAll(r.Body)
		if e != nil {
			t.Error(e)
		}
		if body != "" {
			var got, want interface{}
			if e = json.Unmarshal(raw, &got); e != nil {
				t.Error(e)
			}
			if e = json.Unmarshal([]byte(body), &want); e != nil {
				t.Fatal(e)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("body keys or values differ")
			}
		} else if len(raw) > 0 {
			t.Error("unexpected body")
		}
		fmt.Fprint(w, `{"code":200,"data":`+response+`}`)
	}))
	defer server.Close()
	call(NewV3Client("test-key", server.URL))
}
func TestListSecurityGroupsContract(t *testing.T) {
	drop3Contract(t, "GET", "/cloud/networking/security-groups", ``, `{"data":[{"securityGroupId":2,"name":"web","description":"tier","isDefault":true,"defaultForVpcId":1,"memberCount":3,"rules":[{"ruleId":9}],"references":[{"securityGroupId":8,"ruleId":7}]}]}`, func(c *V3Client) {
		got, err := c.ListSecurityGroups()
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].SecurityGroupID != 2 || got[0].MemberCount != 3 || !got[0].IsDefault {
			t.Fatalf("decoded: %#v", got)
		}
	})
}
func TestListVPCSecurityGroupsContract(t *testing.T) {
	drop3Contract(t, "GET", "/vpcs/1/security-groups", ``, `{"data":[{"securityGroupId":2,"name":"web","description":"tier","isDefault":true,"defaultForVpcId":1,"memberCount":3,"rules":[{"ruleId":9}],"references":[{"securityGroupId":8,"ruleId":7}]}]}`, func(c *V3Client) {
		got, err := c.ListVPCSecurityGroups(1)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].SecurityGroupID != 2 || got[0].MemberCount != 3 || !got[0].IsDefault {
			t.Fatalf("decoded: %#v", got)
		}
	})
}
func TestGetSecurityGroupContract(t *testing.T) {
	drop3Contract(t, "GET", "/cloud/networking/security-groups/2", ``, `{"securityGroupId":2,"name":"web","description":"tier","isDefault":true,"defaultForVpcId":1,"memberCount":3,"rules":[{"ruleId":9}],"references":[{"securityGroupId":8,"ruleId":7}]}`, func(c *V3Client) {
		got, err := c.GetSecurityGroup(2)
		if err != nil {
			t.Fatal(err)
		}
		if got.Name != "web" || got.Rules[0].RuleID != 9 || got.References[0].RuleID != 7 {
			t.Fatalf("decoded: %#v", got)
		}
	})
}
func TestCreateSecurityGroupContract(t *testing.T) {
	drop3Contract(t, "POST", "/cloud/networking/security-groups", `{"name":"web","description":"tier"}`, `{"securityGroupId":2}`, func(c *V3Client) {
		got, err := c.CreateSecurityGroup(&CreateSecurityGroupRequest{Name: "web", Description: "tier"})
		if err != nil {
			t.Fatal(err)
		}
		if got != 2 {
			t.Fatal(got)
		}
	})
}
func TestUpdateSecurityGroupContract(t *testing.T) {
	drop3Contract(t, "PATCH", "/cloud/networking/security-groups/2", `{"name":"web","description":""}`, `{}`, func(c *V3Client) {
		name, desc := "web", ""
		err := c.UpdateSecurityGroup(2, &UpdateSecurityGroupRequest{Name: &name, Description: &desc})
		if err != nil {
			t.Fatal(err)
		}
	})
}
func TestDeleteSecurityGroupContract(t *testing.T) {
	drop3Contract(t, "DELETE", "/cloud/networking/security-groups/2", ``, `{}`, func(c *V3Client) {
		err := c.DeleteSecurityGroup(2)
		if err != nil {
			t.Fatal(err)
		}
	})
}
func TestListSecurityGroupRulesContract(t *testing.T) {
	drop3Contract(t, "GET", "/cloud/networking/security-groups/2/rules", ``, `{"data":[{"ruleId":3,"ipVersion":4,"direction":"inbound","protocol":"tcp","portStart":80,"portEnd":443,"network":"10.0.0.0/24","securityGroupId":4,"description":"web"}]}`, func(c *V3Client) {
		got, err := c.ListSecurityGroupRules(2)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].RuleID != 3 || got[0].PortEnd != 443 || got[0].Network != "10.0.0.0/24" {
			t.Fatalf("decoded: %#v", got)
		}
	})
}
func TestCreateSecurityGroupRuleContract(t *testing.T) {
	drop3Contract(t, "POST", "/cloud/networking/security-groups/2/rules", `{"ipVersion":4,"direction":"inbound","protocol":"tcp","port":{"start":80,"end":443},"network":"10.0.0.0/24","description":"web"}`, `{"ruleId":3}`, func(c *V3Client) {
		got, err := c.CreateSecurityGroupRule(2, &CreateSecurityGroupRuleRequest{IPVersion: 4, Direction: "inbound", Protocol: "tcp", Port: &SecurityGroupPort{Start: 80, End: 443}, Network: "10.0.0.0/24", Description: "web"})
		if err != nil {
			t.Fatal(err)
		}
		if got != 3 {
			t.Fatal(got)
		}
	})
}
func TestDeleteSecurityGroupRuleContract(t *testing.T) {
	drop3Contract(t, "DELETE", "/cloud/networking/security-groups/2/rules/3", ``, `{}`, func(c *V3Client) {
		err := c.DeleteSecurityGroupRule(2, 3)
		if err != nil {
			t.Fatal(err)
		}
	})
}
func TestListSecurityGroupMembersContract(t *testing.T) {
	drop3Contract(t, "GET", "/cloud/networking/security-groups/2/members", ``, `{"data":[{"memberId":3,"vpcId":1,"mbpkgid":8,"nodePoolId":9,"enforcementState":"enforced","createdOn":"today"}]}`, func(c *V3Client) {
		got, err := c.ListSecurityGroupMembers(2)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].MemberID != 3 || got[0].EnforcementState != "enforced" || got[0].NodePoolID != 9 {
			t.Fatalf("decoded: %#v", got)
		}
	})
}
func TestCreateSecurityGroupMemberContract(t *testing.T) {
	drop3Contract(t, "POST", "/cloud/networking/security-groups/2/members", `{"vpcId":1,"mbpkgid":8}`, `{"memberId":3}`, func(c *V3Client) {
		got, err := c.CreateSecurityGroupMember(2, &CreateSecurityGroupMemberRequest{VPCID: 1, MBPkgID: 8})
		if err != nil {
			t.Fatal(err)
		}
		if got != 3 {
			t.Fatal(got)
		}
	})
}
func TestDeleteSecurityGroupMemberContract(t *testing.T) {
	drop3Contract(t, "DELETE", "/cloud/networking/security-groups/2/members/3", ``, `{}`, func(c *V3Client) {
		err := c.DeleteSecurityGroupMember(2, 3)
		if err != nil {
			t.Fatal(err)
		}
	})
}
func TestGetVPCSecurityGroupEnforcementContract(t *testing.T) {
	drop3Contract(t, "GET", "/vpcs/1/security-groups/enforcement", ``, `{"vpcId":1,"enabled":false,"enabledOn":"today"}`, func(c *V3Client) {
		got, err := c.GetVPCSecurityGroupEnforcement(1)
		if err != nil {
			t.Fatal(err)
		}
		if got.VPCID != 1 || got.Enabled || got.EnabledOn != "today" {
			t.Fatal(got)
		}
	})
}
func TestSetVPCSecurityGroupEnforcementContract(t *testing.T) {
	drop3Contract(t, "PUT", "/vpcs/1/security-groups/enforcement", `{"enabled":false}`, `{"vpcId":1,"enabled":false,"enabledOn":"today"}`, func(c *V3Client) {
		got, err := c.SetVPCSecurityGroupEnforcement(1, &SetVPCSecurityGroupEnforcementRequest{Enabled: false})
		if err != nil {
			t.Fatal(err)
		}
		if got.VPCID != 1 || got.Enabled || got.EnabledOn != "today" {
			t.Fatal(got)
		}
	})
}
func TestListVPCSubnetsContract(t *testing.T) {
	drop3Contract(t, "GET", "/vpcs/1/subnets", ``, `{"data":[{"subnetId":2,"network":"10.1.0.0/24","prefixLength":24,"label":"app","isDefault":false,"state":"ready","readyOn":"today","deleteRequestedOn":"later","deletedOn":"last"}]}`, func(c *V3Client) {
		got, err := c.ListVPCSubnets(1)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 {
			t.Fatal(got)
		}
		r := got[0]
		if r.SubnetID != 2 || r.Network != "10.1.0.0/24" || r.DeletedOn != "last" {
			t.Fatal(got)
		}
	})
}
func TestGetVPCSubnetContract(t *testing.T) {
	drop3Contract(t, "GET", "/vpcs/1/subnets/2", ``, `{"subnetId":2,"network":"10.1.0.0/24","prefixLength":24,"label":"app","isDefault":false,"state":"ready","readyOn":"today","deleteRequestedOn":"later","deletedOn":"last"}`, func(c *V3Client) {
		got, err := c.GetVPCSubnet(1, 2)
		if err != nil {
			t.Fatal(err)
		}
		r := got
		if r.SubnetID != 2 || r.Network != "10.1.0.0/24" || r.DeletedOn != "last" {
			t.Fatal(got)
		}
	})
}
func TestCreateVPCSubnetContract(t *testing.T) {
	drop3Contract(t, "POST", "/vpcs/1/subnets", `{"network":"10.1.0.0/24","label":"app"}`, `{"subnetId":2}`, func(c *V3Client) {
		got, err := c.CreateVPCSubnet(1, &CreateVPCSubnetRequest{Network: "10.1.0.0/24", Label: "app"})
		if err != nil {
			t.Fatal(err)
		}
		if got != 2 {
			t.Fatal(got)
		}
	})
}
func TestDeleteVPCSubnetContract(t *testing.T) {
	drop3Contract(t, "DELETE", "/vpcs/1/subnets/2", ``, `{}`, func(c *V3Client) {
		err := c.DeleteVPCSubnet(1, 2)
		if err != nil {
			t.Fatal(err)
		}
	})
}
func TestListVPCServersContract(t *testing.T) {
	drop3Contract(t, "GET", "/vpcs/1/servers", ``, `{"data":[{"serverId":3,"mbpkgid":8,"fqdn":"web","membership":"attached","subnetId":2,"nicId":9,"state":"attached","address":"10.1.0.4"}]}`, func(c *V3Client) {
		got, err := c.ListVPCServers(1)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 1 || got[0].MBPkgID != 8 || got[0].Address != "10.1.0.4" || got[0].Membership != "attached" {
			t.Fatal(got)
		}
	})
}
func TestAttachVPCServerContract(t *testing.T) {
	drop3Contract(t, "POST", "/vpcs/1/servers", `{"mbpkgid":8,"subnetId":2}`, `{"serverId":3,"mbpkgid":8,"fqdn":"web","membership":"attached","subnetId":2,"nicId":9,"state":"attached","address":"10.1.0.4"}`, func(c *V3Client) {
		got, err := c.AttachVPCServer(1, &AttachVPCServerRequest{MBPkgID: 8, SubnetID: 2})
		if err != nil {
			t.Fatal(err)
		}
		if got.ServerID != 3 || got.NICID != 9 || got.State != "attached" {
			t.Fatal(got)
		}
	})
}
func TestDetachVPCServerContract(t *testing.T) {
	drop3Contract(t, "DELETE", "/vpcs/1/servers/8", ``, `{}`, func(c *V3Client) {
		err := c.DetachVPCServer(1, 8)
		if err != nil {
			t.Fatal(err)
		}
	})
}
func TestPutVPCAttachmentBGPContract(t *testing.T) {
	drop3Contract(t, "PUT", "/vpcs/1/attachments/2/bgp", `{"peerAsn":65001,"maxPrefix":100,"password":"test-password"}`, `{"attachmentId":2}`, func(c *V3Client) {
		err := c.PutVPCAttachmentBGP(1, 2, &VPCAttachmentBGPRequest{PeerASN: 65001, MaxPrefix: 100, Password: "test-password"})
		if err != nil {
			t.Fatal(err)
		}
	})
}
func TestDeleteVPCAttachmentBGPContract(t *testing.T) {
	drop3Contract(t, "DELETE", "/vpcs/1/attachments/2/bgp", ``, `{"attachmentId":2}`, func(c *V3Client) {
		err := c.DeleteVPCAttachmentBGP(1, 2)
		if err != nil {
			t.Fatal(err)
		}
	})
}
func TestDrop3AttachmentBGPDecode(t *testing.T) {
	for _, list := range []bool{false, true} {
		row := `{"id":2,"bgp":{"peerAsn":65001,"maxPrefix":100,"hasPassword":true,"gatewayAsn":4200000123,"gatewayAddress":"10.0.0.1","peerAddress":"10.0.0.2","state":"applied","stateMessage":"ok","appliedOn":"today"}}`
		response := row
		path := "/vpcs/1/attachments/2"
		if list {
			response = `{"data":[` + row + `]}`
			path = "/vpcs/1/attachments"
		}
		drop3Contract(t, "GET", path, "", response, func(c *V3Client) {
			var r *VPCAttachment
			var err error
			if list {
				var rows []VPCAttachment
				rows, err = c.ListVPCAttachments(1)
				if len(rows) > 0 {
					r = &rows[0]
				}
			} else {
				r, err = c.GetVPCAttachment(1, 2)
			}
			if err != nil {
				t.Fatal(err)
			}
			if r == nil || r.BGP == nil || r.BGP.GatewayASN != 4200000123 || !r.BGP.HasPassword || r.BGP.AppliedOn != "today" {
				t.Fatal(r)
			}
		})
	}
}
func TestDrop3BGPCreateAndDebugRedaction(t *testing.T) {
	var logs bytes.Buffer
	old := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(old)
	drop3Contract(t, "POST", "/vpcs/1/attachments", `{"type":"customer-router","accept":["192.0.2.0/24"],"bgp":{"peerAsn":65001,"password":"test-password"}}`, `{"attachmentId":2}`, func(c *V3Client) {
		c.debug = true
		id, err := c.CreateVPCAttachment(1, &CreateVPCAttachmentRequest{Type: "customer-router", Accept: []string{"192.0.2.0/24"}, BGP: &VPCAttachmentBGPRequest{PeerASN: 65001, Password: "test-password"}})
		if err != nil || id != 2 {
			t.Fatal(err)
		}
	})
	drop3Contract(t, "PUT", "/vpcs/1/attachments/2/bgp", `{"peerAsn":65001,"password":"test-password"}`, `{"attachmentId":2}`, func(c *V3Client) {
		c.debug = true
		if err := c.PutVPCAttachmentBGP(1, 2, &VPCAttachmentBGPRequest{PeerASN: 65001, Password: "test-password"}); err != nil {
			t.Fatal(err)
		}
	})
	if strings.Contains(logs.String(), "test-password") || strings.Contains(logs.String(), "test-key") {
		t.Fatal("secret reached log")
	}
}

func TestCreateSecurityGroupRuleSourceContract(t *testing.T) {
	drop3Contract(t, "POST", "/cloud/networking/security-groups/2/rules", `{"ipVersion":4,"direction":"inbound","securityGroupId":9}`, `{"ruleId":3}`, func(c *V3Client) {
		id, err := c.CreateSecurityGroupRule(2, &CreateSecurityGroupRuleRequest{IPVersion: 4, Direction: "inbound", SecurityGroupID: 9})
		if err != nil || id != 3 {
			t.Fatal(err)
		}
	})
}
func TestCreateSecurityGroupMemberNodePoolContract(t *testing.T) {
	drop3Contract(t, "POST", "/cloud/networking/security-groups/2/members", `{"vpcId":1,"nodePoolId":9}`, `{"memberId":3}`, func(c *V3Client) {
		id, err := c.CreateSecurityGroupMember(2, &CreateSecurityGroupMemberRequest{VPCID: 1, NodePoolID: 9})
		if err != nil || id != 3 {
			t.Fatal(err)
		}
	})
}
func TestAttachVPCServerDefaultSubnetContract(t *testing.T) {
	drop3Contract(t, "POST", "/vpcs/1/servers", `{"mbpkgid":8}`, `{"mbpkgid":8,"subnetId":null,"state":"attached","nicId":9}`, func(c *V3Client) {
		r, err := c.AttachVPCServer(1, &AttachVPCServerRequest{MBPkgID: 8})
		if err != nil {
			t.Fatal(err)
		}
		if r.SubnetID != 0 || r.NICID != 9 {
			t.Fatal(r)
		}
	})
}
func TestPutVPCAttachmentBGPOmitsPasswordContract(t *testing.T) {
	drop3Contract(t, "PUT", "/vpcs/1/attachments/2/bgp", `{"peerAsn":65001}`, `{"attachmentId":2}`, func(c *V3Client) {
		if err := c.PutVPCAttachmentBGP(1, 2, &VPCAttachmentBGPRequest{PeerASN: 65001}); err != nil {
			t.Fatal(err)
		}
	})
}
