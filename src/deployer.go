package src

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go.uber.org/zap"
)

type DockerDeployer struct {
	docker *client.Client
	logger *zap.Logger

	network string
}

func NewDeployer(docker *client.Client) *DockerDeployer {
	return &DockerDeployer{
		docker: docker,
		logger: DefaultLogger.With(
			zap.String("component", "DockerDeployer"),
		),
		network: Network,
	}
}

func (d *DockerDeployer) CreateNetwork(ctx context.Context) error {
	d.logger.Info("Creating network if not exists")
	networks, err := d.docker.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list networks: %w", err)
	}

	for _, nw := range networks {
		if nw.Name == d.network {
			d.logger.Info("Network already exists")
			return nil
		}
	}

	d.logger.Info("Creating network")
	if _, err := d.docker.NetworkCreate(ctx, d.network, network.CreateOptions{}); err != nil {
		return fmt.Errorf("failed to create network: %w", err)
	}
	return nil
}

func (d *DockerDeployer) DeployMysql(ctx context.Context) error {
	d.logger.Info("Deploying MySQL container")

	containers, err := d.docker.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	for _, c := range containers {
		if strings.Contains(c.Names[0], MysqlContainerName) {
			d.logger.Info("Mysql container already exists, skip deployment")
			return nil
		}
	}

	if err := d.pullImageIfNotExist(ctx, "mysql:latest"); err != nil {
		return fmt.Errorf("failed to pull MySQL image: %w", err)
	}

	d.logger.Info("Creating MySQL container")
	response, err := d.docker.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: "mysql:latest",
			ExposedPorts: nat.PortSet{
				"3306/tcp": struct{}{},
			},
			Env: []string{
				EncodeDockerEnv("MYSQL_ROOT_PASSWORD", "rootpwd"),
				EncodeDockerEnv("MYSQL_DATABASE", DbSchema),
				EncodeDockerEnv("MYSQL_USER", DbUser),
				EncodeDockerEnv("MYSQL_PASSWORD", DbPassword),
			},
		},
		nil,
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				d.network: {},
			},
		},
		nil,
		MysqlContainerName,
	)
	if err != nil {
		return fmt.Errorf("failed to create MySQL container: %w", err)
	}

	d.logger.Info("Starting MySQL container")
	if err := d.docker.ContainerStart(ctx, response.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start MySQL container: %w", err)
	}
	return nil
}

func (d *DockerDeployer) pullImageIfNotExist(ctx context.Context, imageName string) error {
	images, err := d.docker.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list images: %w", err)
	}

	found := false
	for _, img := range images {
		if len(img.RepoTags) > 0 && img.RepoTags[0] == imageName {
			found = true
			d.logger.Info("Mysql image already exists, skip pulling")
		}
	}

	if !found {
		d.logger.Info("Mysql image not exist, pulling")
		reader, err := d.docker.ImagePull(ctx, imageName, image.PullOptions{})
		if err != nil {
			return fmt.Errorf("failed to pull image %s: %w", imageName, err)
		}
		defer reader.Close()

		_, err = io.Copy(io.Discard, reader)
		return fmt.Errorf("failed to read image pull response: %w", err)
	}
	return nil
}
