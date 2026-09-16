package gona

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// Service is a top-level account service record.
type Service struct {
	ID          int             `json:"id,omitempty"`
	Description string          `json:"description,omitempty"`
	Raw         json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete service payload while exposing common fields.
func (s *Service) UnmarshalJSON(data []byte) error {
	type Alias Service
	aux := (*Alias)(s)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0], data...)
	return nil
}

// ServiceListOptions filters service list endpoints by service id.
type ServiceListOptions struct {
	ServiceID *int
}

// ColocationService is a colocation service record.
type ColocationService struct {
	ID             int             `json:"id,omitempty"`
	ServiceID      int             `json:"service_id,omitempty"`
	DatacenterID   int             `json:"datacenter_id,omitempty"`
	RackIdentifier string          `json:"rack_identifier,omitempty"`
	PowerDetails   string          `json:"power_details,omitempty"`
	Description    string          `json:"description,omitempty"`
	Raw            json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete colocation service payload.
func (s *ColocationService) UnmarshalJSON(data []byte) error {
	type Alias ColocationService
	aux := (*Alias)(s)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0], data...)
	return nil
}

// IPTransitService is an IP transit service record.
type IPTransitService struct {
	ID           int             `json:"id,omitempty"`
	ServiceID    int             `json:"service_id,omitempty"`
	DatacenterID int             `json:"datacenter_id,omitempty"`
	BGPGroupID   int             `json:"bgp_group_id,omitempty"`
	Description  string          `json:"description,omitempty"`
	Raw          json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete IP transit service payload.
func (s *IPTransitService) UnmarshalJSON(data []byte) error {
	type Alias IPTransitService
	aux := (*Alias)(s)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0], data...)
	return nil
}

