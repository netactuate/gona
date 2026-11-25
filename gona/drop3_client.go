package gona

import (
	"errors"
	"net/url"
)

// Password-bearing requests must not reach request-body debug logging.
func (c *V3Client) drop3Request(method, path string, body interface{}) (*V3APIResponse, error) {
	client := *c
	client.debug = false
	r, err := client.doRequest(method, path, body)
	var transport *url.Error
	if errors.As(err, &transport) {
		return r, transport.Err
	}
	return r, err
}
