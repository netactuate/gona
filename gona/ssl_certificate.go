package gona

import (
	"encoding/json"
	"fmt"
)

type SSLCertificate struct {
	SSLCertificateID int                 `json:"sslCertificateId"`
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	Fingerprint      string              `json:"fingerprint"`
	Domains          []string            `json:"domains"`
	IsActive         bool                `json:"isActive"`
	Status           string              `json:"status"`
	Dates            *SSLCertificateDates `json:"dates,omitempty"`
}

type SSLCertificateDates struct {
	Created    string `json:"created,omitempty"`
	Updated    string `json:"updated,omitempty"`
	NotBefore  string `json:"notBefore,omitempty"`
	Expiration string `json:"expiration,omitempty"`
}

type CreateSSLCertificateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
}

type UpdateSSLCertificateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Certificate string `json:"certificate,omitempty"`
	PrivateKey  string `json:"privateKey,omitempty"`
}

type CreateSSLCertificateResponse struct {
	SSLCertificateID int `json:"sslCertificateId"`
}

func (c *V3Client) CreateSSLCertificate(req *CreateSSLCertificateRequest) (*CreateSSLCertificateResponse, error) {
	resp, err := c.post("/ssl-certificates", req)
	if err != nil {
		return nil, fmt.Errorf("create SSL certificate: %w", err)
	}
	var result CreateSSLCertificateResponse
	if err := json.Unmarshal(resp.Data, &result); err != nil {
		return nil, fmt.Errorf("create SSL certificate unmarshal: %w", err)
	}
	return &result, nil
}

func (c *V3Client) GetSSLCertificate(id int) (*SSLCertificate, error) {
	path := fmt.Sprintf("/ssl-certificates/%d", id)
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get SSL certificate %d: %w", id, err)
	}
	var cert SSLCertificate
	if err := json.Unmarshal(resp.Data, &cert); err != nil {
		return nil, fmt.Errorf("get SSL certificate unmarshal: %w", err)
	}
	return &cert, nil
}

func (c *V3Client) UpdateSSLCertificate(id int, req *UpdateSSLCertificateRequest) error {
	path := fmt.Sprintf("/ssl-certificates/%d", id)
	_, err := c.patch(path, req)
	if err != nil {
		return fmt.Errorf("update SSL certificate %d: %w", id, err)
	}
	return nil
}

func (c *V3Client) DeleteSSLCertificate(id int) error {
	path := fmt.Sprintf("/ssl-certificates/%d", id)
	_, err := c.del(path)
	if err != nil {
		return fmt.Errorf("delete SSL certificate %d: %w", id, err)
	}
	return nil
}