// IPTransitIPAddress is an IP address assigned to an IP transit service.
type IPTransitIPAddress struct {
	ID                 int             `json:"id,omitempty"`
	ServiceIPTransitID int             `json:"service_iptransit_id,omitempty"`
	IP                 string          `json:"ip,omitempty"`
	Raw                json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete IP transit IP address payload.
func (a *IPTransitIPAddress) UnmarshalJSON(data []byte) error {
	type Alias IPTransitIPAddress
	aux := (*Alias)(a)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	a.Raw = append(a.Raw[:0], data...)
	return nil
}

// IPTransitPort is a port assigned to an IP transit service.
type IPTransitPort struct {
	ID                 int             `json:"id,omitempty"`
	ServiceIPTransitID int             `json:"service_iptransit_id,omitempty"`
	Name               string          `json:"name,omitempty"`
	Raw                json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete IP transit port payload.
func (p *IPTransitPort) UnmarshalJSON(data []byte) error {
	type Alias IPTransitPort
	aux := (*Alias)(p)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	p.Raw = append(p.Raw[:0], data...)
	return nil
}

// TransportService is a transport service record.
type TransportService struct {
	ID           int             `json:"id,omitempty"`
	ServiceID    int             `json:"service_id,omitempty"`
	DatacenterID int             `json:"datacenter_id,omitempty"`
	Description  string          `json:"description,omitempty"`
	Raw          json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete transport service payload.
func (s *TransportService) UnmarshalJSON(data []byte) error {
	type Alias TransportService
	aux := (*Alias)(s)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	s.Raw = append(s.Raw[:0], data...)
	return nil
}

// TransportPort is a port assigned to a transport service.
type TransportPort struct {
	ID                 int             `json:"id,omitempty"`
	ServiceTransportID int             `json:"service_transport_id,omitempty"`
	Name               string          `json:"name,omitempty"`
	Raw                json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete transport port payload.
func (p *TransportPort) UnmarshalJSON(data []byte) error {
	type Alias TransportPort
	aux := (*Alias)(p)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	p.Raw = append(p.Raw[:0], data...)
	return nil
}

// GetServices returns all account services.
func (c *Client) GetServices() ([]Service, error) {
	var services []Service
	if err := c.get(context.Background(), "services", &services); err != nil {
		return nil, fmt.Errorf("get services: %w", err)
	}
	return services, nil
}

// GetColocationServices returns colocation services, optionally filtered by service id.
func (c *Client) GetColocationServices(opts ServiceListOptions) ([]ColocationService, error) {
	var services []ColocationService
	if err := c.get(context.Background(), serviceListPath("services/colocation", opts.ServiceID), &services); err != nil {
		return nil, fmt.Errorf("get colocation services: %w", err)
	}
	return services, nil
}

// GetColocationService returns a colocation service by id.
func (c *Client) GetColocationService(id int) (*ColocationService, error) {
	var service ColocationService
	if err := c.get(context.Background(), "services/colocation/"+strconv.Itoa(id), &service); err != nil {
		return nil, fmt.Errorf("get colocation service %d: %w", id, err)
	}
	return &service, nil
}

// GetIPTransitServices returns IP transit services, optionally filtered by service id.
func (c *Client) GetIPTransitServices(opts ServiceListOptions) ([]IPTransitService, error) {
	var services []IPTransitService
	if err := c.get(context.Background(), serviceListPath("services/iptransit", opts.ServiceID), &services); err != nil {
		return nil, fmt.Errorf("get IP transit services: %w", err)
	}
	return services, nil
}

// GetIPTransitService returns an IP transit service by id.
func (c *Client) GetIPTransitService(id int) (*IPTransitService, error) {
	var service IPTransitService
	if err := c.get(context.Background(), "services/iptransit/"+strconv.Itoa(id), &service); err != nil {
		return nil, fmt.Errorf("get IP transit service %d: %w", id, err)
	}
	return &service, nil
}

// GetIPTransitIPAddresses returns IP addresses, optionally filtered by IP transit service id.
func (c *Client) GetIPTransitIPAddresses(serviceIPTransitID *int) ([]IPTransitIPAddress, error) {
	var addresses []IPTransitIPAddress
	path := serviceListPathWithParam("services/iptransit/ips", "service_iptransit_id", serviceIPTransitID)
	if err := c.get(context.Background(), path, &addresses); err != nil {
		return nil, fmt.Errorf("get IP transit IP addresses: %w", err)
	}
	return addresses, nil
}

// GetIPTransitPorts returns ports, optionally filtered by IP transit service id.
func (c *Client) GetIPTransitPorts(serviceIPTransitID *int) ([]IPTransitPort, error) {
	var ports []IPTransitPort
	path := serviceListPathWithParam("services/iptransit/ports", "service_iptransit_id", serviceIPTransitID)
	if err := c.get(context.Background(), path, &ports); err != nil {
		return nil, fmt.Errorf("get IP transit ports: %w", err)
	}
	return ports, nil
}

// GetTransportServices returns transport services, optionally filtered by service id.
func (c *Client) GetTransportServices(opts ServiceListOptions) ([]TransportService, error) {
	var services []TransportService
	if err := c.get(context.Background(), serviceListPath("services/transport", opts.ServiceID), &services); err != nil {
		return nil, fmt.Errorf("get transport services: %w", err)
	}
	return services, nil
}

// GetTransportService returns a transport service by id.
func (c *Client) GetTransportService(id int) (*TransportService, error) {
	var service TransportService
	if err := c.get(context.Background(), "services/transport/"+strconv.Itoa(id), &service); err != nil {
		return nil, fmt.Errorf("get transport service %d: %w", id, err)
	}
	return &service, nil
}

// GetTransportPorts returns ports, optionally filtered by transport service id.
func (c *Client) GetTransportPorts(serviceTransportID *int) ([]TransportPort, error) {
	var ports []TransportPort
	path := serviceListPathWithParam("services/transport/ports", "service_transport_id", serviceTransportID)
	if err := c.get(context.Background(), path, &ports); err != nil {
		return nil, fmt.Errorf("get transport ports: %w", err)
	}
	return ports, nil
}

func serviceListPath(path string, serviceID *int) string {
	return serviceListPathWithParam(path, "service_id", serviceID)
}

func serviceListPathWithParam(path, name string, value *int) string {
	if value == nil {
		return path
	}
	values := url.Values{}
	values.Set(name, strconv.Itoa(*value))
	return path + "?" + values.Encode()
}
