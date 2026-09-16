package gona

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
)

type PlatformStatusService struct {
	Service     string                   `json:"service"`
	ComponentID string                   `json:"component_id"`
	Locations   []PlatformStatusLocation `json:"locations"`
}

type PlatformStatusLocation struct {
	Location    string `json:"location"`
	ContainerID string `json:"container_id"`
	Status      string `json:"status"`
	LastUpdated string `json:"last_updated"`
}

type platformStatusServiceWire struct {
	ComponentID string                            `json:"component_id"`
	Locations   map[string]PlatformStatusLocation `json:"locations"`
}

func (c *Client) GetPlatformStatus() ([]PlatformStatusService, error) {
	var raw map[string]platformStatusServiceWire
	if err := c.get(context.Background(), "platform/status", &raw); err != nil {
		return nil, fmt.Errorf("get platform status: %w", err)
	}

	services := make([]string, 0, len(raw))
	for service := range raw {
		services = append(services, service)
	}
	sort.Strings(services)

	result := make([]PlatformStatusService, len(services))
	for i, service := range services {
		wire := raw[service]
		locations := make([]string, 0, len(wire.Locations))
		for location := range wire.Locations {
			locations = append(locations, location)
		}
		sort.Strings(locations)

		statusLocations := make([]PlatformStatusLocation, len(locations))
		for j, location := range locations {
			statusLocation := wire.Locations[location]
			statusLocation.Location = location
			statusLocations[j] = statusLocation
		}

		result[i] = PlatformStatusService{
			Service:     service,
			ComponentID: wire.ComponentID,
			Locations:   statusLocations,
		}
	}

	return result, nil
}

type PlatformChangeLogEntry struct {
	ChangeLogID      string          `json:"change_log_id"`
	Title            string          `json:"title"`
	ShortDescription string          `json:"short_description"`
	Status           string          `json:"status"`
	EntryJSON        json.RawMessage `json:"entry_json"`
}

func (e *PlatformChangeLogEntry) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	*e = PlatformChangeLogEntry{EntryJSON: append(e.EntryJSON[:0], data...)}
	if err := unmarshalStringField(raw, "id", &e.ChangeLogID); err != nil {
		return err
	}
	if err := unmarshalStringField(raw, "title", &e.Title); err != nil {
		return err
	}
	if err := unmarshalStringField(raw, "short_description", &e.ShortDescription); err != nil {
		return err
	}
	if err := unmarshalStringField(raw, "status", &e.Status); err != nil {
		return err
	}
	return nil
}

