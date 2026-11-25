package gona

import (
	"encoding/json"
	"fmt"
)

type NKEStorageAvailability struct {
	LocationID       int  `json:"locationId"`
	StorageAvailable bool `json:"storageAvailable"`
}

// ListNKEStorageAvailability optionally filters by location; zero lists all locations.
func (c *V3Client) ListNKEStorageAvailability(locationID int) ([]NKEStorageAvailability, error) {
	path := "/nke/storage-availability"
	if locationID > 0 {
		path += fmt.Sprintf("?locationId=%d", locationID)
	}
	r, err := c.getList(path)
	if err != nil {
		return nil, fmt.Errorf("list NKE storage availability: %w", err)
	}
	var rows []NKEStorageAvailability
	if err = json.Unmarshal(r.Data, &rows); err != nil {
		return nil, fmt.Errorf("decode NKE storage availability: %w", err)
	}
	return rows, nil
}
