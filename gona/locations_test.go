package gona

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetLocations(t *testing.T) {
	tests := map[string]struct {
		mockResponse   any
		mockStatusCode int
		wantErr        bool
		wantLocations  []Location
	}{
		"successful response with multiple locations": {
			mockResponse: []Location{
				{
					ID:        1,
					Name:      "Los Angeles",
					IATACode:  "LAX",
					Continent: "North America",
					Flag:      "us",
				},
				{
					ID:        2,
					Name:      "Amsterdam",
					IATACode:  "AMS",
					Continent: "Europe",
					Flag:      "nl",
				},
			},
			mockStatusCode: http.StatusOK,
			wantLocations: []Location{
				{
					ID:        1,
					Name:      "Los Angeles",
					IATACode:  "LAX",
					Continent: "North America",
					Flag:      "us",
				},
				{
					ID:        2,
					Name:      "Amsterdam",
					IATACode:  "AMS",
					Continent: "Europe",
					Flag:      "nl",
				},
			},
		},
		"empty locations list": {
			mockStatusCode: http.StatusOK,
		},
		// Currently is not filtered.
		"disabled location": {
			mockResponse: []Location{
				{
					ID:        3,
					Name:      "Tokyo",
					IATACode:  "TYO",
					Continent: "Asia",
					Flag:      "jp",
					Disabled:  1,
				},
			},
			mockStatusCode: http.StatusOK,
			wantLocations: []Location{
				{
					ID:        3,
					Name:      "Tokyo",
					IATACode:  "TYO",
					Continent: "Asia",
					Flag:      "jp",
					Disabled:  1,
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the request path
				if r.URL.Path != "/api/cloud/locations" {
					t.Errorf("Expected path /api/cloud/locations, got %s", r.URL.Path)
				}

				// Verify API key is present
				if r.URL.Query().Get("key") == "" {
					t.Error("Expected API key in query parameters")
				}

				// Send mock response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				resp := mockAPIResponse(tt.mockResponse)
				json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			// Create client with mock server URL
			client := NewClientCustom("test-api-key", server.URL+"/api/")

			// Call GetLocations
			ctx := context.Background()
			locations, err := client.GetLocations(ctx)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLocations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check result
			if !tt.wantErr {
				if len(locations) != len(tt.wantLocations) {
					t.Errorf("GetLocations() returned %d locations, want %d", len(locations), len(tt.wantLocations))
					return
				}

				for i, loc := range locations {
					t.Run(fmt.Sprint("location", i), func(t *testing.T) {
						checkLocation(t, loc, tt.wantLocations[i])
					})
				}
			}
		})
	}
}

func TestGetLocations_ErrorResponse(t *testing.T) {
	// Create mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := apiResponse{
			Result:  "error",
			Message: "Internal server error",
			Data:    nil,
			Code:    http.StatusInternalServerError,
			Fields:  nil,
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClientCustom("test-api-key", server.URL+"/api/")

	// Call GetLocations
	ctx := context.Background()
	locations, err := client.GetLocations(ctx)

	// Should return an error
	if err == nil {
		t.Error("GetLocations() expected error for 500 response, got nil")
	}

	// Locations should be nil on error
	if locations != nil {
		t.Errorf("GetLocations() expected nil locations on error, got %v", locations)
	}
}

func TestLocation_JSONSerialization(t *testing.T) {
	location := Location{
		ID:        1,
		Name:      "Los Angeles",
		IATACode:  "LAX",
		Continent: "North America",
		Flag:      "us",
	}

	// Test marshaling
	data, err := json.Marshal(location)
	if err != nil {
		t.Fatalf("Failed to marshal Location: %v", err)
	}

	// Test unmarshaling
	var decoded Location
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal Location: %v", err)
	}
	checkLocation(t, decoded, location)
}

// Verify all fields
func checkLocation(t *testing.T, actual, wanted Location) {
	if actual.ID != wanted.ID {
		t.Errorf("ID = %d, want %d", actual.ID, wanted.ID)
	}
	if actual.Name != wanted.Name {
		t.Errorf("Name = %s, want %s", actual.Name, wanted.Name)
	}
	if actual.IATACode != wanted.IATACode {
		t.Errorf("IATACode = %s, want %s", actual.IATACode, wanted.IATACode)
	}
	if actual.Continent != wanted.Continent {
		t.Errorf("Continent = %s, want %s", actual.Continent, wanted.Continent)
	}
	if actual.Flag != wanted.Flag {
		t.Errorf("Flag = %s, want %s", actual.Flag, wanted.Flag)
	}
	if actual.Disabled != wanted.Disabled {
		t.Errorf("Disabled = %d, want %d", actual.Disabled, wanted.Disabled)
	}
}

func TestGetLocations_QueryParameterEncoding(t *testing.T) {
	var capturedURL *url.URL

	// Create mock server that captures the request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := mockAPIResponse([]Location{})
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClientCustom("my-test-key", server.URL+"/api/")
	ctx := context.Background()

	// Call GetLocations
	_, err := client.GetLocations(ctx)
	if err != nil {
		t.Fatalf("GetLocations() error = %v", err)
	}

	// Verify the API key was properly encoded in the URL
	if capturedURL == nil {
		t.Fatal("URL was not captured")
	}

	apiKey := capturedURL.Query().Get("key")
	if apiKey != "my-test-key" {
		t.Errorf("API key = %s, want my-test-key", apiKey)
	}
}
