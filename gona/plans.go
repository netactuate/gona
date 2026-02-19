package gona

// Plan struct defines the purchaceable plans/packages
type Plan struct {
	ID        int    `json:"plan_id"`
	Name      string `json:"plan"`
	RAM       string `json:"ram"`
	Disk      string `json:"disk"`
	Transfer  string `json:"transfer"`
	Price     string `json:"price"`
	Available float64 `json:"available"`
}

// GetPlans external method on Client to list available Plans
func (c *Client) GetPlans() ([]Plan, error) {

	var planList []Plan

	if err := c.get("cloud/sizes", &planList); err != nil {
		return nil, err
	}

	return planList, nil
}
