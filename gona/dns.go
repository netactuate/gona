package gona

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// FlexibleTTL is a TTL that the API returns as a string from GET /dns/zone/{id} and as a
// number from GET /dns/zones. Modelled the
// way FlexibleIPv4 in router.go already handles the same disagreement, rather than inventing
// a second idiom for it.
type FlexibleTTL string

// Int returns the TTL as an integer, or 0 if it is empty or unparseable. Callers that
// need the schema's int type use this rather than converting at each site.
func (f FlexibleTTL) Int() int {
	n, err := strconv.Atoi(string(f))
	if err != nil {
		return 0
	}
	return n
}

func (f *FlexibleTTL) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexibleTTL(s)
		return nil
	}
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		*f = FlexibleTTL(strconv.FormatInt(i, 10))
		return nil
	}
	if string(data) == "null" {
		*f = ""
		return nil
	}
	return fmt.Errorf("ttl must be either a string or a number, got %s", string(data))
}

// DNSZoneSOA is the zone's start-of-authority block. The API returns it as an object on
// GET /dns/zone/{id}, and every numeric field inside it comes back as a STRING.
type DNSZoneSOA struct {
	Primary    string `json:"primary"`
	Hostmaster string `json:"hostmaster"`
	Serial     string `json:"serial"`
	Refresh    string `json:"refresh"`
	Retry      string `json:"retry"`
	Expire     string `json:"expire"`
	DefaultTTL string `json:"default_ttl"`
}

// DNSZone models a zone as the API actually returns it, which is not what the list and the
// single GET have in common:
//
//	soa  is an OBJECT, not a string
//	ns   is an ARRAY of NS records, not a string
//	ttl  is a STRING ("3600"), not an integer
//	master appears on the LIST response only; the single GET returns ns and soa instead
//
// Getting any of these wrong makes the whole response fail to unmarshal, which is how the
// first live run of the DNS acceptance tests failed.
type DNSZone struct {
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	IP      string      `json:"ip,omitempty"`
	Master  string      `json:"master,omitempty"`
	NS      []DNSRecord `json:"ns,omitempty"`
	SOA     *DNSZoneSOA `json:"soa,omitempty"`
	TTL     FlexibleTTL `json:"ttl,omitempty"`
	Records []DNSRecord `json:"records,omitempty"`
}

type CreateDNSZoneRequest struct {
	Name string `json:"name"`
	Type string `json:"type"`
	IP   string `json:"ip,omitempty"`
}

type DNSRecord struct {
	ID       int         `json:"id"`
	ZoneID   int         `json:"domain_id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Content  string      `json:"content"`
	TTL      FlexibleTTL `json:"ttl"`
	Priority int         `json:"prio"`
}

type CreateDNSRecordRequest struct {
	ZoneID        int    `json:"domain_id"`
	Name          string `json:"name,omitempty"`
	Type          string `json:"type"`
	TTL           int    `json:"ttl,omitempty"`
	Priority      int    `json:"prio,omitempty"`
	RecordContent string `json:"record_content"`
}

type UpdateDNSRecordRequest struct {
	ID            int    `json:"id"`
	ZoneID        int    `json:"domain_id"`
	Name          string `json:"name,omitempty"`
	Type          string `json:"type"`
	TTL           int    `json:"ttl,omitempty"`
	Priority      int    `json:"prio,omitempty"`
	RecordContent string `json:"record_content"`
}

func (c *Client) ListZones(zoneType string) ([]DNSZone, error) {
	var zones []DNSZone
	values := url.Values{}
	values.Add("type", zoneType)
	if err := c.get(context.Background(), "dns/zones?"+values.Encode(), &zones); err != nil {
		return nil, err
	}
	return zones, nil
}

func (c *Client) GetZone(id int) (*DNSZone, error) {
	var zone DNSZone
	if err := c.get(context.Background(), "dns/zone/"+strconv.Itoa(id), &zone); err != nil {
		return nil, err
	}
	return &zone, nil
}

func (c *Client) CreateZone(req *CreateDNSZoneRequest) (*DNSZone, error) {
	values := url.Values{}
	values.Add("name", req.Name)
	values.Add("type", req.Type)
	if req.IP != "" {
		values.Add("ip", req.IP)
	}

	var zone DNSZone
	if err := c.post(context.Background(), "dns/zone", []byte(values.Encode()), &zone); err != nil {
		return nil, err
	}
	return &zone, nil
}

func (c *Client) DeleteZone(id int) error {
	return c.delete(context.Background(), "dns/zone/"+strconv.Itoa(id), nil, nil)
}

func (c *Client) CreateRecord(req *CreateDNSRecordRequest) (*DNSRecord, error) {
	values := url.Values{}
	values.Add("domain_id", strconv.Itoa(req.ZoneID))
	values.Add("name", req.Name)
	values.Add("type", req.Type)
	values.Add("record_content", req.RecordContent)
	if req.TTL != 0 {
		values.Add("ttl", strconv.Itoa(req.TTL))
	}
	if req.Priority != 0 {
		values.Add("prio", strconv.Itoa(req.Priority))
	}

	var record DNSRecord
	if err := c.post(context.Background(), "dns/record", []byte(values.Encode()), &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (c *Client) GetRecord(id int) (*DNSRecord, error) {
	var record DNSRecord
	if err := c.get(context.Background(), "dns/record/"+strconv.Itoa(id), &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (c *Client) UpdateRecord(req *UpdateDNSRecordRequest) (*DNSRecord, error) {
	values := url.Values{}
	values.Add("id", strconv.Itoa(req.ID))
	values.Add("domain_id", strconv.Itoa(req.ZoneID))
	values.Add("name", req.Name)
	values.Add("type", req.Type)
	values.Add("record_content", req.RecordContent)
	if req.TTL != 0 {
		values.Add("ttl", strconv.Itoa(req.TTL))
	}
	if req.Priority != 0 {
		values.Add("prio", strconv.Itoa(req.Priority))
	}

	var record DNSRecord
	if err := c.put(context.Background(), "dns/record", []byte(values.Encode()), &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (c *Client) DeleteRecord(id int) error {
	return c.delete(context.Background(), "dns/record/"+strconv.Itoa(id), nil, nil)
}

func (c *Client) ListRecords(zoneID int) ([]DNSRecord, error) {
	var records []DNSRecord
	if err := c.get(context.Background(), "dns/records/"+strconv.Itoa(zoneID), &records); err != nil {
		return nil, err
	}
	return records, nil
}
