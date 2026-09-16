package gona

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// Tag is a first-class NetActuate tag object. Tags are global and shared; a tag
// is associated with resources (servers, NKE clusters, ...) many-to-many via the
// assign-resource / remove-resource endpoints. GET /tags returns every tag with
// its full Resources slice embedded.
type Tag struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	Icon           string        `json:"icon"`
	Color          string        `json:"color"`
	IsDefault      int           `json:"is_default"`
	IsFavorite     int           `json:"is_favorite"`
	IsLocked       int           `json:"is_locked"`
	ShowDashboard  int           `json:"show_dashboard"`
	CreatedAt      string        `json:"created_at"`
	MbID           int           `json:"mb_id"`
	ResourcesCount int           `json:"resources_count"`
	Resources      []TagResource `json:"resources"`
}

// TagResource is a single tag<->resource assignment row (Tag.Resources[]).
// Identifier is the resource's primary numeric id (server mbpkgid, nke clusterId).
type TagResource struct {
	ID            int    `json:"id"`
	ResourceTagID int    `json:"resource_tag_id"`
	ResourceName  string `json:"resource_name"`
	Identifier    int    `json:"identifier"`
	CreatedAt     string `json:"created_at"`
}

// TagLog is a log entry for a tag.
type TagLog struct {
	ID        int             `json:"id,omitempty"`
	TagID     int             `json:"tag_id,omitempty"`
	Action    string          `json:"action,omitempty"`
	Message   string          `json:"message,omitempty"`
	CreatedAt string          `json:"created_at,omitempty"`
	Raw       json.RawMessage `json:"-"`
}

// UnmarshalJSON records the complete tag log payload while exposing common fields.
func (l *TagLog) UnmarshalJSON(data []byte) error {
	type Alias TagLog
	aux := (*Alias)(l)
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	l.Raw = append(l.Raw[:0], data...)
	return nil
}

// CreateTagRequest is the POST /tags body. Per the API, the create call only
// accepts name/description/icon/color; the is_* flags are PUT-only.
type CreateTagRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Color       string `json:"color,omitempty"`
}

// UpdateTagRequest is the PUT /tags/{id} body. Flags are 0/1 integers.
type UpdateTagRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Icon          string `json:"icon,omitempty"`
	Color         string `json:"color,omitempty"`
	IsDefault     int    `json:"is_default"`
	IsFavorite    int    `json:"is_favorite"`
	IsLocked      int    `json:"is_locked"`
	ShowDashboard int    `json:"show_dashboard"`
}

type tagResourceRef struct {
	ResourceName string `json:"resource_name"`
	// The API request schema types identifier as a string even though it is
	// stored/returned as a number; send it as a string for spec-compliance.
	Identifier string `json:"identifier"`
}

// GetTags lists every tag (each with its Resources slice embedded).
func (c *Client) GetTags() ([]Tag, error) {
	var tags []Tag
	if err := c.get(context.Background(), "tags", &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

// GetTag returns a single tag by id. There is no single-GET endpoint, so this
// filters the list (the same source the portal uses).
func (c *Client) GetTag(id int) (*Tag, error) {
	tags, err := c.GetTags()
	if err != nil {
		return nil, err
	}
	for i := range tags {
		if tags[i].ID == id {
			return &tags[i], nil
		}
	}
	return nil, fmt.Errorf("tag %d not found", id)
}

// CreateTag creates a tag object. The caller should read back by id/name to get
// the authoritative record (the create response shape is not contractually fixed).
func (c *Client) CreateTag(r *CreateTagRequest) (*Tag, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	var tag Tag
	if err := c.postJSON(context.Background(), "tags", body, &tag); err != nil {
		return nil, err
	}
	return &tag, nil
}

// UpdateTag updates a tag object (name/description/icon/color and the is_* flags).
func (c *Client) UpdateTag(id int, r *UpdateTagRequest) (*Tag, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	var tag Tag
	if err := c.putJSON(context.Background(), "tags/"+strconv.Itoa(id), body, &tag); err != nil {
		return nil, err
	}
	return &tag, nil
}

// DeleteTag deletes a tag object. The API returns 412 when the tag is locked or
// still assigned; that surfaces as an error from do().
func (c *Client) DeleteTag(id int) error {
	return c.delete(context.Background(), "tags/"+strconv.Itoa(id), nil, nil)
}

// AssignTagResource attaches tagID to a resource (idempotent server-side).
func (c *Client) AssignTagResource(tagID int, resourceName string, identifier int) error {
	body, err := json.Marshal(tagResourceRef{
		ResourceName: resourceName,
		Identifier:   strconv.Itoa(identifier),
	})
	if err != nil {
		return err
	}
	return c.postJSON(context.Background(), "tags/"+strconv.Itoa(tagID)+"/assign-resource", body, nil)
}

// RemoveTagResource detaches tagID from a resource. It removes only the
// association; the tag object itself is never deleted.
func (c *Client) RemoveTagResource(tagID int, resourceName string, identifier int) error {
	body, err := json.Marshal(tagResourceRef{
		ResourceName: resourceName,
		Identifier:   strconv.Itoa(identifier),
	})
	if err != nil {
		return err
	}
	return c.postJSON(context.Background(), "tags/"+strconv.Itoa(tagID)+"/remove-resource", body, nil)
}

// GetResourceTags returns the tags currently assigned to a given resource.
func (c *Client) GetResourceTags(resourceName string, resourceID int) ([]Tag, error) {
	var tags []Tag
	path := fmt.Sprintf("tags/resource/%s/id/%d", url.PathEscape(resourceName), resourceID)
	if err := c.get(context.Background(), path, &tags); err != nil {
		return nil, fmt.Errorf("get tags for %s %d: %w", resourceName, resourceID, err)
	}
	return tags, nil
}

// GetTagResources returns the resources currently assigned to a tag.
func (c *Client) GetTagResources(tagID int) ([]TagResource, error) {
	var resources []TagResource
	path := fmt.Sprintf("tags/%d/resources", tagID)
	if err := c.get(context.Background(), path, &resources); err != nil {
		return nil, fmt.Errorf("get resources for tag %d: %w", tagID, err)
	}
	return resources, nil
}

// GetTagLogs returns log entries for a tag.
func (c *Client) GetTagLogs(tagID int) ([]TagLog, error) {
	var logs []TagLog
	path := fmt.Sprintf("tags/%d/logs", tagID)
	if err := c.get(context.Background(), path, &logs); err != nil {
		return nil, fmt.Errorf("get logs for tag %d: %w", tagID, err)
	}
	return logs, nil
}
