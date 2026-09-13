package gona

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

// SSHKey Struct 
type SSHKey struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"ssh_key"`
	Fingerprint string `json:"fingerprint"`
}

// GetSSHKeys will list all SSH Keys installed for the account
func (c *Client) GetSSHKeys() (keys []SSHKey, err error) {
	var sshkeyList []SSHKey
	if err := c.get(context.Background(), "account/ssh_keys", &sshkeyList); err != nil {
		return nil, err
	}
	return sshkeyList, nil
}

// GetSSHKey will list the information on a specific key
// GetSSHKey returns one SSH key.
//
// An SSH key that no longer exists does NOT 404 and does NOT 422. vAPI2 answers
// GET account/ssh_key/{id} with 200 and a null data member, verified live 2026-09-10
// against id 99999999:
//
//	{"result":"success","message":null,"meta":[],"data":null,"code":200}
//
// Unmarshalled into a struct that is silently the zero value, so without the check below
// this returns SSHKey{} and a nil error, and resourceSshKeyRead then writes an empty name
// and an empty key into Terraform state instead of removing the resource. That is B-08's
// silent state corruption, which is worse than a loud error because a plan afterwards looks
// plausible and proposes to "restore" a key that is gone.
//
// A real key always has a non zero id, so a zero id after a successful call means absent.
func (c *Client) GetSSHKey(id int) (sshkey SSHKey, err error) {
	if err := c.get(context.Background(), "account/ssh_key/"+strconv.Itoa(id), &sshkey); err != nil {
		return SSHKey{}, err
	}
	if sshkey.ID == 0 {
		return SSHKey{}, &NotFoundError{
			Method:     "GET",
			URL:        "account/ssh_key/" + strconv.Itoa(id),
			StatusCode: 200,
			Code:       200,
			Message:    "ssh key " + strconv.Itoa(id) + " returned a null body, so it does not exist",
		}
	}
	return sshkey, nil
}

// CreateSSHKey creates a key
func (c *Client) CreateSSHKey(name, key string) (sshkey SSHKey, err error) {
	values := url.Values{}
	values.Add("ssh_key", key)
	values.Add("name", name)

	if err := c.post(context.Background(), "account/ssh_key", []byte(values.Encode()), &sshkey); err != nil {
		return SSHKey{}, err
	}

	return sshkey, nil
}

// UpdateSSHKey updates a key's name and/or content
func (c *Client) UpdateSSHKey(id int, name, key string) (SSHKey, error) {
	body := map[string]string{
		"name":    name,
		"ssh_key": key,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return SSHKey{}, err
	}

	var sshkey SSHKey
	if err := c.patch(context.Background(), "account/ssh_key/"+strconv.Itoa(id), jsonData, &sshkey); err != nil {
		return SSHKey{}, err
	}
	return sshkey, nil
}

// DeleteSSHKey deletes a key
func (c *Client) DeleteSSHKey(id int) error {
	if err := c.delete(context.Background(), "account/ssh_key/"+strconv.Itoa(id), nil, nil); err != nil {
		return err
	}
	return nil
}
