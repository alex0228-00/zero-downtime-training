package src

import (
	"fmt"
	"os"
)

const (
	EnvDBHost     = "DB_HOST"
	EnvDBPort     = "DB_PORT"
	EnvDBUser     = "DB_USER_NAME"
	EnvDBPassword = "DB_USER_PASSWORD"

	EnvDBSchema   = "DB_SCHEMA"
	EnvServerPort = "SERVER_PORT"

	DockerImageName    = "zero-downtime-training"
	Network            = "zero-downtime-training"
	MysqlContainerName = "zero-downtime-training-mysql"

	DbUser     = "testuser"
	DbPassword = "testpassword"
	DbSchema   = "zero-downtime-training"
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

func EncodeDockerEnv(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}
