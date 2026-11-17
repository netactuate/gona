package gona

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

func idempotencyHeaders(key string) http.Header {
	if key == "" {
		return nil
	}
	return http.Header{"Idempotency-Key": []string{key}}
}

func credentialKey(key, path string) bool {
	switch strings.ToLower(key) {
	case "password", "root_password", "kubeconfig", "token", "access_key", "secret_key", "secret", "private_key", "api_key", "key", "idempotency-key":
		return true
	case "value":
		return strings.Contains(strings.ToLower(path), "secrets")
	}
	return false
}

// Redaction works on a copy so diagnostic logging cannot alter the wire payload.
func redactBody(path string, body []byte) string {
	var value interface{}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if json.Valid(body) && decoder.Decode(&value) == nil {
		maskJSON(path, value)
		encoded, err := json.Marshal(value)
		if err == nil {
			return string(encoded)
		}
	}
	// A kubeconfig or other credential document is not JSON, but its base64 padding makes it
	// parse as form data, which would log it whole. On those paths, anything that is not JSON
	// is withheld before the form parse is tried.
	lower := strings.ToLower(path)
	for _, sensitive := range []string{"kubeconfig", "secrets", "storage"} {
		if strings.Contains(lower, sensitive) {
			return "[body redacted]"
		}
	}
	// ParseQuery also accepts arbitrary text as a key; require an assignment before
	// treating a body as form data.
	if strings.Contains(string(body), "=") {
		if values, err := url.ParseQuery(string(body)); err == nil {
			for key := range values {
				if credentialKey(key, path) {
					values[key] = []string{"[redacted]"}
				}
			}
			return strings.ReplaceAll(values.Encode(), "%5Bredacted%5D", "[redacted]")
		}
	}
	return string(body)
}

func maskJSON(path string, value interface{}) {
	switch node := value.(type) {
	case map[string]interface{}:
		for key, child := range node {
			if credentialKey(key, path) {
				node[key] = "[redacted]"
			} else {
				maskJSON(path, child)
			}
		}
	case []interface{}:
		for _, child := range node {
			maskJSON(path, child)
		}
	}
}
