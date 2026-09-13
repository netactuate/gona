package gona

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type OIDCBool bool

func (b *OIDCBool) UnmarshalJSON(data []byte) error {
	var boolValue bool
	if err := json.Unmarshal(data, &boolValue); err == nil {
		*b = OIDCBool(boolValue)
		return nil
	}

	var numberValue json.Number
	if err := json.Unmarshal(data, &numberValue); err == nil {
		i, err := strconv.Atoi(numberValue.String())
		if err != nil {
			return err
		}
		*b = OIDCBool(i != 0)
		return nil
	}

	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		switch stringValue {
		case "true", "1":
			*b = true
			return nil
		case "false", "0", "":
			*b = false
			return nil
		}
	}

	return fmt.Errorf("OIDCBool must be a bool or numeric bool, got %s", string(data))
}

func (b OIDCBool) Bool() bool {
	return bool(b)
}

type OIDCClient struct {
	ClientID         int         `json:"clientId"`
	CreatedOn        string      `json:"createdOn"`
	LastUsedOn       *string     `json:"lastUsedOn"`
	Label            string      `json:"label"`
	Description      string      `json:"description"`
	JWKSURI          *string     `json:"jwksUri"`
	AccountDefault   bool        `json:"accountDefault"`
	DefaultAudience  string      `json:"defaultAudience"`
	TTL              int         `json:"ttl"`
	EnforceAllowList OIDCBool    `json:"enforceAllowList"`
	Tenant           json.Number `json:"-"`
	Keys             []OIDCClientKey
	AuthLogs         []OIDCClientAuthLog
	ChangeLogs       []OIDCClientChangeLog
}

type oidcClientListRow struct {
	ClientID         int      `json:"clientId"`
	CreatedOn        string   `json:"createdOn"`
	LastUsedOn       *string  `json:"lastUsedOn"`
	Label            string   `json:"label"`
	Description      string   `json:"description"`
	JWKSURI          *string  `json:"jwksUri"`
	JWKSHTTPSURL     *string  `json:"jwksHttpsUrl"`
	AccountDefault   bool     `json:"accountDefault"`
	DefaultAudience  string   `json:"defaultAudience"`
	TTL              int      `json:"ttl"`
	EnforceAllowList OIDCBool `json:"enforceAllowList"`
}

func (r oidcClientListRow) toClient(tenant json.Number) OIDCClient {
	jwksURI := r.JWKSURI
	// The create and single-client APIs use jwksUri, but the list API returns
	// the same value as jwksHttpsUrl. Keep both mapped so reads do not drift.
	if jwksURI == nil {
		jwksURI = r.JWKSHTTPSURL
	}
	return OIDCClient{
		ClientID:         r.ClientID,
		CreatedOn:        r.CreatedOn,
		LastUsedOn:       r.LastUsedOn,
		Label:            r.Label,
		Description:      r.Description,
		JWKSURI:          jwksURI,
		AccountDefault:   r.AccountDefault,
		DefaultAudience:  r.DefaultAudience,
		TTL:              r.TTL,
		EnforceAllowList: r.EnforceAllowList,
		Tenant:           tenant,
	}
}

type oidcClientsResponse struct {
	Tenant  json.Number `json:"tenant"`
	Clients V3ListData  `json:"clients"`
}

type oidcClientDetailResponse struct {
	Metadata oidcClientMetadata `json:"metadata"`
	Keys     V3ListData         `json:"keys"`
	Logs     oidcClientLogs     `json:"logs"`
}

type oidcClientMetadata struct {
	CreatedOn   string  `json:"createdOn"`
	LastUsedOn  *string `json:"lastUsedOn"`
	Label       string  `json:"label"`
	Description string  `json:"description"`
	JWKSURI     *string `json:"jwksUri"`
}

type oidcClientLogs struct {
	Changes V3ListData `json:"changes"`
	Auth    V3ListData `json:"auth"`
}

type CreateOIDCClientRequest struct {
	Label            string  `json:"label,omitempty"`
	Description      string  `json:"description,omitempty"`
	JWKSURI          *string `json:"jwksUri,omitempty"`
	AccountDefault   bool    `json:"accountDefault"`
	EnforceAllowList bool    `json:"enforceAllowList"`
	TTL              int     `json:"ttl,omitempty"`
	DefaultAudience  string  `json:"defaultAudience,omitempty"`
}

