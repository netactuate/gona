package gona

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// LocationByCurrentIP is the location detected for the caller's current IP.
type LocationByCurrentIP struct {
	IP       string          `json:"ip,omitempty"`
	Location string          `json:"location,omitempty"`
	Raw      json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete location response while exposing common fields.
func (l *LocationByCurrentIP) UnmarshalJSON(data []byte) error {
	type Alias LocationByCurrentIP
	aux := (*Alias)(l)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	l.Raw = append(l.Raw[:0], data...)
	return nil
}

// GraphQuery identifies the switch port graph and time range to fetch.
type GraphQuery struct {
	Port int
	Time string
}

// Graph is the graph payload returned for a switch port.
type Graph struct {
	Raw json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete graph payload.
func (g *Graph) UnmarshalJSON(data []byte) error {
	g.Raw = append(g.Raw[:0], data...)
	return nil
}

// GetLocationByCurrentIP returns the platform location detected from the caller's IP address.
func (c *Client) GetLocationByCurrentIP() (*LocationByCurrentIP, error) {
	var location LocationByCurrentIP
	if err := c.get(context.Background(), "location", &location); err != nil {
		return nil, fmt.Errorf("get location by current IP: %w", err)
	}
	return &location, nil
}

// GetGraph returns graph data for a switch port and time range. Time must be one
// of daily, weekly, monthly or yearly.
func (c *Client) GetGraph(q GraphQuery) (*Graph, error) {
	values := url.Values{}
	values.Set("port", strconv.Itoa(q.Port))
	values.Set("time", q.Time)

	var graph Graph
	if err := c.get(context.Background(), "graphs/graph?"+values.Encode(), &graph); err != nil {
		return nil, fmt.Errorf("get graph for port %d: %w", q.Port, err)
	}
	return &graph, nil
}
