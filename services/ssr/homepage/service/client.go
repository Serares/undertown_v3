package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/Serares/undertown_v3/repositories/repository/lite"
)

type FeaturedPropertiesResponse struct {
	Results []lite.ListFeaturedPropertiesRow `json:"Results"`
}

type SSRClient struct {
	Log    *slog.Logger
	Client *http.Client
}

func NewHomeClient(log *slog.Logger) *SSRClient {
	return &SSRClient{
		Log: log,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (hc *SSRClient) ListFeaturedProperties(url string) ([]lite.ListFeaturedPropertiesRow, error) {
	r, err := hc.Client.Get(url)
	if err != nil {
		hc.Log.Error("error requesting featured properties", "error", err)
		return nil, fmt.Errorf("error trying to query the url: %s error")
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			hc.Log.Error("error requesting featured properties", "error", err, "body", msg)
			return nil, fmt.Errorf("cannot read body: %w", err)
		}
	}
	var resp FeaturedPropertiesResponse
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}

	return resp.Results, nil
}
