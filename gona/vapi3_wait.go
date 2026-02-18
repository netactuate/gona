package gona

import (
	"fmt"
	"time"
)

type V3WaitConfig struct {
	Interval time.Duration // time between polls
	Timeout  time.Duration // max total wait time
}

var (
	StorageWaitConfig = V3WaitConfig{Interval: 10 * time.Second, Timeout: 2 * time.Minute}
	VPCWaitConfig     = V3WaitConfig{Interval: 1 * time.Minute, Timeout: 10 * time.Minute}  //1min
	NKEWaitConfig     = V3WaitConfig{Interval: 15 * time.Second, Timeout: 10 * time.Minute}
	RouterWaitConfig  = V3WaitConfig{Interval: 10 * time.Second, Timeout: 5 * time.Minute}
)

func (c *V3Client) waitForCondition(checkFn func() (bool, error), config V3WaitConfig) error {
	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()

	deadline := time.After(config.Timeout)

	ready, err := checkFn()
	if err != nil {
		return err
	}
	if ready {
		return nil
	}

	for {
		select {
		case <-deadline:
			return fmt.Errorf("timeout waiting for condition after %v", config.Timeout)
		case <-ticker.C:
			ready, err := checkFn()
			if err != nil {
				return err
			}
			if ready {
				return nil
			}
			c.debugLog("waiting for condition... (interval=%v, timeout=%v)", config.Interval, config.Timeout)
		}
	}
}

func (c *V3Client) waitForRouterSettled(routerID int, maxSuccessfulGETs int, interval time.Duration) error {
	if maxSuccessfulGETs <= 0 {
		maxSuccessfulGETs = 3
	}

	path := fmt.Sprintf("/cloud-routing/routers/%d", routerID)
	successCount := 0

	timeout := time.Duration(maxSuccessfulGETs*2) * interval
	deadline := time.After(timeout)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-deadline:
			return fmt.Errorf("timeout waiting for router %d to settle after %v (%d/%d successful GETs)",
				routerID, timeout, successCount, maxSuccessfulGETs)
		case <-ticker.C:
			_, err := c.get(path)
			if err != nil {
				successCount = 0
				c.debugLog("router %d settle check failed (resetting count): %v", routerID, err)
				continue
			}
			successCount++
			c.debugLog("router %d settle check: %d/%d successful GETs", routerID, successCount, maxSuccessfulGETs)
			if successCount >= maxSuccessfulGETs {
				return nil
			}
		}
	}
}
