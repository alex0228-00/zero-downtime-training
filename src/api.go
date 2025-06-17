package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Asset struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Source string `json:"source"`
}

type ApiClient struct {
	Host string
	Port string
}

func (c *ApiClient) url(path string) string {
	return fmt.Sprintf("http://%s:%s%s", c.Host, c.Port, path)
}

func (c *ApiClient) HealthCheck() error {
	resp, err := http.Get(c.url("/health"))
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %s", resp.Status)
	}
	return nil
}

func (c *ApiClient) CreateAsset(asset *Asset) error {
	body, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal asset: %w", err)
	}

	resp, err := http.Post(c.url("/api/asset"), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("create asset failed with status: %s", resp.Status)
	}
	return nil
}

func (c *ApiClient) ReadAsset(id string) (*Asset, error) {
	resp, err := http.Get(c.url("/api/asset/" + id))
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // asset not found
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get asset failed with status: %s", resp.Status)
	}

	var asset Asset
	err = json.NewDecoder(resp.Body).Decode(&asset)
	if err != nil {
		return nil, fmt.Errorf("failed to decode asset: %w", err)
	}

	return &asset, nil
}

func (c *ApiClient) DeleteAsset(id string) error {
	req, err := http.NewRequest(http.MethodDelete, c.url("/api/asset/"+id), nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete asset failed with status: %s", resp.Status)
	}
	return nil
}

func (c *ApiClient) UpdateAssetSourceByID(id, source string) error {
	// First get the existing asset
	asset, err := c.ReadAsset(id)
	if err != nil {
		return err
	}
	if asset == nil {
		return fmt.Errorf("asset with ID %s not found", id)
	}

	// Update source
	asset.Source = source

	body, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal updated asset: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, c.url("/api/asset/"+id), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create update request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update asset failed with status: %s", resp.Status)
	}
	return nil
}