type UpdateOIDCClientRequest struct {
	Label            string  `json:"label,omitempty"`
	Description      string  `json:"description,omitempty"`
	JWKSURI          *string `json:"jwksUri,omitempty"`
	AccountDefault   *bool   `json:"accountDefault,omitempty"`
	EnforceAllowList *bool   `json:"enforceAllowList,omitempty"`
	TTL              int     `json:"ttl,omitempty"`
	DefaultAudience  string  `json:"defaultAudience,omitempty"`
}

type oidcClientCreateResponse struct {
	ClientID int `json:"clientId"`
}

func (c *V3Client) CreateOIDCClient(req *CreateOIDCClientRequest) (int, error) {
	resp, err := c.post("/oidc/clients", req)
	if err != nil {
		return 0, fmt.Errorf("create OIDC client: %w", err)
	}
	var created oidcClientCreateResponse
	if err := json.Unmarshal(resp.Data, &created); err != nil {
		return 0, fmt.Errorf("create OIDC client unmarshal: %w", err)
	}
	return created.ClientID, nil
}

func (c *V3Client) GetOIDCClients() ([]OIDCClient, error) {
	return c.getOIDCClientsPage("/oidc/clients?limit=1000")
}

func (c *V3Client) getOIDCClientsPage(path string) ([]OIDCClient, error) {
	resp, err := c.get(path)
	if err != nil {
		return nil, fmt.Errorf("get OIDC clients: %w", err)
	}

	var wrapped oidcClientsResponse
	if err := json.Unmarshal(resp.Data, &wrapped); err != nil {
		return nil, fmt.Errorf("get OIDC clients unmarshal: %w", err)
	}

	var rows []oidcClientListRow
	if err := json.Unmarshal(wrapped.Clients.Data, &rows); err != nil {
		return nil, fmt.Errorf("get OIDC clients rows unmarshal: %w", err)
	}

	clients := make([]OIDCClient, 0, len(rows))
	for _, row := range rows {
		clients = append(clients, row.toClient(wrapped.Tenant))
	}

	meta := wrapped.Clients.Meta
	if meta.Total > 0 && meta.Limit > 0 && meta.Offset+meta.Limit < meta.Total {
		nextOffset := meta.Offset + meta.Limit
		nextPath, err := v3ListPagePath(path, nextOffset, meta.Limit)
		if err != nil {
			return nil, err
		}
		nextClients, err := c.getOIDCClientsPage(nextPath)
		if err != nil {
			return nil, err
		}
		clients = append(clients, nextClients...)
	}

	return clients, nil
}

func (c *V3Client) GetOIDCClient(clientID int) (*OIDCClient, error) {
	resp, err := c.get(fmt.Sprintf("/oidc/clients/%d", clientID))
	if err != nil {
		return nil, fmt.Errorf("get OIDC client %d: %w", clientID, err)
	}

	var detail oidcClientDetailResponse
	if err := json.Unmarshal(resp.Data, &detail); err != nil {
		return nil, fmt.Errorf("get OIDC client %d unmarshal: %w", clientID, err)
	}

	clients, err := c.GetOIDCClients()
	if err != nil {
		return nil, err
	}

	var client *OIDCClient
	for i := range clients {
		if clients[i].ClientID == clientID {
			client = &clients[i]
			break
		}
	}
	if client == nil {
		return nil, &V3NotFoundError{StatusCode: 404, Body: fmt.Sprintf("OIDC client %d not found in client list", clientID)}
	}

	client.CreatedOn = detail.Metadata.CreatedOn
	client.LastUsedOn = detail.Metadata.LastUsedOn
	client.Label = detail.Metadata.Label
	client.Description = detail.Metadata.Description
	client.JWKSURI = detail.Metadata.JWKSURI
	// Absent and malformed must not look alike here. An empty sub-list is normal, so a nil or
	// empty Data is skipped, but anything present that fails to unmarshal is returned rather
	// than discarded. ListStorageBlockVolumes swallowed exactly this and returned empty structs
	// with no error, which took a live run to notice.
	for _, sub := range []struct {
		name string
		raw  json.RawMessage
		into interface{}
	}{
		{"keys", detail.Keys.Data, &client.Keys},
		{"auth logs", detail.Logs.Auth.Data, &client.AuthLogs},
		{"change logs", detail.Logs.Changes.Data, &client.ChangeLogs},
	} {
		if len(sub.raw) == 0 {
			continue
		}
		if err := json.Unmarshal(sub.raw, sub.into); err != nil {
			return nil, fmt.Errorf("get OIDC client %d %s unmarshal: %w", clientID, sub.name, err)
		}
	}
	return client, nil
}

