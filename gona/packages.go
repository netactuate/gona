package gona

import "strconv"
import (
	//"github.com/google/go-querystring/query"
    "bytes"
    "encoding/json"
)

// Package struct stores the purchaced package values
type Package struct {
	ID        int    `json:"mbpkgid,string"`
	Status    string `json:"package_status"`
	Locked    string `json:"locked"`
	PlanName  string `json:"name"`
	Installed int    `json:"installed,string"`
}

type CancelRequest struct {
    MBPKGID     int     `json:"mbpkgid"`
    DomUPackage *string `json:"domU_package,omitempty"`
    Comments    *string `json:"comments,omitempty"`
    CancelType  string  `json:"cancel_type"`
    Agree       int  `json:"agree"`
    Password    *string `json:"password,omitempty"`
}

// GetPackages external method on Client that returns a
// list of Package object from the API
func (c *Client) GetPackages() ([]Package, error) {

	var packageList []Package

	if err := c.get("cloud/packages", &packageList); err != nil {
		return nil, err
	}

	return packageList, nil
}

// GetPackage external method on Client that takes an id (int) as it's sole
// argument and returns a single Package object
func (c *Client) GetPackage(id int) (pkg Package, err error) {
	if err := c.get("/cloud/package/"+strconv.Itoa(id), &pkg); err != nil {
		return Package{Installed: 0}, err
	}
	return pkg, nil
}

func (c *Client) CancelPackage(req *CancelRequest) (result interface{}, err error) {
    jsonBody, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }
    //call with JSON in case not work properly
    reqHTTP, err := c.newRequest("POST", "cloud/package/cancel/", bytes.NewBuffer(jsonBody))
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
