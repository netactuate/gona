package gona

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// mockAPIResponse creates an apiResponse with the given data
func mockAPIResponse(data any) apiResponse {
	return apiResponse{
		Result: "success",
		Data:   data,
		Code:   http.StatusOK,
	}
}

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.apiKey != apiKey {
		t.Errorf("apiKey = %q, want %q", client.apiKey, apiKey)
	}
	if client.endPoint.String() != BaseEndpoint {
		t.Errorf("endPoint = %q, want %q", client.endPoint.String(), BaseEndpoint)
	}
	if !strings.HasPrefix(client.userAgent, "gona/") {
		t.Errorf("userAgent = %q, want prefix 'gona/'", client.userAgent)
	}
}

func TestNewClientCustom(t *testing.T) {
	apiKey := "custom-api-key"
	apiURL := "https://custom.api.com/v1/"
	client := NewClientCustom(apiKey, apiURL)

	if client == nil {
		t.Fatal("NewClientCustom() returned nil")
	}
	if client.apiKey != apiKey {
		t.Errorf("apiKey = %q, want %q", client.apiKey, apiKey)
	}
	if client.endPoint.String() != apiURL {
		t.Errorf("endPoint = %q, want %q", client.endPoint.String(), apiURL)
	}
}

func TestGetKeyFromEnv(t *testing.T) {
	// Save original value
	original := os.Getenv("NA_API_KEY")
	defer os.Setenv("NA_API_KEY", original)

	tests := map[string]struct {
		envValue string
		want     string
	}{
		"key set":     {"test-key-123", "test-key-123"},
		"key empty":   {"", ""},
		"key not set": {"", ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.envValue == "" && name == "key not set" {
				os.Unsetenv("NA_API_KEY")
			} else {
				os.Setenv("NA_API_KEY", tt.envValue)
			}

			got := GetKeyFromEnv()
			if got != tt.want {
				t.Errorf("GetKeyFromEnv() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestApiKeyPath(t *testing.T) {
	tests := map[string]struct {
		path   string
		apiKey string
		want   string
	}{
		"simple path":            {"cloud/servers", "key123", "cloud/servers?key=key123"},
		"path with query":        {"cloud/servers?filter=active", "key123", "cloud/servers?filter=active&key=key123"},
		"key needs escaping":     {"cloud/servers", "key with spaces", "cloud/servers?key=key+with+spaces"},
		"special characters":     {"cloud/servers", "key&special=value", "cloud/servers?key=key%26special%3Dvalue"},
		"path with existing key": {"cloud/servers?existing=param", "mykey", "cloud/servers?existing=param&key=mykey"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := apiKeyPath(tt.path, tt.apiKey)
			if got != tt.want {
				t.Errorf("apiKeyPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if !r.URL.Query().Has("key") {
			t.Error("Expected API key in query parameters")
		}
		if r.Header.Get("Accept") != ContentType {
			t.Errorf("Accept header = %q, want %q", r.Header.Get("Accept"), ContentType)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(
			mockAPIResponse(map[string]string{"test": "data"}),
		)
	}))
	defer server.Close()

	client := NewClientCustom("test-key", server.URL+"/")
	ctx := context.Background()
	var result map[string]string

	err := client.get(ctx, "test/path", &result)
	if err != nil {
		t.Fatalf("get() error = %v", err)
	}
	if result["test"] != "data" {
		t.Errorf("result = %v, want map with test:data", result)
	}
}

func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("Content-Type = %q, want application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(apiResponse{
			Result: "success",
			Code:   http.StatusOK,
		})
	}))
	defer server.Close()

	client := NewClientCustom("test-key", server.URL+"/")
	ctx := context.Background()

	err := client.post(ctx, "test/path", []byte("param=value"), nil)
	if err != nil {
		t.Fatalf("post() error = %v", err)
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	tests := map[string]struct {
		statusCode int
		response   apiResponse
		wantErr    bool
	}{
		"success 200": {
			statusCode: http.StatusOK,
			response:   apiResponse{Result: "success", Code: http.StatusOK},
		},
		"error 500": {
			statusCode: http.StatusInternalServerError,
			response:   apiResponse{Result: "error", Code: http.StatusInternalServerError, Message: "Internal error"},
			wantErr:    true,
		},
		"error 422 with fields": {
			statusCode: 422,
			response:   apiResponse{Result: "error", Code: 422, Message: "Validation error", Fields: map[string]any{"name": "required"}},
			wantErr:    true,
		},
		"422 with mbpkgid ignored": {
			statusCode: 422,
			response:   apiResponse{Result: "error", Code: 422, Fields: map[string]any{"mbpkgid": "invalid"}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client := NewClientCustom("test-key", server.URL+"/")
			ctx := context.Background()

			err := client.get(ctx, "test/path", nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UserAgent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		if !strings.HasPrefix(userAgent, "gona/") {
			t.Errorf("User-Agent = %q, want prefix 'gona/'", userAgent)
		}
		if !strings.Contains(userAgent, Version) {
			t.Errorf("User-Agent = %q, want to contain version %q", userAgent, Version)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(apiResponse{Code: http.StatusOK})
	}))
	defer server.Close()

	client := NewClientCustom("test-key", server.URL+"/")
	ctx := context.Background()
	client.get(ctx, "test", nil)
}
