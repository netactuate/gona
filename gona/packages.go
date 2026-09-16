package gona

import "strconv"
import (
	//"github.com/google/go-querystring/query"
	"bytes"
	"context"
	"encoding/json"
)

// Package struct stores the purchased package values.
//
// ID, Locked and Installed are json.Number rather than int or string. The API does not quote
// those values: cloud/packages and cloud/package/{id} both send them as unquoted numbers, so every call
// failed outright with
//
//	json: invalid use of ,string struct tag, trying to unmarshal unquoted value into int
//
// Locked was worse: declared a string against a numeric field. GetPackages and GetPackage
// could therefore never have worked against the live API, and nothing called them until the
// netactuate_packages data source did.
//
// json.Number accepts a quoted or an unquoted number, which is the point. This platform is
// demonstrably inconsistent about value shapes across endpoints, so a type that tolerates both
// removes the class rather than the instance.
type Package struct {
	ID        json.Number `json:"mbpkgid"`
	Status    string      `json:"package_status"`
	Locked    json.Number `json:"locked"`
	PlanName  string      `json:"name"`
	Installed json.Number `json:"installed"`
}

type CancelRequest struct {
	MBPKGID     int     `json:"mbpkgid"`
	DomUPackage *string `json:"domU_package,omitempty"`
	Comments    *string `json:"comments,omitempty"`
	CancelType  string  `json:"cancel_type"`
	Agree       int     `json:"agree"`
	Password    *string `json:"password,omitempty"`
}

// GetPackages external method on Client that returns a
// list of Package object from the API
func (c *Client) GetPackages() ([]Package, error) {

	var packageList []Package

	if err := c.get(context.Background(), "cloud/packages", &packageList); err != nil {
		return nil, err
	}

	return packageList, nil
}

// GetPackage external method on Client that takes an id (int) as it's sole
// argument and returns a single Package object
func (c *Client) GetPackage(id int) (pkg Package, err error) {
	if err := c.get(context.Background(), "cloud/package/"+strconv.Itoa(id), &pkg); err != nil {
		return Package{}, err
	}
	return pkg, nil
}

func (c *Client) CancelPackage(req *CancelRequest) (result interface{}, err error) {
	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	//call with JSON in case not work properly
	reqHTTP, err := c.newRequest(context.Background(), "POST", "cloud/package/cancel/", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	reqHTTP.Header.Set("Content-Type", "application/json")

	result = make(map[string]interface{})
	if err := c.do(reqHTTP, &result); err != nil {
		return nil, err
	}

	return result, nil
}
