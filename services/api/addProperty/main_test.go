package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/repositories/repository/utils"
	"github.com/Serares/undertown_v3/services/api/addProperty/handler"
)

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()
	mockProperty, err := os.ReadFile("testdata/postProperty.json")
	if err != nil {
		t.Error("error reading the mock property file")
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// create db connection
	dbUrl := utils.CreatePsqlUrl()
	db, err := repository.NewPropertiesRepo(dbUrl)
	if err != nil {
		log.Error("error on initializing the db")
	}
	hh := handler.New(log, db)

	ts := httptest.NewServer(hh)

	r, err := http.Post(ts.URL+"/", "application/json", bytes.NewBuffer(mockProperty))
	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode != http.StatusCreated {
		t.Fatalf("failed to create the mock property: %d", r.StatusCode)
	}
	return ts.URL, func() {
		ts.Close()
	}
}

// currently using psql with docker to test locally
func TestPost(t *testing.T) {
	mockProperty, err := os.ReadFile("testdata/postProperty.json")
	if err != nil {
		t.Error("error reading the mock property file")
	}
	// unmarshal the property to be able to change it's fields
	var parsedProperty handler.POSTProperty
	if err = json.NewDecoder(bytes.NewBuffer(mockProperty)).Decode(&parsedProperty); err != nil {
		t.Error("failed to decode the mocked json")
	}

	cases := []struct {
		name           string
		method         string
		URL            string
		expectedStatus int
		expectedError  string
		propertyTitle  string
	}{
		{
			name:           "Testing POST",
			method:         "POST",
			URL:            "/",
			expectedStatus: http.StatusCreated,
			propertyTitle:  "Testing property!",
			expectedError:  "",
		},
		{
			name:           "Testing GET",
			method:         "GET",
			URL:            "/",
			expectedStatus: http.StatusMethodNotAllowed,
			propertyTitle:  "This method is not allowed",
			expectedError:  handler.ErrorMethodNotSupported,
		},
	}

	url, cleanup := setupAPI(t)
	defer cleanup()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parsedProperty.Title = tc.propertyTitle
			var (
				body    bytes.Buffer
				errBody []byte
				err     error
			)
			// marshal the property to json
			if err := json.NewEncoder(&body).Encode(parsedProperty); err != nil {
				t.Fatal("error trying to marshal the request body")
			}

			req, err := http.NewRequest(tc.method, url, &body)
			defer req.Body.Close()
			req.Header.Set("Content-Type", "application/json")
			if err != nil {
				t.Error(err)
			}
			client := &http.Client{}
			resp, err := client.Do(req)
			defer resp.Body.Close()
			if err != nil {
				t.Fatal("Error sending request:", err)
			}

			if resp.StatusCode != tc.expectedStatus {
				t.Fatalf("Expected %q, got %q.", http.StatusText(tc.expectedStatus),
					http.StatusText(resp.StatusCode))
			}

			switch {
			case resp.Header.Get("Content-Type") == "application/json":
				if err = json.NewDecoder(resp.Body).Decode(&resp); err != nil {
					t.Error(err)
				}
				// TODO the success response is defined as struct SuccessCreateResponse
				// test agains that
			case strings.Contains(resp.Header.Get("Content-Type"), "text/plain"):
				if errBody, err = io.ReadAll(resp.Body); err != nil {
					t.Error(err)
				}
				if !strings.Contains(string(errBody), tc.expectedError) {
					t.Errorf("Expected %q, got %q.", tc.expectedError,
						string(errBody))
				}
			default:
				t.Fatalf("Unsupported Content-Type: %q", resp.Header.Get("Content-Type"))
			}
		})
	}
}
