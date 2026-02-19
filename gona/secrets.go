package gona

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

type SecretList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SecretListValue struct {
	ID           int    `json:"id"`
	SecretListID int    `json:"secret_list_id"`
	SecretKey    string `json:"secret_key"`
	SecretValue  string `json:"secret_value"`
}

// Convert id str -> int to handle both
func (v *SecretListValue) UnmarshalJSON(data []byte) error {
	type Alias SecretListValue
	aux := &struct {
		SecretListID json.Number `json:"secret_list_id"`
		*Alias
	}{
		Alias: (*Alias)(v),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	id, err := strconv.Atoi(aux.SecretListID.String())
	if err != nil {
		return fmt.Errorf("invalid secret_list_id %q: %w", aux.SecretListID, err)
	}
	v.SecretListID = id
	return nil
}

func (c *Client) GetSecretLists() ([]SecretList, error) {
	var lists []SecretList
	if err := c.get("secrets/lists", &lists); err != nil {
		return nil, err
	}
	return lists, nil
}

// GetSecretList returns a specific secret list by ID
func (c *Client) GetSecretList(id int) (SecretList, error) {
	var list SecretList
	if err := c.get("secrets/lists/"+strconv.Itoa(id), &list); err != nil {
		return SecretList{}, err
	}
	return list, nil
}

func (c *Client) CreateSecretList(name string) (SecretList, error) {
	values := url.Values{}
	values.Add("name", name)

	var list SecretList
	if err := c.post("secrets/lists", []byte(values.Encode()), &list); err != nil {
		return SecretList{}, err
	}
	return list, nil
}

// UpdateSecretList updates a secret list name
func (c *Client) UpdateSecretList(id int, name string) (SecretList, error) {
	values := url.Values{}
	values.Add("name", name)

	var list SecretList
	if err := c.post("secrets/lists/"+strconv.Itoa(id), []byte(values.Encode()), &list); err != nil {
		return SecretList{}, err
	}
	return list, nil
}

func (c *Client) DeleteSecretList(id int) error {
	if err := c.delete("secrets/lists/"+strconv.Itoa(id), nil, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetSecretListValues(listID int) ([]SecretListValue, error) {
	var values []SecretListValue
	path := fmt.Sprintf("secrets/lists/%d/values", listID)
	if err := c.get(path, &values); err != nil {
		return nil, err
	}
	return values, nil
}

func (c *Client) GetSecretListValue(listID, valueID int) (SecretListValue, error) {
	var value SecretListValue
	path := fmt.Sprintf("secrets/lists/%d/values/%d", listID, valueID)
	if err := c.get(path, &value); err != nil {
		return SecretListValue{}, err
	}
	return value, nil
}

func (c *Client) CreateSecretListValue(listID int, key, val string) (SecretListValue, error) {
	formValues := url.Values{}
	formValues.Add("secret_key", key)
	formValues.Add("secret_value", val)

	var value SecretListValue
	path := fmt.Sprintf("secrets/lists/%d/values", listID)
	if err := c.post(path, []byte(formValues.Encode()), &value); err != nil {
		return SecretListValue{}, err
	}
	return value, nil
}

func (c *Client) UpdateSecretListValue(listID, valueID int, key, val string) (SecretListValue, error) {
	formValues := url.Values{}
	formValues.Add("secret_key", key)
	formValues.Add("secret_value", val)

	var value SecretListValue
	path := fmt.Sprintf("secrets/lists/%d/values/%d", listID, valueID)
	if err := c.post(path, []byte(formValues.Encode()), &value); err != nil {
		return SecretListValue{}, err
	}
	return value, nil
}

func (c *Client) DeleteSecretListValue(listID, valueID int) error {
	path := fmt.Sprintf("secrets/lists/%d/values/%d", listID, valueID)
	if err := c.delete(path, nil, nil); err != nil {
		return err
	}
	return nil
}
