package gona

import "context"

type Size struct {
	PlanID    int     `json:"plan_id"`
	Plan      string  `json:"plan"`
	RAM       string  `json:"ram"`
	Disk      string  `json:"disk"`
	Transfer  string  `json:"transfer"`
	Price     string  `json:"price"`
	CPU       int     `json:"cpu"`
	Port      string  `json:"port"`
	Available float64 `json:"available"`
}

func (c *Client) GetSizes() ([]Size, error) {
	var sizes []Size
	if err := c.get(context.Background(), "cloud/sizes", &sizes); err != nil {
		return nil, err
	}
	return sizes, nil
}
