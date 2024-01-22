package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Serares/undertown_v3/repositories/repository"
	"github.com/Serares/undertown_v3/services/api/addProperty/handler"
	"github.com/Serares/undertown_v3/services/api/addProperty/service"
	"github.com/Serares/undertown_v3/services/api/addProperty/types"
	"github.com/joho/godotenv"
)

func setupAPI(t *testing.T) (string, func()) {
	t.Helper()
	err := godotenv.Load(".env.local")
	if err != nil {
		t.Error("error loading the .env.local file")
	}
	mockImage, err := os.Open("testdata/mockImage.jpg")
	if err != nil {
		t.Error("error reading the mock property file")
	}
	defer mockImage.Close()
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := repository.NewPropertiesRepo()
	if err != nil {
		log.Error("error on initializing the db")
	}
	service := service.NewSubmitService(log, db)
	hh := handler.New(log, service)
	ts := httptest.NewServer(hh)

	titleField, err := multipartWriter.CreateFormField("title")
	if err != nil {
		t.Errorf("error adding title field to the form")
	}

	_, err = titleField.Write([]byte("mock_title"))
	if err != nil {
		t.Error("error adding mock value to title")
	}

	imageWriter, err := multipartWriter.CreateFormFile("images", mockImage.Name())
	if err != nil {
		t.Error("error trying to create the file for form")
	}
	_, err = io.Copy(imageWriter, mockImage)
	if err != nil {
		t.Error("error trying to copy the image to the form file writer")
	}

	multipartWriter.Close()
	request, err := http.NewRequest(http.MethodPost, ts.URL+"/", &requestBody)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Errorf("error sending the request %v", err)
	}
	defer response.Body.Close()
	return ts.URL, func() {
		err := db.CloseDbConnection(context.Background())
		if err != nil {
			t.Errorf("failed to close the db connection: %v", err)
		}
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
	var parsedProperty types.POSTProperty
	if err = json.NewDecoder(bytes.NewBuffer(mockProperty)).Decode(&parsedProperty); err != nil {
		t.Error("failed to decode the mocked json")
	}

	// cases := []struct {
	// 	name           string
	// 	method         string
	// 	URL            string
	// 	expectedStatus int
	// 	expectedError  string
	// 	propertyTitle  string
	// }{
	// 	{
	// 		name:           "Testing POST",
	// 		method:         "POST",
	// 		URL:            "/",
	// 		expectedStatus: http.StatusCreated,
	// 		propertyTitle:  "Testing property!",
	// 		expectedError:  "",
	// 	},
	// 	{
	// 		name:           "Testing GET",
	// 		method:         "GET",
	// 		URL:            "/",
	// 		expectedStatus: http.StatusMethodNotAllowed,
	// 		propertyTitle:  "This method is not allowed",
	// 		expectedError:  types.ErrorMethodNotSupported,
	// 	},
	// }

	_, cleanup := setupAPI(t)
	defer cleanup()

	// for _, tc := range cases {
	// 	t.Run(tc.name, func(t *testing.T) {
	// 		parsedProperty.Title = tc.propertyTitle
	// 		var (
	// 			body    bytes.Buffer
	// 			errBody []byte
	// 			err     error
	// 		)
	// 		// marshal the property to json
	// 		if err := json.NewEncoder(&body).Encode(parsedProperty); err != nil {
	// 			t.Fatal("error trying to marshal the request body")
	// 		}

	// 		req, err := http.NewRequest(tc.method, url, &body)
	// 		defer req.Body.Close()
	// 		req.Header.Set("Content-Type", "application/json")
	// 		if err != nil {
	// 			t.Error(err)
	// 		}
	// 		client := &http.Client{}
	// 		resp, err := client.Do(req)
	// 		defer resp.Body.Close()
	// 		if err != nil {
	// 			t.Fatal("Error sending request:", err)
	// 		}

	// 		if resp.StatusCode != tc.expectedStatus {
	// 			t.Fatalf("Expected %q, got %q.", http.StatusText(tc.expectedStatus),
	// 				http.StatusText(resp.StatusCode))
	// 		}

	// 		switch {
	// 		case resp.Header.Get("Content-Type") == "application/json":
	// 			if err = json.NewDecoder(resp.Body).Decode(&resp); err != nil {
	// 				t.Error(err)
	// 			}
	// 			// TODO the success response is defined as struct SuccessCreateResponse
	// 			// test agains that
	// 		case strings.Contains(resp.Header.Get("Content-Type"), "text/plain"):
	// 			if errBody, err = io.ReadAll(resp.Body); err != nil {
	// 				t.Error(err)
	// 			}
	// 			if !strings.Contains(string(errBody), tc.expectedError) {
	// 				t.Errorf("Expected %q, got %q.", tc.expectedError,
	// 					string(errBody))
	// 			}
	// 		default:
	// 			t.Fatalf("Unsupported Content-Type: %q", resp.Header.Get("Content-Type"))
	// 		}
	// 	})
	// }
}