func unmarshalStringField(raw map[string]json.RawMessage, key string, out *string) error {
	val, ok := raw[key]
	if !ok || string(val) == "null" {
		return nil
	}
	var s string
	if err := json.Unmarshal(val, &s); err == nil {
		*out = s
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(val, &n); err == nil {
		*out = n.String()
		return nil
	}
	return fmt.Errorf("%s is not a string or number", key)
}

func (c *Client) GetPlatformChangeLog() ([]PlatformChangeLogEntry, error) {
	var entries []PlatformChangeLogEntry
	if err := c.get(context.Background(), "platform/change-log", &entries); err != nil {
		return nil, fmt.Errorf("get platform change log: %w", err)
	}
	return entries, nil
}

// GetPlatformChangeLogEntry returns a single platform change log entry by id.
func (c *Client) GetPlatformChangeLogEntry(id int) (*PlatformChangeLogEntry, error) {
	var entry PlatformChangeLogEntry
	path := fmt.Sprintf("platform/change-log/%d", id)
	if err := c.get(context.Background(), path, &entry); err != nil {
		return nil, fmt.Errorf("get platform change log entry %d: %w", id, err)
	}
	return &entry, nil
}

// GetPlatformDatacenters returns datacenters for a location.
func (c *Client) GetPlatformDatacenters(location string) ([]Datacenter, error) {
	var datacenters []Datacenter
	path := "platform/datacenters/" + url.PathEscape(location)
	if err := c.get(context.Background(), path, &datacenters); err != nil {
		return nil, fmt.Errorf("get platform datacenters for %s: %w", location, err)
	}
	return datacenters, nil
}

// PlatformLookingGlassInit is the looking glass initialization payload.
type PlatformLookingGlassInit struct {
	Raw json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete looking glass initialization payload.
func (i *PlatformLookingGlassInit) UnmarshalJSON(data []byte) error {
	i.Raw = append(i.Raw[:0], data...)
	return nil
}

// PlatformLookingGlassExecuteOptions contains optional looking glass execution parameters.
type PlatformLookingGlassExecuteOptions struct {
	Action   *string
	Target   *string
	Location *string
	Full     *int
}

// PlatformLookingGlassResult is the output returned by a looking glass action.
type PlatformLookingGlassResult struct {
	Raw json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete looking glass execution payload.
func (r *PlatformLookingGlassResult) UnmarshalJSON(data []byte) error {
	r.Raw = append(r.Raw[:0], data...)
	return nil
}

// PlatformMaintenanceInfo is a platform maintenance detail payload.
type PlatformMaintenanceInfo struct {
	Raw json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete platform maintenance detail payload.
func (i *PlatformMaintenanceInfo) UnmarshalJSON(data []byte) error {
	i.Raw = append(i.Raw[:0], data...)
	return nil
}

// GetPlatformLookingGlassInit returns the options used to initialize looking glass calls.
func (c *Client) GetPlatformLookingGlassInit() (*PlatformLookingGlassInit, error) {
	var init PlatformLookingGlassInit
	if err := c.get(context.Background(), "platform/looking-glass/init", &init); err != nil {
		return nil, fmt.Errorf("get platform looking glass init: %w", err)
	}
	return &init, nil
}

// ExecutePlatformLookingGlass executes a looking glass action with optional query parameters.
func (c *Client) ExecutePlatformLookingGlass(opts PlatformLookingGlassExecuteOptions) (*PlatformLookingGlassResult, error) {
	values := url.Values{}
	if opts.Action != nil {
		values.Set("action", *opts.Action)
	}
	if opts.Target != nil {
		values.Set("target", *opts.Target)
	}
	if opts.Location != nil {
		values.Set("location", *opts.Location)
	}
	if opts.Full != nil {
		values.Set("full", strconv.Itoa(*opts.Full))
	}

	path := "platform/looking-glass/execute"
	if encoded := values.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var result PlatformLookingGlassResult
	if err := c.get(context.Background(), path, &result); err != nil {
		return nil, fmt.Errorf("execute platform looking glass: %w", err)
	}
	return &result, nil
}

// GetPlatformMaintenanceInfo returns maintenance detail by id.
func (c *Client) GetPlatformMaintenanceInfo(id int) (*PlatformMaintenanceInfo, error) {
	var info PlatformMaintenanceInfo
	path := fmt.Sprintf("platform/maintenance-info/%d", id)
	if err := c.get(context.Background(), path, &info); err != nil {
		return nil, fmt.Errorf("get platform maintenance info %d: %w", id, err)
	}
	return &info, nil
}

type PlatformEvents struct {
	Active   []PlatformEvent `json:"active"`
	Upcoming []PlatformEvent `json:"upcoming"`
	Historic []PlatformEvent `json:"historic"`
}

type PlatformEvent struct {
	EventID    string   `json:"event_id"`
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Status     string   `json:"status"`
	StartTime  string   `json:"start_time"`
	EndTime    string   `json:"end_time"`
	Components []string `json:"components"`
	Containers []string `json:"containers"`
}

func (c *Client) GetPlatformIncidents(location string) (PlatformEvents, error) {
	return c.getPlatformEvents("platform/incidents/"+location, "get platform incidents")
}

func (c *Client) GetPlatformIncidentHistory(location string) (PlatformEvents, error) {
	return c.getPlatformEvents("platform/incidents/history/"+location, "get platform incident history")
}

func (c *Client) GetPlatformMaintenance(location string) (PlatformEvents, error) {
	return c.getPlatformEvents("platform/maintenance/"+location, "get platform maintenance")
}

func (c *Client) GetPlatformMaintenanceHistory(location string) (PlatformEvents, error) {
	return c.getPlatformEvents("platform/maintenance/history/"+location, "get platform maintenance history")
}

func (c *Client) getPlatformEvents(path, operation string) (PlatformEvents, error) {
	var events PlatformEvents
	if err := c.get(context.Background(), path, &events); err != nil {
		return PlatformEvents{}, fmt.Errorf("%s: %w", operation, err)
	}
	return events, nil
}
