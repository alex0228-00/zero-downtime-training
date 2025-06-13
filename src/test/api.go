package test

import (
	"fmt"
	"net/http"

	"zero-downtime-training/src"
)

type ApiClient struct {
	Host string
	Port string
}

func (c *ApiClient) HealthCheck() error {
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/health", c.Host, c.Port))
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %s", resp.Status)
	}
	defer resp.Body.Close()
	return nil
}

func (c *ApiClient) CreateAsset(asset *src.Asset) error {
	return nil
}

func (c *ApiClient) ReadAsset(id string) (*src.Asset, error) {
	return nil, nil
}

func (c *ApiClient) DeleteAsset(id string) error {
	return nil
}

func (c *ApiClient) UpdateAssetSourceByID(id, name string) error {
	return nil
}