func (c *V3Client) UpdateOIDCClient(clientID int, req *UpdateOIDCClientRequest) error {
	_, err := c.patch(fmt.Sprintf("/oidc/clients/%d", clientID), req)
	if err != nil {
		return fmt.Errorf("update OIDC client %d: %w", clientID, err)
	}
	return nil
}

func (c *V3Client) DeleteOIDCClient(clientID int) error {
	_, err := c.del(fmt.Sprintf("/oidc/clients/%d", clientID))
	if err != nil {
		return fmt.Errorf("delete OIDC client %d: %w", clientID, err)
	}
	return nil
}

type OIDCClientKey struct {
	KeyID       int    `json:"keyId"`
	Label       string `json:"label"`
	Description string `json:"description"`
	ProvidedOn  string `json:"providedOn"`
	Type        string `json:"type"`
	Value       string `json:"value"`
	PublicKey   string `json:"publicKey"`
}

type CreateOIDCClientKeyRequest struct {
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
	PublicKey   string `json:"publicKey"`
}

type createOIDCClientKeysRequest struct {
	Keys []CreateOIDCClientKeyRequest `json:"keys"`
}

type createOIDCClientKeysResponse struct {
	Keys []OIDCClientKey `json:"keys"`
}

type UpdateOIDCClientKeyRequest struct {
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
}

func (c *V3Client) CreateOIDCClientKeys(clientID int, keys []CreateOIDCClientKeyRequest) ([]OIDCClientKey, error) {
	resp, err := c.post(fmt.Sprintf("/oidc/clients/%d/keys", clientID), &createOIDCClientKeysRequest{Keys: keys})
	if err != nil {
		return nil, fmt.Errorf("create OIDC client keys for client %d: %w", clientID, err)
	}
	var created createOIDCClientKeysResponse
	if err := json.Unmarshal(resp.Data, &created); err != nil {
		return nil, fmt.Errorf("create OIDC client keys unmarshal: %w", err)
	}
	return created.Keys, nil
}

func (c *V3Client) GetOIDCClientKeys(clientID int) ([]OIDCClientKey, error) {
	listData, err := c.getListUnder(fmt.Sprintf("/oidc/clients/%d/keys?limit=1000", clientID), "keys")
	if err != nil {
		return nil, fmt.Errorf("get OIDC client %d keys: %w", clientID, err)
	}
	var keys []OIDCClientKey
	if err := json.Unmarshal(listData.Data, &keys); err != nil {
		return nil, fmt.Errorf("get OIDC client %d keys unmarshal: %w", clientID, err)
	}
	return keys, nil
}

func (c *V3Client) UpdateOIDCClientKey(clientID, keyID int, req *UpdateOIDCClientKeyRequest) error {
	_, err := c.patch(fmt.Sprintf("/oidc/clients/%d/keys/%d", clientID, keyID), req)
	if err != nil {
		return fmt.Errorf("update OIDC client %d key %d: %w", clientID, keyID, err)
	}
	return nil
}

func (c *V3Client) DeleteOIDCClientKey(clientID, keyID int) error {
	_, err := c.del(fmt.Sprintf("/oidc/clients/%d/keys/%d", clientID, keyID))
	if err != nil {
		return fmt.Errorf("delete OIDC client %d key %d: %w", clientID, keyID, err)
	}
	return nil
}

type OIDCClientVM struct {
	MBPkgID int             `json:"mbpkgid"`
	Raw     json.RawMessage `json:"-"`
}

