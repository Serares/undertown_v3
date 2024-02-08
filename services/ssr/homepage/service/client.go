package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository/lite"
	"github.com/Serares/undertown_v3/utils"
	"github.com/golang-jwt/jwt/v5"
)

type FeaturedPropertiesResponse struct {
	Results []lite.ListFeaturedPropertiesRow `json:"Results"`
}

type GetPropertiesResponse struct {
	Results []lite.Property `json:"Results"`
}

type GetPropertiesByTransactionTypeResponse struct {
	Results []lite.ListPropertiesByTransactionTypeRow `json:"Results"`
}

type GetPropertyResponse struct {
	// this should contain only one result
	Results []lite.Property `json:"Results"`
}

type SSRClient struct {
	Log    *slog.Logger
	Client *http.Client
}

func NewClient(log *slog.Logger) *SSRClient {
	return &SSRClient{
		Log: log,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// This should return the response body of the request
// ❗
// Have to create a jwt with the secret string to authorize SSR lambda to access the api gateway :)
// Other option is to get the SSRLambda arn and allow the request by checking the ARN in the authorizer
// But this is not possible because the API stack has to be created before the SSR stack
func (ssrc *SSRClient) sendRequest(url, method, contentType string,
	expStatus int, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return []byte{}, err
	}
	if contentType == "" {
		contentType = "applicaiton/json"
	}
	isLocal := os.Getenv("IS_LOCAL")
	if isLocal != "true" {
		token, err := ssrc.generateJwt()
		if err != nil {
			return []byte{}, err
		}
		req.Header.Set("Authorization", token)
	}
	req.Header.Set("Content-Type", contentType)
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

// ❗tokens are base64 encoded
// TODO check if the error of 'key is of invalid type: ECDSA sign expects *ecsda.PrivateKey' is still occuring
func (ssrc *SSRClient) generateJwt() (string, error) {
	claims := utils.JWTClaims{IsSsr: true}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("the jwt secret is empty")
	}
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error generating the jwt %w", err)
	}

	// base64Token := base64.RawStdEncoding.EncodeToString([]byte(tokenString))
	return tokenString, nil
}

func (ssrc *SSRClient) ListFeaturedProperties(url string) ([]lite.ListFeaturedPropertiesRow, error) {
	response, err := ssrc.sendRequest(url, http.MethodGet, "application/json", http.StatusOK, nil)
	if err != nil {
		ssrc.Log.Error("error requesting featured properties", "error", err, "url", url)
		return []lite.ListFeaturedPropertiesRow{}, fmt.Errorf("error trying to query the url: %s error", err.Error())
	}

	var resp FeaturedPropertiesResponse
	if err := json.Unmarshal(response, &resp); err != nil {
		return []lite.ListFeaturedPropertiesRow{}, err
	}

	return resp.Results, nil
}

func (ssrc *SSRClient) GetProperties(url string) ([]lite.Property, error) {
	r, err := ssrc.sendRequest(url, http.MethodGet, "application/json", http.StatusAccepted, nil)
	if err != nil {
		ssrc.Log.Error("error requesting property", "error", err)
		return []lite.Property{}, fmt.Errorf("error trying to query the url: %s error", err.Error())
	}

	var resp GetPropertyResponse
	if err := json.Unmarshal(r, &resp); err != nil {
		return []lite.Property{}, err
	}

	return resp.Results, nil
}

func (ssrc *SSRClient) GetPropertiesByTransactionType(url string) ([]lite.ListPropertiesByTransactionTypeRow, error) {
	r, err := ssrc.sendRequest(url, http.MethodGet, "application/json", http.StatusAccepted, nil)
	if err != nil {
		ssrc.Log.Error("error requesting properties by transaction type", "error", err, "url", url)
		return []lite.ListPropertiesByTransactionTypeRow{}, fmt.Errorf("error trying to query the url: %s error", err.Error())
	}

	var resp GetPropertiesByTransactionTypeResponse
	if err := json.Unmarshal(r, &resp); err != nil {
		return []lite.ListPropertiesByTransactionTypeRow{}, err
	}

	return resp.Results, nil
}

func (ssrc *SSRClient) GetProperty(url string) (lite.Property, error) {
	r, err := ssrc.sendRequest(url, http.MethodGet, "application/json", http.StatusAccepted, nil)
	if err != nil {
		ssrc.Log.Error("error requesting property", "error", err)
		return lite.Property{}, fmt.Errorf("error trying to query the url: %s error", err.Error())
	}

	var resp GetPropertyResponse
	if err := json.Unmarshal(r, &resp); err != nil {
		return lite.Property{}, err
	}

	return resp.Results[0], nil
}
