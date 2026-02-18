package gona

import (
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
	if err := c.get("account/ssh_keys", &sshkeyList); err != nil {
		return nil, err
	}
	return sshkeyList, nil
}

// GetSSHKey will list the information on a specific key
func (c *Client) GetSSHKey(id int) (sshkey SSHKey, err error) {
	if err := c.get("account/ssh_key/"+strconv.Itoa(id), &sshkey); err != nil {
		return SSHKey{}, err
	}
	return sshkey, nil
}

// CreateSSHKey creates a key
func (c *Client) CreateSSHKey(name, key string) (sshkey SSHKey, err error) {
	values := url.Values{}
	values.Add("ssh_key", key)
	values.Add("name", name)

	if err := c.post("account/ssh_key", []byte(values.Encode()), &sshkey); err != nil {
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
	if err := c.patch("account/ssh_key/"+strconv.Itoa(id), jsonData, &sshkey); err != nil {
		return SSHKey{}, err
	}
	return sshkey, nil
}

// DeleteSSHKey deletes a key
func (c *Client) DeleteSSHKey(id int) error {
	if err := c.delete("account/ssh_key/"+strconv.Itoa(id), nil, nil); err != nil {
		return err
	}
	return nil
}
