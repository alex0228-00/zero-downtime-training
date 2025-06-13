package src

import "os"

const (
	EnvDBHost     = "DB_HOST"
	EnvDBPort     = "DB_PORT"
	EnvDBUser     = "DB_USER_NAME"
	EnvDBPassword = "DB_USER_PASSWORD"

	EnvDBSchema   = "DB_SCHEMA"
	EnvServerPort = "SERVER_PORT"
)

func GetEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic("Environment variable " + key + " is not set")
	}
	return value
}

func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

type Asset struct {
	ID     string
	Name   string
	Source string
}

type AssetManager interface {
	CreateAsset(asset *Asset) error
	ReadAsset(id string) (*Asset, error)
	DeleteAsset(id string) error
	UpdateAssetSourceByID(id, name string) error
}
