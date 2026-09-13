// Package gona provides a simple golang interface to the NetActuate
// Rest API at https://vapi2.netactuate.com/
package gona

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Version, BaseEndpoint, ContentType constants
const (
	Version      = "0.4.0"
	BaseEndpoint = "https://vapi2.netactuate.com/api/"
	ContentType  = "application/json"
)

// Client is the main object (struct) to which we attach most
// methods/functions.
// It has the following fields:
// (client, userAgent, endPoint, apiKey)
type Client struct {
	client    *http.Client
	userAgent string
	endPoint  *url.URL
	apiKey    string
}

// GetKeyFromEnv is a simple function to grab the value for
// "NA_API_KEY" from the environment
func GetKeyFromEnv() string {
	return os.Getenv("NA_API_KEY")
}

// NewClientCustom is the main entrypoint for instantiating a Client struct.
// It takes your API Key as it's sole argument
// and returns the Client struct ready to talk to the API
func NewClientCustom(apikey string, apiurl string) *Client {
	useragent := "gona/" + Version
	client := &http.Client{
		Transport: &http.Transport{
			TLSNextProto: make(
				map[string]func(string, *tls.Conn) http.RoundTripper,
			),
		},
	}
	endpoint, err := url.Parse(apiurl)
	if err != nil {
		panic(fmt.Sprintf("invalid API URL: %v", err))
	}

	return &Client{
		userAgent: useragent,
		client:    client,
		endPoint:  endpoint,
		apiKey:    apikey,
	}
}

// NewClient takes an apikey and calls NewClientCustom with the hardcoded
// BaseEndpoint constant API URL
func NewClient(apikey string) *Client {
	return NewClientCustom(apikey, BaseEndpoint)
}

// apiKeyPath is just a short internal function for appending the key to the url
func apiKeyPath(path, apiKey string) string {
	if strings.Contains(path, "?") {
		return path + "&key=" + apiKey
	}
	return path + "?key=" + apiKey
}

func (c *Client) debugLog(format string, v ...interface{}) {
	if os.Getenv("NA_API_DEBUG") == "" {
		return
	}
	log.Printf("[DEBUG] "+format, v...)
}

// get internal method on Client struct for providing the HTTP GET call
func (c *Client) get(ctx context.Context, path string, data interface{}) error {
	req, err := c.newRequest(ctx, "GET", path, nil)
	if err != nil {
		return err
	}
	return c.do(req, data)
}

// post internal method on Client struct for providing the HTTP POST call
func (c *Client) post(ctx context.Context, path string, values []byte, data interface{}) error {
	c.debugLog("POST data for %s: %s", path, string(values))

	req, err := c.newRequest(ctx, "POST", path, bytes.NewBuffer(values))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req, data)
}

func (c *Client) patch(ctx context.Context, path string, values []byte, data interface{}) error {
	c.debugLog("PATCH data for %s: %s", path, string(values))

	req, err := c.newRequest(ctx, "PATCH", path, bytes.NewBuffer(values))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.do(req, data)
}

func (c *Client) put(ctx context.Context, path string, values []byte, data interface{}) error {
	c.debugLog("PUT data for %s: %s", path, string(values))

	req, err := c.newRequest(ctx, "PUT", path, bytes.NewBuffer(values))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.do(req, data)
}

func (c *Client) postJSON(ctx context.Context, path string, values []byte, data interface{}) error {
	c.debugLog("POST JSON data for %s: %s", path, string(values))

	req, err := c.newRequest(ctx, "POST", path, bytes.NewBuffer(values))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.do(req, data)
}

func (c *Client) putJSON(ctx context.Context, path string, values []byte, data interface{}) error {
	c.debugLog("PUT JSON data for %s: %s", path, string(values))

	req, err := c.newRequest(ctx, "PUT", path, bytes.NewBuffer(values))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	return c.do(req, data)
}

// delete internal method on Client struct for providing the HTTP DELETE call
func (c *Client) delete(ctx context.Context, path string, values url.Values, data interface{}) error {
	req, err := c.newRequest(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}
	return c.do(req, data)
}

// Two functions (newRequest, do) below are used by the http method name functions above
// newRequest internal method on Client struct to be wrapped inside the above http method
// named functions for doing the actual work of the get/post/put/patch/delete methods
func (c *Client) newRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error) {
	relPath, err := url.Parse(apiKeyPath(path, c.apiKey))

	if err != nil {
		return nil, err

	}

	url := c.endPoint.ResolveReference(relPath)

	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return nil, err

	}

	req.Header.Add("User-Agent", c.userAgent)
	req.Header.Add("Accept", ContentType)

	c.debugLog("making a %s request to %s", method, url)
	return req, nil
}

