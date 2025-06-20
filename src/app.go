package src

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	HealthCheckN        = 5
	HealthCheckInterval = 5
)

type App struct {
	tag  string
	port string

	client *ApiClient

	docker      *client.Client
	containerID string
}

func NewApp(port, tag string, docker *client.Client) *App {
	return &App{
		port: port,
		tag:  tag,
		client: &ApiClient{
			Host: "localhost",
			Port: port,
		},
		docker: docker,
	}
}

func (app *App) Deploy() error {
	ctx := context.Background()

	if err := app.deployContainer(ctx); err != nil {
		return fmt.Errorf("failed to deploy container: %w", err)
	}

	if err := app.healthCheck(); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	return nil
}

func (app *App) Stop() error {
	if app.containerID == "" {
		return fmt.Errorf("no container to stop")
	}
	ctx := context.Background()
	if err := app.docker.ContainerStop(ctx, app.containerID, container.StopOptions{}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}
	if err := app.docker.ContainerRemove(ctx, app.containerID, container.RemoveOptions{}); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}
	app.containerID = ""
	return nil
}

func (app *App) healthCheck() error {
	for range HealthCheckN {
		if err := app.client.HealthCheck(); err == nil {
			return nil
		}
		time.Sleep(time.Second * HealthCheckInterval)
	}
	return fmt.Errorf("timeout for health check")
}

func (app *App) deployContainer(ctx context.Context) error {
	image := fmt.Sprintf("%s:%s", DockerImageName, app.tag)

	response, err := app.docker.ContainerCreate(
		context.Background(),
		&container.Config{
			ExposedPorts: nat.PortSet{
				"80/tcp": struct{}{},
			},
			Image: image,
			Env: []string{
				EncodeDockerEnv(EnvDBHost, MysqlContainerName),
				EncodeDockerEnv(EnvDBPort, "3306"),
				EncodeDockerEnv(EnvDBUser, DbUser),
				EncodeDockerEnv(EnvDBPassword, DbPassword),
				EncodeDockerEnv(EnvDBSchema, DbSchema),
				EncodeDockerEnv(EnvServerPort, "80"),
			},
		},
		&container.HostConfig{
			// AutoRemove: true,
			PortBindings: nat.PortMap{
				"80/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: app.port,
					},
				},
			},
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				Network: {},
			},
		},
		nil,
		"",
	)
	if err != nil {
		return fmt.Errorf("failed to create docker container: %w", err)
	}

	if err := app.docker.ContainerStart(ctx, response.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start docker container: %w", err)
	}

	app.containerID = response.ID
	return nil
}
