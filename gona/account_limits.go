package gona

import (
	"encoding/json"
	"fmt"
)

type AccountLimit struct {
	Used         int      `json:"used"`
	Max          int      `json:"max"`
	AllowedPlans []string `json:"allowedPlans"`
}

func (c *V3Client) GetAccountLimits() (map[string]AccountLimit, error) {
	resp, err := c.get("/account-limits")
	if err != nil {
		return nil, fmt.Errorf("get account limits: %w", err)
	}
	var limits map[string]AccountLimit
	if err := json.Unmarshal(resp.Data, &limits); err != nil {
		return nil, fmt.Errorf("get account limits: %w", err)
	}
	return limits, nil
}
