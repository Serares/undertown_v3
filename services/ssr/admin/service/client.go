package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository/lite"
)

type ListPropertiesResponse struct {
	Results []lite.Property `json:"Results"`
}

type GetPropertyResponse struct {
	Results []lite.Property `json:"Results"` // always one property even though it's returning an array
}

type SSRAdminClient struct {
	Log    *slog.Logger
	Client *http.Client
}

func NewAdminClient(log *slog.Logger) *SSRAdminClient {
	return &SSRAdminClient{
		Log: log,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (ssrc *SSRAdminClient) sendRequest(url, method, contentType, token string,
	expStatus int, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return []byte{}, err
	}
	if contentType == "" {
		contentType = "applicaiton/json"
	}
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", token)
	r, err := ssrc.Client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer r.Body.Close()
	msg, err := io.ReadAll(r.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("cannot read body: %w", err)
	}
	if r.StatusCode != expStatus {
		err = fmt.Errorf("error invalid response")
		if r.StatusCode == http.StatusNotFound {
			err = fmt.Errorf("not found")
		}
		return []byte{}, fmt.Errorf("%w: %s", err, msg)
	}

	return msg, nil
}

func (ssrc *SSRAdminClient) AddProperty(body io.Reader, url, authToken, contentType string) error {
	// the body should already come as multipart from the client
	_, err := ssrc.sendRequest(url, "POST", contentType, authToken, http.StatusOK, body)
	if err != nil {
		return err
	}
	// ssrc.Log.Debug("add property", slog.String("response", string(msg)))
	return nil
}

func (ssrc *SSRAdminClient) EditProperty(body io.Reader, url, authToken, contentType string) error {
	_, err := ssrc.sendRequest(url, "POST", contentType, authToken, http.StatusOK, body)
	if err != nil {
		return err
	}
	return nil
}

func (ssrc *SSRAdminClient) Login(email, password, url string) (string, error) {
	// creat a io.Reader body
	creds := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(creds); err != nil {
		return "", err
	}
	msg, err := ssrc.sendRequest(url, "POST", "application/json", "", http.StatusAccepted, &body)
	if err != nil {
		return "", err
	}
	loginResponse := struct {
		JWTToken string `json:"jwtToken,omitempty"` // this is returned in base64
		Error    string `json:"error,omitempty"`
	}{}
	if err := json.Unmarshal(msg, &loginResponse); err != nil {
		return "", err
	}
	if loginResponse.JWTToken != "" {
		return loginResponse.JWTToken, nil
	}
	return "", fmt.Errorf("no token")
}

func (ssrc *SSRAdminClient) List(url, authToken string) ([]lite.Property, error) {
	resp, err := ssrc.sendRequest(url, http.MethodGet, "", authToken, http.StatusOK, nil)
	if err != nil {
		return []lite.Property{}, err
	}
	var response ListPropertiesResponse

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return []lite.Property{}, fmt.Errorf("error deconding the listing properties response %w", err)
	}
	return response.Results, nil
}

func (ssrc *SSRAdminClient) GetProperty(url, humanReadableId, authToken string) (lite.Property, error) {
	resp, err := ssrc.sendRequest(url, http.MethodGet, "", authToken, http.StatusOK, nil)
	if err != nil {
		return lite.Property{}, err
	}
	var response GetPropertyResponse

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return lite.Property{}, fmt.Errorf("error deconding the listing properties response %w", err)
	}
	return response.Results[0], nil
}