// apiResponse is a message returned by the API that is used both for successful
// responses and for some error responses.
type apiResponse struct {
	Result  string                 `json:"result"`
	Message string                 `json:"message"`
	Data    json.RawMessage        `json:"data"`
	Code    int                    `json:"code"`
	Fields  map[string]interface{} `json:"fields"`
}

type NotFoundError struct {
	Method     string
	URL        string
	StatusCode int
	Code       int
	Message    string
	Body       string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("not found on %s %s: code %d / %d, response: %s / %s", e.Method, e.URL, e.StatusCode, e.Code, e.Message, e.Body)
}

// redactURL removes the API key from a URL before it reaches an error message or a log.
//
// The V2 API takes the key as a query parameter, so req.URL carries a live credential.
// Every error below is returned to Terraform, which prints it to the operator's terminal,
// their CI logs and any support ticket or screenshot that follows. Marking the provider's
// api_key field Sensitive does not help: that governs plan and state rendering, not the
// text of an error. So the redaction has to happen here, where the URL is stringified.
func redactURL(u fmt.Stringer) string {
	if u == nil {
		return ""
	}
	s := u.String()
	parsed, err := url.Parse(s)
	if err != nil {
		// Cannot parse it, so cannot prove it is safe. Say nothing rather than guess.
		return "[url redacted]"
	}
	q := parsed.Query()
	if q.Get("key") == "" {
		return s
	}
	q.Set("key", "REDACTED")
	parsed.RawQuery = q.Encode()
	return parsed.String()
}

func IsNotFound(err error) bool {
	var notFound *NotFoundError
	return errors.As(err, &notFound)
}