type oidcClientVMsResponse struct {
	VMs []json.RawMessage `json:"vms"`
}

type addOIDCClientVMsRequest struct {
	VMs []oidcClientVMRequest `json:"vms"`
}

type oidcClientVMRequest struct {
	MBPkgID int `json:"mbpkgid"`
}

func (c *V3Client) AddOIDCClientVMs(clientID int, mbpkgids []int) error {
	vms := make([]oidcClientVMRequest, len(mbpkgids))
	for i, mbpkgid := range mbpkgids {
		vms[i] = oidcClientVMRequest{MBPkgID: mbpkgid}
	}
	_, err := c.post(fmt.Sprintf("/oidc/clients/%d/allow-list/vms", clientID), &addOIDCClientVMsRequest{VMs: vms})
	if err != nil {
		return fmt.Errorf("add OIDC client %d VMs: %w", clientID, err)
	}
	return nil
}

func (c *V3Client) RemoveOIDCClientVM(clientID, mbpkgid int) error {
	_, err := c.del(fmt.Sprintf("/oidc/clients/%d/allow-list/vms/%d", clientID, mbpkgid))
	if err != nil {
		return fmt.Errorf("remove OIDC client %d VM %d: %w", clientID, mbpkgid, err)
	}
	return nil
}

func (c *V3Client) GetOIDCClientVMs(clientID int) ([]OIDCClientVM, error) {
	resp, err := c.get(fmt.Sprintf("/oidc/clients/%d/allow-list/vms", clientID))
	if err != nil {
		return nil, fmt.Errorf("get OIDC client %d VMs: %w", clientID, err)
	}
	var wrapped oidcClientVMsResponse
	if err := json.Unmarshal(resp.Data, &wrapped); err != nil {
		return nil, fmt.Errorf("get OIDC client %d VMs unmarshal: %w", clientID, err)
	}
	vms := make([]OIDCClientVM, 0, len(wrapped.VMs))
	for _, raw := range wrapped.VMs {
		vm := OIDCClientVM{Raw: raw}
		var asObject struct {
			MBPkgID int `json:"mbpkgid"`
		}
		if err := json.Unmarshal(raw, &asObject); err == nil {
			vm.MBPkgID = asObject.MBPkgID
		} else {
			var asInt int
			if err := json.Unmarshal(raw, &asInt); err == nil {
				vm.MBPkgID = asInt
			}
		}
		vms = append(vms, vm)
	}
	return vms, nil
}

type OIDCClientAuthLog struct {
	LogID     int    `json:"id"`
	IssuedOn  string `json:"issuedOn"`
	ExpiresOn string `json:"expiresOn"`
	JTI       string `json:"jti"`
}

func (c *V3Client) GetOIDCClientAuthLogs(clientID int) ([]OIDCClientAuthLog, error) {
	listData, err := c.getListUnder(fmt.Sprintf("/oidc/clients/%d/auth-logs?limit=1000", clientID), "logs")
	if err != nil {
		return nil, fmt.Errorf("get OIDC client %d auth logs: %w", clientID, err)
	}
	var logs []OIDCClientAuthLog
	if err := json.Unmarshal(listData.Data, &logs); err != nil {
		return nil, fmt.Errorf("get OIDC client %d auth logs unmarshal: %w", clientID, err)
	}
	return logs, nil
}

type OIDCClientChangeLog struct {
	KeyID      int    `json:"keyId"`
	RecordedOn string `json:"recordedOn"`
	Type       string `json:"type"`
}

func (c *V3Client) GetOIDCClientChangeLogs(clientID int) ([]OIDCClientChangeLog, error) {
	listData, err := c.getListUnder(fmt.Sprintf("/oidc/clients/%d/change-logs?limit=1000", clientID), "logs")
	if err != nil {
		return nil, fmt.Errorf("get OIDC client %d change logs: %w", clientID, err)
	}
	var logs []OIDCClientChangeLog
	if err := json.Unmarshal(listData.Data, &logs); err != nil {
		return nil, fmt.Errorf("get OIDC client %d change logs unmarshal: %w", clientID, err)
	}
	return logs, nil
}