// do internal method on Client struct for making the HTTP calls
func (c *Client) do(req *http.Request, data interface{}) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	c.debugLog("got a response: %s", string(body))

	r := &apiResponse{}
	if err := json.Unmarshal(body, r); err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return &NotFoundError{
				Method:     req.Method,
				URL:        redactURL(req.URL),
				StatusCode: resp.StatusCode,
				Body:       string(body),
			}
		}
		return fmt.Errorf("could not unmarshal response %q: %w", string(body), err)
	}

	if resp.StatusCode == http.StatusNotFound || r.Code == http.StatusNotFound {
		return &NotFoundError{
			Method:     req.Method,
			URL:        redactURL(req.URL),
			StatusCode: resp.StatusCode,
			Code:       r.Code,
			Message:    r.Message,
			Body:       string(r.Data),
		}
	}

	// Deleting a BGP session that is already gone (a retry, or one cleared out
	// of band) 422s with this message instead of a real error -- treat it as
	// success so DeleteBGPSession is idempotent.
	if (resp.StatusCode == 422 || r.Code == 422) && r.Fields != nil {
		if msgs, ok := r.Fields["id"].([]interface{}); ok {
			for _, msg := range msgs {
				if str, _ := msg.(string); str == "The bgp id could not be found" {
					return nil
				}
			}
		}
	}

	// A DNS zone or record that no longer exists 422s with a field error rather than
	// returning 404. Verified live 2026-09-10: GET /dns/zone/{id} on a zone deleted
	// moments earlier returns
	//   422 {"fields":{"id":["The id must be a valid zone id"]}}
	//
	// Without this, a zone deleted out of band produces a hard error on every refresh
	// and the resource can never be reconciled or removed. That is exactly the B-01
	// defect, and it would have shipped in brand new code, so it is caught here rather
	// than left for a customer to find.
	//
	// 2026-09-10: the SAME idiom appears on a THIRD field. A firewall set that no longer
	// exists answers GET /firewall/sets/{id} with
	//   422 {"fields":{"firewall_set_id":["The firewall set must be a valid"]}}
	// Verified live against both a deleted set and an id that never existed. Without this,
	// a firewall set deleted out of band hard errors on every refresh, which is B-01 again
	// on a security resource. Found by the W12 acceptance suite when its own CheckDestroy
	// could not tell "gone" from "broken".
	//
	// The field name differs per resource, so the check is keyed on the field AND the
	// message, never on the message alone. A blanket "must be a valid" match would swallow
	// genuine validation errors on create and turn a bad request into a silent no-op.
	if (resp.StatusCode == 422 || r.Code == 422) && r.Fields != nil {
		// One field can carry several not-found messages, so this is field to messages.
		// Keyed on the pair, never on the message alone: a blanket "must be a valid"
		// match would swallow genuine validation errors on create and turn a bad request
		// into a silent no-op. Every entry below was verified live on 2026-09-10.
		for field, wants := range map[string][]string{
			"id": {
				"must be a valid zone id",   // GET dns/zone/{id}
				"must be a valid record id", // GET dns/record/{id}
				"must be a valid Image ID",  // GET cloud/images/{id}
			},
			"firewall_set_id": {
				"The firewall set must be a valid", // GET firewall/sets/{id}
			},
			"secret_list_id": {
				"The secret list id must be a valid", // GET secrets/lists/{id}
			},
			"secret_list_value_id": {
				"The secret list value id must be a valid", // GET secrets/lists/{id}/values/{id}
			},
		} {
			msgs, ok := r.Fields[field].([]interface{})
			if !ok {
				continue
			}
			for _, msg := range msgs {
				str, _ := msg.(string)
				for _, want := range wants {
					if strings.Contains(str, want) {
						return &NotFoundError{
							Method:     req.Method,
							URL:        redactURL(req.URL),
							StatusCode: resp.StatusCode,
							Code:       r.Code,
							Message:    str,
							Body:       string(body),
						}
					}
				}
			}
		}
	}

	// B-08. A server that no longer exists answers GET cloud/server?mbpkgid=N with
	//   422 {"fields":{"mbpkgid":["The mbpkgid must be a valid mbpkgid"]}}
	// The blanket mbpkgid swallow below then returns nil AND a zero valued Server, so
	// resourceServerRead hydrates state with empty strings and zeros instead of removing
	// the resource. That is silent state corruption, and it is worse than the loud error
	// B-01 produced. Confirmed live 2026-09-10 against a server destroyed minutes before.
	//
	// Narrow, deliberately: only the "not a valid mbpkgid" shape becomes a not-found.
	// Any other mbpkgid 422 keeps the historical swallow, because that swallow was added
	// to make the provider work and removing it wholesale would be a blind change.
	if (resp.StatusCode == 422 || r.Code == 422) && r.Fields != nil {
		if msgs, ok := r.Fields["mbpkgid"].([]interface{}); ok {
			for _, msg := range msgs {
				str, _ := msg.(string)
				if strings.Contains(str, "must be a valid mbpkgid") ||
					strings.Contains(str, "must be a valid package") {
					return &NotFoundError{
						Method:     req.Method,
						URL:        redactURL(req.URL),
						StatusCode: resp.StatusCode,
						Code:       r.Code,
						Message:    str,
						Body:       string(body),
					}
				}
			}
		}
	}

	// Error Handling - This currently ignores invalid mbpkdgid errors to enable the Terraform Provider
	if (resp.StatusCode == 422 || r.Code == 422) && (r.Fields != nil && r.Fields["mbpkgid"] == nil) {
		fieldStr := ""
		for key, value := range r.Fields {
			fieldStr = fieldStr + fmt.Sprintf("%s: %v, ", key, value)
		}
		return fmt.Errorf("got an ERROR response on %s %s: code %d / %d, response: %s / %s", req.Method, redactURL(req.URL), resp.StatusCode, r.Code, r.Message, fieldStr)
	}

	if (resp.StatusCode != http.StatusOK && resp.StatusCode != 422) || (r.Code != http.StatusOK && r.Code != 422) {
		return fmt.Errorf("got an error response on %s %s: code %d / %d, response: %s / %s", req.Method, redactURL(req.URL), resp.StatusCode, r.Code, r.Message, string(r.Data))
	}

	// Unmarshal the data field into the caller's typed struct only on success
	if data != nil && len(r.Data) > 0 {
		if err := json.Unmarshal(r.Data, data); err != nil {
			// Do NOT put the raw body in the error. A dedicated server response carries
			// ipmi_cxuser and ipmi_cxpass, and its build debug blob carries the BMC
			// password in clear text, so an unmarshal failure was printing live IPMI
			// credentials into terraform output, CI logs and anything scraping them.
			// Observed 2026-09-10 on GET dedicated/servers/{id}. Same class as the API
			// key leak, and worse, because a BMC is out of band access to the machine.
			return fmt.Errorf("could not unmarshal response data (%d bytes, redacted): %w",
				len(r.Data), err)
		}
	}

	return nil
}
